package marshal

import (
	"fmt"
	"github.com/scim2/tools/structs"
	"testing"
)

type testUnmarshalInterface struct {
	Name string
}

func (t *testUnmarshalInterface) UnmarshalSCIM(resource structs.Resource) error {
	t.Name = resource["name"].(string)
	return nil
}

func ExampleUnmarshal() {
	var r testUnmarshalInterface
	_ = Unmarshal(structs.Resource{
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
	s := structs.Resource{
		"name": "Quint",
		"nil":  nil,
		"last": structs.Resource{
			"name": "Daenen",
		},
		"nickNames": []interface{}{
			structs.Resource{
				"name": "quint",
			},
			structs.Resource{
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

	resourceMap := structs.Resource{
		"userName": "di-wu",
		"name": structs.Resource{
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
