package repositories

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterMatchingConfigRepository interface {
	LoadMasterMatchingConfig() ([]models.MasterMatchingConfig, error)
}

type sqlMasterMatchingConfigRepository struct {
	DB *sql.DB
}

func NewSQLMasterMatchingConfigRepository(db *sql.DB) MasterMatchingConfigRepository {
	return &sqlMasterMatchingConfigRepository{DB: db}
}

func (r sqlMasterMatchingConfigRepository) LoadMasterMatchingConfig() ([]models.MasterMatchingConfig, error) {
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
