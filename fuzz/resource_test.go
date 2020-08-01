package fuzz

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/elimity-com/scim/schema"
	"github.com/google/gofuzz"
)

func TestReferenceSchemaNeverEmpty(t *testing.T) {
	var s ReferenceSchema
	raw, _ := schema.CoreUserSchema().MarshalJSON()
	_ = json.Unmarshal(raw, &s)

	// displayName, name.givenName and emails.value can never be empty.
	s.NeverEmpty("displayName", "name.givenName", "emails.value")

	f := fuzz.New().Funcs(
		NewResourceFuzzer(s),
	)

	for i := 0; i < 100; i++ {
		var resource Resource
		f.Fuzz(&resource)
		fmt.Println(resource)

		if _, ok := resource["userName"]; !ok {
			t.Errorf("userName not present")
		}
		if _, ok := resource["displayName"]; !ok {
			t.Errorf("displayName not present")
		}

		if nameMap, ok := resource["name"]; !ok {
			t.Errorf("name not present")
		} else {
			name, ok := nameMap.(map[string]interface{})
			if !ok {
				t.Errorf("name not a complex attribute")
			}
			if _, ok := name["givenName"]; !ok {
				t.Errorf("name.givenName not present")
			}
		}

		if emailsMap, ok := resource["emails"]; !ok {
			t.Errorf("emails not present")
		} else {
			emails, ok := emailsMap.([]map[string]interface{})
			if !ok {
				t.Errorf("email not a complex multi valued attribute")
			}
			if len(emails) == 0 {
				t.Errorf("emails is empty")
				break
			}
			for _, email := range emails {
				if _, ok := email["value"]; !ok {
					t.Errorf("email.value is not present")
				}
			}
		}
	}
}
