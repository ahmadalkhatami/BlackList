package models

import (
	"time"
)

type MasterLocalBlacklist struct {
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

func (m *MasterLocalBlacklist) GetID() int64                { return m.Id }
func (m *MasterLocalBlacklist) GetNama() string             { return m.Nama }
func (m *MasterLocalBlacklist) GetTempatLahir() *string     { return m.TempatLahir }
func (m *MasterLocalBlacklist) GetTanggalLahir() *time.Time { return m.TanggalLahir }
func (m *MasterLocalBlacklist) GetKTP() *string             { return m.KTP }
func (m *MasterLocalBlacklist) GetNPWP() *string            { return m.NPWP }
func (m *MasterLocalBlacklist) GetNoPaspor() *string        { return m.NoPaspor }
func (m *MasterLocalBlacklist) GetCreatedAt() time.Time     { return m.CreatedAt }
func (m *MasterLocalBlacklist) GetIsActive() bool           { return m.IsActive }
func (m *MasterLocalBlacklist) GetType() string             { return m.Type }	
func (m *MasterLocalBlacklist) GetSource() string           { return "MASTER_LOCAL_BLACKLIST" }
