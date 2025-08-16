package repository

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterWMDRepository interface {
	LoadMasterWMD() ([]models.MasterWMD, error)
}

type SQLMasterWMDRepository struct {
	DB *sql.DB
}

func NewSQLMasterWMDRepository(db *sql.DB) MasterWMDRepository {
	return &SQLMasterWMDRepository{DB: db}
}

func (r SQLMasterWMDRepository) LoadMasterWMD() ([]models.MasterWMD, error) {
	rows, err := r.DB.Query("SELECT [Id], [Nama], [Alias1], [Alias2], [Alias3], [Alias4], [Alias5], [Alias6], [Alias7], [Alias8], [Alias9], [Alias10], [Type], [TempatLahir], [TanggalLahir], [KTP], [NPWP], [NoPaspor], [CreatedAt], [UpdatedAt], [IsActive] FROM [dbo].[MASTER_WMD] WHERE [IsActive] = 1;")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var records []models.MasterWMD
	for rows.Next() {
		var rec models.MasterWMD
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
