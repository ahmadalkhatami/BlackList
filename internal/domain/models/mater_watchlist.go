package models

import "time"

type MasterWatchlist struct {
	ID   int64
	Nama string
	// Aliases      map[string]interface{}
	Aliases      []string
	TempatLahir  string
	TanggalLahir time.Time
	KTP          string
	NPWP         string
	NoPaspor     string
	Source       string // "DTTOT", "WMD", or "LOCAL_BLACKLIST"
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	IsActive     bool
}
