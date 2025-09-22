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

	startTime := time.Now()

	debug := utils.IsDebugMode()
	fmt.Printf("DEBUG MODE: %t\n", debug)

	trigeredBy := utils.TrigeredBy()
	fmt.Printf("Triggered By: %d\n", trigeredBy)

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

	ctx := context.Background()

	container := NewContainer(sqlDB, ctx)
	if err != nil {
		fmt.Printf("❌ Error saat init container: %v\n", err)
		return err
	}

	err = db.EnsureSequences(ctx, sqlDB)
	if err != nil {
		fmt.Printf("❌ Error saat membuat sequeence: %v\n", err)
		return nil
	}

	err = container.Match.RunMatch(ctx)
	if err != nil {
		fmt.Printf("❌ Error saat matching: %v\n", err)
		return nil
	}

	duration := time.Since(startTime)

	results := container.Match.GetResults()

	if debug {
		fmt.Println("\n--- MATCH RESULTS ---")
		for _, r := range results.MatchResult {
			fmt.Printf(
				"Id=%d | BatchId=%v | CifNumber=%v | CustomerName=%v | WatchlistId=%v | WatchlistSource=%v | SimilarityScore=%v | Status=%v | ProcessDate=%v | ProcessTime=%v | CreatedAt=%v\n",
				r.Id,
				func() interface{} {
					if r.BatchId != nil {
						return *r.BatchId
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if r.CifNumber != nil {
						return *r.CifNumber
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if r.CustomerName != nil {
						return *r.CustomerName
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if r.WatchlistId != nil {
						return *r.WatchlistId
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if r.WatchlistSource != nil {
						return *r.WatchlistSource
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if r.SimilarityScore != nil {
						return *r.SimilarityScore
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if r.Status != nil {
						return *r.Status
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if r.ProcessDate != nil {
						return r.ProcessDate.Format("2006-01-02 15:04:05")
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if r.ProcessTime != nil {
						return r.ProcessTime.Format("2006-01-02 15:04:05")
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if r.CreatedAt != nil {
						return r.CreatedAt.Format("2006-01-02 15:04:05")
					} else {
						return "<nil>"
					}
				}(),
			)
		}

		fmt.Println("\n--- MATCH DETAILS ---")
		for _, d := range results.MatchDetail {
			fmt.Printf(
				"Id=%d | MatchingResultId=%v | FieldName=%v | CustomerValue=%v | WatchlistValue=%v | FieldScore=%v | FieldWeight=%v | AlgorithmUsed=%v\n",
				d.Id,
				func() interface{} {
					if d.MatchingResultId != nil {
						return *d.MatchingResultId
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if d.FieldName != nil {
						return *d.FieldName
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if d.CustomerValue != nil {
						return *d.CustomerValue
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if d.WatchlistValue != nil {
						return *d.WatchlistValue
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if d.FieldScore != nil {
						return *d.FieldScore
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if d.FieldWeight != nil {
						return *d.FieldWeight
					} else {
						return "<nil>"
					}
				}(),
				func() interface{} {
					if d.AlgorithmUsed != nil {
						return *d.AlgorithmUsed
					} else {
						return "<nil>"
					}
				}(),
			)
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
