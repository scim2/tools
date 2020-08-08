package marshal

import (
	"errors"
	"fmt"
	"github.com/di-wu/scim-tools/structs"
	"reflect"
)

var unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

func Unmarshal(resource structs.Resource, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("value not a ptr or is nil")
	}

	if rv.Type().Implements(unmarshalerType) {
		m, ok := v.(Unmarshaler)
		if !ok {
			return errors.New("value does not implement unmarshaler")
		}
		return m.UnmarshalSCIM(resource)
	}

	ptr := rv.Elem()
	t := ptr.Type()

	switch ptr.Kind() {
	case reflect.Struct:
		for i := 0; i < ptr.NumField(); i++ {
			tag := parseTags(t.Field(i))
			if tag.ignore {
				continue
			}
		}
		return nil
	default:
		return errors.New(fmt.Sprintf("unsupported kind %s", rv.Kind()))
	}
}

// Unmarshaler is the interface implemented by types that can unmarshal a SCIM resource of themselves.
type Unmarshaler interface {
	UnmarshalSCIM(structs.Resource) error
}
