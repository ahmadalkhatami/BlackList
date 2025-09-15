package repositories

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/utils"
	"context"
	"database/sql"
)

// Filter struct
type SystemConfigFilter struct {
	Key  *string
	Type *string
}

type SystemConfigOption func(*SystemConfigFilter)

func WithConfigKey(key string) SystemConfigOption {
	return func(f *SystemConfigFilter) {
		f.Key = &key
	}
}

func WithConfigType(t string) SystemConfigOption {
	return func(f *SystemConfigFilter) {
		f.Type = &t
	}
}

type SystemConfigRepository interface {
	Load(ctx context.Context, opts ...SystemConfigOption) ([]models.SystemConfig, error)
	LoadOne(ctx context.Context, opts ...SystemConfigOption) (models.SystemConfig, error)
}

type sqlSystemConfigRepository struct {
	DB *sql.DB
}

func NewSQLSystemConfigRepository(db *sql.DB) SystemConfigRepository {
	return &sqlSystemConfigRepository{DB: db}
}

func (r *sqlSystemConfigRepository) Load(ctx context.Context, opts ...SystemConfigOption) ([]models.SystemConfig, error) {
	filter := &SystemConfigFilter{}
	for _, opt := range opts {
		opt(filter)
	}

	fb := utils.NewQueryBuilder()

	if filter.Key != nil {
		fb.Add("ConfigKey", *filter.Key)
	}

	if filter.Type != nil {
		fb.Add("ConfigType", *filter.Type)
	}

	query := `
		SELECT [Id], [ConfigKey], [ConfigValue], [ConfigType], [Description], [CreatedBy], [CreatedAt], [UpdatedAt]
		FROM [dbo].[SYSTEM_CONFIG]
	`
	query, args := fb.Build(query)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.SystemConfig
	for rows.Next() {
		var rec models.SystemConfig
		if err := rows.Scan(
			&rec.ID,
			&rec.ConfigKey,
			&rec.ConfigValue,
			&rec.ConfigType,
			&rec.Description,
			&rec.CreatedBy,
			&rec.CreatedAt,
			&rec.UpdatedAt,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func (r *sqlSystemConfigRepository) LoadOne(ctx context.Context, opts ...SystemConfigOption) (models.SystemConfig, error) {
	recs, err := r.Load(ctx, opts...)
	if err != nil {
		return models.SystemConfig{}, err
	}
	if len(recs) == 0 {
		return models.SystemConfig{}, sql.ErrNoRows
	}
	return recs[0], nil
}

/*
ctx := context.Background()

// Ambil semua config
configs, _ := repo.Load(ctx)

// Ambil semua config dengan type tertentu
configs, _ = repo.Load(ctx, WithConfigType("MATCHING_ALGORITHM"))

// Ambil satu config by key
cfg, _ := repo.LoadOne(ctx, WithConfigKey("SELECT"))

*/
