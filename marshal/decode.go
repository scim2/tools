package marshal

import (
	"errors"
	"fmt"
	"github.com/scim2/tools/structs"
	"reflect"
)

var unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

func Unmarshal(data structs.Resource, value interface{}) error {
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
				s := reflect.ValueOf(fV)
				switch t := fV.(type) {
				case []interface{}:
					if v.Kind() != reflect.Slice {
						break
					}
					field := reflect.MakeSlice(v.Type(), len(t), len(t))
					for i, v := range t {
						switch t := v.(type) {
						case structs.Resource:
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
				case structs.Resource:
					field := reflect.New(v.Type())
					initializeStruct(v.Type(), field.Elem())
					if err := Unmarshal(t, field.Interface()); err != nil {
						return err
					}
					v.Set(field.Elem())
					continue
				}
				if s.Type() != v.Type() {
					return fmt.Errorf(
						"types of %q do not match: got %s, want %s",
						name, s.Type(), v.Type(),
					)
				}
				v.Set(reflect.ValueOf(fV))
			}
		}
	}
	return nil
}

// Unmarshaler is the interface implemented by types that can unmarshal a SCIM description of themselves.
type Unmarshaler interface {
	UnmarshalSCIM(structs.Resource) error
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
