package services

import (
	"BlackListWorker/models"
	"database/sql"
)

// LoadMasterTeroris mengambil data dari tabel MASTER_TERORIS
func LoadMasterTeroris(db *sql.DB) ([]models.MasterWatchlist, error) {
	rows, err := db.Query(`
		SELECT Id, Nama, Alias1, Alias2, Alias3, Alias4, TempatLahir, TanggalLahir, Ktp, Npwp, NoPaspor, CreatedAt, UpdatedAt, IsActive 
		FROM MASTER_TERORIS
		WHERE IsActive = 1
		AND Nama = 'Prince Eddie Reinger'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MasterWatchlist
	for rows.Next() {
		var m models.MasterWatchlist
		var alias1, alias2, alias3, alias4 sql.NullString

		err := rows.Scan(
			&m.ID, &m.Nama,
			&alias1, &alias2, &alias3, &alias4,
			&m.TempatLahir, &m.TanggalLahir, &m.KTP, &m.NPWP, &m.NoPaspor,
			&m.CreatedAt, &m.UpdatedAt, &m.IsActive,
		)
		if err != nil {
			return nil, err
		}

		// Gabungkan alias
		m.Alias = combineAliases(alias1, alias2, alias3, alias4)
		m.Source = "MASTER_TERORIS" // ini penting untuk pemrosesan di matching

		result = append(result, m)
	}

	return result, nil
}

// Fungsi bantu untuk gabungkan alias
func combineAliases(a1, a2, a3, a4 sql.NullString) []string {
	aliases := []string{}
	for _, a := range []sql.NullString{a1, a2, a3, a4} {
		if a.Valid && a.String != "" {
			aliases = append(aliases, a.String)
		}
	}
	return aliases
}
