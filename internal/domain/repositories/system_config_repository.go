package repositories

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type SystemConfigRepository interface {
	Load(cfgKey string) (models.SystemConfig, error)
}

type sqlSystemConfigRepository struct {
	DB *sql.DB
}

func NewSQLSystemConfigRepository(db *sql.DB) SystemConfigRepository {
	return &sqlSystemConfigRepository{DB: db}
}

func (r sqlSystemConfigRepository) Load(cfgKey string) (models.SystemConfig, error) {

	var rec models.SystemConfig

	query := "SELECT TOP 1 [Id], [ConfigKey], [ConfigValue], [ConfigType], [Description], [CreatedBy], [CreatedAt], [UpdatedAt] FROM [dbo].[SYSTEM_CONFIG] WHERE [ConfigKey] = @p1;"

	err := r.DB.QueryRow(query, sql.Named("p1", cfgKey)).Scan(
		&rec.ID,
		&rec.ConfigKey,
		&rec.ConfigValue,
		&rec.ConfigType,
		&rec.Description,
		&rec.CreatedBy,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)

	if err != nil {
		return models.SystemConfig{}, err
	}

	return rec, nil
}
