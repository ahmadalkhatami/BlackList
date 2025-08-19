package services

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

// LoadMatchingConfig mengambil data konfigurasi pencocokan dari tabel MATCHING_CONFIG
func LoadMatchingConfig(db *sql.DB) ([]models.MatchingConfig, error) {
	rows, err := db.Query(`
		;WITH BASECONFIG AS (
			SELECT MM.Id, WatchlistSource, FieldName, FieldWeight, MatchingAlgorithm, mmc.IsActive, CreatedBy, CreatedAt 
			FROM MASTER_MATCHING mm WITH(NOLOCK)
			INNER JOIN MASTER_MATCHING_CONFIG mmc WITH(NOLOCK) ON mm.Id = mmc.MatchingId AND mmc.IsActive = 1 AND mm.IsActive = 1
		)
		SELECT * FROM BASECONFIG
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
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, cfg)
	}

	return configs, nil
}
