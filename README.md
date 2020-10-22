# SCIM Tools
This repository/module contains various utility functions to make it easier to work with SCIM servers and clients.

**!** most packages are a wip

## Fuzzer
Build on top of [gofuzz](https://github.com/google/gofuzz/).

```go
var refSchema ReferenceSchema
// define the reference schema yourself or use unmarshal json

resource := New(refSchema).
    // multi valued fields have one value.
    NumElements(1, 1).
    // displayName and name.givenName can never be empty.
    NeverEmpty("displayName", "name.givenName").
    // other fields are empty.
    EmptyChance(1).
    // create fuzzed resource.
    Fuzz()

// OUTPUT: map[displayName:vWKdUsVprh name:map[givenName:ieVkQrrcKL] userName:RFlLpsMnBW]
```

## Encoder
A simple encoder that converts structs to maps based on their tags.

##### Tags
- `multiValued` (or `mV`) \
  Makes the attribute multi valued.

```go
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

resource, _ := Marshal(resourceStruct)

// OUTPUT: map[name:map[familyName:Daenen givenName:Quint] userName:di-wu]
```

## Decoder
A simple decoder that fills structs with maps.

**!** no pointers and tags supported

```go
resourceMap := structs.Resource{
	"userName": "di-wu",
	"name": structs.Resource{
		"firstName": "Quint",
		"lastName":  "Daenen",
	},
}

var resource ResourceStruct
_ = Unmarshal(resourceMap, &resource)

// OUTPUT: {di-wu {Quint Daenen}}
```

## Struct Generator
Converts a schema to a structure representing the resource described in that schema.

```go
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
//     ExternalId string
//     Id         string
//     UserName   string
// }
```
