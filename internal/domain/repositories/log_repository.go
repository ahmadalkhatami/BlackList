package repositories

import (
	"BlackListWorker/internal/domain/models"
	"database/sql"

	mssql "github.com/denisenkom/go-mssqldb"
)

type ProcessLogRepository interface {
	LoadProcessLog() ([]models.ProcessingLog, error)
	SaveErrorProcessLog(processLog []models.ProcessingLog) error
}

type sqlProcessLogRepository struct {
	DB *sql.DB
}

func NewSQLProcessLogRepository(db *sql.DB) ProcessLogRepository {
	return &sqlProcessLogRepository{DB: db}
}

func (r sqlProcessLogRepository) LoadProcessLog() ([]models.ProcessingLog, error) {
	return nil, nil
}

func (r sqlProcessLogRepository) SaveErrorProcessLog(processLog []models.ProcessingLog) error {
	if len(processLog) == 0 {
		return nil
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(mssql.CopyIn(
		"PROCESSING_LOG",
		mssql.BulkOptions{KeepNulls: true},
		"BatchId",
		"LogLevel",
		"LogMessage",
		"ErrorDetails",
		"LogTime",
	))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, res := range processLog {
		_, err = stmt.Exec(
			res.BatchID,
			res.LogLevel,
			res.LogMessage,
			res.ErrorDetail,
			res.LogTime,
		)
		if err != nil {
			return err
		}
	}

	if _, err = stmt.Exec(); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return nil
	}

	return nil
}
