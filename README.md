# SCIM Tools
This repository/module contains various utility functions to make it easier to work with SCIM servers and clients.

## Fuzzer
Build on top of [gofuzz](https://github.com/google/gofuzz/).

> **!** the current implementation only fuzzes required fields.

```go
var refSchema ReferenceSchema
// define the reference schema yourself or use unmarshal json

// displayName and name.givenName can never be empty.
s.NeverEmpty("displayName", "name.givenName")

f := fuzz.New().Funcs(
    NewResourceFuzzer(s), 
)

var resource Resource
f.Fuzz(&resource)

// OUTPUT: map[displayName:vWKdUsVprh userName:RFlLpsMnBW]
```

