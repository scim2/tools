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
			if option == "multiValued" || option == "mV" {
				t.multiValued = true
			}
			if strings.HasPrefix(option, "index=") || strings.HasPrefix(option, "i=") {
				option = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(option, "i"), "ndex"), "=")
				for _, v := range strings.Split(option, ";") {
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
