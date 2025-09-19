package bootstrap

import (
	"BlackListWorker/config"
	"BlackListWorker/internal/db"
	"BlackListWorker/internal/utils"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type app struct{}

func NewApp() *app {
	return &app{}
}

func (a *app) Start() error {

	config.LoadEnv()

	cfg := config.Load()
	conn := db.NewSQLServerConnector(cfg.DBServer, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	sqlDB, err := conn.Connect()
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	defer sqlDB.Close()

	if err := printCurrentDB(sqlDB); err != nil {
		return err
	}

	// 2️⃣ Buat container (service & matchers)
	container := NewContainer(sqlDB)

	startTime := time.Now()

	// 3️⃣ Jalankan match service
	ctx := context.Background()
	err = container.Match.RunMatch(ctx)
	if err != nil {
		fmt.Printf("❌ Error saat matching: %v\n", err)
		return nil
	}

	duration := time.Since(startTime)

	// 4️⃣ Ambil hasil matching
	results := container.Match.GetResults()
	fmt.Printf("✅ Total match: %d | Waktu proses: %s\n", len(results.MatchResult), duration)

	debug := utils.IsDebugMode()
	if debug {
		for _, r := range results.MatchResult {
			fmt.Printf("CIF: %s | Watchlist: %d | Score: %.2f\n", *r.GetCifNumber(), *r.GetWatchlistId(), *r.GetSimilarityScore())
		}
	}

	fmt.Printf("⏱️ Matching selesai dalam: %s\n", duration)

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
