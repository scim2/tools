package marshal

import (
	"reflect"
	"testing"
)

func TestParseTags(t *testing.T) {
	type _tags struct {
		IgnoreJSON      string `json:"ignore"`
		IgnoreOtherTags string `json:"ignoreOther" scim:"ignore"`
		Colon           string `scim:"c:o.l.o:n"`
	}

	tagRefs := []tag{
		{name: "ignoreJSON"},
		{name: "ignore"},
		{name: "c:o.l.o:n"},
	}

	v := reflect.TypeOf(_tags{})
	for i := 0; i < v.NumField(); i++ {
		if tags := parseTags(v.Field(i)); !reflect.DeepEqual(tagRefs[i], tags) {
			t.Errorf("expected %v, got %v", tagRefs[i], tags)
		}
	}
}
