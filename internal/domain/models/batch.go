package models

import "time"

type BatchProcessing struct {
	Id               int64
	ProcessType      *string
	Status           *string
	TotalRecords     *int
	ProcessedRecords *int
	MatchedRecords   *int
	StartTime        *time.Time
	EndTime          *time.Time
	ErrorMessage     *string
	InitiatedBy      *int
	FilePath         *string
}

func (b BatchProcessing) GetID() int64              { return b.Id }
func (b BatchProcessing) GetProcessType() *string   { return b.ProcessType }
func (b BatchProcessing) GetStatus() *string        { return b.Status }
func (b BatchProcessing) GetTotalRecords() *int     { return b.TotalRecords }
func (b BatchProcessing) GetProcessedRecords() *int { return b.ProcessedRecords }
func (b BatchProcessing) GetMatchedRecords() *int   { return b.MatchedRecords }
func (b BatchProcessing) GetStartTime() *time.Time  { return b.StartTime }
func (b BatchProcessing) GetEndTime() *time.Time    { return b.EndTime }
func (b BatchProcessing) GetErrorMessage() *string  { return b.ErrorMessage }
func (b BatchProcessing) GetInitiatedBy() *int      { return b.InitiatedBy }
func (b BatchProcessing) GetFilePath() *string      { return b.FilePath }
