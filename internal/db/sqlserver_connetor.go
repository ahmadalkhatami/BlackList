package db

import (
	"database/sql"
	"fmt"
)

// private means accessible only within the same package

// Abstraksi
// upper case prefix for public
type DBConnector interface {
	Connect() (*sql.DB, error)
}

// SQLServerConnector implements DBConnector for SQL Server
// lower case prefix for private
type sQLServerConnector struct {
	server   string
	user     string
	password string
	database string
}

func NewSQLServerConnector(server, user, password, database string) *sQLServerConnector {
	return &sQLServerConnector{
		server:   server,
		user:     user,
		password: password,
		database: database,
	}
}

func (c sQLServerConnector) Connect() (*sql.DB, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", c.server, c.user, c.password, c.database)
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

/*
Implementation :

func SomeFunction() {

	var connector db.DBConnector
	connector = db.NewSQLServerConnector("localhost", "sa", "password", "BlackList")

	sqlDB, err := connector.Connect()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer sqlDB.Close()

	fmt.Println("Connected to SQL Server!")
}
*/
