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
	Type         string // INDIVIDU | CORPORATE
	CreatedAt    time.Time
	IsActive     bool
}

type Watchlistable interface {
	GetID() int64
	GetNama() string
	GetAliases() []string
	GetTempatLahir() *string
	GetTanggalLahir() *time.Time
	GetKTP() *string
	GetNPWP() *string
	GetNoPaspor() *string
	GetSource() string
	GetType() string
	GetCreatedAt() time.Time
	GetIsActive() bool
}

func (m *MasterWatchlist) GetID() int64                { return m.ID }
func (m *MasterWatchlist) GetNama() string             { return m.Nama }
func (m *MasterWatchlist) GetAliases() []string        { return m.Aliases }
func (m *MasterWatchlist) GetTempatLahir() *string     { return m.TempatLahir }
func (m *MasterWatchlist) GetTanggalLahir() *time.Time { return m.TanggalLahir }
func (m *MasterWatchlist) GetKTP() *string             { return m.KTP }
func (m *MasterWatchlist) GetNPWP() *string            { return m.NPWP }
func (m *MasterWatchlist) GetNoPaspor() *string        { return m.NoPaspor }
func (m *MasterWatchlist) GetSource() string           { return m.Source }
func (m *MasterWatchlist) GetType() string             { return m.Type }
func (m *MasterWatchlist) GetCreatedAt() time.Time     { return m.CreatedAt }
func (m *MasterWatchlist) GetIsActive() bool           { return m.IsActive }

/* Getter untuk masing-masing alias
func (m *MasterWatchlist) GetAlias1() string { return getAlias(m.Aliases, 0) }
func (m *MasterWatchlist) GetAlias2() string { return getAlias(m.Aliases, 1) }
func (m *MasterWatchlist) GetAlias3() string { return getAlias(m.Aliases, 2) }
func (m *MasterWatchlist) GetAlias4() string { return getAlias(m.Aliases, 3) }
func (m *MasterWatchlist) GetAlias5() string { return getAlias(m.Aliases, 4) }
func (m *MasterWatchlist) GetAlias6() string { return getAlias(m.Aliases, 5) }
func (m *MasterWatchlist) GetAlias7() string { return getAlias(m.Aliases, 6) }
func (m *MasterWatchlist) GetAlias8() string { return getAlias(m.Aliases, 7) }
func (m *MasterWatchlist) GetAlias9() string { return getAlias(m.Aliases, 8) }
func (m *MasterWatchlist) GetAlias10() string { return getAlias(m.Aliases, 9) }

func getAlias(aliases []string, idx int) string {
	if idx >= 0 && idx < len(aliases) {
		return aliases[idx]
	}
	return ""
} */
