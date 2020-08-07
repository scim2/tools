package marshal

import "unicode"

func lowerFirstRune(s string) string {
	for i, v := range s {
		return string(unicode.ToLower(v)) + s[i+1:]
	}
	return s
}
