package models

import "time"

type MasterTeroris struct {
	Id           int64 // use NEWSEQUENTIALID() in sqlserver for unique like uuid but more sequential for better indexing
	Nama         string
	Alias1       *string
	Alias2       *string
	Alias3       *string
	Alias4       *string
	Type         string // INDIVIDU | CORPORATE
	TempatLahir  *string
	TanggalLahir *time.Time
	KTP          *string
	NPWP         *string
	NoPaspor     *string
	CreatedAt    time.Time
	IsActive     bool
}

func (m *MasterTeroris) GetID() int64                { return m.Id }
func (m *MasterTeroris) GetNama() string             { return m.Nama }
func (m *MasterTeroris) GetTempatLahir() *string     { return m.TempatLahir }
func (m *MasterTeroris) GetTanggalLahir() *time.Time { return m.TanggalLahir }
func (m *MasterTeroris) GetKTP() *string             { return m.KTP }
func (m *MasterTeroris) GetNPWP() *string            { return m.NPWP }
func (m *MasterTeroris) GetNoPaspor() *string        { return m.NoPaspor }
func (m *MasterTeroris) GetCreatedAt() time.Time     { return m.CreatedAt }
func (m *MasterTeroris) GetIsActive() bool           { return m.IsActive }
func (m *MasterTeroris) GetType() string             { return m.Type }
func (m *MasterTeroris) GetSource() string           { return "MASTER_TERORIS" }
