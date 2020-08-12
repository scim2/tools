package marshal

import "unicode"

// ensureMapInMap return the sub map with given key, creates it if not present.
func ensureMapInMap(key string, m map[string]map[string]interface{}) map[string]interface{} {
	sub, ok := m[key]
	if !ok {
		m[key] = make(map[string]interface{})
		sub = m[key]
	}
	return sub
}

// lowerFirstRune lowers the first rune of a string.
// e.g. "UserName" into "userName"
func lowerFirstRune(s string) string {
	for i, v := range s {
		return string(unicode.ToLower(v)) + s[i+1:]
	}
	return s
}
