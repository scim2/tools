package gen_test

import (
	"fmt"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
	gen "github.com/scim2/tools/structs/generate"
	"testing"
)

func TestGenerateStruct(t *testing.T) {
	if _, err := gen.NewStructGenerator(schema.Schema{}); err == nil {
		t.Error("error expected, got none")
	}
}
func ExampleStructGenerator_Generate_empty() {
	g, _ := gen.NewStructGenerator(schema.Schema{
		Name:        optional.NewString("User"),
		Description: optional.NewString("User Account"),
	})
	fmt.Print(g.Generate())

	// Output:
	// // User Account
	// type User struct {}
}

func ExampleStructGenerator_Generate_minimal() {
	g, _ := gen.NewStructGenerator(schema.Schema{
		Name: optional.NewString("User"),
		Attributes: []schema.CoreAttribute{
			schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
				Name:       "userName",
				Required:   true,
				Uniqueness: schema.AttributeUniquenessServer(),
			})),
		},
	})
	fmt.Print(g.Generate())

	// Output:
	// type User struct {
	//     UserName string
	// }
}

func ExampleStructGenerator_Generate_user() {
	g, _ := gen.NewStructGenerator(schema.CoreUserSchema())
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
	// // or they MAY return just the individual component attributes using the other sub-attributes, or they MAY return both. If
	// // variants are returned, they SHOULD be describing the same name, with the formatted name indicating how the component
	// // should be combined.
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
	// type UserEmails struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // Phone numbers for the User. The value SHOULD be canonicalized by the service provider according to the format specified
	// // RFC 3966, e.g., 'tel:+1-201-555-0123'. Canonical type values of 'work', 'home', 'mobile', 'fax', 'pager', and 'other'.
	// type UserPhoneNumbers struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // Instant messaging addresses for the User.
	// type UserIms struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // URLs of photos of the User.
	// type UserPhotos struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // A physical mailing address for this User. Canonical type values of 'work', 'home', and 'other'. This attribute is a
	// // type with the following sub-attributes.
	// type UserAddresses struct {
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
	// type UserGroups struct {
	//     Value   *string
	//     $Ref    *string
	//     Display *string
	//     Type    *string
	// }
	//
	// // A list of entitlements for the User that represent a thing the User has.
	// type UserEntitlements struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // A list of roles for the User that collectively represent who the User is, e.g., 'Student', 'Faculty'.
	// type UserRoles struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
	//
	// // A list of certificates issued to the User.
	// type UserX509Certificates struct {
	//     Value   *string
	//     Display *string
	//     Type    *string
	//     Primary *bool
	// }
}
