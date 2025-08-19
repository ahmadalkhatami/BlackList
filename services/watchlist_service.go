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
		WHERE IsActive = 1 AND Nama = 'AHMADElijah Corkery'
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

func LoadMasterWMD(db *sql.DB) ([]models.MasterWatchlist, error) {
	rows, err := db.Query(`
		SELECT Id, Nama, Alias1, Alias2, Alias3, Alias4, Alias5, Alias6, Alias7, Alias8, Alias9, Alias10,
		TempatLahir, TanggalLahir, Ktp, Npwp, NoPaspor, CreatedAt, UpdatedAt, IsActive 
		FROM MASTER_WMD
		WHERE IsActive = 1 AND Nama = 'AHMADElijah Corkery'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MasterWatchlist
	for rows.Next() {
		var m models.MasterWatchlist
		var aliases [10]sql.NullString

		err := rows.Scan(
			&m.ID, &m.Nama,
			&aliases[0], &aliases[1], &aliases[2], &aliases[3],
			&aliases[4], &aliases[5], &aliases[6], &aliases[7],
			&aliases[8], &aliases[9],
			&m.TempatLahir, &m.TanggalLahir, &m.KTP, &m.NPWP, &m.NoPaspor,
			&m.CreatedAt, &m.UpdatedAt, &m.IsActive,
		)
		if err != nil {
			return nil, err
		}

		m.Alias = WMDAliases(aliases[0], aliases[1], aliases[2], aliases[3],
		                      aliases[4], aliases[5], aliases[6], aliases[7],
		                      aliases[8], aliases[9])
		m.Source = "WMD" // ini penting untuk pemrosesan di matching

		result = append(result, m)
	}

	return result, nil
}

func LoadMasterLocalBlacklist(db *sql.DB) ([]models.MasterWatchlist, error) {
	rows, err := db.Query(`
		SELECT Id, Nama, Alias1, Alias2, Alias3, Alias4, TempatLahir, TanggalLahir, Ktp, Npwp, NoPaspor, CreatedAt, UpdatedAt, IsActive 
		FROM MASTER_LOCAL_BLACKLIST
		WHERE IsActive = 1 AND Nama = 'AHMADElijah Corkery'
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

		m.Alias = combineAliases(alias1, alias2, alias3, alias4)
		m.Source = "LOCAL_BLACKLIST" // ini penting untuk pemrosesan di matching

		result = append(result, m)
	}

	return result, nil
}

// LoadAllWatchlists menggabungkan data dari MASTER_TERORIS, MASTER_WMD, dan MASTER_LOCAL_BLACKLIST
func LoadAllWatchlists(db *sql.DB) ([]models.MasterWatchlist, error) {
	var all []models.MasterWatchlist

	// Ambil dari MASTER_TERORIS
	teroris, err := LoadMasterTeroris(db)
	if err != nil {
		return nil, err
	}
	all = append(all, teroris...)

	// Ambil dari MASTER_WMD
	wmd, err := LoadMasterWMD(db)
	if err != nil {
		return nil, err
	}
	all = append(all, wmd...)

	// Ambil dari MASTER_LOCAL_BLACKLIST
	local, err := LoadMasterLocalBlacklist(db)
	if err != nil {
		return nil, err
	}
	all = append(all, local...)

	return all, nil
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

func WMDAliases(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10 sql.NullString) []string {
	aliases := []string{}
	for _, a := range []sql.NullString{a1, a2, a3, a4, a5, a6, a7, a8, a9, a10} {
		if a.Valid && a.String != "" {
			aliases = append(aliases, a.String)
		}
	}
	return aliases
}
