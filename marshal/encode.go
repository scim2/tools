package marshal

import (
	"errors"
	"fmt"
	"github.com/scim2/tools/structs"
	"reflect"
)

var marshalerType = reflect.TypeOf((*Marshaler)(nil)).Elem()

func Marshal(value interface{}) (structs.Resource, error) {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return nil, errors.New("value is invalid")
	}

	t := v.Type()
	if t.Implements(marshalerType) {
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return nil, errors.New("ptr is nil")
		}
		m, ok := v.Interface().(Marshaler)
		if !ok {
			return nil, errors.New("value does not implement marshaler")
		}
		return m.MarshalSCIM()
	}

	switch t.Kind() {
	case reflect.Interface:
		if v.IsNil() {
			return nil, errors.New("interface is nil")
		}
		return Marshal(v.Elem().Interface())
	case reflect.Ptr:
		ptr := v.Elem()
		return Marshal(ptr.Elem().Interface())
	case reflect.Struct:
		resource := make(structs.Resource)

		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			tag := parseTags(t.Field(i))
			if tag.ignore {
				continue
			}

			field := v.Field(i)
			if !tag.allowZero && field.IsZero() {
				continue
			}
			if err := structEncoder(resource, field, tag); err != nil {
				return nil, err
			}
		}
		return resource, nil
	default:
		return unsupportedTypeEncoder(v)
	}
}

func structEncoder(resource structs.Resource, field reflect.Value, tag tag) error {
	if tag.sub == nil {
		if tag.multiValued {
			return structEncoderSimpleMultiValued(resource, field, tag)
		}
		return structEncoderSimple(resource, field, tag)
	} else {
		if tag.multiValued {
			return structEncoderComplexMultiValued(resource, field, tag)
		}
		return structEncoderComplex(resource, field, tag)
	}
}

func structEncoderComplex(resource structs.Resource, field reflect.Value, tag tag) error {
	subResource := resource.EnsureComplexAttribute(tag.name)
	if subResource.Exists(tag.sub.name) {
		return errors.New(fmt.Sprintf("duplicate names: %s", tag.sub.name))
	}
	return structEncoderSimple(subResource, field, *tag.sub)
}

func structEncoderComplexMultiValued(resource structs.Resource, field reflect.Value, tag tag) error {
	switch field.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			value := make(structs.Resource)
			if err := structEncoder(value, field.Index(i), *tag.sub); err != nil {
				return err
			}
			resource.EnsureComplexMultiValuedAttribute(tag.name, 0)
			if err := resource.AppendComplexMultiValuedAttribute(tag.name, value)
				err != nil {
				return err
			}
		}
	case reflect.Ptr, reflect.Interface:
		return structEncoderComplexMultiValued(resource, field.Elem(), tag)
	default:
		value := make(structs.Resource)
		if err := structEncoder(value, field, *tag.sub); err != nil {
			return err
		}
		resource.EnsureComplexMultiValuedAttribute(tag.name, tag.max())
		if err := resource.AppendComplexMultiValuedAttribute(tag.name, value); err != nil {
			return err
		}
	}
	return nil
}

