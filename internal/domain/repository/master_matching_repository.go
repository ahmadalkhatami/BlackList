package repository

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MasterMatchingRepository interface {
	LoadMasterMatching(query string) ([]models.MasterMatching, error)
}

type SQLMasterMatchingRepository struct {
	DB *sql.DB
}

func (r SQLMasterMatchingRepository) LoadMasterMatching(query string) ([]models.MasterMatching, error) {
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
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
