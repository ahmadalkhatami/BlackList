package repositories

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterWMDRepository interface {
	Load() ([]models.MasterWMD, error)
	LoadIndividu() ([]models.MasterWMD, error)
	LoadCorporate() ([]models.MasterWMD, error)
}

type sqlMasterWMDRepository struct {
	DB *sql.DB
}

func NewSQLMasterWMDRepository(db *sql.DB) MasterWMDRepository {
	return &sqlMasterWMDRepository{DB: db}
}

func (r sqlMasterWMDRepository) Load() ([]models.MasterWMD, error) {
	rows, err := r.DB.Query("SELECT [Id], [Nama], [Alias1], [Alias2], [Alias3], [Alias4], [Alias5], [Alias6], [Alias7], [Alias8], [Alias9], [Alias10], [Type], [TempatLahir], [TanggalLahir], [KTP], [NPWP], [NoPaspor], [CreatedAt], [IsActive] FROM [dbo].[MASTER_WMD] WHERE [IsActive] = 1;")
	if err != nil {
		return []models.MasterWMD{}, err
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
			&rec.Alias5,
			&rec.Alias6,
			&rec.Alias7,
			&rec.Alias8,
			&rec.Alias9,
			&rec.Alias10,
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

func (r sqlMasterWMDRepository) LoadIndividu() ([]models.MasterWMD, error) {
	rows, err := r.DB.Query("SELECT [Id], [Nama], [Alias1], [Alias2], [Alias3], [Alias4], [Alias5], [Alias6], [Alias7], [Alias8], [Alias9], [Alias10], [Type], [TempatLahir], [TanggalLahir], [KTP], [NPWP], [NoPaspor], [CreatedAt], [IsActive] FROM [dbo].[MASTER_WMD] WHERE [IsActive] = 1 AND TOLOWER(Type) = 'individu';")
	if err != nil {
		return []models.MasterWMD{}, err
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
			&rec.Alias5,
			&rec.Alias6,
			&rec.Alias7,
			&rec.Alias8,
			&rec.Alias9,
			&rec.Alias10,
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

func (r sqlMasterWMDRepository) LoadCorporate() ([]models.MasterWMD, error) {
	rows, err := r.DB.Query("SELECT [Id], [Nama], [Alias1], [Alias2], [Alias3], [Alias4], [Alias5], [Alias6], [Alias7], [Alias8], [Alias9], [Alias10], [Type], [TempatLahir], [TanggalLahir], [KTP], [NPWP], [NoPaspor], [CreatedAt], [IsActive] FROM [dbo].[MASTER_WMD] WHERE [IsActive] = 1 AND TOLOWER(Type) = 'korporasi';")
	if err != nil {
		return []models.MasterWMD{}, err
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
			&rec.Alias5,
			&rec.Alias6,
			&rec.Alias7,
			&rec.Alias8,
			&rec.Alias9,
			&rec.Alias10,
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
