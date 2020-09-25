package gen_test

import (
	"encoding/json"
	"fmt"
	scimSchema "github.com/elimity-com/scim/schema"
	"github.com/scim2/tools/schema"
	gen "github.com/scim2/tools/structs/generate"
	"testing"
)

func TestGenerateStruct(t *testing.T) {
	if _, err := gen.NewStructGenerator(schema.ReferenceSchema{}); err == nil {
		t.Error("error expected, got none")
	}
}
func ExampleStructGenerator_Generate_empty() {
	g, _ := gen.NewStructGenerator(schema.ReferenceSchema{
		Name:        "User",
		Description: "User Account",
	})
	fmt.Print(g.Generate())

	// Output:
	// // User Account
	// type User struct {}
}

func ExampleStructGenerator_Generate_minimal() {
	g, _ := gen.NewStructGenerator(schema.ReferenceSchema{
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
	//     UserName string
	// }
}

func ExampleStructGenerator_Generate_user() {
	var s schema.ReferenceSchema
	raw, _ := scimSchema.CoreUserSchema().MarshalJSON()
	_ = json.Unmarshal(raw, &s)

	g, _ := gen.NewStructGenerator(s)
	fmt.Print(g.UsePtr(true).Generate())

	// Output:
	// // User Account
	// type User struct {
	//     UserName          string
	//     Name              *UserName
	//     DisplayName       *string
	//     NickName          *string
	//     ProfileUrl        *string
	//     Title             *string
	//     UserType          *string
	//     PreferredLanguage *string
	//     Locale            *string
	//     Timezone          *string
	//     Active            *bool
	//     Password          *string
	//     Emails            []UserEmail
	//     PhoneNumbers      []UserPhoneNumber
	//     Ims               []UserIm
	//     Photos            []UserPhoto
	//     Addresses         []UserAddresse
	//     Groups            []UserGroup
	//     Entitlements      []UserEntitlement
	//     Roles             []UserRole
	//     X509Certificates  []UserX509Certificate
	// }
	//
	// // The components of the user's real name. Providers MAY return just the full name as a single string in the formatted
	// // or they MAY return just the individual component attributes using the other sub-attributes, or they MAY return both.
	// // both variants are returned, they SHOULD be describing the same name, with the formatted name indicating how the
	// // attributes should be combined.
	// type UserName struct {
	//     Formatted       *string
	//     FamilyName      *string
	//     GivenName       *string
	//     MiddleName      *string
	//     HonorificPrefix *string
	//     HonorificSuffix *string
	// }
	//
	// // Email addresses for the user. The value SHOULD be canonicalized by the service provider, e.g., 'bjensen@example.com'
	// // of 'bjensen@EXAMPLE.COM'. Canonical type values of 'work', 'home', and 'other'.
	// type UserEmail struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // Phone numbers for the User. The value SHOULD be canonicalized by the service provider according to the format
	// // in RFC 3966, e.g., 'tel:+1-201-555-0123'. Canonical type values of 'work', 'home', 'mobile', 'fax', 'pager', and
	// type UserPhoneNumber struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // Instant messaging addresses for the User.
	// type UserIm struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // URLs of photos of the User.
	// type UserPhoto struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // A physical mailing address for this User. Canonical type values of 'work', 'home', and 'other'. This attribute is a
	// // type with the following sub-attributes.
	// type UserAddresse struct {
	//     Formatted     *string
	//     StreetAddress *string
	//     Locality      *string
	//     Region        *string
	//     PostalCode    *string
	//     Country       *string
	//     Type          *string
	// }
	//
	// // A list of groups to which the user belongs, either through direct membership, through nested groups, or dynamically
	// type UserGroup struct {
	//     Value   *string
	//     Ref     *string
	//     Display *string
	//     Type    *string
	// }
	//
	// // A list of entitlements for the User that represent a thing the User has.
	// type UserEntitlement struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // A list of roles for the User that collectively represent who the User is, e.g., 'Student', 'Faculty'.
	// type UserRole struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // A list of certificates issued to the User.
	// type UserX509Certificate struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
}
