package repositories

import (
	"BlackListWorker/internal/domain/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type BatchProcessingRepository interface {
	Create(batch *models.BatchProcessing) (int64, error)
	UpdateStatus(batch *models.BatchProcessing) error
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
            (@Id, @ProcessType, @Status, @TotalRecords, 0, 0, GETDATE(), @InitiatedBy, @FilePath);
    `

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		sql.Named("Id", batchID),
		sql.Named("ProcessType", batch.ProcessType),
		sql.Named("Status", batch.Status),
		sql.Named("TotalRecords", batch.TotalRecords),
		sql.Named("InitiatedBy", batch.InitiatedBy),
		sql.Named("FilePath", batch.FilePath),
	)
	if err != nil {

		if strings.Contains(err.Error(), "PRIMARY KEY constraint") {

			syncErr := r.ResetSequenceId(ctx)
			if syncErr != nil {
				return 0, fmt.Errorf("failed to sync sequence: %w", syncErr)
			}

			batchID, err = r.GetSequencedId(ctx)
			if err != nil {
				return 0, fmt.Errorf("failed to get sequenced id after sync: %w", err)
			}

			_, err = stmt.ExecContext(ctx,
				sql.Named("Id", batchID),
				sql.Named("ProcessType", batch.ProcessType),
				sql.Named("Status", batch.Status),
				sql.Named("TotalRecords", batch.TotalRecords),
				sql.Named("InitiatedBy", batch.InitiatedBy),
				sql.Named("FilePath", batch.FilePath),
			)
			if err != nil {
				return 0, fmt.Errorf("failed to insert batch after sync: %w", err)
			}

			return batchID, nil
		}

		return 0, fmt.Errorf("failed to insert batch: %w", err)
	}

	return batchID, nil
}

func (r *sqlBatchProcessingRepository) UpdateStatus(b *models.BatchProcessing) error {
	ctx := context.Background()

	query := `
		UPDATE BATCH_PROCESSING
		SET Status = @Status,
			ProcessedRecords = @ProcessedRecords,
			MatchedRecords = @MatchedRecords,
			TotalRecords = @TotalRecords,
			EndTime = GETDATE(),
			ErrorMessage = @ErrorMessage,
			FilePath = @FilePath
		WHERE Id = @Id
	`

	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		sql.Named("Status", b.Status),
		sql.Named("TotalRecords", b.TotalRecords),
		sql.Named("ProcessedRecords", b.ProcessedRecords),
		sql.Named("MatchedRecords", b.MatchedRecords),
		sql.Named("ErrorMessage", b.ErrorMessage),
		sql.Named("FilePath", b.FilePath),
		sql.Named("Id", b.Id),
	)
	if err != nil {
		return fmt.Errorf("failed to update batch status: %w", err)
	}

	return nil
}

func (r *sqlBatchProcessingRepository) GetLastId(ctx context.Context) (int64, error) {
	var lastID sql.NullInt64

	query := `SELECT ISNULL(MAX(Id), 0) FROM BATCH_PROCESSING;`
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
	return err
}

func (r *sqlBatchProcessingRepository) GetSequencedId(ctx context.Context) (int64, error) {
	var newBatchId sql.NullInt64

	query := `SELECT NEXT VALUE FOR Seq_MatchingBatch AS NewBatchId;`
	err := r.DB.QueryRowContext(ctx, query).Scan(&newBatchId)
	if err != nil {
		return 0, err
	}

	if !newBatchId.Valid {
		return 0, nil
	}

	return newBatchId.Int64, nil
}
