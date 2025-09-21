package repositories

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/utils"
	"context"
	"database/sql"
)

type LocalBlacklistFilter struct {
	IsActive *bool
	Type     *string
}

type LocalBlacklistOption func(*LocalBlacklistFilter)

func WithLBActive(active bool) LocalBlacklistOption {
	return func(f *LocalBlacklistFilter) {
		f.IsActive = &active
	}
}

func WithLBType(t string) LocalBlacklistOption {
	return func(f *LocalBlacklistFilter) {
		f.Type = &t
	}
}

type MasterLocalBlacklistRepository interface {
	Load(ctx context.Context, opts ...LocalBlacklistOption) ([]models.MasterLocalBlacklist, error)
}

type sqlMasterLocalBlacklistRepository struct {
	DB *sql.DB
}

func NewSQLMasterLocalBlacklistRepository(db *sql.DB) MasterLocalBlacklistRepository {
	return &sqlMasterLocalBlacklistRepository{DB: db}
}

func (r sqlMasterLocalBlacklistRepository) Load(ctx context.Context, opts ...LocalBlacklistOption) ([]models.MasterLocalBlacklist, error) {
	filter := LocalBlacklistFilter{}
	for _, opt := range opts {
		opt(&filter)
	}

	fb := utils.NewQueryBuilder()

	fb.Add("IsActive", filter.IsActive)
	fb.Add("Type", filter.Type)

	query := `
		SELECT [Id],[Nama],[Alias1],[Alias2],[Alias3],[Alias4],
			[Type],[TempatLahir],[TanggalLahir],[KTP],[NPWP],
			[NoPaspor],[CreatedAt],[IsActive]
		FROM [dbo].[MASTER_LOCAL_BLACKLIST]
	`
	query, args := fb.Build(query)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MasterLocalBlacklist
	for rows.Next() {
		var rec models.MasterLocalBlacklist
		if err := rows.Scan(
			&rec.ID,
			&rec.Nama,
			&rec.Alias1,
			&rec.Alias2,
			&rec.Alias3,
			&rec.Alias4,
			&rec.Type,
			&rec.TempatLahir,
			&rec.TanggalLahir,
			&rec.KTP,
			&rec.NPWP,
			&rec.NoPaspor,
			&rec.CreatedAt,
			&rec.IsActive,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

/*
ctx := context.Background()
repo := NewSQLMasterLocalBlacklistRepository(db)

// tanpa filter
all, _ := repo.Load(ctx)

// filter aktif saja
active, _ := repo.Load(ctx, WithLBActive(true))

// filter aktif + individu
individu, _ := repo.Load(ctx, WithLBActive(true), WithLBType("individu"))

// filter aktif + corporate
corporate, _ := repo.Load(ctx, WithLBActive(true), WithLBType("korporasi"))
*/
