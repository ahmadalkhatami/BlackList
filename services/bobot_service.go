package services

import (
	"BlackListWorker/models"
	"database/sql"
)

// LoadMatchingConfig mengambil data konfigurasi pencocokan dari tabel MATCHING_CONFIG
func LoadMatchingConfig(db *sql.DB) ([]models.MatchingConfig, error) {
	rows, err := db.Query(`
		SELECT id, watchlist_source, field_name, field_weight, matching_algorithm, is_active, created_by, created_at, updated_at
		FROM MATCHING_CONFIG
		WHERE is_active = 1
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []models.MatchingConfig
	for rows.Next() {
		var cfg models.MatchingConfig
		err := rows.Scan(
			&cfg.ID,
			&cfg.WatchlistSource,
			&cfg.FieldName,
			&cfg.FieldWeight,
			&cfg.MatchingAlgorithm,
			&cfg.IsActive,
			&cfg.CreatedBy,
			&cfg.CreatedAt,
			&cfg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, cfg)
	}

	return configs, nil
}
