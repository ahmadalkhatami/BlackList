package services

import (
	"BlackListWorker/config"
	"BlackListWorker/internal/db"
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repository"
	"BlackListWorker/internal/services"
	"errors"
	"fmt"
)

var ErrUnknownAlgorithm = errors.New("Unknown Algorithm")

type MatchService interface {
	RunMatch(configID int) ([]models.MatchingResult, error)
}

type MatchServiceImpl struct {
	ConfigRepo               repository.MasterMatchingConfigRepository
	MasterTerorisRepo        repository.MasterTerorisRepository
	MasterWMDRepo            repository.MasterWMDRepository
	MasterLocalBlacklistRepo repository.MasterLocalBlacklistRepository
	MasterNasabahRepo        repository.MasterNasabahRepository
	NumWorker                int
}

func NewMatchService(
	cfg repository.MasterMatchingConfigRepository,
	masterTeroris repository.MasterTerorisRepository,
	masterWMD repository.MasterWMDRepository,
	masterLocalBlacklist repository.MasterLocalBlacklistRepository,
	masterNasabah repository.MasterNasabahRepository,
	numWorker int) MatchService {

	if numWorker < 1 {
		numWorker = 1
	}

	return &MatchServiceImpl{
		ConfigRepo:               cfg,
		MasterTerorisRepo:        masterTeroris,
		MasterWMDRepo:            masterWMD,
		MasterLocalBlacklistRepo: masterLocalBlacklist,
		MasterNasabahRepo:        masterNasabah,
		NumWorker:                numWorker,
	}
}

func (s *MatchServiceImpl) RunMatch(configID int) ([]models.MatchingResult, error) {

	config.LoadEnv()
	dbConfig := config.Load()

	connector := db.NewSQLServerConnector(dbConfig.DBServer, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName)
	sqlDB, err := connector.Connect()
	if err != nil {
		return []models.MatchingResult{}, fmt.Errorf("Failed connect to database : %w", err)
	}
	defer sqlDB.Close()

	masterMatchingRepo := repository.NewSQLMasterMatchingRepository(sqlDB)
	masterMatchingConfigRepo := repository.NewSQLMasterMatchingConfigRepository(sqlDB)
	systemConfigRepo := repository.NewSQLSystemConfigRepository(sqlDB)

	configService := services.NewConfigService(masterMatchingRepo, masterMatchingConfigRepo, systemConfigRepo)

	matchConfig, err := configService.GetJoinedMatchingConfig()
	if err != nil {
		return nil, err
	}

	//Load Config
	// cfg, err := s.ConfigRepo.LoadMasterMatchingConfig()
	// if err != nil {
	// 	return nil, err
	// }

	// var calc similarity.Calculator
	// switch cfg[0].MatchingAlgorithm {
	// case "jaro_winkler":
	// 	calc = &similarity.JaroWingklerCalculator{
	// 		BoostThreshold: cfg[0].FieldWeight,
	// 		PrefixSize:     4, // Default prefix size, can be adjusted
	// 	}
	// case "levenshtein":
	// 	calc = &similarity.LevenshteinCalculator{
	// 		MaxDistance: 100,
	// 	}
	// default:
	// 	return nil, ErrUnknownAlgorithm
	// }

	//Load Data
	// masterTeroris, err := s.MasterTerorisRepo.LoadMasterTeroris()
	// if err != nil {
	// 	return nil, err
	// }

	// masterWMD, err := s.MasterWMDRepo.LoadMasterWMD()
	// if err != nil {
	// 	return nil, err
	// }

	// masterLocalBlacklist, err := s.MasterLocalBlacklistRepo.LoadMasterLocalBlacklist()
	// if err != nil {
	// 	return nil, err
	// }

	// masterNasabah, err := s.MasterNasabahRepo.LoadMasterNasabah()
	// if err != nil {
	// 	return nil, err
	// }

	// // run worker manager
	// manager := &worker.WorkerManager{
	// 	NumWorker: s.NumWorker,
	// 	Threshold: cfg[0].FieldWeight,
	// 	Calc:      calc,
	// 	// Calculator: calc,
	// 	// MasterTeroris: masterTeroris,
	// 	// MasterWMD: masterWMD,
	// 	// MasterLocalBlacklist: masterLocalBlacklist,
	// 	// MasterNasabah: masterNasabah,
	// }

	// resultTeroris := manager.Run(masterTeroris, masterNasabah)
	// return resultTeroris, nil

	return nil, nil
}
