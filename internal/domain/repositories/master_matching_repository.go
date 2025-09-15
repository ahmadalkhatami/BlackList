package repositories

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/utils"
	"context"
	"database/sql"
)

type MasterMatchingFilter struct {
	IsActive        *bool
	IsIndividual    *bool
	WatchlistSource *string
}

type MasterMatchingOption func(*MasterMatchingFilter)

func WithActive(isActive bool) MasterMatchingOption {
	return func(f *MasterMatchingFilter) {
		f.IsActive = &isActive
	}
}

func WithIndividual(isIndividual bool) MasterMatchingOption {
	return func(f *MasterMatchingFilter) {
		f.IsIndividual = &isIndividual
	}
}

func WithSource(source string) MasterMatchingOption {
	return func(f *MasterMatchingFilter) {
		f.WatchlistSource = &source
	}
}

type MasterMatchingRepository interface {
	Load(ctx context.Context, opts ...MasterMatchingOption) ([]models.MasterMatching, error)
}

type sqlMasterMatchingRepository struct {
	DB *sql.DB
}

func NewSQLMasterMatchingRepository(db *sql.DB) MasterMatchingRepository {
	return &sqlMasterMatchingRepository{DB: db}
}

func (r sqlMasterMatchingRepository) Load(ctx context.Context, opts ...MasterMatchingOption) ([]models.MasterMatching, error) {
	filter := MasterMatchingFilter{}
	for _, opt := range opts {
		opt(&filter)
	}

	fb := utils.NewQueryBuilder()

	if filter.IsIndividual != nil {
		fb.Add("IsIndividual", filter.IsIndividual)
	}

	if filter.IsActive != nil {
		fb.Add("IsActive", filter.IsActive)
	}

	if filter.WatchlistSource != nil {
		fb.Add("WatchlistSource", filter.WatchlistSource)
	}

	query := `
		SELECT [Id], [WatchlistSource], [IsIndividual], [IsActive]
		FROM [dbo].[MASTER_MATCHING]
	`
	query, args := fb.Build(query)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MasterMatching
	for rows.Next() {
		var rec models.MasterMatching
		if err := rows.Scan(
			&rec.Id,
			&rec.WatchlistSource,
			&rec.Type,
			&rec.IsActive,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

/*
contoh penggunaan:

ctx := context.Background()
repo := NewSQLMasterMatchingRepository(db)

// ambil semua
all, _ := repo.Load(ctx)

// hanya active
active, _ := repo.Load(ctx, WithActive(true))

// active + individu
individu, _ := repo.Load(ctx, WithActive(true), WithIndividual(true))

// active + corporate
corporate, _ := repo.Load(ctx, WithActive(true), WithIndividual(false))

// filter by source
dttot, _ := repo.Load(ctx, WithSource("MASTER_TERORIS"))

// filter by source + corporate
corpDttot, _ := repo.Load(ctx, WithActive(true), WithIndividual(false), WithSource("MASTER_TERORIS"))
*/
