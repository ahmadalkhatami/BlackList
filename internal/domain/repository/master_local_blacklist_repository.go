package repository

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterLocalBlacklistRepository interface {
	LoadMasterLocalBlacklist(query string) ([]models.MasterLocalBlacklist, error)
}

type SQLMasterLocalBlacklistRepository struct {
	DB *sql.DB
}

func (r SQLMasterLocalBlacklistRepository) LoadMasterLocalBlacklist(query string) ([]models.MasterLocalBlacklist, error) {
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var records []models.MasterLocalBlacklist
	for rows.Next() {
		var rec models.MasterLocalBlacklist
		if err := rows.Scan(
			&rec.Id,
			&rec.Nama,
			&rec.Alias1,
			&rec.Alias2,
			&rec.Alias3,
			&rec.Alias4,
			&rec.Type,
			&rec.TempatLahir,
			&rec.TanggalLahir,
			&rec.KTP,
			&rec.NPWP,
			&rec.NoPaspor,
			&rec.CreatedAt,
			&rec.IsActive); err != nil {
			return nil, err
		}

		records = append(records, rec)
	}

	return records, nil
}
