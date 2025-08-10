package repository

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MatchingResultRepository interface {
	LoadMatchingResult(query string) ([]models.MatchingResult, error)
	SaveMatchingResult(results []models.MatchingResult) error
}

type SQLMatchingResultRepository struct {
	DB *sql.DB
}

func NewSQLMatchingResultRepository(db *sql.DB) MatchingResultRepository {
	return &SQLMatchingResultRepository{DB: db}
}

func (r SQLMatchingResultRepository) LoadMatchingResult(query string) ([]models.MatchingResult, error) {
	rows, err := r.DB.Query(query)
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
			&rec.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, nil
}

func (r *SQLMatchingResultRepository) SaveMatchingResult(results []models.MatchingResult) error {
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
