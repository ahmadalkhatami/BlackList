package similarity

import "BlackListWorker/internal/utils"

// ExactCalculator menghitung similarity exact value setelah normalisasi
type exactCalculator struct{}

func newExactCalculator() *exactCalculator {
	return &exactCalculator{}
}

func (c *exactCalculator) Calculate(a, b string) float64 {
	a = utils.Normalize(a)
	b = utils.Normalize(b)

	if a == b {
		return 1
	}
	return 0
}
