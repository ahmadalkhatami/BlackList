package similarity

import "errors"

var ErrUnknownAlgorithm = errors.New("unknown similarity algorithm")

// Enum2Man
const (
	AlgoLevenshtein Algorithm = "levenshtein"
	AlgoJaroWinkler Algorithm = "jaro_winkler"
	AlgoFuzzyWuzzy  Algorithm = "fuzzywuzzy"
)

type Algorithm string

// internal config untuk parameter opsional
type options struct {
	maxDistance    int
	boostThreshold float64
	prefixSize     int
}

// Option adalah fungsi yang bisa mengubah konfigurasi
type Option func(*options)

// Default values
func defaultOptions() *options {
	return &options{
		maxDistance:    0,
		boostThreshold: 0.7,
		prefixSize:     4,
	}
}

// Functional options
func WithMaxDistance(d int) Option {
	return func(o *options) { o.maxDistance = d }
}
func WithBoostThreshold(b float64) Option {
	return func(o *options) { o.boostThreshold = b }
}
func WithPrefixSize(p int) Option {
	return func(o *options) { o.prefixSize = p }
}

// NewCalculator membuat instance calculator sesuai algoritma yang dipilih.
// Bisa diberikan Option opsional untuk override default config.
func NewCalculator(algo Algorithm, opts ...Option) (Calculator, error) {
	cfg := defaultOptions()
	for _, o := range opts {
		o(cfg)
	}

	switch algo {
	case AlgoLevenshtein:
		return NewLevenshteinCalculator(cfg.maxDistance), nil
	case AlgoJaroWinkler:
		return NewJaroWinklerCalculator(cfg.boostThreshold, cfg.prefixSize), nil
	case AlgoFuzzyWuzzy:
		return NewFuzzyWuzzyCalculator(), nil
	default:
		return nil, ErrUnknownAlgorithm
	}
}
