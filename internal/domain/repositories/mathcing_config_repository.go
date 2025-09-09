package repositories

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterMatchingConfigRepository interface {
	Load() ([]models.MasterMatchingConfig, error)
	LoadIndividu() ([]models.MasterMatchingConfig, error)
	LoadCorporate() ([]models.MasterMatchingConfig, error)
}

type sqlMasterMatchingConfigRepository struct {
	DB *sql.DB
}

func NewSQLMasterMatchingConfigRepository(db *sql.DB) MasterMatchingConfigRepository {
	return &sqlMasterMatchingConfigRepository{DB: db}
}

func (r sqlMasterMatchingConfigRepository) Load() ([]models.MasterMatchingConfig, error) {
	rows, err := r.DB.Query("SELECT [Id], [MatchingId], [FieldName], [FieldWeight], [MatchingAlgorithm], [IsActive] FROM [dbo].[MASTER_MATCHING_CONFIG] WHERE [IsActive] = 1;")
	if err != nil {
		return []models.MasterMatchingConfig{}, err
	}

	defer rows.Close()

	var records []models.MasterMatchingConfig
	for rows.Next() {
		var rec models.MasterMatchingConfig
		if err := rows.Scan(
			&rec.Id,
			&rec.MatchingId,
			&rec.FieldName,
			&rec.FieldWeight,
			&rec.MatchingAlgorithm,
			&rec.IsActive); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func (r sqlMasterMatchingConfigRepository) LoadIndividu() ([]models.MasterMatchingConfig, error) {
	return nil, nil
}
func (r sqlMasterMatchingConfigRepository) LoadCorporate() ([]models.MasterMatchingConfig, error) {
	return nil, nil
}
