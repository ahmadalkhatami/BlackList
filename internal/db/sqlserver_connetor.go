package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

// private means accessible only within the same package

// Abstraksi
// upper case prefix for public
type DBConnector interface {
	Connect() (*sql.DB, error)
}

// SQLServerConnector implements DBConnector for SQL Server
// lower case prefix for private
type sqlServerConnector struct {
	server   string
	user     string
	password string
	database string
}

func NewSQLServerConnector(server, user, password, database string) *sqlServerConnector {
	return &sqlServerConnector{
		server:   server,
		user:     user,
		password: password,
		database: database,
	}
}

func (c *sqlServerConnector) Connect() (*sql.DB, error) {
	// connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
	// 	c.server, c.user, c.password, c.database)

	// connString := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s",
	// 	c.user, c.password, c.server, c.database)

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable",
		c.server, c.user, c.password, c.database)

	log.Printf("Connecting with: sqlserver://%s:****@%s?database=%s",
		c.user, c.server, c.database)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal ping ke SQL Server: %w", err)
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
