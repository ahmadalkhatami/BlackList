package models

import "time"

type ProcessingLog struct {
	ID          int64
	BatchID     int64
	LogLevel    string // INFO, WARN, ERROR
	LogMessage  string
	ErrorDetail string
	LogTime     time.Time
}
