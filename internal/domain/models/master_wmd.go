package models

import "time"

type MasterWMD struct {
	Id           int64 // use NEWSEQUENTIALID() in sqlserver for unique like uuid but more sequential for better indexing
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
	Type         string // INDIVIDU | CORPORATE
	TempatLahir  *string
	TanggalLahir *time.Time
	KTP          *string
	NPWP         *string
	NoPaspor     *string
	CreatedAt    time.Time
	IsActive     bool
}

func (m MasterWMD) GetID() int64                { return m.Id }
func (m MasterWMD) GetNama() string             { return m.Nama }
func (m MasterWMD) GetTempatLahir() *string     { return m.TempatLahir }
func (m MasterWMD) GetTanggalLahir() *time.Time { return m.TanggalLahir }
func (m MasterWMD) GetKTP() *string             { return m.KTP }
func (m MasterWMD) GetNPWP() *string            { return m.NPWP }
func (m MasterWMD) GetNoPaspor() *string        { return m.NoPaspor }
func (m MasterWMD) GetCreatedAt() time.Time     { return m.CreatedAt }
func (m MasterWMD) GetIsActive() bool           { return m.IsActive }
func (m MasterWMD) GetSource() string           { return "MASTER_WMD" }
