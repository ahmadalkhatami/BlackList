package similarity

import "errors"

var ErrUnknownAlgorithm = errors.New("unknown similarity algorithm")

type Algorithm string

const (
	AlgoLevenshtein Algorithm = "levenshtein"
	AlgoJaroWinkler Algorithm = "jaro_winkler"
	AlgoFuzzyWuzzy  Algorithm = "fuzzywuzzy"
	AlgoCosine      Algorithm = "cosine"
)

type similarityOptions struct {
	maxDistance    int
	boostThreshold float64
	prefixSize     int
}

type Option func(*similarityOptions)

func defaultOptions() *similarityOptions {
	return &similarityOptions{
		maxDistance:    0,
		boostThreshold: 0.7,
		prefixSize:     4,
	}
}

func WithMaxDistance(d int) Option {
	return func(o *similarityOptions) { o.maxDistance = d }
}
func WithBoostThreshold(b float64) Option {
	return func(o *similarityOptions) { o.boostThreshold = b }
}
func WithPrefixSize(p int) Option {
	return func(o *similarityOptions) { o.prefixSize = p }
}

func NewCalculator(algo Algorithm, opts ...Option) (Calculator, error) {
	cfg := defaultOptions()
	for _, o := range opts {
		o(cfg)
	}

	switch algo {
	case AlgoLevenshtein:
		return newLevenshteinCalculator(cfg.maxDistance), nil
	case AlgoJaroWinkler:
		return newJaroWinklerCalculator(cfg.boostThreshold, cfg.prefixSize), nil
	case AlgoFuzzyWuzzy:
		return newFuzzyWuzzyCalculator(), nil
	case AlgoCosine:
		return newCosineCalculator(cfg.boostThreshold), nil
	default:
		return nil, ErrUnknownAlgorithm
	}
}
