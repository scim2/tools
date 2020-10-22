package marshal

import "unicode"

// lowerFirstRune lowers the first rune of a string.
// e.g. "UserName" into "userName"
func lowerFirstRune(s string) string {
	for i, v := range s {
		return string(unicode.ToLower(v)) + s[i+1:]
	}
	return s
}
