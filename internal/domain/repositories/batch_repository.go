package repositories

import (
	"BlackListWorker/internal/domain/models"
	"context"
	"database/sql"
	"fmt"
)

type BatchProcessingRepository interface {
	Create(batch *models.BatchProcessing) (int64, error)
	UpdateStatus(batchID int64, status string, processed, matched int, errMsg string) error
	GetLastId(ctx context.Context) (int64, error)
	ResetSequenceId(ctx context.Context) error
	GetSequencedId(ctx context.Context) (int64, error)
}

type sqlBatchProcessingRepository struct {
	DB *sql.DB
}

func NewSQLBatchProcessingRepository(db *sql.DB) BatchProcessingRepository {
	return &sqlBatchProcessingRepository{DB: db}
}

func (r *sqlBatchProcessingRepository) Create(batch *models.BatchProcessing) (int64, error) {
	query := `
		DECLARE @NewId BIGINT = NEXT VALUE FOR Seq_BatchProcessing;
		INSERT INTO BATCH_PROCESSING
			(Id, ProcessType, Status, TotalRecords, ProcessedRecords, MatchedRecords, StartTime, InitiatedBy, FilePath)
		VALUES
			(@NewId, ?, 'running', ?, 0, 0, GETDATE(), ?, ?);
		SELECT @NewId;
	`

	var id int64
	err := r.DB.QueryRow(query,
		batch.ProcessType,
		batch.TotalRecords,
		batch.InitiatedBy,
		batch.FilePath,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert batch: %w", err)
	}
	return id, nil
}

// Update batch status → completed / failed
func (r *sqlBatchProcessingRepository) UpdateStatus(batchID int64, status string, processed, matched int, errMsg string) error {
	query := `
		UPDATE BATCH_PROCESSING
		SET Status = ?,
			ProcessedRecords = ?,
			MatchedRecords = ?,
			EndTime = GETDATE(),
			ErrorMessage = ?
		WHERE Id = ?
	`
	_, err := r.DB.Exec(query, status, processed, matched, errMsg, batchID)
	if err != nil {
		return fmt.Errorf("failed to update batch status: %w", err)
	}
	return nil
}

func (r *sqlBatchProcessingRepository) GetLastId(ctx context.Context) (int64, error) {
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

func (r *sqlBatchProcessingRepository) ResetSequenceId(ctx context.Context) error {

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

func (r *sqlBatchProcessingRepository) GetSequencedId(ctx context.Context) (int64, error) {
	var NewBatchId sql.NullInt64

	query := `SELECT NEXT VALUE FOR Seq_MatchingBatch AS NewBatchId;`
	err := r.DB.QueryRowContext(ctx, query).Scan(&NewBatchId)
	if err != nil {
		return 0, err
	}

	if !NewBatchId.Valid {
		return 0, nil
	}

	return NewBatchId.Int64, nil
}
