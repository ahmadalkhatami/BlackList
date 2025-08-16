package models

type MasterMatchingConfig struct {
	// use NEWSEQUENTIALID() in sqlserver for unique like uuid but more sequential for better indexing
	Id                string
	MatchingId        string
	FieldName         string
	FieldWeight       float64
	MatchingAlgorithm string
	IsActive          bool
}
