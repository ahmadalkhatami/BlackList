package bootstrap

import (
	"BlackListWorker/internal/domain/repositories"
	"BlackListWorker/internal/services"
	"database/sql"
)

type Container struct {
	DB        *sql.DB
	Config    services.SystemConfigService
	CIF       services.MasterNasabahService
	WatchList services.WatchlistService
	Match     services.MatchService
}

func NewContainer(db *sql.DB) *Container {

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

	// ------------------ matchers ------------------
	dttotMatcher := services.NewDTTOTMatcher(cifSvc, watchlistSvc, configSvc)
	wmdMatcher := services.NewWMDMatcher(cifSvc, watchlistSvc, configSvc)
	localMatcher := services.NewLocalBlacklistMatcher(cifSvc, watchlistSvc, configSvc)

	// ------------------ match service ------------------
	matchSvc := services.NewMatchService(dttotMatcher, wmdMatcher, localMatcher)

	return &Container{
		DB:        db,
		Config:    configSvc,
		CIF:       cifSvc,
		WatchList: watchlistSvc,
		Match:     matchSvc,
	}
}
