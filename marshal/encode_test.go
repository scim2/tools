package marshal

import (
	"fmt"
	"testing"

	"github.com/di-wu/scim-tools/structs"
)

func TestMarshal(t *testing.T) {
	t.Run("map", func(t *testing.T) {
		type Roles []struct {
			Value string `scim:"value"`
		}

		type Name struct {
			Full string `scim:"fullName"`
		}

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
			Roles        Roles    `scim:"roles,multiValued"`
			FullName     Name     `scim:"name"`
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
			Roles: Roles{
				{Value: "Author"},
				{Value: "Member"},
			},
			FullName: Name{Full: "Quint Daenen"},
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
				"fullName":   "Quint Daenen",
			},
			"roles": []map[string]interface{}{
				{"value": "Author"},
				{"value": "Member"},
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

	type (
		s1 struct {
			String string
			Int    int
			Bool   bool `scim:",0"`

			IgnoreMe []interface{} `scim:",!"`
		}

		s2 struct {
			Slice    []interface{} `scim:",mV"`
			Slice3   interface{}   `scim:"slice,mV,i=0;2"`
			InvalidS []interface{} // must be multi valued or ignored
		}
	)

	for _, test := range []struct {
		v        interface{}
		errMsg   string
		resource string
	}{
		{v: "", errMsg: "unsupported type string"},
		{v: 0, errMsg: "unsupported type int"},
		{v: 0.0, errMsg: "unsupported type float64"},
		{v: true, errMsg: "unsupported type bool"},
		{v: []interface{}{}, errMsg: "unsupported type []interface {}"},
		{v: &[]interface{}{}, errMsg: "unsupported type []interface {}"},
		{v: make([]string, 10, 10), errMsg: "unsupported type []string"},
		{v: map[string]interface{}{}, resource: "map[]"},
		{v: &map[string]interface{}{}, resource: "map[]"},
		{v: map[string]interface{}{"_": 0}, resource: "map[_:0]"},
		{v: &map[string]interface{}{"_": 0}, resource: "map[_:0]"},
		{v: map[int]interface{}{0: "_"}, errMsg: "key of map is not a string"},

		{v: s1{String: "_", Int: 0}, resource: "map[bool:false string:_]"},
		{v: s1{String: "_", Int: 1, Bool: true}, resource: "map[bool:true int:1 string:_]"},
		{v: s1{IgnoreMe: []interface{}{"_"}}, resource: "map[bool:false]"},

		{v: s2{Slice: []interface{}{"_"}}, resource: "map[slice:[_]]"},
		{v: s2{Slice3: "_"}, resource: "map[slice:[_ <nil> _]]"},
		{v: s2{InvalidS: []interface{}{"_"}}, errMsg: "invalid simple attribute: slice"},
	} {
		r, err := Marshal(test.v)
		if test.errMsg != "" {
			if err == nil {
				t.Errorf("error expected: %q, got none", test.errMsg)
			} else if test.errMsg != err.Error() {
				t.Errorf("expected %q, got %q", test.errMsg, err)
			}
		} else if test.errMsg == "" && err != nil {
			t.Errorf("no error expected: got %q", err)
		}

		resource := fmt.Sprintf("%v", r)
		if test.resource != "" && test.resource != resource {
			t.Errorf("expected %s, got %s", test.resource, resource)
		}
	}
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
