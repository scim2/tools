package marshal

import (
	"reflect"
	"strconv"
	"strings"
)

type attributeType string

var (
	complexMultiValued attributeType = "c mV"
	complex            attributeType = "c"
	simpleMultiValued  attributeType = "s mV"
	simple             attributeType = "s"
)

type tag struct {
	name, sub   string
	multiValued bool
	allowZero   bool
	ignore      bool
	indexes     []int
}

func (t tag) attrType() attributeType {
	if t.sub == "" {
		if t.multiValued {
			return simpleMultiValued
		}
		return simple
	}
	if t.multiValued {
		return complexMultiValued
	}
	return complex
}

func (t tag) max() int {
	var max int
	for _, i := range t.indexes {
		if i > max {
			max = i
		}
	}
	return max
}

func parseTags(field reflect.StructField) tag {
	var t tag

	tag := field.Tag.Get("scim")
	tags := strings.Split(tag, ",")
	if tag == "" {
		t.name = lowerFirstRune(field.Name)
	} else {
		parts := strings.Split(tags[0], ".")
		if parts[0] != "" {
			t.name = parts[0]
		} else {
			t.name = lowerFirstRune(field.Name)
		}

		if len(parts) > 1 {
			t.sub = parts[1]
		}
	}

	if len(tags) > 1 {
		for _, option := range tags[1:] {
			if strings.HasPrefix(option, "index=") || strings.HasPrefix(option, "i=") {
				option = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(option, "i"), "ndex"), "=")
				for _, v := range strings.Split(option, ";") {
					if i, err := strconv.Atoi(v); err == nil {
						t.indexes = append(t.indexes, i)
					}
				}
			}
			switch option {
			case "multiValued", "mV":
				t.multiValued = true
			case "zero", "0":
				t.allowZero = true
			case "ignore", "!":
				t.ignore = true
			}
		}
	}

	return t
}
