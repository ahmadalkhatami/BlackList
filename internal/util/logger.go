package util

import (
	"database/sql"
	"fmt"
)

// func Info(msg string, args ...interface{}) {
// 	log.Printf("Info: "+msg, args...)
// }

// func Error(msg string, args ...interface{}) {
// 	log.Printf("Error: "+msg, args...)
// }

type DBLogger struct {
	db *sql.DB
}

// Constructor
func New(db *sql.DB) *DBLogger {
	return &DBLogger{db: db}
}

// Core function
func (l *DBLogger) Log(batchID, processName, level, message string) {
	_, err := l.db.Exec(`
		EXEC InsertProcessLog @BatchId = ?, @ProcessName = ?, @LogLevel = ?, @LogMessage = ?`,
		batchID, processName, level, message,
	)
	if err != nil {
		// fallback kalau gagal insert log → tampilkan di console
		fmt.Printf("[DB-LOGGING ERROR] %v\n", err)
	}
}

// Helper functions
func (l *DBLogger) Info(batchID, processName, message string) {
	l.Log(batchID, processName, "INFO", message)
}

func (l *DBLogger) Warn(batchID, processName, message string) {
	l.Log(batchID, processName, "WARN", message)
}

func (l *DBLogger) Error(batchID, processName, message string) {
	l.Log(batchID, processName, "ERROR", message)
}

/*
Example Use:
type sqlMatchingResultRepository struct {
	DB     *sql.DB
	Logger *dblogger.DBLogger
}

func (r *sqlMatchingResultRepository) SaveMatchingResultBatch(batchID string, results []models.MatchingResult) error {
	if len(results) == 0 {
		r.Logger.Info(batchID, "SaveMatchingResultBatch", "No results to insert")
		return nil
	}

	tx, err := r.DB.Begin()
	if err != nil {
		r.Logger.Error(batchID, "SaveMatchingResultBatch", fmt.Sprintf("Begin transaction failed: %v", err))
		return err
	}
	defer tx.Rollback()

	// ... proses CopyIn dll

	if err = tx.Commit(); err != nil {
		r.Logger.Error(batchID, "SaveMatchingResultBatch", fmt.Sprintf("Commit failed: %v", err))
		return err
	}

	r.Logger.Info(batchID, "SaveMatchingResultBatch", fmt.Sprintf("Successfully inserted %d rows", len(results)))
	return nil
}

*/
