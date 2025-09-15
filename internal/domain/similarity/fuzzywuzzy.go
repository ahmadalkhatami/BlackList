package similarity

import (
	"BlackListWorker/internal/utils"

	"github.com/agext/levenshtein"
)

// FuzzyWuzzyCalculator menghitung similarity ala FuzzyWuzzy (token-based).
type fuzzyWuzzyCalculator struct{}

func newFuzzyWuzzyCalculator() *fuzzyWuzzyCalculator {
	return &fuzzyWuzzyCalculator{}
}

func (c *fuzzyWuzzyCalculator) Calculate(a, b string) float64 {
	a = utils.Normalize(a)
	b = utils.Normalize(b)

	if a == "" || b == "" {
		return 0
	}

	// Hitung Levenshtein similarity
	d := levenshtein.Distance(a, b, nil)
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	if maxLen == 0 {
		return 1
	}
	return 1 - float64(d)/float64(maxLen)
}
