package marshal

import (
	"testing"

	"github.com/di-wu/scim-tools/structs"
)

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

func testMarshal(t *testing.T, test interface{}) {
	resource, err := Marshal(test)
	if err != nil {
		t.Error(err)
		return
	}

	testName(t, resource)
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
