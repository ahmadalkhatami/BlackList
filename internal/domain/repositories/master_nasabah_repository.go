package repositories

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/utils"
	"context"
	"database/sql"
)

type MasterNasabahRepository interface {
	Load(ctx context.Context, opts ...NasabahOption) ([]models.MasterNasabah, error)
	LoadAll(ctx context.Context) ([]models.MasterNasabah, error)
}

type sqlMasterNasabahRepository struct {
	DB *sql.DB
}

func NewSQLMasterNasabahRepository(db *sql.DB) MasterNasabahRepository {
	return &sqlMasterNasabahRepository{DB: db}
}

type nasabahFilter struct {
	Status *string
}

type NasabahOption func(*nasabahFilter)

func WithStatus(status string) NasabahOption {
	return func(f *nasabahFilter) {
		f.Status = &status
	}
}

func (r sqlMasterNasabahRepository) Load(ctx context.Context, opts ...NasabahOption) ([]models.MasterNasabah, error) {

	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	filter := &nasabahFilter{}
	for _, opt := range opts {
		opt(filter)
	}

	fb := utils.NewQueryBuilder()

	fb.Add("StatusNasabah", filter.Status)

	query := `
		SELECT [Id], [CIFNumber], [NamaNasabah], [TempatLahir], [TanggalLahir],
		       [KTP], [NPWP], [NoPaspor], [StatusNasabah], [CreatedAt]
		FROM [dbo].[MASTER_NASABAH]
	`
	query, args := fb.Build(query)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MasterNasabah
	for rows.Next() {
		var rec models.MasterNasabah
		if err := rows.Scan(
			&rec.ID,
			&rec.CifNumber,
			&rec.NamaNasabah,
			&rec.TempatLahir,
			&rec.TanggalLahir,
			&rec.KTP,
			&rec.NPWP,
			&rec.NoPaspor,
			&rec.StatusNasabah,
			&rec.CreatedAt,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func (r sqlMasterNasabahRepository) LoadAll(ctx context.Context) ([]models.MasterNasabah, error) {
	return r.Load(ctx, WithStatus("ACTIVE"))
}

/*
ctx := context.Background()
repo := repositories.NewSQLMasterNasabahRepository(db)
nasabah1, _ := repo.LoadAll(ctx)
nasabah2, _ := repo.Load(ctx)
nasabah3, _ := repo.Load(ctx, repositories.WithStatus("ACTIVE"), repositories.WithLimit(50))
*/
