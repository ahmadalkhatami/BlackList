package db

import (
	"context"
	"database/sql"
	"fmt"
)

func EnsureSequences(ctx context.Context, db *sql.DB) error {
	sequences := []string{
		"Seq_MatchingResult",
		"Seq_MatchingDetail",
		"Seq_MatchingBatch",
	}

	for _, seq := range sequences {
		query := fmt.Sprintf(`
		IF NOT EXISTS (SELECT * FROM sys.sequences WHERE name = '%s')
		BEGIN
			CREATE SEQUENCE %s
				START WITH 1
				INCREMENT BY 1;
		END
		`, seq, seq)

		_, err := db.ExecContext(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to create sequence %s: %w", seq, err)
		}
	}

	return nil
}
