package models

import "time"

type MatchingResult struct {
	ID              int64
	BatchID         int64
	CIFNumber       string
	CustomerName    string
	WatchlistID     int64
	WatchlistSource string
	SimilarityScore float64
	Status          string // PENDING, REVIEWED, etc.
	ProcessDate     time.Time
	ProcessTime     time.Time
	CreatedAt       time.Time
}
