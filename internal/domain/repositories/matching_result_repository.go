package repositories

import (
	"BlackListWorker/internal/domain/models"
	"context"
	"database/sql"

	mssql "github.com/denisenkom/go-mssqldb"
)

type MatchingResultRepository interface {
	Load(ctx context.Context, query *string) ([]models.MatchingResult, error)
	Save(ctx context.Context, results []models.MatchingResult) error
	SaveBatch(ctx context.Context, batchID string, results []models.MatchingResult) error
}

type sqlMatchingResultRepository struct {
	DB *sql.DB
}

func NewSQLMatchingResultRepository(db *sql.DB) MatchingResultRepository {
	return &sqlMatchingResultRepository{DB: db}
}

func (r *sqlMatchingResultRepository) Load(ctx context.Context, query *string) ([]models.MatchingResult, error) {
	defaultQuery := `
		SELECT [Id], [BatchId], [CIFNumber], [CustomerName],
		       [WatchlistId], [WatchlistSource], [SimilarityScore],
		       [Status], [ProcessDate], [ProcessTime], [CreatedAt]
		FROM [dbo].[MATCHING_RESULTS]`

	if query == nil || *query == "" {
		query = &defaultQuery
	}

	rows, err := r.DB.QueryContext(ctx, *query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MatchingResult
	for rows.Next() {
		var rec models.MatchingResult
		if err := rows.Scan(
			&rec.ID,
			&rec.BatchID,
			&rec.CIFNumber,
			&rec.CustomerName,
			&rec.WatchlistID,
			&rec.WatchlistSource,
			&rec.SimilarityScore,
			&rec.Status,
			&rec.ProcessDate,
			&rec.ProcessTime,
			&rec.CreatedAt,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func (r *sqlMatchingResultRepository) Save(ctx context.Context, results []models.MatchingResult) error {
	if len(results) == 0 {
		return nil
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const q = `
		INSERT INTO dbo.MATCHING_RESULTS
			(BatchId, CIFNumber, CustomerName, WatchlistId, WatchlistSource, 
			 SimilarityScore, Status, ProcessDate, ProcessTime)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9)`

	stmt, err := tx.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, res := range results {
		if _, err := stmt.ExecContext(
			ctx,
			res.BatchID,
			res.CIFNumber,
			res.CustomerName,
			res.WatchlistID,
			res.WatchlistSource,
			res.SimilarityScore,
			res.Status,
			res.ProcessDate,
			res.ProcessTime,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *sqlMatchingResultRepository) SaveBatch(ctx context.Context, batchID string, results []models.MatchingResult) error {
	if len(results) == 0 {
		return nil
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(mssql.CopyIn(
		"MATCHING_RESULTS",
		mssql.BulkOptions{KeepNulls: true},
		"BatchId",
		"CIFNumber",
		"CustomerName",
		"WatchlistId",
		"WatchlistSource",
		"SimilarityScore",
		"Status",
		"ProcessDate",
		"ProcessTime",
	))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, res := range results {
		var customerName sql.NullString
		if res.CustomerName != "" {
			customerName = sql.NullString{String: res.CustomerName, Valid: true}
		}

		if _, err := stmt.Exec(
			batchID,
			res.CIFNumber,
			customerName,
			res.WatchlistID,
			res.WatchlistSource,
			res.SimilarityScore,
			res.Status,
			res.ProcessDate,
			res.ProcessTime,
		); err != nil {
			return err
		}
	}

	if _, err := stmt.Exec(); err != nil {
		return err
	}

	return tx.Commit()
}
