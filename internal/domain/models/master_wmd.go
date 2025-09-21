package models

import "time"

type MasterWMD struct {
	ID           int64
	Nama         string
	Alias1       *string
	Alias2       *string
	Alias3       *string
	Alias4       *string
	Alias5       *string
	Alias6       *string
	Alias7       *string
	Alias8       *string
	Alias9       *string
	Alias10      *string
	Type         *string
	TempatLahir  *string
	TanggalLahir *time.Time
	KTP          *string
	NPWP         *string
	NoPaspor     *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	IsActive     *bool
}

// Implement Watchlistable
func (m *MasterWMD) GetID() int64    { return m.ID }
func (m *MasterWMD) GetNama() string { return m.Nama }
func (m *MasterWMD) GetAliases() []string {
	var aliases []string
	for _, a := range []*string{
		m.Alias1, m.Alias2, m.Alias3, m.Alias4, m.Alias5,
		m.Alias6, m.Alias7, m.Alias8, m.Alias9, m.Alias10,
	} {
		if a != nil && *a != "" {
			aliases = append(aliases, *a)
		}
	}
	return aliases
}
func (m *MasterWMD) GetTempatLahir() *string     { return m.TempatLahir }
func (m *MasterWMD) GetTanggalLahir() *time.Time { return m.TanggalLahir }
func (m *MasterWMD) GetKTP() *string             { return m.KTP }
func (m *MasterWMD) GetNPWP() *string            { return m.NPWP }
func (m *MasterWMD) GetNoPaspor() *string        { return m.NoPaspor }
func (m *MasterWMD) GetSource() *string {
	src := "MASTER_WMD"
	return &src
}
func (m *MasterWMD) GetType() *string         { return m.Type }
func (m *MasterWMD) GetCreatedAt() *time.Time { return m.CreatedAt }
func (m *MasterWMD) GetIsActive() *bool       { return m.IsActive }
