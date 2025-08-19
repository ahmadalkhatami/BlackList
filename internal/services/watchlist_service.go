package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repository"
	"database/sql"
	"log"
)

type WatchlistServiceInterface interface {
	LoadDTTOT() ([]models.MasterTeroris, error)
	LoadWMD() ([]models.MasterWMD, error)
	LoadLocalBlacklist() ([]models.MasterLocalBlacklist, error)
	CombineAliases(a1, a2, a3, a4 sql.NullString) ([]string, error)
}

type WatchlistServiceImpl struct {
	MasterDTTOT          repository.MasterTerorisRepository
	MasterWMD            repository.MasterWMDRepository
	MasterLocalblacklist repository.MasterLocalBlacklistRepository
}

func NewWatchlistService(
	masterDTTOT repository.MasterTerorisRepository,
	masterWMD repository.MasterWMDRepository,
	masterLocalBlacklist repository.MasterLocalBlacklistRepository,
) WatchlistServiceInterface {
	return &WatchlistServiceImpl{
		MasterDTTOT:          masterDTTOT,
		MasterWMD:            masterWMD,
		MasterLocalblacklist: masterLocalBlacklist,
	}
}

func (w *WatchlistServiceImpl) LoadDTTOT() ([]models.MasterTeroris, error) {

	records, err := w.MasterDTTOT.LoadMasterTeroris()
	if err != nil {
		// return []models.MasterTeroris{}, err
		log.Fatalf("Failed to load Master Teroris: %v", err)
	}

	return records, nil
}

func (w *WatchlistServiceImpl) LoadWMD() ([]models.MasterWMD, error) {

	records, err := w.MasterWMD.LoadMasterWMD()
	if err != nil {
		// return []models.MasterWMD{}, err
		log.Fatalf("Failed to load Master WMD: %v", err)
	}

	return records, nil
}

func (w *WatchlistServiceImpl) LoadLocalBlacklist() ([]models.MasterLocalBlacklist, error) {
	records, err := w.MasterLocalblacklist.LoadMasterLocalBlacklist()
	if err != nil {
		// return []models.MasterLocalBlacklist{}, err
		log.Fatalf("Failed to load Master Local Blacklist: %v", err)
	}

	return records, nil
}

func (w *WatchlistServiceImpl) CombineAliases(a1, a2, a3, a4 sql.NullString) ([]string, error) {
	aliases := []string{}
	for _, a := range []sql.NullString{a1, a2, a3, a4} {
		if a.Valid && a.String != "" {
			aliases = append(aliases, a.String)
		}
	}
	return aliases, nil
}
