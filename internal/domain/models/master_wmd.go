package models

import "time"

type MasterWMD struct {
	Id           string // use NEWSEQUENTIALID() in sqlserver for unique like uuid but more sequential for better indexing
	Nama         string
	Alias1       string
	Alias2       string
	Alias3       string
	Alias4       string
	Alias5       string
	Alias6       string
	Alias7       string
	Alias8       string
	Alias9       string
	Alias10      string
	Type         string // INDIVIDU | CORPORATE
	TempatLahir  string
	TanggalLahir time.Time
	KTP          string
	NPWP         string
	NoPaspor     string
	CreatedAt    time.Time
	IsActive     bool
}
