package models

type JoinedMatchingConfig struct {
	Id                string
	WatchlistSource   string
	Type              bool
	MatchingId        string
	FieldName         string
	FieldWeight       float64
	MatchingAlgorithm string
	IsActive          bool
}
