package generate_test

import (
	"fmt"
	"testing"

	"github.com/scim2/tools/generate"
	"github.com/scim2/tools/schema"
)

func TestGenerateStruct(t *testing.T) {
	if _, err := generate.NewStructGenerator(schema.ReferenceSchema{}); err == nil {
		t.Error("error expected, got none")
	}
}

func ExampleStructGenerator_AddTags() {
	g, _ := generate.NewStructGenerator(schema.ReferenceSchema{
		Name:        "User",
		Description: "User Account",
	})
	g.AddTags(func(a *schema.Attribute) (tags map[string]string) {
		tags = make(map[string]string)
		if a.Required {
			tags["x"] = "required"
		}
		if a.Uniqueness == schema.Server {
			x, ok := tags["x"]
			if !ok {
				tags["x"] = "unique"
			} else {
				tags["x"] = x + ",unique"
			}
		}
		return tags
	})
	fmt.Print(g.Generate())

	// Output:
	// // User Account
	// type User struct {
	//     ExternalId string
	//     Id         string `x:"required,unique"`
	// }
}

func ExampleStructGenerator_CustomTypes() {
	ref := schema.ReferenceSchema{
		Name:       "User",
		Attributes: []*schema.Attribute{schema.MetaAttribute},
	}
	g, _ := generate.NewStructGenerator(ref)
	g.CustomTypes([]generate.CustomType{
		{
			PkgPrefix: "uuid",
			AttrName:  "id",
			TypeName:  "UUID",
		},
		{
			AttrName: "meta",
			TypeName: "Meta",
		},
	})
	fmt.Print(g.Generate())

	// Output:
	// type User struct {
	//     ExternalId string
	//     Id         uuid.UUID
	//     Meta       Meta
	// }
}

func ExampleStructGenerator_Generate_empty() {
	g, _ := generate.NewStructGenerator(schema.ReferenceSchema{
		Name:        "User",
		Description: "User Account",
	})
	fmt.Print(g.Generate())

	// Output:
	// // User Account
	// type User struct {
	//     ExternalId string
	//     Id         string
	// }
}

func ExampleStructGenerator_Generate_minimal() {
	g, _ := generate.NewStructGenerator(schema.ReferenceSchema{
		Name: "User",
		Attributes: []*schema.Attribute{
			{
				Name:       "userName",
				Required:   true,
				Uniqueness: schema.Server,
			},
		},
	})
	fmt.Print(g.Generate())

	// Output:
	// type User struct {
	//     ExternalId string
	//     Id         string
	//     UserName   string
	// }
}
