package repository

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"

	mssql "github.com/denisenkom/go-mssqldb"
)

type MatchingResultRepository interface {
	LoadMatchingResult(query string) ([]models.MatchingResult, error)
	SaveMatchingResult(results []models.MatchingResult) error
	SaveMatchingResultBatch(batchID string, results []models.MatchingResult) error
}

type sqlMatchingResultRepository struct {
	DB *sql.DB
}

func NewSQLMatchingResultRepository(db *sql.DB) MatchingResultRepository {
	return &sqlMatchingResultRepository{DB: db}
}

func (r sqlMatchingResultRepository) LoadMatchingResult(query string) ([]models.MatchingResult, error) {
	rows, err := r.DB.Query(query)
	if err != nil {
		return []models.MatchingResult{}, err
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
			&rec.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func (r *sqlMatchingResultRepository) SaveMatchingResult(results []models.MatchingResult) error {
	trx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer trx.Rollback()

	stmt, err := trx.Prepare(
		`INSERT INTO MATCHING_RESULTS (
			BatchId,
			CIFNumber,
			CustomerName,
			WatchlistId,
			WatchlistSource,
			SimilarityScore,
			Status,
			ProcessDate,
			ProcessTime
		) VALUES (
			@BatchId,
			@CIFNumber,
			@CustomerName,
			@WatchlistID,
			@WatchlistSource,
			@SimilarityScore,
			@Status,
			@ProcessDate,
			@ProcessTime
	)`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, res := range results {
		if _, err := stmt.Exec(
			res.BatchID,
			res.CIFNumber,
			res.CustomerName,
			res.WatchlistID,
			res.WatchlistSource,
			res.SimilarityScore,
			res.Status,
			res.ProcessDate,
			res.ProcessTime); err != nil {
			return err
		}
	}

	return trx.Commit()
}

func (r *sqlMatchingResultRepository) SaveMatchingResultBatch(batchID string, results []models.MatchingResult) error {
	if len(results) == 0 {
		return nil
	}

	tx, err := r.DB.Begin()
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
		// kalau ada kolom yang bisa NULL, gunakan sql.NullXxx agar aman
		var customerName sql.NullString
		if res.CustomerName != "" {
			customerName = sql.NullString{String: res.CustomerName, Valid: true}
		} else {
			customerName = sql.NullString{Valid: false}
		}

		_, err = stmt.Exec(
			batchID,
			res.CIFNumber,
			customerName, // sudah ter-handle NULL
			res.WatchlistID,
			res.WatchlistSource,
			res.SimilarityScore,
			res.Status,
			res.ProcessDate,
			res.ProcessTime,
		)
		if err != nil {
			return err
		}
	}

	// flush buffer ke SQL Server
	if _, err = stmt.Exec(); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
