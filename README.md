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

# Marshaller
A simple marshaller that converts structs to maps based on their tags.


##### Tags
- `index=0;1;5` (or `i=1`) \
  Assigns the value to all elements with the given index in the multi valued attribute.
- `multiValued` (or `mV`) \
  Makes the attribute multi valued.
- `zero` (or `0`) \
  Allow (Go) zero values (0 for integers for example).
- `ignore` (or `!`) \
  Skips the field entirely.

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

