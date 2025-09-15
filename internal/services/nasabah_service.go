package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repositories"
	"context"
)

type MasterNasabahService interface {
	Load(ctx context.Context) ([]models.MasterNasabah, error)
}

type MasterNasabahImpl struct {
	MasterNasabah repositories.MasterNasabahRepository
}

type MasterNasabahOption func(*MasterNasabahImpl)

func WithMasterNasabah(r repositories.MasterNasabahRepository) MasterNasabahOption {
	return func(m *MasterNasabahImpl) {
		m.MasterNasabah = r
	}
}

func NewMasterNasabahService(opts ...MasterNasabahOption) MasterNasabahService {
	svc := &MasterNasabahImpl{}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func (m *MasterNasabahImpl) Load(ctx context.Context) ([]models.MasterNasabah, error) {
	if m.MasterNasabah == nil {
		return nil, nil
	}
	return m.MasterNasabah.Load(ctx, repositories.WithStatus("ACTIVE"))
}
