package models

import "time"

type MasterMatchingConfig struct {
	ID                int64
	MatchingID        int64
	FieldName         *string
	FieldWeight       *float64
	MatchingAlgorithm *string
	IsActive          *bool
	CreatedBy         *int
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
}

func (m *MasterMatchingConfig) GetID() int64                  { return m.ID }
func (m *MasterMatchingConfig) GetMatchingID() int64          { return m.MatchingID }
func (m *MasterMatchingConfig) GetFieldName() *string         { return m.FieldName }
func (m *MasterMatchingConfig) GetFieldWeight() *float64      { return m.FieldWeight }
func (m *MasterMatchingConfig) GetMatchingAlgorithm() *string { return m.MatchingAlgorithm }
func (m *MasterMatchingConfig) GetIsActive() *bool            { return m.IsActive }
func (m *MasterMatchingConfig) GetCreatedBy() *int            { return m.CreatedBy }
func (m *MasterMatchingConfig) GetCreatedAt() *time.Time      { return m.CreatedAt }
func (m *MasterMatchingConfig) GetUpdatedAt() *time.Time      { return m.UpdatedAt }
