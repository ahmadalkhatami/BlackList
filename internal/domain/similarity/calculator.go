package similarity

type Calculator interface {
	Calculate(a, b string) float64
}

/* How To Use
calc, _ := similarity.NewCalculator("levenshtein", 5)
score := calc.Calculate("Jokowi", "Joko Widodo")
fmt.Println("Levenshtein score:", score)

calc2, _ := similarity.NewCalculator("jaro_winkler", 0.8, 4)
score2 := calc2.Calculate("Trump", "Tramp")
fmt.Println("Jaro-Winkler score:", score2)
*/
