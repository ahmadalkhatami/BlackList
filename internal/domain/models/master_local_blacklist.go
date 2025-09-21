package models

import "time"

type MasterLocalBlacklist struct {
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
func (m *MasterLocalBlacklist) GetID() int64    { return m.ID }
func (m *MasterLocalBlacklist) GetNama() string { return m.Nama }
func (m *MasterLocalBlacklist) GetAliases() []string {
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
func (m *MasterLocalBlacklist) GetTempatLahir() *string     { return m.TempatLahir }
func (m *MasterLocalBlacklist) GetTanggalLahir() *time.Time { return m.TanggalLahir }
func (m *MasterLocalBlacklist) GetKTP() *string             { return m.KTP }
func (m *MasterLocalBlacklist) GetNPWP() *string            { return m.NPWP }
func (m *MasterLocalBlacklist) GetNoPaspor() *string        { return m.NoPaspor }
func (m *MasterLocalBlacklist) GetSource() *string {
	src := "MASTER_LOCAL_BLACKLIST"
	return &src
}
func (m *MasterLocalBlacklist) GetType() *string         { return m.Type }
func (m *MasterLocalBlacklist) GetCreatedAt() *time.Time { return m.CreatedAt }
func (m *MasterLocalBlacklist) GetIsActive() *bool       { return m.IsActive }
