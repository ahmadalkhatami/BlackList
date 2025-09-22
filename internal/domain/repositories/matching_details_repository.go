package repositories

import (
	"BlackListWorker/internal/domain/models"
	"context"
	"database/sql"
	"fmt"

	mssql "github.com/denisenkom/go-mssqldb"
)

type MatchingDetailsRepository interface {
	Load(ctx context.Context) ([]models.MatchingDetail, error)
	Save(ctx context.Context, details []models.MatchingDetail) error
	SaveBatch(ctx context.Context, batchID int64, details []models.MatchingDetail) error
	GetLastId(ctx context.Context) (int64, error)
	ResetSequenceId(ctx context.Context) error
}

type sqlMatchingDetailsRepository struct {
	DB *sql.DB
}

func NewSQLMatchingDetailsRepository(db *sql.DB) MatchingDetailsRepository {
	return &sqlMatchingDetailsRepository{DB: db}
}

// Query sudah fixed di sini, nggak perlu parametris
func (r *sqlMatchingDetailsRepository) Load(ctx context.Context) ([]models.MatchingDetail, error) {

	query := `
		SELECT 
			ID,
			MatchingResultId,
			FieldName,
			CustomerValue,
			WatchlistValue,
			FieldScore,
			FieldWeight,
			AlgorithmUsed
		FROM dbo.MATCHING_DETAILS
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MatchingDetail
	for rows.Next() {
		var rec models.MatchingDetail
		if err := rows.Scan(
			&rec.Id,
			&rec.MatchingResultId,
			&rec.FieldName,
			&rec.CustomerValue,
			&rec.WatchlistValue,
			&rec.FieldScore,
			&rec.FieldWeight,
			&rec.AlgorithmUsed,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func (r *sqlMatchingDetailsRepository) Save(ctx context.Context, details []models.MatchingDetail) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	const q = `
		INSERT INTO dbo.MATCHING_DETAILS
		(MatchingResultId, FieldName, CustomerValue, WatchlistValue, FieldScore, FieldWeight, AlgorithmUsed)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)
	`

	stmt, err := tx.PrepareContext(ctx, q)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, detail := range details {
		_, err := stmt.ExecContext(
			ctx,
			detail.MatchingResultId,
			detail.FieldName,
			detail.CustomerValue,
			detail.WatchlistValue,
			detail.FieldScore,
			detail.FieldWeight,
			detail.AlgorithmUsed,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *sqlMatchingDetailsRepository) SaveBatch(ctx context.Context, batchID int64, details []models.MatchingDetail) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(mssql.CopyIn(
		"MATCHING_DETAILS",
		mssql.BulkOptions{},
		"MatchingResultId",
		"FieldName",
		"CustomerValue",
		"WatchlistValue",
		"FieldScore",
		"FieldWeight",
		"AlgorithmUsed",
	))
	if err != nil {
		return err
	}

	for _, detail := range details {
		_, err = stmt.Exec(
			detail.MatchingResultId,
			detail.FieldName,
			detail.CustomerValue,
			detail.WatchlistValue,
			detail.FieldScore,
			detail.FieldWeight,
			detail.AlgorithmUsed,
		)
		if err != nil {
			stmt.Close()
			tx.Rollback()
			return err
		}
	}

	if _, err := stmt.Exec(); err != nil {
		stmt.Close()
		tx.Rollback()
		return err
	}

	if err := stmt.Close(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *sqlMatchingDetailsRepository) GetLastId(ctx context.Context) (int64, error) {
	var lastID sql.NullInt64

	query := `SELECT MAX(Id) FROM MATCHING_DETAILS;`
	err := r.DB.QueryRowContext(ctx, query).Scan(&lastID)
	if err != nil {
		return 0, err
	}

	if !lastID.Valid {
		return 0, nil
	}

	return lastID.Int64, nil
}

func (r *sqlMatchingDetailsRepository) ResetSequenceId(ctx context.Context) error {

	lastID, err := r.GetLastId(ctx)
	if err != nil {
		return err
	}

	nextID := lastID + 1
	query := fmt.Sprintf(`ALTER SEQUENCE Seq_MatchingDetail RESTART WITH %d;`, nextID)

	_, err = r.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
