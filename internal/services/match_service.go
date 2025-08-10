package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repository"
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
	cfg, err := s.ConfigRepo.LoadMasterMatchingConfig("SELECT * FROM ")
}