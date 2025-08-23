package repositories

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterMatchingRepository interface {
	LoadMasterMatching() ([]models.MasterMatching, error)
}

type sqlMasterMatchingRepository struct {
	DB *sql.DB
}

func NewSQLMasterMatchingRepository(db *sql.DB) MasterMatchingRepository {
	return &sqlMasterMatchingRepository{DB: db}
}

func (r sqlMasterMatchingRepository) LoadMasterMatching() ([]models.MasterMatching, error) {
	// rows, err := r.DB.Query("SELECT [Id] ,[Name] ,[WatchlistSource] ,[IsIndividual] ,[Description] ,[IsActive] FROM [dbo].[MASTER_MATCHING] WHERE [IsActive] = 1;")
	rows, err := r.DB.Query("SELECT [Id], [WatchlistSource], [IsIndividual], [IsActive] FROM [dbo].[MASTER_MATCHING] WHERE [IsActive] = 1;")
	if err != nil {
		return []models.MasterMatching{}, err
	}

	defer rows.Close()

	var records []models.MasterMatching
	for rows.Next() {
		var rec models.MasterMatching
		if err := rows.Scan(&rec.Id, &rec.WatchlistSource, &rec.Type, &rec.IsActive); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}
