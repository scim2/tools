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

> **!** the current implementation does not marshal recursively.

```go
type ResourceStruct struct {
	Name string `scim:"userName"`
}

resourceStruct := ResourceStruct{
    Name: "di-wu",
}

resource, _ := Marshal(resourceStruct)

// OUTPUT: map[userName:di-wu]
```

