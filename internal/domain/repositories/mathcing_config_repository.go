package repositories

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/utils"
	"context"
	"database/sql"
)

type MasterMatchingConfigFilter struct {
	IsActive   *bool
	MatchingID *int64
}

type MasterMatchingConfigOption func(*MasterMatchingConfigFilter)

func WithActiveConfig(isActive bool) MasterMatchingConfigOption {
	return func(f *MasterMatchingConfigFilter) {
		f.IsActive = &isActive
	}
}

func WithMatchingID(id int64) MasterMatchingConfigOption {
	return func(f *MasterMatchingConfigFilter) {
		f.MatchingID = &id
	}
}

type MasterMatchingConfigRepository interface {
	Load(ctx context.Context, opts ...MasterMatchingConfigOption) ([]models.MasterMatchingConfig, error)
}

type sqlMasterMatchingConfigRepository struct {
	DB *sql.DB
}

func NewSQLMasterMatchingConfigRepository(db *sql.DB) MasterMatchingConfigRepository {
	return &sqlMasterMatchingConfigRepository{DB: db}
}

func (r *sqlMasterMatchingConfigRepository) Load(ctx context.Context, opts ...MasterMatchingConfigOption) ([]models.MasterMatchingConfig, error) {
	filter := &MasterMatchingConfigFilter{}
	for _, opt := range opts {
		opt(filter)
	}

	fb := utils.NewQueryBuilder()

	if filter.IsActive != nil {
		fb.Add("IsActive", filter.IsActive)
	}

	if filter.MatchingID != nil {
		fb.Add("MatchingId", *filter.MatchingID)
	}

	query := `
		SELECT [Id], [MatchingId], [FieldName], [FieldWeight], [MatchingAlgorithm], [IsActive]
		FROM [dbo].[MASTER_MATCHING_CONFIG]
	`
	query, args := fb.Build(query)

	rows, err := r.DB.QueryContext(ctx, query, args...)
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
			&rec.IsActive,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func buildArgs(f *MasterMatchingConfigFilter) []interface{} {
	var args []interface{}
	if f.MatchingID != nil {
		args = append(args, *f.MatchingID)
	}
	return args
}

/*
ctx := context.Background()

// Ambil semua yg aktif
configs, _ := repo.Load(ctx, WithActive(true))

// Ambil semua untuk MatchingID tertentu
configs, _ = repo.Load(ctx, WithMatchingID(123))

// Ambil yg aktif dan untuk MatchingID tertentu
configs, _ = repo.Load(ctx, WithActive(true), WithMatchingID(123))
*/
