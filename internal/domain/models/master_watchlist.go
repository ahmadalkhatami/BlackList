package models

import (
	"time"
)

// Struct MasterWatchlist
type MasterWatchlist struct {
	ID           int64
	Nama         string
	Aliases      []string
	TempatLahir  *string
	TanggalLahir *time.Time
	KTP          *string
	NPWP         *string
	NoPaspor     *string
	Source       string
	CreatedAt    time.Time
	IsActive     bool
}

type Watchlistable interface {
	GetID() int64
	GetNama() string
	GetTempatLahir() *string
	GetTanggalLahir() *time.Time
	GetKTP() *string
	GetNPWP() *string
	GetNoPaspor() *string
	GetSource() string
	GetCreatedAt() time.Time
	GetIsActive() bool
}
