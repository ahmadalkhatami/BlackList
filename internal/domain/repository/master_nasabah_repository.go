package repository

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterNasabahRepository interface {
	LoadMasterNasabah(query string) ([]models.MasterNasabah, error)
}

type SQLMasterNasabahRepository struct {
	DB *sql.DB
}

func (r SQLMasterNasabahRepository) LoadMasterNasabah(query string) ([]models.MasterNasabah, error) {
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var records []models.MasterNasabah
	for rows.Next() {
		var rec models.MasterNasabah
		if err := rows.Scan(
			&rec.Id,
			&rec.CIFNumber,
			&rec.NamaNasabah,
			&rec.TempatLahir,
			&rec.TanggalLahir,
			&rec.KTP,
			&rec.NPWP,
			&rec.NoPaspor,
			&rec.StatusNasabah,
			&rec.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}
