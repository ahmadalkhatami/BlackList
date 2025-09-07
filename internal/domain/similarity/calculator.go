package similarity

type Calculator interface {
	Calculate(a, b string) float64
}

/* How To Use

func main() {
	// 1. Default Levenshtein
	calc1, err := similarity.NewCalculator(similarity.AlgoLevenshtein)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Levenshtein (default):", calc1.Calculate("hello world", "helo wrld"))

	// 2. Levenshtein dengan override maxDistance
	calc2, _ := similarity.NewCalculator(similarity.AlgoLevenshtein, similarity.WithMaxDistance(2))
	fmt.Println("Levenshtein (maxDistance=2):", calc2.Calculate("hello world", "helooooo wrld"))

	// 3. Jaro-Winkler default
	calc3, _ := similarity.NewCalculator(similarity.AlgoJaroWinkler)
	fmt.Println("Jaro-Winkler (default):", calc3.Calculate("martha", "marhta"))

	// 4. Jaro-Winkler dengan override
	calc4, _ := similarity.NewCalculator(similarity.AlgoJaroWinkler,
		similarity.WithBoostThreshold(0.85),
		similarity.WithPrefixSize(6),
	)
	fmt.Println("Jaro-Winkler (override):", calc4.Calculate("martha", "marhta"))

	// 5. FuzzyWuzzy
	calc5, _ := similarity.NewCalculator(similarity.AlgoFuzzyWuzzy)
	fmt.Println("FuzzyWuzzy:", calc5.Calculate("hello world", "helo wrld"))
}

*/
