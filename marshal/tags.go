package marshal

import (
	"reflect"
	"strconv"
	"strings"
)

type tag struct {
	name, sub   string
	multiValued bool
	indexes     []int
}

func parseTags(field reflect.StructField) tag {
	var t tag

	tag := field.Tag.Get("scim")
	tags := strings.Split(tag, ",")
	if tag == "" {
		t.name = lowerFirstRune(field.Name)
	} else {
		parts := strings.Split(tags[0], ".")
		t.name = parts[0]
		if len(parts) > 1 {
			t.sub = parts[1]
		}
	}

	if len(tags) > 1 {
		for _, option := range tags[1:] {
			if option == "multiValued" {
				t.multiValued = true
			}
			if strings.HasPrefix(option, "index=") {
				for _, v := range strings.Split(strings.TrimPrefix(option, "index="), ";") {
					i, err := strconv.Atoi(v)
					if err == nil {
						t.indexes = append(t.indexes, i)
					}
				}
			}
		}
	}

	return t
}
