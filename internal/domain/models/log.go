package models

import "time"

type ProcessingLog struct {
	Id           int64
	BatchId      *int64
	LogLevel     *string
	LogMessage   *string
	ErrorDetails *string
	LogTime      *time.Time
}

func (p ProcessingLog) GetID() int64             { return p.Id }
func (p ProcessingLog) GetBatchID() *int64       { return p.BatchId }
func (p ProcessingLog) GetLogLevel() *string     { return p.LogLevel }
func (p ProcessingLog) GetLogMessage() *string   { return p.LogMessage }
func (p ProcessingLog) GetErrorDetails() *string { return p.ErrorDetails }
func (p ProcessingLog) GetLogTime() *time.Time   { return p.LogTime }
