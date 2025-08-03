package models

import "time"

type SystemConfig struct {
	ID          int64
	ConfigKey   string
	ConfigValue string
	ConfigType  string // DECIMAL, INTEGER, STRING
	Description string
	CreatedBy   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
