package repository

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterMatchingConfigRepository interface {
	LoadMasterMatchingConfig(query string) ([]models.MasterMatchingConfig, error)
}

type SQLMasterMatchingConfigRepository struct {
	DB *sql.DB
}

func (r SQLMasterMatchingConfigRepository) LoadMasterMatchingConfig(query string) ([]models.MasterMatchingConfig, error) {
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
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
