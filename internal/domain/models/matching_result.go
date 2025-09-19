package models

import "time"

type MatchingResult struct {
	Id              int64
	BatchId         *int64
	CifNumber       *string
	CustomerName    *string
	WatchlistId     *int64
	WatchlistSource *string
	SimilarityScore *float64
	Status          *string
	ProcessDate     *time.Time
	ProcessTime     *time.Time
	CreatedAt       *time.Time
}

func (m MatchingResult) GetID() int64                 { return m.Id }
func (m MatchingResult) GetBatchID() *int64           { return m.BatchId }
func (m MatchingResult) GetCifNumber() *string        { return m.CifNumber }
func (m MatchingResult) GetCustomerName() *string     { return m.CustomerName }
func (m MatchingResult) GetWatchlistId() *int64       { return m.WatchlistId }
func (m MatchingResult) GetWatchlistSource() *string  { return m.WatchlistSource }
func (m MatchingResult) GetSimilarityScore() *float64 { return m.SimilarityScore }
func (m MatchingResult) GetStatus() *string           { return m.Status }
func (m MatchingResult) GetProcessDate() *time.Time   { return m.ProcessDate }
func (m MatchingResult) GetProcessTime() *time.Time   { return m.ProcessTime }
func (m MatchingResult) GetCreatedAt() *time.Time     { return m.CreatedAt }
