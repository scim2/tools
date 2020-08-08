package marshal

import (
	"errors"
	"fmt"
	"github.com/di-wu/scim-tools/structs"
	"reflect"
)

var marshalerType = reflect.TypeOf((*Marshaler)(nil)).Elem()

func Marshal(v interface{}) (structs.Resource, error) {
	return reflectValue(reflect.ValueOf(v))
}

func interfaceEncoder(v reflect.Value) (map[string]interface{}, error) {
	if v.IsNil() {
		return nil, errors.New("interface is nil")
	}
	return reflectValue(v.Elem())
}

func marshalerEncoder(v reflect.Value) (map[string]interface{}, error) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil, errors.New("ptr is nil")
	}
	m, ok := v.Interface().(Marshaler)
	if !ok {
		return nil, errors.New("value does not implement marshaller")
	}
	return m.MarshalSCIM()
}

func mapEncoder(v reflect.Value) (map[string]interface{}, error) {
	t := v.Type()
	if t.Key().Kind() != reflect.String {
		return nil, errors.New("key of map is not a string")
	}

	elementEncoder := typeEncoder(t.Elem())
	resource := make(structs.Resource)
	keys := v.MapKeys()

	for _, k := range keys {
		if k.Kind() != reflect.String {
			return nil, errors.New("invalid type")
		}

		value := v.MapIndex(k)
		if value.Kind() == reflect.Interface {
			value = value.Elem()
		}

		switch value.Kind() {
		case reflect.Map:
			sub, err := elementEncoder(value)
			if err != nil {
				return nil, err
			}
			resource[k.String()] = sub
		default:
			resource[k.String()] = value.Interface()
		}
	}
	return resource, nil
}

func structEncoder(v reflect.Value) (map[string]interface{}, error) {
	var (
		resource    = make(structs.Resource)
		mVToAdd     = make(map[string][]interface{})
		mVMapsToAdd = make(map[string][]map[string]interface{})
		mapsToAdd   = make(map[string]map[string]interface{})
	)

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

		switch tag.attrType() {
		case simple:
			switch field.Kind() {
			case reflect.Struct:
				sub, err := structEncoder(field)
				if err != nil {
					return nil, err
				}
				m := ensureMapInMap(tag.name, mapsToAdd)
				for k, v := range sub {
					if _, ok := m[k]; ok {
						return nil, errors.New(fmt.Sprintf("duplicate names: %s", tag.sub))
					}
					m[k] = v
				}
			case reflect.Slice:
				return nil, errors.New(fmt.Sprintf("invalid simple attribute: %s", field.Kind()))
			default:
				resource[tag.name] = field.Interface()
			}
		case simpleMultiValued:
			mv, ok := mVToAdd[tag.name]
			if !ok {
				mVToAdd[tag.name] = make([]interface{}, 0)
				mv = mVToAdd[tag.name]
			}

			for _, i := range tag.indexes {
				for len(mv) < i+1 {
					mv = append(mv, nil)
				}
			}

			switch field.Kind() {
			case reflect.Slice:
				for i := 0; i < field.Len(); i++ {
					fieldIndex := field.Index(i)
					switch fieldIndex.Kind() {
					case reflect.Struct:
						sub, err := structEncoder(fieldIndex)
						if err != nil {
							return nil, err
						}
						mv = append(mv, sub)
					default:
						value := fieldIndex.Interface()
						mv = append(mv, value)
					}
				}
			default:
				if len(tag.indexes) == 0 {
					mv = append(mv, field.Interface())
				} else {
					for _, i := range tag.indexes {
						mv[i] = field.Interface()
					}
				}
			}
			mVToAdd[tag.name] = mv
		case complex:
			m := ensureMapInMap(tag.name, mapsToAdd)
			if _, ok := m[tag.sub]; ok {
				return nil, errors.New(fmt.Sprintf("duplicate names: %s", tag.sub))
			}
			m[tag.sub] = field.Interface()
		case complexMultiValued:
			mv, ok := mVMapsToAdd[tag.name]
			if !ok {
				mVMapsToAdd[tag.name] = make([]map[string]interface{}, 0)
				mv = mVMapsToAdd[tag.name]
			}

			for _, i := range tag.indexes {
				for len(mv) < i+1 {
					mv = append(mv, make(map[string]interface{}))
				}
			}

			if len(tag.indexes) == 0 {
				var added bool
				for _, m := range mv {
					_, ok := m[tag.sub]
					if !ok {
						m[tag.sub] = field.Interface()
						added = true
						break
					}
				}
				if !added {
					mv = append(mv, map[string]interface{}{
						tag.sub: field.Interface(),
					})
				}
			} else {
				for _, i := range tag.indexes {
					mv[i][tag.sub] = field.Interface()
				}
			}
			mVMapsToAdd[tag.name] = mv
		}
	}

	for k, v := range mVToAdd {
		resource[k] = v
	}

	for k, v := range mVMapsToAdd {
		resource[k] = v
	}

	for k, v := range mapsToAdd {
		resource[k] = v
	}

	return resource, nil
}

func reflectValue(v reflect.Value) (structs.Resource, error) {
	if !v.IsValid() {
		return nil, errors.New("value is invalid")
	}
	return typeEncoder(v.Type())(v)
}

func unsupportedTypeEncoder(v reflect.Value) (map[string]interface{}, error) {
	return nil, errors.New(fmt.Sprintf("unsupported type %s", v.Type()))
}

// Marshaler is the interface implemented by types that can marshal themselves into SCIM resources.
type Marshaler interface {
	MarshalSCIM() (structs.Resource, error)
}

type encoderFunc func(v reflect.Value) (map[string]interface{}, error)

func typeEncoder(t reflect.Type) encoderFunc {
	if t.Implements(marshalerType) {
		return marshalerEncoder
	}

	switch t.Kind() {
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Map:
		return mapEncoder
	case reflect.Ptr:
		ptr := t.Elem()
		return typeEncoder(ptr.Elem())
	case reflect.Struct:
		return structEncoder
	default:
		return unsupportedTypeEncoder
	}
}
