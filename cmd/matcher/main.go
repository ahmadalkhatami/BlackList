package main

import (
	"BlackListWorker/config"
	"BlackListWorker/internal/db"
	"BlackListWorker/internal/domain/repository"
	"log"
)

func main() {
	cfg := config.Load()

	connector := db.SQLServerConnector{
		Server:  cfg.DBServer,
		User:    cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
	}

	sqlDB, err := connector.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	defer sqlDB.Close()

	masterMatchingRepo := repository.NewSQLMasterMatchingRepository(sqlDB)
	masterMatchingConfig := repository.NewSQLMatchingConfigRepository(sqlDB)
	masterTerorisRepo := repository.NewSQLMasterTerorisRepository(sqlDB)
}