package attributes

import (
    "fmt"
    "strings"
)

var (
    errNotFound = func(id string) error {
        return fmt.Errorf("could not find %q in attributes", id)
    }
    errInvalid = func(id, typ string) error {
        return fmt.Errorf("attribute %q is not a %s", id, typ)
    }
)

// Contains checks whether the given map contains the id. This check is case insensitive!
func Contains(id string, a map[string]interface{}) (interface{}, bool) {
    id = strings.ToLower(id)
    for k, v := range a {
        if id == strings.ToLower(k) {
            return v, true
        }
    }
    return nil, false
}

// GetBool searches the given map for a boolean that matches the given id.
func GetBool(id string, a map[string]interface{}) (bool, error) {
    i, found := Contains(id, a)
    if !found {
        return false, errNotFound(id)
    }
    t, ok := i.(bool)
    if !ok {
        return false, errInvalid(id, "bool")
    }
    return t, nil
}

// GetFloat searches the given map for a float that matches the given id.
func GetFloat(id string, a map[string]interface{}) (float64, error) {
    i, found := Contains(id, a)
    if !found {
        return 0, errNotFound(id)
    }
    f, ok := i.(float64)
    if !ok {
        return 0, errInvalid(id, "float64")
    }
    return f, nil
}

// GetFloatAsInt searches the given map for a float that matches the given id and converts it to an int if possible.
func GetFloatAsInt(id string, a map[string]interface{}) (int, error) {
    f, err := GetFloat(id, a)
    if err != nil {
        return 0, err
    }
    if f != float64(int64(f)) {
        return 0, errInvalid(id, "int")
    }
    return int(f), nil
}

// GetMap searches the given map for a map that matches the given id.
func GetMap(id string, a map[string]interface{}) (map[string]interface{}, error) {
    i, found := Contains(id, a)
    if !found {
        return nil, errNotFound(id)
    }
    m, ok := i.(map[string]interface{})
    if !ok {
        return nil, errInvalid(id, "map[string]interface{}")
    }
    return m, nil
}

// GetString searches the given map for a string that matches the given id.
func GetString(id string, a map[string]interface{}) (string, error) {
    i, found := Contains(id, a)
    if !found {
        return "", errNotFound(id)
    }
    str, ok := i.(string)
    if !ok {
        return "", errInvalid(id, "string")
    }
    return str, nil
}

// GetStringInSubMap searches the given map for a string with key sID in the map matching the given mID.
func GetStringInSubMap(mID, sID string, a map[string]interface{}) (string, error) {
    m, err := GetMap(mID, a)
    if err != nil {
        return "", err
    }
    str, err := GetString(sID, m)
    if err != nil {
        return "", err
    }
    return str, nil
}
