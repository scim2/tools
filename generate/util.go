package generate

import (
	"log"
	"regexp"
	"strings"
)

func keepAlpha(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(s, "")
}

func singular(s string) string {
	if strings.HasSuffix(s, "ses") {
		return strings.TrimSuffix(s, "es")
	}
	if strings.HasSuffix(s, "s") {
		return strings.TrimSuffix(s, "s")
	}
	return s
}

// cap capitalizes every word and removes the spaces.
func cap(s string) string {
	return strings.Replace(strings.Title(s), " ", "", -1)
}

// comment adds two slashed at the beginning of every line.
func comment(s string) string {
	s = strings.TrimSpace(s)
	return "// " + strings.Replace(s, "\n", "\n// ", -1)
}

// wrap splits the given strings in lines of maximum lw characters.
func wrap(s string, lw int) string {
	var wrapped string
	for _, line := range strings.Split(s, "\n") {
		if len(line) <= lw {
			wrapped += line
			continue
		}

		var wrappedLine string
		for _, word := range strings.Split(line, " ") {
			if len(wrappedLine)+len(word) <= lw {
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
