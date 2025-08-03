package models

import "time"

type MatchingConfig struct {
	ID                int64
	WatchlistSource   string
	FieldName         string
	FieldWeight       float64
	MatchingAlgorithm string // e.g., "jaro_winkler"
	IsActive          bool
	CreatedBy         int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
