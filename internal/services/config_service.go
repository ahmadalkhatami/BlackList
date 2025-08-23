package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repository"
	"strconv"
)

const defaultThreshold = 1

type SystemConfigInterface interface {
	GetSystemConfig(cfgKey string) (models.SystemConfig, error)
	GetJoinedMatchingConfig() ([]models.JoinedMatchingConfig, error)
	GetThresholdFromConfig(cfgKey string) (float64, error)
}

type SystemConfigImpl struct {
	MasterMatching       repository.MasterMatchingRepository
	MasterMatchingConfig repository.MasterMatchingConfigRepository
	SystemConfig         repository.SystemConfigRepository
}

type ConfigOption func(*SystemConfigImpl)

func WithMasterMatching(r repository.MasterMatchingRepository) ConfigOption {
	return func(s *SystemConfigImpl) {
		s.MasterMatching = r
	}
}

func WithMasterMatchingConfig(r repository.MasterMatchingConfigRepository) ConfigOption {
	return func(s *SystemConfigImpl) {
		s.MasterMatchingConfig = r
	}
}

func WithSystemConfig(r repository.SystemConfigRepository) ConfigOption {
	return func(s *SystemConfigImpl) {
		s.SystemConfig = r
	}
}

func NewConfigService(opts ...ConfigOption) SystemConfigInterface {
	svc := &SystemConfigImpl{}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
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

func (s *SystemConfigImpl) GetThresholdFromConfig(cfgKey string) (float64, error) {
	cfg, err := s.GetSystemConfig(cfgKey)
	if err != nil || cfg.ConfigKey == "" {
		return defaultThreshold, nil
	}

	if cfg.ConfigType == "DECIMAL" || cfg.ConfigType == "INTEGER" {
		if val, err := strconv.ParseFloat(cfg.ConfigValue, 64); err == nil {
			return val, nil
		}
	}

	return defaultThreshold, nil
}
