package marshal

import (
	"fmt"
	"testing"

	"github.com/di-wu/scim-tools/structs"
)

func TestMarshal(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		type User struct {
			UserName     string
			DisplayName  string
			FirstName    string   `scim:"name.givenName"`
			LastName     string   `scim:"name.familyName"`
			Email        string   `scim:"emails.value,multiValued,index=0;2"`
			WorkEmail    string   `scim:"emails.value,multiValued"`
			EmailType    string   `scim:"emails.type,multiValued,index=0;1"`
			Primary      bool     `scim:"emails.primary,multiValued,index=2"`
			RandomValues []string `scim:"values,multiValued"`
			RandomValue  string   `scim:"values,multiValued,index=0"`
		}

		user := User{
			UserName:     "di-wu",
			DisplayName:  "quint",
			FirstName:    "Quint",
			LastName:     "Daenen",
			Email:        "me@di-wu.be",
			WorkEmail:    "quint@elimity.com",
			EmailType:    "work",
			Primary:      true,
			RandomValues: []string{"replaced", "random", "values"},
			RandomValue:  "some",
		}

		resource, err := Marshal(user)
		if err != nil {
			t.Error(err)
		}

		ref := map[string]interface{}{
			"displayName": "quint",
			"emails": []map[string]interface{}{
				{
					"type":  "work",
					"value": "me@di-wu.be",
				},
				{
					"type":  "work",
					"value": "quint@elimity.com",
				},
				{
					"primary": true,
					"value":   "me@di-wu.be",
				},
			},
			"name": map[string]interface{}{
				"familyName": "Daenen",
				"givenName":  "Quint",
			},
			"userName": "di-wu",
			"values": []string{
				"some", "random", "values",
			},
		}
		if fmt.Sprintf("%v", ref) != fmt.Sprintf("%v", resource) {
			t.Errorf("resources do not match:\n%v\n%v", resource, ref)
		}
	})
}

func TestMarshalMap(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		testMarshal(t, map[string]interface{}{
			"name": "test",
		})
	})

	t.Run("string", func(t *testing.T) {
		testMarshal(t, map[string]string{
			"name": "test",
		})
	})

	_, err := Marshal([]string{"test"})
	if err == nil {
		t.Error("error expected")
	}
}

func TestMarshalStruct(t *testing.T) {
	test := testMarshaller{
		Name: "test",
	}

	t.Run("func", func(t *testing.T) {
		testMarshal(t, test)
	})

	t.Run("ptr", func(t *testing.T) {
		testMarshal(t, &test)
	})

	t.Run("interface", func(t *testing.T) {
		resource, err := test.MarshalSCIM()
		if err != nil {
			t.Error(err)
			return
		}
		testName(t, resource)
	})

	t.Run("tags", func(t *testing.T) {
		testMarshal(t, testTagMarshaller{
			UserName: "test",
		})
	})
}

func testMarshal(t *testing.T, test interface{}) structs.Resource {
	resource, err := Marshal(test)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resource
}

func testName(t *testing.T, resource structs.Resource) {
	v, ok := resource["name"]
	if !ok {
		t.Errorf("could nog find \"name\" in map")
		return
	}
	if v.(string) != "test" {
		t.Errorf("expected \"test\", got %v", v)
	}
}

type testMarshaller struct {
	Name string
}

func (m testMarshaller) MarshalSCIM() (structs.Resource, error) {
	return structs.Resource{
		"name": m.Name,
	}, nil
}

type testTagMarshaller struct {
	UserName string `scim:"name"`
}
