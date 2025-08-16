package repository

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterNasabahRepository interface {
	LoadMasterNasabah() ([]models.MasterNasabah, error)
}

type SQLMasterNasabahRepository struct {
	DB *sql.DB
}

func NewSQLMasterNasabahRepository(db *sql.DB) MasterNasabahRepository {
	return &SQLMasterNasabahRepository{DB: db}
}

func (r SQLMasterNasabahRepository) LoadMasterNasabah() ([]models.MasterNasabah, error) {
	rows, err := r.DB.Query("SELECT [Id], [CIFNumber], [NamaNasabah], [TempatLahir], [TanggalLahir], [KTP], [NPWP], [NoPaspor], [StatusNasabah], [CreatedAt], [UpdatedAt] FROM [dbo].[MASTER_NASABAH] WHERE StatusNasabah = 'AKTIF';")
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
