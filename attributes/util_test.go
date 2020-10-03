package attributes_test

import (
    "fmt"
    "github.com/scim2/tools/attributes"
)

func ExampleContains() {
    attrs := map[string]interface{}{
        "x": 0,
    }

    fmt.Println(attributes.Contains("x", attrs))
    fmt.Println(attributes.Contains("X", attrs))
    fmt.Println(attributes.Contains("y", attrs))

    // Output:
    // 0 true
    // 0 true
    // <nil> false
}

func ExampleGetBool() {
    attrs := map[string]interface{}{
        "x": true,
        "y": 0,
    }

    fmt.Println(attributes.GetBool("x", attrs))
    fmt.Println(attributes.GetBool("y", attrs))
    fmt.Println(attributes.GetBool("z", attrs))

    // Output:
    // true <nil>
    // false attribute "y" is not a bool
    // false could not find "z" in attributes
}

func ExampleGetFloat() {
    attrs := map[string]interface{}{
        "x": 0.1,
        "y": false,
    }

    fmt.Println(attributes.GetFloat("x", attrs))
    fmt.Println(attributes.GetFloat("y", attrs))
    fmt.Println(attributes.GetFloat("z", attrs))

    // Output:
    // 0.1 <nil>
    // 0 attribute "y" is not a float64
    // 0 could not find "z" in attributes
}

func ExampleGetFloatAsInt() {
    attrs := map[string]interface{}{
        "x": float64(1),
        "y": 0.1,
    }

    fmt.Println(attributes.GetFloatAsInt("x", attrs))
    fmt.Println(attributes.GetFloatAsInt("y", attrs))
    fmt.Println(attributes.GetFloatAsInt("z", attrs))

    // Output:
    // 1 <nil>
    // 0 attribute "y" is not a int
    // 0 could not find "z" in attributes
}

func ExampleGetMap() {
    attrs := map[string]interface{}{
        "x": map[string]interface{}{
            "z": true,
        },
        "y": 0,
    }

    fmt.Println(attributes.GetMap("x", attrs))
    fmt.Println(attributes.GetMap("y", attrs))
    fmt.Println(attributes.GetMap("z", attrs))

    // Output:
    // map[z:true] <nil>
    // map[] attribute "y" is not a map[string]interface{}
    // map[] could not find "z" in attributes
}

func ExampleGetString() {
    attrs := map[string]interface{}{
        "x": "x",
        "y": 0,
    }

    fmt.Println(attributes.GetString("x", attrs))
    fmt.Println(attributes.GetString("y", attrs))
    fmt.Println(attributes.GetString("z", attrs))

    // Output:
    // x <nil>
    //  attribute "y" is not a string
    //  could not find "z" in attributes
}

func ExampleGetStringInSubMap() {
    attrs := map[string]interface{}{
        "x": map[string]interface{}{
            "x": "x",
            "y": 0,
        },
    }

    fmt.Println(attributes.GetStringInSubMap("x", "x", attrs))
    fmt.Println(attributes.GetStringInSubMap("x", "y", attrs))
    fmt.Println(attributes.GetStringInSubMap("x", "z", attrs))
    fmt.Println(attributes.GetStringInSubMap("y", "x", attrs))

    // Output:
    // x <nil>
    //  attribute "y" is not a string
    //  could not find "z" in attributes
    //  could not find "y" in attributes
}
