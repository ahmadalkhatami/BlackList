package repositories

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"
)

type MatchingDetailsRepository interface {
	LoadMatchingDetails(query string) ([]models.MatchingDetail, error)
	SaveMatchingDetails(details []models.MatchingDetail) error
	SaveMatchingDetailsBatch(batchID string, details []models.MatchingDetail) error
}

type sqlMatchingDetailsRepository struct {
	DB *sql.DB
}

func NewSQLMatchingDetailsRepository(db *sql.DB) MatchingDetailsRepository {
	return &sqlMatchingDetailsRepository{DB: db}
}

func (r sqlMatchingDetailsRepository) LoadMatchingDetails(query string) ([]models.MatchingDetail, error) {
	rows, err := r.DB.Query(query)
	if err != nil {
		return []models.MatchingDetail{}, err
	}

	defer rows.Close()

	var records []models.MatchingDetail
	for rows.Next() {
		var rec models.MatchingDetail
		if err := rows.Scan(
			&rec.ID,
			&rec.MatchingResultID,
			&rec.FieldName,
			&rec.CustomerValue,
			&rec.WatchlistValue,
			&rec.FieldScore,
			&rec.FieldWeight,
			&rec.AlgorithmUsed); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func (r sqlMatchingDetailsRepository) SaveMatchingDetails(details []models.MatchingDetail) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO [dbo].[MATCHING_DETAILS] (MatchingResultId, FieldName, CustomerValue, WatchlistValue, FieldScore, FieldWeight, AlgorithmUsed) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, detail := range details {
		_, err := stmt.Exec(
			detail.MatchingResultID,
			detail.FieldName,
			detail.CustomerValue,
			detail.WatchlistValue,
			detail.FieldScore,
			detail.FieldWeight,
			detail.AlgorithmUsed)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r sqlMatchingDetailsRepository) SaveMatchingDetailsBatch(batchID string, details []models.MatchingDetail) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO [dbo].[MATCHING_DETAILS] (MatchingResultId, FieldName, CustomerValue, WatchlistValue, FieldScore, FieldWeight, AlgorithmUsed) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, detail := range details {
		_, err := stmt.Exec(
			detail.MatchingResultID,
			detail.FieldName,
			detail.CustomerValue,
			detail.WatchlistValue,
			detail.FieldScore,
			detail.FieldWeight,
			detail.AlgorithmUsed)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
