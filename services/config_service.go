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
		SELECT id, config_key, config_value, config_type, description
		FROM SYSTEM_CONFIG
		WHERE config_key = @p1
	`

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
