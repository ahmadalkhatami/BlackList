package repositories

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterTerorisRepository interface {
	Load() ([]models.MasterTeroris, error)
	LoadIndividu() ([]models.MasterTeroris, error)
	LoadCorporate() ([]models.MasterTeroris, error)
}

type sqlMasterTerorisRepository struct {
	DB *sql.DB
}

func NewSQLMasterTerorisRepository(db *sql.DB) MasterTerorisRepository {
	return &sqlMasterTerorisRepository{DB: db}
}

func (r sqlMasterTerorisRepository) Load() ([]models.MasterTeroris, error) {
	rows, err := r.DB.Query("SELECT [Id], [Nama], [Alias1], [Alias2], [Alias3], [Alias4], [Type], [TempatLahir], [TanggalLahir], [KTP], [NPWP], [NoPaspor], [CreatedAt], [IsActive] FROM [dbo].[MASTER_TERORIS] WHERE [IsActive] = 1;")
	if err != nil {
		return []models.MasterTeroris{}, err
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

func (r sqlMasterTerorisRepository) LoadIndividu() ([]models.MasterTeroris, error) {
	rows, err := r.DB.Query("SELECT [Id], [Nama], [Alias1], [Alias2], [Alias3], [Alias4], [Type], [TempatLahir], [TanggalLahir], [KTP], [NPWP], [NoPaspor], [CreatedAt], [IsActive] FROM [dbo].[MASTER_TERORIS] WHERE [IsActive] = 1 AND TOLOWER(Type) = 'individu';")
	if err != nil {
		return []models.MasterTeroris{}, err
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

func (r sqlMasterTerorisRepository) LoadCorporate() ([]models.MasterTeroris, error) {
	rows, err := r.DB.Query("SELECT [Id], [Nama], [Alias1], [Alias2], [Alias3], [Alias4], [Type], [TempatLahir], [TanggalLahir], [KTP], [NPWP], [NoPaspor], [CreatedAt], [IsActive] FROM [dbo].[MASTER_TERORIS] WHERE [IsActive] = 1 AND TOLOWER(Type) = 'korporasi';")
	if err != nil {
		return []models.MasterTeroris{}, err
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
