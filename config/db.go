package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

func ConnectDB() (*sql.DB, error) {
	connString := os.Getenv("DATABASE_URL")
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db.Ping error: %w", err)
	}

	return db, nil
}
