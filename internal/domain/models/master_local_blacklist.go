package models

import "time"

type MasterLocalBlacklist struct {
	Id           string // use NEWSEQUENTIALID() in sqlserver for unique like uuid but more sequential for better indexing
	Nama         string
	Alias1       string
	Alias2       string
	Alias3       string
	Alias4       string
	Type         string // INDIVIDU | CORPORATE
	TempatLahir  string
	TanggalLahir time.Time
	KTP          string
	NPWP         string
	NoPaspor     string
	CreatedAt    time.Time
	IsActive     bool
}
