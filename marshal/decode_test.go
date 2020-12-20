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

func ExampleUnmarshal_extension() {
	type EnterpriseUserExtension struct {
		EmployeeNumber string
	}

	type User struct {
		UserName       string
		EnterpriseUser EnterpriseUserExtension `scim:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"`
	}

	var user User
	_ = Unmarshal(map[string]interface{}{
		"userName": "di-wu",
		"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
			"employeeNumber": "0001",
		},
	}, &user)
	fmt.Println(user)

	// Output:
	// {di-wu {0001}}

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

type ResourceAttributes map[string]interface{}
type Attributes []interface{}
type String string

func TestUnmarshal_TypeAliases(t *testing.T) {
	s := ResourceAttributes{
		"name": String("Quint"),
		"nil":  nil,
		"last": ResourceAttributes{
			"name": String("Daenen"),
		},
		"nickNames": Attributes{
			ResourceAttributes{
				"name": String("quint"),
			},
			ResourceAttributes{
				"name": String("di-wu"),
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
			"givenName":  "Quint",
			"familyName": "Daenen",
		},
	}

	var resource ResourceStruct
	_ = Unmarshal(resourceMap, &resource)
	fmt.Println(resource)

	// Output:
	// map[name:map[familyName:Daenen givenName:Quint] userName:di-wu]
	// {di-wu {Quint Daenen}}
}
