package similarity

import (
	"math"
	"strings"
)

type cosineCalculator struct {
	boostThreshold float64
}

func newCosineCalculator(boostThreshold float64) Calculator {
	return &cosineCalculator{
		boostThreshold: boostThreshold,
	}
}

func computeTFWithBoost(text string, boost float64) map[string]float64 {
	tf := make(map[string]float64)
	words := strings.Fields(strings.ToLower(text))
	for _, w := range words {
		tf[w] += 1.0
	}
	for k := range tf {
		tf[k] /= float64(len(words))
		if tf[k] >= boost {
			tf[k] *= 1.2 // boost 20%
		}
	}
	return tf
}

func (c *cosineCalculator) Calculate(a, b string) float64 {
	tfA := computeTFWithBoost(a, c.boostThreshold)
	tfB := computeTFWithBoost(b, c.boostThreshold)

	var dot, magA, magB float64
	for k, vA := range tfA {
		dot += vA * tfB[k]
		magA += vA * vA
	}
	for _, v := range tfB {
		magB += v * v
	}
	if magA == 0 || magB == 0 {
		return 0
	}
	return dot / (math.Sqrt(magA) * math.Sqrt(magB))
}
