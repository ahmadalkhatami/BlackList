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
		query := `
		IF NOT EXISTS (SELECT * FROM sys.sequences WHERE name = ?)
		BEGIN
			CREATE SEQUENCE ` + seq + `
				START WITH 1
				INCREMENT BY 1;
		END
		`

		_, err := db.ExecContext(ctx, query, seq)
		if err != nil {
			return fmt.Errorf("failed to create sequence %s: %w", seq, err)
		}
	}

	return nil
}
