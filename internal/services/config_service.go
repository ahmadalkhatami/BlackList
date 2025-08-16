package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repository"
)

type GetSystemConfigByKey interface {
	GetSystemConfig(configKey string) (models.SystemConfig, error)
	GetJoinedMatchingConfig() ([]JoinedMatchingConfig, error)
}

type JoinedMatchingConfig struct {
	Id                string
	WatchlistSource   string
	Type              bool
	MatchingId        string
	FieldName         string
	FieldWeight       float64
	MatchingAlgorithm string
	IsActive          bool
}

type SystemConfigImpl struct {
	MasterMatching       repository.MasterMatchingRepository
	MasterMatchingConfig repository.MasterMatchingConfigRepository
	SystemConfig         repository.SystemConfigRepository
}

func NewConfigService(
	masterMatching repository.MasterMatchingRepository,
	masterMatchingConfig repository.MasterMatchingConfigRepository,
	systemConfig repository.SystemConfigRepository) GetSystemConfigByKey {
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

func (s *SystemConfigImpl) GetJoinedMatchingConfig() ([]JoinedMatchingConfig, error) {

	matchingList, err := s.MasterMatching.LoadMasterMatching()
	if err != nil {
		return nil, err
	}

	configList, err := s.MasterMatchingConfig.LoadMasterMatchingConfig()
	if err != nil {
		return nil, err
	}

	var result []JoinedMatchingConfig

	for _, m := range matchingList {
		for _, c := range configList {
			if m.Id == c.MatchingId {
				result = append(result, JoinedMatchingConfig{
					Id:                c.Id,
					WatchlistSource:   m.WatchlistSource,
					Type:              m.Type,
					MatchingId:        c.MatchingId,
					FieldName:         c.FieldName,
					FieldWeight:       c.FieldWeight,
					MatchingAlgorithm: c.MatchingAlgorithm,
					IsActive:          c.IsActive,
				})
			}
		}
	}

	return result, nil
}
