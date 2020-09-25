package gen_test

import (
	"fmt"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
	gen "github.com/scim2/tools/structs/generate"
	"testing"
)

func TestGenerateStruct(t *testing.T) {
	if _, err := gen.GenerateStruct(schema.Schema{}); err == nil {
		t.Error("error expected, got none")
	}
}

func ExampleGenerateStruct() {
	b, err := gen.GenerateStruct(schema.Schema{
		Name: optional.NewString("User"),
	})
	fmt.Print(b)
	fmt.Print(err)

	// Output:
	// type User struct {}
	// <nil>
}
