package models

type MatchingDetail struct {
	Id               int64
	MatchingResultId *int64
	FieldName        *string
	CustomerValue    *string
	WatchlistValue   *string
	FieldScore       *float64
	FieldWeight      *float64
	AlgorithmUsed    *string
}

func (m MatchingDetail) GetID() int64                { return m.Id }
func (m MatchingDetail) GetMatchingResultID() *int64 { return m.MatchingResultId }
func (m MatchingDetail) GetFieldName() *string       { return m.FieldName }
func (m MatchingDetail) GetCustomerValue() *string   { return m.CustomerValue }
func (m MatchingDetail) GetWatchlistValue() *string  { return m.WatchlistValue }
func (m MatchingDetail) GetFieldScore() *float64     { return m.FieldScore }
func (m MatchingDetail) GetFieldWeight() *float64    { return m.FieldWeight }
func (m MatchingDetail) GetAlgorithmUsed() *string   { return m.AlgorithmUsed }
