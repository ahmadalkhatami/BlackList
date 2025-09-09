package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repositories"
	"BlackListWorker/internal/utils"
	"strconv"
)

const defaultThreshold = 1

type SystemConfigInterface interface {
	Get(cfgKey string) (models.SystemConfig, error)
	GetJoinedMatchingConfig() ([]models.JoinedMatchingConfig, error)
	GetThreshold(cfgKey string) (float64, error)
}

type SystemConfigImpl struct {
	MasterMatching       repositories.MasterMatchingRepository
	MasterMatchingConfig repositories.MasterMatchingConfigRepository
	SystemConfig         repositories.SystemConfigRepository
}

type ConfigOption func(*SystemConfigImpl)

func WithMasterMatching(r repositories.MasterMatchingRepository) ConfigOption {
	return func(s *SystemConfigImpl) {
		s.MasterMatching = r
	}
}

func WithMasterMatchingConfig(r repositories.MasterMatchingConfigRepository) ConfigOption {
	return func(s *SystemConfigImpl) {
		s.MasterMatchingConfig = r
	}
}

func WithSystemConfig(r repositories.SystemConfigRepository) ConfigOption {
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

func (s *SystemConfigImpl) Get(cfgKey string) (models.SystemConfig, error) {
	record, err := s.SystemConfig.Load(cfgKey)
	if err != nil {
		return models.SystemConfig{}, err
	}
	return record, nil
}

func (s *SystemConfigImpl) GetJoinedMatchingConfig() ([]models.JoinedMatchingConfig, error) {

	configList, err := s.MasterMatchingConfig.Load()
	if err != nil {
		return nil, err
	}

	matchingList, err := s.MasterMatching.Load()
	if err != nil {
		return nil, err
	}

	matchingMap := make(map[string]models.MasterMatching, len(matchingList))
	for _, m := range matchingList {
		matchingMap[m.Id] = m
	}

	// generic MapSlice2
	result := utils.MapSlice2(
		configList,
		matchingMap,
		func(c models.MasterMatchingConfig) string { return c.MatchingId },
		MapConfigAndMatching,
	)

	return result, nil
}

func MapConfigAndMatching(c models.MasterMatchingConfig, m models.MasterMatching) models.JoinedMatchingConfig {
	return models.JoinedMatchingConfig{
		Id:                c.Id,
		MatchingId:        c.MatchingId,
		FieldName:         c.FieldName,
		FieldWeight:       c.FieldWeight,
		WatchlistSource:   m.WatchlistSource,
		Type:              m.Type,
		MatchingAlgorithm: c.MatchingAlgorithm,
		IsActive:          c.IsActive && m.IsActive,
	}
}

func (s *SystemConfigImpl) GetThreshold(cfgKey string) (float64, error) {
	cfg, err := s.Get(cfgKey)
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
