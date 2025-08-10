package models

type MatchingDetail struct {
	ID               int64
	MatchingResultID int64
	FieldName        string
	CustomerValue    string
	WatchlistValue   string
	FieldScore       float64
	FieldWeight      float64
	AlgorithmUsed    string
}
