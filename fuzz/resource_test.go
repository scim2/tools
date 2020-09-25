package fuzz

import (
	"encoding/json"
	schema2 "github.com/scim2/tools/schema"
	"testing"

	"github.com/elimity-com/scim/schema"
)

func TestReferenceSchemaNeverEmpty(t *testing.T) {
	var s schema2.ReferenceSchema
	raw, _ := schema.CoreUserSchema().MarshalJSON()
	_ = json.Unmarshal(raw, &s)

	f := New(s).
		EmptyChance(1).
		NumElements(1, 1).
		NeverEmpty("displayName", "name.givenName", "emails")

	for i := 0; i < 100; i++ {
		resource := f.Fuzz()

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
			emails, ok := emailsMap.([]interface{})
			if !ok {
				t.Errorf("email not a multi valued attribute")
			}
			if len(emails) == 0 {
				t.Errorf("emails is empty")
				break
			}
			for _, e := range emails {
				if email, ok := e.(map[string]interface{}); !ok {
					t.Errorf("emails is not a complex attribute")
				} else {
					if _, ok := email["display"]; !ok {
						t.Errorf("emails.value is not present")
					}
					if _, ok := email["primary"]; !ok {
						t.Errorf("emails.value is not present")
					}
					if _, ok := email["type"]; !ok {
						t.Errorf("emails.value is not present")
					}
					if _, ok := email["value"]; !ok {
						t.Errorf("emails.value is not present")
					}
				}
			}
		}
	}
}
