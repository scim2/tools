package safe

import "github.com/scim2/tools/attributes"

// GetBool searches the given map for a boolean that matches the given id.
// Returns false if not found.
func GetBool(id string, a map[string]interface{}) bool {
    t, _ := attributes.GetBool(id, a)
    return t
}

// GetFloat searches the given map for a float that matches the given id.
// Returns 0.0 if not found.
func GetFloat(id string, a map[string]interface{}) float64 {
    f, _ := attributes.GetFloat(id, a)
    return f
}

// GetFloatAsInt searches the given map for a float that matches the given id and converts it to an int if possible.
// Returns 0 if the float is not a whole number or not found.
func GetFloatAsInt(id string, a map[string]interface{}) int {
    i, _ := attributes.GetFloatAsInt(id, a)
    return i
}

// GetMap searches the given map for a map that matches the given id.
// Returns nil if not found.
func GetMap(id string, a map[string]interface{}) map[string]interface{} {
    m, _ := attributes.GetMap(id, a)
    return m
}

// GetString searches the given map for a string that matches the given id.
// Returns an empty string if not found.
func GetString(id string, a map[string]interface{}) string {
    str, _ := attributes.GetString(id, a)
    return str
}

// GetStringInSubMap searches the given map for a string with key sID in the map matching the given mID.
// Returns an empty string if the map or string is not found.
func GetStringInSubMap(mID, sID string, a map[string]interface{}) string {
    str, _ := attributes.GetStringInSubMap(mID, sID, a)
    return str
}
