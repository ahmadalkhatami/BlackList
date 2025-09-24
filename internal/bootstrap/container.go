package bootstrap

import (
	"BlackListWorker/internal/domain/repositories"
	"BlackListWorker/internal/services"
	"context"
	"database/sql"
	"fmt"
)

type Container struct {
	DB            *sql.DB
	Config        services.SystemConfigService
	CIF           services.MasterNasabahService
	WatchList     services.WatchlistService
	Match         services.MatchService
	BookService   *services.IDService
	ResultService services.MatchingResultService
}

func NewContainer(db *sql.DB, ctx context.Context) *Container {

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

	idSvc := services.NewIDService(resultRepo, detailRepo, batchRepo)
	if err := idSvc.InitAtomicIDs(ctx); err != nil {
		panic(fmt.Sprintf("failed to init IDService: %v", err))
	}

	dttotMatcher := services.NewDTTOTMatcher(cifSvc, watchlistSvc, configSvc, idSvc)
	wmdMatcher := services.NewWMDMatcher(cifSvc, watchlistSvc, configSvc, idSvc)
	localMatcher := services.NewLocalBlacklistMatcher(cifSvc, watchlistSvc, configSvc, idSvc)

	matchSvc := services.NewMatchService(batchRepo, dttotMatcher, wmdMatcher, localMatcher)

	resultSvc := services.NewMatchingResultService(resultRepo, detailRepo)

	return &Container{
		DB:            db,
		Config:        configSvc,
		CIF:           cifSvc,
		WatchList:     watchlistSvc,
		Match:         matchSvc,
		BookService:   idSvc,
		ResultService: resultSvc,
	}
}
