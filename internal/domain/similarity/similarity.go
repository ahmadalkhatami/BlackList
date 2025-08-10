package similarity

type Calculator interface {
	Calculate(a, b string) float64
}
