package similarity

import (
	"BlackListWorker/internal/utils"
	"strings"
)

// TokenCalculator menghitung similarity berbasis token (kata-per-kata).
type tokenCalculator struct{}

func newTokenCalculator() *tokenCalculator {
	return &tokenCalculator{}
}

func (c *tokenCalculator) Calculate(a, b string) float64 {
	a = utils.Normalize(a)
	b = utils.Normalize(b)

	if a == "" || b == "" {
		return 0
	}

	tokensA := strings.Fields(a)
	tokensB := strings.Fields(b)

	if len(tokensA) == 0 || len(tokensB) == 0 {
		return 0
	}

	// Hitung jumlah token yang sama
	match := 0
	tokenSetB := make(map[string]struct{})
	for _, t := range tokensB {
		tokenSetB[t] = struct{}{}
	}
	for _, t := range tokensA {
		if _, ok := tokenSetB[t]; ok {
			match++
		}
	}

	// Skor = jumlah token yang sama / rata-rata jumlah token
	avgLen := float64(len(tokensA)+len(tokensB)) / 2
	return float64(match) / avgLen
}
