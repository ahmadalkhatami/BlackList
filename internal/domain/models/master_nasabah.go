package models

import (
	"time"
)

type MasterNasabah struct {
	ID            int64
	CifNumber     string
	NamaNasabah   string
	TempatLahir   *string
	TanggalLahir  *time.Time
	KTP           *string
	NPWP          *string
	NoPaspor      *string
	StatusNasabah *string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

func (m MasterNasabah) GetID() int64                { return m.ID }
func (m MasterNasabah) GetCifNumber() string        { return m.CifNumber }
func (m MasterNasabah) GetNamaNasabah() string      { return m.NamaNasabah }
func (m MasterNasabah) GetTempatLahir() *string     { return m.TempatLahir }
func (m MasterNasabah) GetTanggalLahir() *time.Time { return m.TanggalLahir }
func (m MasterNasabah) GetKTP() *string             { return m.KTP }
func (m MasterNasabah) GetNPWP() *string            { return m.NPWP }
func (m MasterNasabah) GetNoPaspor() *string        { return m.NoPaspor }
func (m MasterNasabah) GetStatusNasabah() *string   { return m.StatusNasabah }
func (m MasterNasabah) GetCreatedAt() *time.Time    { return m.CreatedAt }
func (m MasterNasabah) GetUpdatedAt() *time.Time    { return m.UpdatedAt }
