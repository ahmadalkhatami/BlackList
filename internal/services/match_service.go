package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repository"
	"BlackListWorker/internal/domain/similarity"
	"errors"
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

	if numWorker <= 0 {
		numWorker = 0
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

	//Load Config
	cfg, err := s.ConfigRepo.LoadMasterMatchingConfig("SELECT * FROM ")
	if err != nil {
		return nil, err
	}

	var calc similarity.Calculator
	switch cfg[0].MatchingAlgorithm {
	case "jaro_winkler":
		calc = similarity.NewJaroWinklerCalculator{
			BoostThreshold: cfg[0].FieldWeight,
			PrefixSize:     4, // Default prefix size, can be adjusted
		}
	case "levenshtein":
		calc = similarity.NewLevenshteinCalculator{
			MaxDistance: cfg[0].FieldWeight,
		}
	default:
		return nil, ErrUnknownAlgorithm
	}
}
