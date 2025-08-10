package repository

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterTerorisRepository interface {
	LoadMasterTeroris(query string) ([]models.MasterTeroris, error)
}

type SQLMasterTerorisRepository struct {
	DB *sql.DB
}

func (r SQLMasterTerorisRepository) LoadMasterTeroris(query string) ([]models.MasterTeroris, error) {
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var records []models.MasterTeroris
	for rows.Next() {
		var rec models.MasterTeroris
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
