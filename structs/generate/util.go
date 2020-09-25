package gen

import "strings"

// cap capitalizes every word and removes the spaces.
func cap(s string) string {
	return strings.Replace(strings.Title(s), " ", "", -1)
}

// comment adds two slashed at the beginning of every line.
func comment(s string) string {
	s = strings.TrimSpace(s)
	return "// " + strings.Replace(s, "\n", "\n// ", -1)
}

// wrap splits the given strings in lines of maximum 120 characters.
func wrap(s string) string {
	var wrapped string
	for _, line := range strings.Split(s, "\n") {
		if len(line) <= 120 {
			wrapped += line
			continue
		}

		var wrappedLine string
		for _, word := range strings.Split(line, " ") {
			if len(wrappedLine)+len(word) <= 120 {
				if wrappedLine != "" {
					wrappedLine += " "
				}
				wrappedLine += word
			} else {
				wrapped += "\n" + wrappedLine
				wrappedLine = ""
			}
		}
		wrapped += "\n" + wrappedLine
	}
	return wrapped
}
