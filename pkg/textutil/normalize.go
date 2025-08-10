package textutil

import (
	"strings"
	"unicode"
)

/* Normalize : trim, lower, remove duplicate spaces, remove non-letter/number */
func Normalize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	/* remove control characters, keep letters/numbers and spaces */
	var b []rune
	lastSpace := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b = append(b, r)
			lastSpace = false
		} else if unicode.IsSpace(r) {
			if !lastSpace {
				b = append(b, ' ')
				lastSpace = true
			}
		}
	}
	return strings.TrimSpace(string(b))
}
