package models

import "time"

type CIF struct {
	ID            int64
	CIFNumber     string
	NamaNasabah   string
	TempatLahir   string
	TanggalLahir  time.Time
	KTP           string
	NPWP          string
	NoPaspor      string
	StatusNasabah string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
