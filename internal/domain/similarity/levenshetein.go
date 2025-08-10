package similarity

import "github.com/agext/levenshtein"

type LevenshteinCalculator struct {
	MaxDistance int
}

func NewLevenstheinCalculator(maxDistance int) Calculator {
	return &LevenshteinCalculator{MaxDistance: maxDistance}
}

func (c *LevenshteinCalculator) Calculate(a, b string) float64 {
	d := levenshtein.Distance(a, b, nil)

	if c.MaxDistance > 0 && d > c.MaxDistance {
		return 0
	}

	maxLen := len(a)

	if len(b) > maxLen {
		maxLen = len(b)
	}

	if maxLen == 0 {
		return 1
	}

	return 1 - float64(d)/float64(maxLen)
}
