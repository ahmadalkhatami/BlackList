package repositories

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/utils"
	"context"
	"database/sql"
	"fmt"
)

type WMDFilter struct {
	IsActive *bool
	Type     *string
}
type WMDOption func(*WMDFilter)

func WithWMDType(t string) WMDOption {
	return func(f *WMDFilter) {
		f.Type = &t
	}
}

func WithWMDActive(active bool) WMDOption {
	return func(f *WMDFilter) {
		f.IsActive = &active
	}
}

type MasterWMDRepository interface {
	Load(ctx context.Context, opts ...WMDOption) ([]models.MasterWMD, error)
}

type sqlMasterWMDRepository struct {
	DB *sql.DB
}

func NewSQLMasterWMDRepository(db *sql.DB) MasterWMDRepository {
	return &sqlMasterWMDRepository{DB: db}
}

func (r sqlMasterWMDRepository) Load(ctx context.Context, opts ...WMDOption) ([]models.MasterWMD, error) {
	filter := WMDFilter{}
	for _, opt := range opts {
		opt(&filter)
	}

	fb := utils.NewQueryBuilder()

	if filter.IsActive != nil {
		fb.Add("IsActive", *filter.IsActive)
	}

	if filter.Type != nil {
		fb.Where("LOWER([Type]) = LOWER(?)", *filter.Type)
	}

	query := `
		SELECT [Id], [Nama], [Alias1], [Alias2], [Alias3], [Alias4], [Alias5], [Alias6], [Alias7], [Alias8], [Alias9], [Alias10],
		       [Type], [TempatLahir], [TanggalLahir], [KTP], [NPWP], [NoPaspor], [CreatedAt], [IsActive]
		FROM [dbo].[MASTER_WMD]
	`
	query, args := fb.Build(query)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MasterWMD
	for rows.Next() {
		var rec models.MasterWMD
		if err := rows.Scan(
			&rec.ID,
			&rec.Nama,
			&rec.Alias1,
			&rec.Alias2,
			&rec.Alias3,
			&rec.Alias4,
			&rec.Alias5,
			&rec.Alias6,
			&rec.Alias7,
			&rec.Alias8,
			&rec.Alias9,
			&rec.Alias10,
			&rec.Type,
			&rec.TempatLahir,
			&rec.TanggalLahir,
			&rec.KTP,
			&rec.NPWP,
			&rec.NoPaspor,
			&rec.CreatedAt,
			&rec.IsActive,
		); err != nil {
			return nil, fmt.Errorf("scan WMD: %w", err)
		}
		records = append(records, rec)
	}

	fmt.Printf("🔍 Load WMD: %d record(s) loaded\n", len(records))
	return records, nil
}

/*
Contoh penggunaan:

ctx := context.Background()
repo := NewSQLMasterWMDRepository(db)

// ambil semua aktif
all, _ := repo.Load(ctx, WithWMDActive(true))

// hanya individu aktif
individu, _ := repo.Load(ctx, WithWMDActive(true), WithWMDType("individu"))

// hanya korporasi aktif
corporate, _ := repo.Load(ctx, WithWMDActive(true), WithWMDType("korporasi"))
*/
