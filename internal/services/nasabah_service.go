package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repository"
)

type MasterNasabahInterface interface {
	Load() ([]models.MasterNasabah, error)
}

type MasterNasabahImpl struct {
	MasterNasabah repository.MasterNasabahRepository
}

func NewMasterNasabah(masterNasabah repository.MasterNasabahRepository) MasterNasabahInterface {
	return &MasterNasabahImpl{
		MasterNasabah: masterNasabah,
	}
}

func (m *MasterNasabahImpl) Load() ([]models.MasterNasabah, error) {

	result, err := m.MasterNasabah.LoadMasterNasabah()
	if err != nil {
		return []models.MasterNasabah{}, err
	}

	return result, nil
}
