package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repositories"
	"BlackListWorker/internal/utils"
	"context"
	"strconv"
)

const defaultThreshold = 1

type SystemConfigService interface {
	Get(ctx context.Context, opts ...repositories.SystemConfigOption) (models.SystemConfig, error)
	GetJoinedMatchingConfig(ctx context.Context) ([]models.JoinedMatchingConfig, error)
	GetThreshold(ctx context.Context) (float64, error)
	GetMatchingAlgorithm(ctx context.Context) (string, error)
	LoadIndividu(ctx context.Context) ([]models.JoinedMatchingConfig, error)
	LoadCorporate(ctx context.Context) ([]models.JoinedMatchingConfig, error)
}

type systemConfigService struct {
	MasterMatching       repositories.MasterMatchingRepository
	MasterMatchingConfig repositories.MasterMatchingConfigRepository
	SystemConfig         repositories.SystemConfigRepository
}

type ConfigOption func(*systemConfigService)

func WithMasterMatching(r repositories.MasterMatchingRepository) ConfigOption {
	return func(s *systemConfigService) { s.MasterMatching = r }
}

func WithMasterMatchingConfig(r repositories.MasterMatchingConfigRepository) ConfigOption {
	return func(s *systemConfigService) { s.MasterMatchingConfig = r }
}

func WithSystemConfig(r repositories.SystemConfigRepository) ConfigOption {
	return func(s *systemConfigService) { s.SystemConfig = r }
}

func NewConfigService(opts ...ConfigOption) SystemConfigService {
	svc := &systemConfigService{}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func (s *systemConfigService) Get(ctx context.Context, opts ...repositories.SystemConfigOption) (models.SystemConfig, error) {
	return s.SystemConfig.LoadOne(ctx, opts...)
}

func (s *systemConfigService) LoadIndividu(ctx context.Context) ([]models.JoinedMatchingConfig, error) {
	return s.loadJoined(ctx, func() ([]models.MasterMatching, error) {
		return s.MasterMatching.Load(ctx, repositories.WithIndividual(true))
	})
}

func (s *systemConfigService) LoadCorporate(ctx context.Context) ([]models.JoinedMatchingConfig, error) {
	return s.loadJoined(ctx, func() ([]models.MasterMatching, error) {
		return s.MasterMatching.Load(ctx, repositories.WithIndividual(false))
	})
}

func (s *systemConfigService) GetJoinedMatchingConfig(ctx context.Context) ([]models.JoinedMatchingConfig, error) {
	return s.loadJoined(ctx, func() ([]models.MasterMatching, error) {
		return s.MasterMatching.Load(ctx)
	})
}

func (s *systemConfigService) loadJoined(
	ctx context.Context,
	loadMatching func() ([]models.MasterMatching, error),
) ([]models.JoinedMatchingConfig, error) {

	configList, err := s.MasterMatchingConfig.Load(ctx)
	if err != nil {
		return nil, err
	}

	matchingList, err := loadMatching()
	if err != nil {
		return nil, err
	}

	matchingMap := make(map[string]models.MasterMatching, len(matchingList))
	for _, m := range matchingList {
		matchingMap[m.Id] = m
	}

	result := utils.MapSlice2(
		configList,
		matchingMap,
		func(c models.MasterMatchingConfig) string { return c.MatchingId },
		MapConfigMatching,
	)

	return result, nil
}

func MapConfigMatching(c models.MasterMatchingConfig, m models.MasterMatching) models.JoinedMatchingConfig {
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

func (s *systemConfigService) GetThreshold(ctx context.Context) (float64, error) {
	cfg, err := s.Get(ctx, repositories.WithConfigKey("MATCHING_THRESHOLD"))
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

func (s *systemConfigService) GetMatchingAlgorithm(ctx context.Context) (string, error) {
	cfg, err := s.Get(ctx, repositories.WithConfigKey("MATCHING_SIMILARITY_METHOD"))
	if err != nil || cfg.ConfigKey == "" {
		return "fuzzywuzzy", nil
	}

	return cfg.ConfigValue, nil
}
