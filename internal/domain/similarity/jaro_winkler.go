package similarity

import "github.com/xrash/smetrics"

type JaroWingklerCalculator struct {
	BoostThreshold float64
	PrefixSize     int
}

func NewJaroWinklerCalculator(boost float64, prefix int) Calculator {
	return &JaroWingklerCalculator{BoostThreshold: boost, PrefixSize: prefix}
}

func (c *JaroWingklerCalculator) Calculate(a, b string) float64 {
	return smetrics.JaroWinkler(a, b, c.BoostThreshold, c.PrefixSize)
}
