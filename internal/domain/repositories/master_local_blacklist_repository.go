package repositories

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterLocalBlacklistRepository interface {
	Load() ([]models.MasterLocalBlacklist, error)
	LoadIndividu() ([]models.MasterLocalBlacklist, error)
	LoadCorporate() ([]models.MasterLocalBlacklist, error)
}

type sqlMasterLocalBlacklistRepository struct {
	DB *sql.DB
}

func NewSQLMasterLocalBlacklistRepository(db *sql.DB) MasterLocalBlacklistRepository {
	return &sqlMasterLocalBlacklistRepository{DB: db}
}

func (r sqlMasterLocalBlacklistRepository) Load() ([]models.MasterLocalBlacklist, error) {
	rows, err := r.DB.Query("SELECT [Id] ,[Nama] ,[Alias1] ,[Alias2] ,[Alias3] ,[Alias4] ,[Type] ,[TempatLahir] ,[TanggalLahir] ,[KTP] ,[NPWP] ,[NoPaspor] ,[CreatedAt] ,[IsActive] FROM [dbo].[MASTER_LOCAL_BLACKLIST] WHERE [IsActive] = 1;")
	if err != nil {
		return []models.MasterLocalBlacklist{}, err
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

func (r sqlMasterLocalBlacklistRepository) LoadIndividu() ([]models.MasterLocalBlacklist, error) {
	rows, err := r.DB.Query("SELECT [Id] ,[Nama] ,[Alias1] ,[Alias2] ,[Alias3] ,[Alias4] ,[Type] ,[TempatLahir] ,[TanggalLahir] ,[KTP] ,[NPWP] ,[NoPaspor] ,[CreatedAt] ,[IsActive] FROM [dbo].[MASTER_LOCAL_BLACKLIST] WHERE [IsActive] = 1 AND TOLOWER(Type) = 'individu';")
	if err != nil {
		return []models.MasterLocalBlacklist{}, err
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

func (r sqlMasterLocalBlacklistRepository) LoadCorporate() ([]models.MasterLocalBlacklist, error) {
	rows, err := r.DB.Query("SELECT [Id] ,[Nama] ,[Alias1] ,[Alias2] ,[Alias3] ,[Alias4] ,[Type] ,[TempatLahir] ,[TanggalLahir] ,[KTP] ,[NPWP] ,[NoPaspor] ,[CreatedAt] ,[IsActive] FROM [dbo].[MASTER_LOCAL_BLACKLIST] WHERE [IsActive] = 1 AND TOLOWER(Type) = 'korporasi';")
	if err != nil {
		return []models.MasterLocalBlacklist{}, err
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
