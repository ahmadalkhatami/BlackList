package repositories

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/utils"
	"context"
	"database/sql"
)

type MasterTerorisFilter struct {
	IsActive *bool
	Type     *string
}

type MasterTerorisOption func(*MasterTerorisFilter)

func WithTerorisActive(active bool) MasterTerorisOption {
	return func(f *MasterTerorisFilter) {
		f.IsActive = &active
	}
}

func WithTerorisType(t string) MasterTerorisOption {
	return func(f *MasterTerorisFilter) {
		f.Type = &t
	}
}

type MasterTerorisRepository interface {
	Load(ctx context.Context, opts ...MasterTerorisOption) ([]models.MasterTeroris, error)
}

type sqlMasterTerorisRepository struct {
	DB *sql.DB
}

func NewSQLMasterTerorisRepository(db *sql.DB) MasterTerorisRepository {
	return &sqlMasterTerorisRepository{DB: db}
}

func (r sqlMasterTerorisRepository) Load(ctx context.Context, opts ...MasterTerorisOption) ([]models.MasterTeroris, error) {
	filter := MasterTerorisFilter{}
	for _, opt := range opts {
		opt(&filter)
	}

	fb := utils.NewQueryBuilder()

	fb.Add("IsActive", *filter.IsActive)

	if filter.Type != nil {
		fb.Where("LOWER([Type]) = LOWER(?)", *filter.Type)
	}

	query := `
		SELECT [Id],[Nama],[Alias1],[Alias2],[Alias3],[Alias4],
		       [Type],[TempatLahir],[TanggalLahir],[KTP],[NPWP],
		       [NoPaspor],[CreatedAt],[IsActive]
		FROM [dbo].[MASTER_TERORIS]
	`
	query, args := fb.Build(query)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MasterTeroris
	for rows.Next() {
		var rec models.MasterTeroris
		if err := rows.Scan(
			&rec.Id,
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
Contoh penggunaan:

ctx := context.Background()
repo := NewSQLMasterTerorisRepository(db)

// ambil semua
all, _ := repo.Load(ctx)

// hanya active
active, _ := repo.Load(ctx, WithTerorisActive(true))

// individu
individu, _ := repo.Load(ctx, WithTerorisActive(true), WithTerorisType("individu"))

// corporate
corporate, _ := repo.Load(ctx, WithTerorisActive(true), WithTerorisType("korporasi"))
*/
