package models

type MasterMatching struct {
	ID              int64
	Name            *string
	WatchlistSource *string
	IsIndividual    *bool
	Description     *string
	IsActive        *bool
}

func (m *MasterMatching) GetID() int64                { return m.ID }
func (m *MasterMatching) GetName() *string            { return m.Name }
func (m *MasterMatching) GetWatchlistSource() *string { return m.WatchlistSource }
func (m *MasterMatching) GetIsIndividual() *bool      { return m.IsIndividual }
func (m *MasterMatching) GetDescription() *string     { return m.Description }
func (m *MasterMatching) GetIsActive() *bool          { return m.IsActive }
