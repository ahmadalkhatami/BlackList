package utils

import (
	"sort"
	"strings"
)

// Normalize : ubah ke lowercase, trim, tokenisasi, lalu urutkan token
// agar robust terhadap perbedaan urutan kata.
func Normalize(s string) string {
	parts := strings.Fields(strings.ToLower(strings.TrimSpace(s)))
	sort.Strings(parts)
	return strings.Join(parts, " ")
}
