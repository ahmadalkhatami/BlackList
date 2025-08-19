package services

import (
	"BlackListWorker/models"
	"database/sql"
	"fmt"
)

// LoadWatchlistGeneric memuat data dari tabel watchlist tertentu
func LoadWatchlistGeneric(db *sql.DB, tableName, source string, aliasCount int) ([]models.MasterWatchlist, error) {
	// Buat list kolom alias dinamis
	aliasCols := ""
	for i := 1; i <= aliasCount; i++ {
		if i > 1 {
			aliasCols += ", "
		}
		aliasCols += fmt.Sprintf("Alias%d", i)
	}

	// Query builder
	query := fmt.Sprintf(`
		SELECT Id, Nama, %s, TempatLahir, TanggalLahir, Ktp, Npwp, NoPaspor, 
		       CreatedAt, UpdatedAt, IsActive 
		FROM %s
		WHERE IsActive = 1 AND Nama = 'AHMADElijah Corkery'
	`, aliasCols, tableName)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MasterWatchlist

	// Loop data hasil query
	for rows.Next() {
		var m models.MasterWatchlist

		// Buat array dynamic untuk Scan()
		aliasVals := make([]sql.NullString, aliasCount)
		scanArgs := []any{&m.ID, &m.Nama}
		for i := range aliasVals {
			scanArgs = append(scanArgs, &aliasVals[i])
		}
		scanArgs = append(scanArgs,
			&m.TempatLahir, &m.TanggalLahir,
			&m.KTP, &m.NPWP, &m.NoPaspor,
			&m.CreatedAt, &m.UpdatedAt, &m.IsActive,
		)

		// Scan data
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		// Gabungkan alias
		m.Alias = collectAliases(aliasVals)
		m.Source = source

		result = append(result, m)
	}

	return result, nil
}

// LoadAllWatchlists menggabungkan semua sumber watchlist
func LoadAllWatchlists(db *sql.DB) ([]models.MasterWatchlist, error) {
	var all []models.MasterWatchlist

	loaders := []struct {
		table      string
		source     string
		aliasCount int
	}{
		{"MASTER_TERORIS", "MASTER_TERORIS", 4},
		{"MASTER_WMD", "WMD", 10},
		{"MASTER_LOCAL_BLACKLIST", "LOCAL_BLACKLIST", 4},
	}

	for _, l := range loaders {
		data, err := LoadWatchlistGeneric(db, l.table, l.source, l.aliasCount)
		if err != nil {
			return nil, err
		}
		all = append(all, data...)
	}

	return all, nil
}

// collectAliases membersihkan alias kosong/null
func collectAliases(aliases []sql.NullString) []string {
	var result []string
	for _, a := range aliases {
		if a.Valid && a.String != "" {
			result = append(result, a.String)
		}
	}
	return result
}