func structEncoderSimple(resource structs.Resource, field reflect.Value, tag tag) error {
	// Ignore invalid fields.
	if !field.IsValid() {
		return nil
	}

	switch field.Kind() {
	// If the simple attribute is a map that means that it is in fact a complex attribute where the name is implicit.
	case reflect.Map:
		t := field.Type()
		if t.Key().Kind() != reflect.String {
			return errors.New("key of map is not a string")
		}

		mapField, err := resource.AddEmptyComplexAttribute(tag.name)
		if err != nil {
			return err
		}

		for _, k := range field.MapKeys() {
			value := field.MapIndex(k)

			// If the value is an interface or ptr, use the underlying element.
			for value.Kind() == reflect.Interface ||
				value.Kind() == reflect.Ptr {
				value = value.Elem()
			}

			fieldInterface, err := validSimpleAttribute(value)
			if err != nil {
				return err
			}
			if err := mapField.Add(k.String(), fieldInterface); err != nil {
				return err
			}
		}
	case reflect.Ptr, reflect.Interface:
		return structEncoderSimple(resource, field.Elem(), tag)
	case reflect.Struct:
		fieldStruct := make(structs.Resource)
		t := field.Type()
		for i := 0; i < field.NumField(); i++ {
			tagIndex := parseTags(t.Field(i))
			if tagIndex.ignore {
				continue
			}

			fieldIndex := field.Field(i)
			if !tagIndex.allowZero && fieldIndex.IsZero() {
				continue
			}
			if err := structEncoder(fieldStruct, fieldIndex, tagIndex); err != nil {
				return err
			}
		}
		if depth := fieldStruct.Depth(); 1 < depth {
			return fmt.Errorf("nested depth exceeded: %d", depth)
		}

		fieldMap := resource.EnsureComplexAttribute(tag.name)
		for k, v := range fieldStruct {
			if err := fieldMap.Add(k, v); err != nil {
				return err
			}
		}
	case reflect.Array, reflect.Slice:
		// Simple attributes can never be an array or a slice.
		return errors.New(fmt.Sprintf("invalid simple attribute: %s", field.Kind()))
	default:
		fieldInterface, err := validSimpleAttribute(field)
		if err != nil {
			return err
		}
		if err := resource.Add(tag.name, fieldInterface); err != nil {
			return err
		}
	}
	return nil
}

func structEncoderSimpleMultiValued(resource structs.Resource, field reflect.Value, tag tag) error {
	switch field.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < field.Len(); i++ {
			value := make(structs.Resource)
			if err := structEncoderSimple(value, field.Index(i), tag); err != nil {
				return err
			}
			for _, v := range value {
				switch field.Index(i).Kind() {
				case reflect.Struct:
					resource.EnsureComplexMultiValuedAttribute(tag.name, tag.max())
					if err := resource.AppendComplexMultiValuedAttribute(tag.name, v.(structs.Resource)); err != nil {
						return err
					}
				default:
					resource.EnsureMultiValuedAttribute(tag.name, tag.max())
					if err := resource.AppendMultiValuedAttribute(tag.name, v); err != nil {
						return err
					}
				}
			}
		}
	case reflect.Map:
		resource.EnsureComplexMultiValuedAttribute(tag.name, tag.max())
		value := make(structs.Resource)
		if err := structEncoderSimple(value, field, tag); err != nil {
			return err
		}
		for _, v := range value {
			if err := resource.AppendComplexMultiValuedAttribute(tag.name, v.(structs.Resource)); err != nil {
				return err
			}
		}
	case reflect.Ptr, reflect.Interface:
		return structEncoderSimpleMultiValued(resource, field.Elem(), tag)
	case reflect.Struct:
		resource.EnsureComplexMultiValuedAttribute(tag.name, tag.max())
		fieldStruct := make(structs.Resource)
		t := field.Type()
		for i := 0; i < field.NumField(); i++ {
			tagIndex := parseTags(t.Field(i))
			if tagIndex.ignore {
				continue
			}

			fieldIndex := field.Field(i)
			if !tagIndex.allowZero && fieldIndex.IsZero() {
				continue
			}
			if err := structEncoder(fieldStruct, fieldIndex, tagIndex); err != nil {
				return err
			}
		}
		if depth := fieldStruct.Depth(); 1 < depth {
			return fmt.Errorf("nested depth exceeded: %d", depth)
		}

		if err := resource.AppendComplexMultiValuedAttribute(tag.name, fieldStruct); err != nil {
			return err
		}
	default:
		resource.EnsureMultiValuedAttribute(tag.name, tag.max())
		value := make(structs.Resource)
		if err := structEncoderSimple(value, field, tag); err != nil {
			return err
		}
		for _, v := range value {
			if err := resource.AppendMultiValuedAttribute(tag.name, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func unsupportedTypeEncoder(v reflect.Value) (map[string]interface{}, error) {
	return nil, errors.New(fmt.Sprintf("unsupported type %s", v.Type()))
}

func validSimpleAttribute(v reflect.Value) (interface{}, error) {
	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil, nil
		}

		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return v.Float(), nil
	case reflect.String:
		return v.String(), nil
	default:
		return nil, fmt.Errorf("invalid simple attribute: %s", v.Kind())
	}
}

// Marshaler is the interface implemented by types that can marshal themselves into SCIM resources.
type Marshaler interface {
	MarshalSCIM() (structs.Resource, error)
}
