package repositories

import (
	"BlackListWorker/internal/domain/models"
	"context"
	"database/sql"
	"fmt"

	mssql "github.com/denisenkom/go-mssqldb"
)

type MatchingResultRepository interface {
	Load(ctx context.Context, query *string) ([]models.MatchingResult, error)
	Save(ctx context.Context, results []models.MatchingResult) error
	SaveBatch(ctx context.Context, batchID string, results []models.MatchingResult) error
	GetLastId(ctx context.Context) (int64, error)
	ResetSequenceId(ctx context.Context) error
	ResetSequenceBatchId(ctx context.Context) error
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
			&rec.Id,
			&rec.BatchId,
			&rec.CifNumber,
			&rec.CustomerName,
			&rec.WatchlistId,
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
			res.BatchId,
			res.CifNumber,
			res.CustomerName,
			res.WatchlistId,
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
		if res.CustomerName != nil {
			customerName = sql.NullString{String: *res.CustomerName, Valid: true}
		}

		if _, err := stmt.Exec(
			batchID,
			res.CifNumber,
			customerName,
			res.WatchlistId,
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

func (r *sqlMatchingResultRepository) GetLastId(ctx context.Context) (int64, error) {
	var lastID sql.NullInt64

	query := `SELECT MAX(Id) FROM MATCHING_RESULTS;`
	err := r.DB.QueryRowContext(ctx, query).Scan(&lastID)
	if err != nil {
		return 0, err
	}

	if !lastID.Valid {
		return 0, nil
	}

	return lastID.Int64, nil
}

func (r *sqlMatchingResultRepository) GetLastBatchId(ctx context.Context) (int64, error) {
	var lastID sql.NullInt64

	query := `SELECT MAX(BatchId) FROM MATCHING_RESULTS;`
	err := r.DB.QueryRowContext(ctx, query).Scan(&lastID)
	if err != nil {
		return 0, err
	}

	if !lastID.Valid {
		return 0, nil
	}

	return lastID.Int64, nil
}

func (r *sqlMatchingResultRepository) ResetSequenceId(ctx context.Context) error {

	lastID, err := r.GetLastId(ctx)
	if err != nil {
		return err
	}

	nextID := lastID + 1
	query := fmt.Sprintf(`ALTER SEQUENCE Seq_MatchingResult RESTART WITH %d;`, nextID)

	_, err = r.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (r *sqlMatchingResultRepository) ResetSequenceBatchId(ctx context.Context) error {

	lastID, err := r.GetLastId(ctx)
	if err != nil {
		return err
	}

	nextID := lastID + 1
	query := fmt.Sprintf(`ALTER SEQUENCE Seq_MatchingBatch RESTART WITH %d;`, nextID)

	_, err = r.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
