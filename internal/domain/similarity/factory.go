package similarity

import "errors"

var ErrUnknownAlgorithm = errors.New("unknown similarity algorithm")

// NewCalculator memilih algoritma berdasarkan nama.
func NewCalculator(algorithm string, opts ...interface{}) (Calculator, error) {
	switch algorithm {
	case "levenshtein":
		max := 0
		if len(opts) > 0 {
			if v, ok := opts[0].(int); ok {
				max = v
			}
		}
		return NewLevenshteinCalculator(max), nil

	case "jaro_winkler":
		boost := 0.7
		prefix := 4
		if len(opts) > 0 {
			if v, ok := opts[0].(float64); ok {
				boost = v
			}
		}
		if len(opts) > 1 {
			if v, ok := opts[1].(int); ok {
				prefix = v
			}
		}
		return NewJaroWinklerCalculator(boost, prefix), nil
	}

	return nil, ErrUnknownAlgorithm
}
