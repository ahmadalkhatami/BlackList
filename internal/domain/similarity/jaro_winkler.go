package similarity

import "github.com/xrash/smetrics"

// JaroWinklerCalculator menghitung similarity berbasis Jaro-Winkler.
type JaroWinklerCalculator struct {
	BoostThreshold float64
	PrefixSize     int
}

func NewJaroWinklerCalculator(boost float64, prefix int) *JaroWinklerCalculator {
	return &JaroWinklerCalculator{BoostThreshold: boost, PrefixSize: prefix}
}

func (c *JaroWinklerCalculator) Calculate(a, b string) float64 {
	return smetrics.JaroWinkler(a, b, c.BoostThreshold, c.PrefixSize)
}
