package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Normalize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

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

func FormatName(input string, lang language.Tag) string {

	str := strings.ReplaceAll(input, "_", " ")

	str = strings.ToLower(str)

	caser := cases.Title(lang)

	return caser.String(str)
}
