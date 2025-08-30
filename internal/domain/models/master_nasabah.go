package models

import (
	"time"
)

type MasterNasabah struct {
	Id            int64
	CIFNumber     string
	NamaNasabah   string
	TempatLahir   *string
	TanggalLahir  *time.Time
	KTP           *string
	NPWP          *string
	NoPaspor      *string
	StatusNasabah string
	CreatedAt     time.Time
}

func (m *MasterNasabah) GetID() int64                { return m.Id }
func (m *MasterNasabah) GetNama() string             { return m.NamaNasabah }
func (m *MasterNasabah) GetTempatLahir() *string     { return m.TempatLahir }
func (m *MasterNasabah) GetTanggalLahir() *time.Time { return m.TanggalLahir }
func (m *MasterNasabah) GetKTP() *string             { return m.KTP }
func (m *MasterNasabah) GetNPWP() *string            { return m.NPWP }
func (m *MasterNasabah) GetNoPaspor() *string        { return m.NoPaspor }
func (m *MasterNasabah) GetCreatedAt() time.Time     { return m.CreatedAt }
func (m *MasterNasabah) GetIsActive() bool           { return m.StatusNasabah == "ACTIVE" }
func (m *MasterNasabah) GetSource() string           { return "NASABAH" }
