package bootstrap

import (
	"BlackListWorker/internal/domain/repositories"
	"BlackListWorker/internal/services"
	"context"
	"database/sql"
	"fmt"
)

type Container struct {
	DB        *sql.DB
	Config    services.SystemConfigService
	CIF       services.MasterNasabahService
	WatchList services.WatchlistService
	Match     services.MatchService
	BookID    *services.IDService
}

func NewContainer(db *sql.DB, ctx context.Context) *Container {

	// repo & service
	cifRepo := repositories.NewSQLMasterNasabahRepository(db)
	cifSvc := services.NewMasterNasabahService(
		services.WithMasterNasabah(cifRepo),
	)
	watchlistSvc := services.NewWatchlistService(
		services.WithMasterLocalBlacklist(repositories.NewSQLMasterLocalBlacklistRepository(db)),
		services.WithMasterTeroris(repositories.NewSQLMasterTerorisRepository(db)),
		services.WithMasterWMD(repositories.NewSQLMasterWMDRepository(db)),
	)
	configSvc := services.NewConfigService(
		services.WithMasterMatching(repositories.NewSQLMasterMatchingRepository(db)),
		services.WithMasterMatchingConfig(repositories.NewSQLMasterMatchingConfigRepository(db)),
		services.WithSystemConfig(repositories.NewSQLSystemConfigRepository(db)),
	)

	resultRepo := repositories.NewSQLMatchingResultRepository(db)
	detailRepo := repositories.NewSQLMatchingDetailsRepository(db)
	batchRepo := repositories.NewSQLBatchProcessingRepository(db)

	// 1. Buat shared IDService
	idSvc := services.NewIDService(resultRepo, detailRepo, batchRepo)
	if err := idSvc.InitAtomicIDs(ctx); err != nil {
		panic(fmt.Sprintf("failed to init IDService: %v", err))
	}

	// 2. Inisialisasi matcher dengan pointer IDService yang sama
	dttotMatcher := services.NewDTTOTMatcher(cifSvc, watchlistSvc, configSvc, idSvc)
	wmdMatcher := services.NewWMDMatcher(cifSvc, watchlistSvc, configSvc, idSvc)
	localMatcher := services.NewLocalBlacklistMatcher(cifSvc, watchlistSvc, configSvc, idSvc)

	// 3. Buat MatchService dengan semua matcher
	matchSvc := services.NewMatchService(dttotMatcher, wmdMatcher, localMatcher)

	return &Container{
		DB:        db,
		Config:    configSvc,
		CIF:       cifSvc,
		WatchList: watchlistSvc,
		Match:     matchSvc,
		BookID:    idSvc, // simpan shared IDService di container
	}
}
