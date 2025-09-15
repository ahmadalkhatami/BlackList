package similarity

import "github.com/xrash/smetrics"

// JaroWinklerCalculator menghitung similarity berbasis Jaro-Winkler.
type jaroWinklerCalculator struct {
	BoostThreshold float64
	PrefixSize     int
}

func newJaroWinklerCalculator(boost float64, prefix int) *jaroWinklerCalculator {
	return &jaroWinklerCalculator{BoostThreshold: boost, PrefixSize: prefix}
}

func (c *jaroWinklerCalculator) Calculate(a, b string) float64 {
	return smetrics.JaroWinkler(a, b, c.BoostThreshold, c.PrefixSize)
}
