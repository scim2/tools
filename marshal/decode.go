package marshal

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
	mapType         = reflect.TypeOf(map[string]interface{}{})
	sliceType       = reflect.TypeOf([]interface{}{})
)

func Unmarshal(data map[string]interface{}, value interface{}) error {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("value is invalid")
	}

	t := v.Type()
	if t.Implements(unmarshalerType) {
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return errors.New("ptr is nil")
		}
		m, ok := v.Interface().(Unmarshaler)
		if !ok {
			return errors.New("value does not implement marshaler")
		}
		return m.UnmarshalSCIM(data)
	}

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	for i := 0; i < t.NumField(); i++ {
		if v := v.Field(i); v.CanAddr() && v.CanSet() {
			f := t.Field(i)
			name := lowerFirstRune(f.Name)
			if fV, ok := data[name]; ok {
				if fV == nil {
					continue
				}
				s := reflect.ValueOf(fV)
				switch s.Kind() {
				case reflect.Array, reflect.Slice:
					t := toDefaultSlice(fV)
					if v.Kind() != reflect.Slice {
						break
					}
					field := reflect.MakeSlice(v.Type(), len(t), len(t))
					for i, v := range t {
						switch reflect.ValueOf(v).Kind() {
						case reflect.Map:
							t := toDefaultMap(v)
							typ := field.Index(i).Type()
							element := reflect.New(typ)
							initializeStruct(typ, element.Elem())
							if err := Unmarshal(t, element.Interface()); err != nil {
								return err
							}
							field.Index(i).Set(element.Elem())
						default:
							field.Index(i).Set(reflect.ValueOf(v))
						}
					}
					v.Set(field)
					continue
				case reflect.Map:
					t := toDefaultMap(fV)
					field := reflect.New(v.Type())
					initializeStruct(v.Type(), field.Elem())
					if err := Unmarshal(t, field.Interface()); err != nil {
						return err
					}
					v.Set(field.Elem())
					continue
				}
				if s.Kind() != v.Kind() {
					return fmt.Errorf(
						"types of %q do not match: got %s, want %s",
						name, s.Type(), v.Type(),
					)
				}

				if s.Type() != v.Type() {
					v.Set(reflect.ValueOf(toType(fV, v.Type())))
				} else {
					v.Set(s)
				}
			}
		}
	}
	return nil
}

// Unmarshaler is the interface implemented by types that can unmarshal a SCIM description of themselves.
type Unmarshaler interface {
	UnmarshalSCIM(map[string]interface{}) error
}

func initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Struct:
			initializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
}

func toDefaultMap(m interface{}) map[string]interface{} {
	if reflect.TypeOf(m) != mapType {
		return toType(m, mapType).(map[string]interface{})
	}
	return m.(map[string]interface{})
}

func toDefaultSlice(m interface{}) []interface{} {
	if reflect.TypeOf(m) != sliceType {
		return toType(m, sliceType).([]interface{})
	}
	return m.([]interface{})
}

func toType(i interface{}, t reflect.Type) interface{} {
	return reflect.
		ValueOf(i).
		Convert(t).
		Interface()
}
