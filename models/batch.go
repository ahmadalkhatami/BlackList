package models

import "time"

type BatchProcessing struct {
	ID              int64
	ProcessType     string // MANUAL or SCHEDULED
	Status          string // RUNNING, COMPLETED, etc.
	TotalRecords    int
	ProcessedRecords int
	MatchedRecords  int
	StartTime       time.Time
	EndTime         time.Time
	ErrorMessage    string
	InitiatedBy     int64
	FilePath        string
}
