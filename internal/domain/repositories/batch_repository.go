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
	ctx := context.Background()

	batchID, err := r.GetSequencedId(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get sequenced id: %w", err)
	}

	query := `
		INSERT INTO BATCH_PROCESSING
			(Id, ProcessType, Status, TotalRecords, ProcessedRecords, MatchedRecords, StartTime, InitiatedBy, FilePath)
		VALUES
			(@Id, @ProcessType, 'running', @TotalRecords, 0, 0, GETDATE(), @InitiatedBy, @FilePath);
	`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		sql.Named("Id", batchID),
		sql.Named("ProcessType", batch.ProcessType),
		sql.Named("TotalRecords", batch.TotalRecords),
		sql.Named("InitiatedBy", batch.InitiatedBy),
		sql.Named("FilePath", batch.FilePath),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert batch: %w", err)
	}

	return batchID, nil
}

func (r *sqlBatchProcessingRepository) UpdateStatus(batchID int64, status string, processed, matched int, errMsg string) error {
	ctx := context.Background()

	query := `
		UPDATE BATCH_PROCESSING
		SET Status = @Status,
			ProcessedRecords = @ProcessedRecords,
			MatchedRecords = @MatchedRecords,
			EndTime = GETDATE(),
			ErrorMessage = @ErrorMessage
		WHERE Id = @Id
	`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		sql.Named("Status", status),
		sql.Named("ProcessedRecords", processed),
		sql.Named("MatchedRecords", matched),
		sql.Named("ErrorMessage", errMsg),
		sql.Named("Id", batchID),
	)
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
