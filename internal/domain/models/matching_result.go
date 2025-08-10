package models

import "time"

type MatchingResult struct {
	ID              string
	BatchID         string
	CIFNumber       string
	CustomerName    string
	WatchlistID     string
	WatchlistSource string
	SimilarityScore float64
	Status          string // PENDING, REVIEWED, etc.
	ProcessDate     time.Time
	ProcessTime     time.Time
	CreatedAt       time.Time
}
