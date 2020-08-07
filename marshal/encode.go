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

func emptyResource() structs.Resource {
	return make(structs.Resource)
}

func interfaceEncoder(v reflect.Value) (map[string]interface{}, error) {
	if v.IsNil() {
		return emptyResource(), errors.New("nil")
	}
	return reflectValue(v.Elem())
}

func marshalerEncoder(v reflect.Value) (map[string]interface{}, error) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return emptyResource(), errors.New("nil")
	}
	m, ok := v.Interface().(Marshaler)
	if !ok {
		return emptyResource(), errors.New("does not implement marshaller")
	}
	return m.MarshalSCIM()
}

func newMapEncoder(v reflect.Value) (map[string]interface{}, error) {
	t := v.Type()
	if t.Key().Kind() != reflect.String {
		return emptyResource(), errors.New("invalid type")
	}

	elementEncoder := typeEncoder(t.Elem())

	resource := make(structs.Resource)

	keys := v.MapKeys()
	for _, k := range keys {
		if k.Kind() != reflect.String {
			return emptyResource(), errors.New("invalid type")
		}

		value := v.MapIndex(k)
		if value.Kind() == reflect.Interface {
			value = value.Elem()
		}

		switch value.Kind() {
		case reflect.Map:
			sub, err := elementEncoder(value)
			if err != nil {
				return emptyResource(), err
			}
			resource[k.String()] = sub
		default:
			resource[k.String()] = value.Interface()
		}
	}
	return resource, nil
}

func newStructEncoder(v reflect.Value) (map[string]interface{}, error) {
	var (
		resource    = make(structs.Resource)
		mVToAdd     = make(map[string][]interface{})
		mVMapsToAdd = make(map[string][]map[string]interface{})
		mapsToAdd   = make(map[string]map[string]interface{})
	)

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.IsZero() {
			continue
		}

		tag := parseTags(t.Field(i))
		if tag.sub == "" && !tag.multiValued {
			resource[tag.name] = field.Interface()
			continue
		}
		if tag.sub != "" {
			if !tag.multiValued {
				m, ok := mapsToAdd[tag.name]
				if !ok {
					mapsToAdd[tag.name] = make(map[string]interface{})
					m = mapsToAdd[tag.name]
				}
				_, ok = m[tag.sub]
				if ok {
					return nil, errors.New("already full")
				}
				m[tag.sub] = field.Interface()
			} else {
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
			continue
		}
		if tag.multiValued {
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
					mv = append(mv, field.Index(i).Interface())
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
		return emptyResource(), errors.New("invalid value")
	}
	return typeEncoder(v.Type())(v)
}

func unsupportedTypeEncoder(v reflect.Value) (map[string]interface{}, error) {
	return emptyResource(), errors.New(fmt.Sprintf("unsupported type %s", v.Type()))
}

// Marshaler is the interface implemented by types that can marshal themselves into valid SCIM resources.
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
		return newMapEncoder
	case reflect.Ptr:
		return typeEncoder(t.Elem())
	case reflect.Struct:
		return newStructEncoder
	default:
		return unsupportedTypeEncoder
	}
}
