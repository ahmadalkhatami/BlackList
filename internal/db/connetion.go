package db

import (
	"database/sql"
	"fmt"
)

type DBConnector interface {
	Connect() (*sql.DB, error)
}

type SQLServerConnector struct {
	Server   string
	User     string
	Password string
	Database string
}

func (c SQLServerConnector) Connect() (*sql.DB, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", c.Server, c.User, c.Password, c.Database)
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
