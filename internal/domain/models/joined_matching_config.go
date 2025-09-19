package models

type JoinedMatchingConfig struct {
	Id                int64
	MatchingId        int64
	WatchlistSource   string
	Type              bool
	FieldName         string
	FieldWeight       float64
	MatchingAlgorithm string
	IsActive          bool
}
