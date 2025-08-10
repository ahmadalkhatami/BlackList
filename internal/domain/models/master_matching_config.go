package models

type MasterMatchingConfig struct {
	Id                string // use NEWSEQUENTIALID() in sqlserver for unique like uuid but more sequential for better indexing
	MatchingId        string
	FieldName         string
	FieldWeight       int16
	MatchingAlgorithm string
	IsActive          bool
}
