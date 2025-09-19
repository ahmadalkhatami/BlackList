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
	UpdatedAt    *time.Time
	IsActive     *bool
}

func (m MasterTeroris) GetID() int64    { return m.ID }
func (m MasterTeroris) GetNama() string { return m.Nama }
func (m MasterTeroris) GetAliases() []string {
	var aliases []string
	for _, a := range []*string{m.Alias1, m.Alias2, m.Alias3, m.Alias4} {
		if a != nil && *a != "" {
			aliases = append(aliases, *a)
		}
	}
	return aliases
}
func (m MasterTeroris) GetTempatLahir() *string     { return m.TempatLahir }
func (m MasterTeroris) GetTanggalLahir() *time.Time { return m.TanggalLahir }
func (m MasterTeroris) GetKTP() *string             { return m.KTP }
func (m MasterTeroris) GetNPWP() *string            { return m.NPWP }
func (m MasterTeroris) GetNoPaspor() *string        { return m.NoPaspor }
func (m MasterTeroris) GetType() *string            { return m.Type }
func (m MasterTeroris) GetCreatedAt() *time.Time    { return m.CreatedAt }
func (m MasterTeroris) GetUpdatedAt() *time.Time    { return m.UpdatedAt }
func (m MasterTeroris) GetIsActive() *bool          { return m.IsActive }
