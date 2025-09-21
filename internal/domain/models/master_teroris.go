package models

import "time"

type MasterTeroris struct {
	ID           int64
	Nama         string
	Alias1       *string
	Alias2       *string
	Alias3       *string
	Alias4       *string
	Type         *string
	TempatLahir  *string
	TanggalLahir *time.Time
	KTP          *string
	NPWP         *string
	NoPaspor     *string
	CreatedAt    *time.Time
	IsActive     *bool
}

// Implement Watchlistable
func (m *MasterTeroris) GetID() int64    { return m.ID }
func (m *MasterTeroris) GetNama() string { return m.Nama }
func (m *MasterTeroris) GetAliases() []string {
	var aliases []string
	if m.Alias1 != nil {
		aliases = append(aliases, *m.Alias1)
	}
	if m.Alias2 != nil {
		aliases = append(aliases, *m.Alias2)
	}
	if m.Alias3 != nil {
		aliases = append(aliases, *m.Alias3)
	}
	if m.Alias4 != nil {
		aliases = append(aliases, *m.Alias4)
	}
	return aliases
}
func (m *MasterTeroris) GetTempatLahir() *string     { return m.TempatLahir }
func (m *MasterTeroris) GetTanggalLahir() *time.Time { return m.TanggalLahir }
func (m *MasterTeroris) GetKTP() *string             { return m.KTP }
func (m *MasterTeroris) GetNPWP() *string            { return m.NPWP }
func (m *MasterTeroris) GetNoPaspor() *string        { return m.NoPaspor }
func (m *MasterTeroris) GetSource() *string {
	src := "MASTER_TERORIS"
	return &src
}
func (m *MasterTeroris) GetType() *string         { return m.Type }
func (m *MasterTeroris) GetCreatedAt() *time.Time { return m.CreatedAt }
func (m *MasterTeroris) GetIsActive() *bool       { return m.IsActive }
