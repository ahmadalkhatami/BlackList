package bootstrap

import (
	"BlackListWorker/config"
	"BlackListWorker/internal/db"
	"BlackListWorker/internal/domain/repository"
	"BlackListWorker/internal/services"
	"fmt"
)

type app struct{}

func NewApp() *app {
	return &app{}
}

func (a *app) Start() error {

	config.LoadEnv()
	dbConfig := config.Load()

	connector := db.NewSQLServerConnector(dbConfig.DBServer, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName)
	sqlDB, err := connector.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer sqlDB.Close()

	/*
		rows, err := sqlDB.Query("SELECT DB_NAME() AS CurrentDB")
		var currentDB string
		for rows.Next() {
			rows.Scan(&currentDB)
		}
		fmt.Println("Current Database:", currentDB) */

	masterMatchingRepo := repository.NewSQLMasterMatchingRepository(sqlDB)
	masterMatchingConfigRepo := repository.NewSQLMasterMatchingConfigRepository(sqlDB)
	systemConfigRepo := repository.NewSQLSystemConfigRepository(sqlDB)

	configService := services.NewConfigService(masterMatchingRepo, masterMatchingConfigRepo, systemConfigRepo)

	matchConfigs, err := configService.GetJoinedMatchingConfig()
	if err != nil {
		return fmt.Errorf("error while getting configs: %w", err)
	}

	for _, j := range matchConfigs {
		fmt.Printf("ID=%s, MatchingID=%s, Source=%s, Field=%s, Weight=%.2f, Type=%t, Algorithm=%s\n",
			j.Id, j.MatchingId, j.WatchlistSource, j.FieldName, j.FieldWeight, j.Type, j.MatchingAlgorithm)
	}

	return nil
}
