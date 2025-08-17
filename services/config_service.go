package services

import (
	"BlackListWorker/models"
	"database/sql"
	"strconv"
	"strings"
)

const defaultThreshold = 0.85

// GetSystemConfigByKey mengambil konfigurasi dari table SYSTEM_CONFIG berdasarkan key
func GetSystemConfigByKey(db *sql.DB, key string) (models.SystemConfig, error) {
	var cfg models.SystemConfig

	query := `
		SELECT Id, ConfigKey, ConfigValue, ConfigType, Description
		FROM SYSTEM_CONFIG
		WHERE ConfigKey = @p1
	`

	// query := `
	// 	;WITH BASECONFIG AS (
	// 	SELECT MM.Id, MMC.MatchingId, WatchlistSource, FieldName, FieldWeight, MatchingAlgorithm, mmc.IsActive, CreatedBy, CreatedAt, UpdatedAt 
	// 	FROM MASTER_MATCHING mm WITH(NOLOCK)
	// 	INNER JOIN MASTER_MATCHING_CONFIG mmc WITH(NOLOCK) ON mm.Id = mmc.MatchingId AND mmc.IsActive = 1 AND mm.IsActive = 1
	// 	)
	// 	SELECT * FROM BASECONFIG
	// `
	

	err := db.QueryRow(query, key).Scan(
		&cfg.ID,
		&cfg.ConfigKey,
		&cfg.ConfigValue,
		&cfg.ConfigType,
		&cfg.Description,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.SystemConfig{}, nil
		}
		return models.SystemConfig{}, err
	}

	// Trim untuk menghindari spasi ekstra
	cfg.ConfigKey = strings.TrimSpace(cfg.ConfigKey)
	cfg.ConfigValue = strings.TrimSpace(cfg.ConfigValue)
	cfg.ConfigType = strings.TrimSpace(cfg.ConfigType)

	return cfg, nil
}

// GetThresholdFromConfig mengambil nilai threshold match dari config
func GetThresholdFromConfig(db *sql.DB, key string) float64 {
	cfg, err := GetSystemConfigByKey(db, key)
	if err != nil || cfg.ConfigKey == "" {
		return defaultThreshold
	}

	if cfg.ConfigType == "DECIMAL" || cfg.ConfigType == "INTEGER" {
		if val, err := strconv.ParseFloat(cfg.ConfigValue, 64); err == nil {
			return val
		}
	}

	return defaultThreshold
}
