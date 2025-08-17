package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repository"
)

type SystemConfigInterface interface {
	GetSystemConfig(configKey string) (models.SystemConfig, error)
	GetJoinedMatchingConfig() ([]models.JoinedMatchingConfig, error)
}

type SystemConfigImpl struct {
	MasterMatching       repository.MasterMatchingRepository
	MasterMatchingConfig repository.MasterMatchingConfigRepository
	SystemConfig         repository.SystemConfigRepository
}

func NewConfigService(
	masterMatching repository.MasterMatchingRepository,
	masterMatchingConfig repository.MasterMatchingConfigRepository,
	systemConfig repository.SystemConfigRepository) SystemConfigInterface {
	return &SystemConfigImpl{
		MasterMatching:       masterMatching,
		MasterMatchingConfig: masterMatchingConfig,
		SystemConfig:         systemConfig,
	}
}

func (s *SystemConfigImpl) GetSystemConfig(cfgKey string) (models.SystemConfig, error) {
	record, err := s.SystemConfig.LoadSystemConfig(cfgKey)
	if err != nil {
		return models.SystemConfig{}, err
	}
	return record, nil
}

func (s *SystemConfigImpl) GetJoinedMatchingConfig() ([]models.JoinedMatchingConfig, error) {

	matchingList, err := s.MasterMatching.LoadMasterMatching()
	if err != nil {
		return nil, err
	}

	configList, err := s.MasterMatchingConfig.LoadMasterMatchingConfig()
	if err != nil {
		return nil, err
	}

	var result []models.JoinedMatchingConfig

	for _, m := range matchingList {
		for _, c := range configList {
			if m.Id == c.MatchingId {
				result = append(result, models.JoinedMatchingConfig{
					Id:                c.Id,
					MatchingId:        c.MatchingId,
					FieldName:         c.FieldName,
					FieldWeight:       c.FieldWeight,
					WatchlistSource:   m.WatchlistSource,
					Type:              m.Type,
					MatchingAlgorithm: c.MatchingAlgorithm,
					IsActive:          c.IsActive && m.IsActive,
				})
			}
		}
	}

	return result, nil
}
