package bootstrap

import (
	"BlackListWorker/config"
	"BlackListWorker/internal/db"
	"database/sql"
	"fmt"
)

type app struct{}

func NewApp() *app {
	return &app{}
}

func (a *app) Start() error {
	// 1. Load environment variables
	config.LoadEnv()

	// 2. Load config & connect to DB
	cfg := config.Load()
	conn := db.NewSQLServerConnector(cfg.DBServer, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	sqlDB, err := conn.Connect()
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	defer sqlDB.Close()

	// 3. Print current DB
	if err := printCurrentDB(sqlDB); err != nil {
		return err
	}

	// // 4. Init matching service
	// svc, err := services.NewMatchingService(sqlDB)
	// if err != nil {
	// 	return fmt.Errorf("create matching service: %w", err)
	// }

	// // 5. Run matching
	// results, details, err := svc.RunAll()
	// if err != nil {
	// 	return fmt.Errorf("run matching service: %w", err)
	// }

	// // 6. Save results
	// return saveResults(sqlDB, results, details)
	return nil
}

func printCurrentDB(sqlDB *sql.DB) error {
	rows, err := sqlDB.Query("SELECT DB_NAME() AS CurrentDB")
	if err != nil {
		return fmt.Errorf("query current db: %w", err)
	}
	defer rows.Close()

	var currentDB string
	for rows.Next() {
		if err := rows.Scan(&currentDB); err != nil {
			return fmt.Errorf("scan current db: %w", err)
		}
	}

	fmt.Println("Current Database:", currentDB)
	return nil
}

// // saveResults sekarang menerima map[int64][]models.MatchingDetail
// func saveResults(sqlDB *sql.DB, results []models.MatchingResult, details map[int64][]models.MatchingDetail) error {
// 	// Get next batch ID
// 	batchID, err := services.GetNextBatchID(sqlDB)
// 	if err != nil {
// 		return fmt.Errorf("get batch id: %w", err)
// 	}

// 	// Assign batchID ke setiap result
// 	for i := range results {
// 		results[i].BatchID = batchID
// 	}

// 	// Insert results dan detail
// 	if err := services.InsertMatchingResults(sqlDB, results, details); err != nil {
// 		return fmt.Errorf("failed to insert matching results: %w", err)
// 	}

// 	fmt.Printf("\nInserted %d results with BatchID %d\n", len(results), batchID)
// 	return nil
// }
