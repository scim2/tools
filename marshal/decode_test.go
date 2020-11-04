package marshal

import (
	"fmt"
	"testing"
)

type testUnmarshalInterface struct {
	Name string
}

func (t *testUnmarshalInterface) UnmarshalSCIM(resource map[string]interface{}) error {
	t.Name = resource["name"].(string)
	return nil
}

func ExampleUnmarshal() {
	var r testUnmarshalInterface
	_ = Unmarshal(map[string]interface{}{
		"name": "Quint",
	}, &r)
	fmt.Println(r)

	// Output:
	// {Quint}
}

type testUnmarshal struct {
	Name string
	Nil  interface{}
	Last struct {
		Name string
	}
	NickNames []struct {
		Name string
	}
}

func TestUnmarshal(t *testing.T) {
	s := map[string]interface{}{
		"name": "Quint",
		"nil":  nil,
		"last": map[string]interface{}{
			"name": "Daenen",
		},
		"nickNames": []interface{}{
			map[string]interface{}{
				"name": "quint",
			},
			map[string]interface{}{
				"name": "di-wu",
			},
		},
	}

	var r testUnmarshal
	if err := Unmarshal(s, &r); err != nil {
		t.Error(err)
	}
	if r.Name != "Quint" {
		t.Error()
	}
	if r.Last.Name != "Daenen" {
		t.Error()
	}
	if r.NickNames[0].Name != "quint" {
		t.Error()
	}
	if r.NickNames[1].Name != "di-wu" {
		t.Error()
	}
}

func Example() {
	type Name struct {
		FirstName string `scim:"givenName"`
		LastName  string `scim:"familyName"`
	}

	type ResourceStruct struct {
		UserName string
		Name     Name
	}

	resourceStruct := ResourceStruct{
		UserName: "di-wu",
		Name: Name{
			FirstName: "Quint",
			LastName:  "Daenen",
		},
	}

	r, _ := Marshal(resourceStruct)
	fmt.Println(r)

	resourceMap := map[string]interface{}{
		"userName": "di-wu",
		"name": map[string]interface{}{
			"firstName": "Quint",
			"lastName":  "Daenen",
		},
	}

	var resource ResourceStruct
	_ = Unmarshal(resourceMap, &resource)
	fmt.Println(resource)

	// Output:
	// map[name:map[familyName:Daenen givenName:Quint] userName:di-wu]
	// {di-wu {Quint Daenen}}
}
