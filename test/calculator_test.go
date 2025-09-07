package test

import (
	"BlackListWorker/internal/domain/similarity"
	"testing"
)

func TestLevenshteinCalculator(t *testing.T) {
	calc := similarity.NewLevenshteinCalculator(0)
	score := calc.Calculate("kitten", "sitting")
	if score <= 0 || score > 1 {
		t.Errorf("unexpected score for levenshtein: %f", score)
	}

	// dengan maxDistance override
	calc2 := similarity.NewLevenshteinCalculator(2)
	score2 := calc2.Calculate("kitten", "sitting")
	if score2 != 0 {
		t.Errorf("expected 0 when distance > maxDistance, got %f", score2)
	}
}

func TestJaroWinklerCalculator(t *testing.T) {
	calc := similarity.NewJaroWinklerCalculator(0.7, 4)
	score := calc.Calculate("martha", "marhta")
	if score <= 0 || score > 1 {
		t.Errorf("unexpected score for jaro-winkler: %f", score)
	}

	// custom boost & prefix
	calc2 := similarity.NewJaroWinklerCalculator(0.9, 6)
	score2 := calc2.Calculate("martha", "marhta")
	if score2 <= 0 {
		t.Errorf("expected positive similarity, got %f", score2)
	}
}

func TestFuzzyWuzzyCalculator(t *testing.T) {
	calc := similarity.NewFuzzyWuzzyCalculator()
	score := calc.Calculate("hello world", "helo wrld")
	if score <= 0 || score > 1 {
		t.Errorf("unexpected score for fuzzywuzzy: %f", score)
	}

	// empty input harus 0
	score2 := calc.Calculate("", "")
	if score2 != 0 {
		t.Errorf("expected 0 for empty strings, got %f", score2)
	}
}

func TestFactory(t *testing.T) {
	// default Levenshtein
	calc1, err := similarity.NewCalculator(similarity.AlgoLevenshtein)
	if err != nil {
		t.Fatal(err)
	}
	if score := calc1.Calculate("a", "b"); score < 0 {
		t.Errorf("unexpected score %f", score)
	}

	// override Levenshtein
	calc2, err := similarity.NewCalculator(similarity.AlgoLevenshtein, similarity.WithMaxDistance(1))
	if err != nil {
		t.Fatal(err)
	}
	if score := calc2.Calculate("abcd", "wxyz"); score != 0 {
		t.Errorf("expected 0 when distance > maxDistance, got %f", score)
	}

	// Jaro-Winkler default
	calc3, err := similarity.NewCalculator(similarity.AlgoJaroWinkler)
	if err != nil {
		t.Fatal(err)
	}
	if score := calc3.Calculate("martha", "marhta"); score <= 0 {
		t.Errorf("unexpected score %f", score)
	}

	// FuzzyWuzzy
	calc4, err := similarity.NewCalculator(similarity.AlgoFuzzyWuzzy)
	if err != nil {
		t.Fatal(err)
	}
	if score := calc4.Calculate("hello", "hallo"); score <= 0 {
		t.Errorf("unexpected score %f", score)
	}

	// unknown algo
	_, err = similarity.NewCalculator("unknown")
	if err == nil {
		t.Errorf("expected error for unknown algorithm")
	}
}
