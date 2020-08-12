package marshal

import (
	"reflect"
	"testing"
)

func TestParseTags(t *testing.T) {
	type _tags struct {
		IgnoreJSON      string `json:"ignore"`
		IgnoreOtherTags string `json:"ignoreOther" scim:"ignore"`

		MV1 bool `scim:",multiValued"`
		MV2 bool `scim:",mV"`

		Zero1 int `scim:",zero"`
		Zero2 int `scim:",0"`

		Ignore1 bool `scim:",ignore"`
		Ignore2 bool `scim:",!"`

		Index1 []int `scim:",index=0"`
		Index2 []int `scim:",i=0"`
		Index3 []int `scim:",i=0;2"`
		Index4 []int `scim:",i=0;2;-1"`
		Index5 []int `scim:",i=all"`
		Index6 []int `scim:",i=0-2"`
		Index7 []int `scim:",i=0-2;5"`

		Sub  rune `scim:"sub,_mV"`
		Sub2 rune `scim:".sub,_mV"`
	}

	tagRefs := []tag{
		{name: "ignoreJSON"},
		{name: "ignore"},

		{name: "mV1", multiValued: true},
		{name: "mV2", multiValued: true},

		{name: "zero1", allowZero: true},
		{name: "zero2", allowZero: true},

		{name: "ignore1", ignore: true},
		{name: "ignore2", ignore: true},

		{name: "index1", indexes: []int{0}},
		{name: "index2", indexes: []int{0}},
		{name: "index3", indexes: []int{0, 2}},
		{name: "index4", indexes: []int{0, 2}},
		{name: "index5", indexes: []int{-1}},
		{name: "index6", indexes: []int{0, 1, 2}},
		{name: "index7", indexes: []int{0, 1, 2, 5}},

		{name: "sub"},
		{name: "sub2", sub: &tag{name: "sub", multiValued: true}},
	}

	v := reflect.TypeOf(_tags{})
	for i := 0; i < v.NumField(); i++ {
		if tags := parseTags(v.Field(i)); !reflect.DeepEqual(tagRefs[i], tags) {
			t.Errorf("expected %v, got %v", tagRefs[i], tags)
		}
	}
}
