package marshal

import (
	"github.com/di-wu/scim-tools/structs"
	"testing"
)

func TestUnmarshalStruct(t *testing.T) {
	var test testMarshaller

	if err := Unmarshal(map[string]interface{}{
		"name": "test",
	}, &test); err != nil {
		t.Error(err)
	}

	if test.Name != "test" {
		t.Errorf("expected \"test\", got %q", test.Name)
	}
}

func (m *testMarshaller) UnmarshalSCIM(resource structs.Resource) error {
	if name, ok := resource["name"]; ok {
		m.Name = name.(string)
	}
	return nil
}
