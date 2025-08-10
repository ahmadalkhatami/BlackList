package models

import "time"

type MasterNasabah struct {
	Id            string // use NEWSEQUENTIALID() in sqlserver for unique like uuid but more sequential for better indexing
	CIFNumber     string
	NamaNasabah   string
	TempatLahir   string
	TanggalLahir  time.Time
	KTP           string
	NPWP          string
	NoPaspor      string
	StatusNasabah string // is different with IsActive ?
	CreatedAt     time.Time
}
