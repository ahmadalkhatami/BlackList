package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repositories"
	"BlackListWorker/pkg/textutil"
)

type WatchlistServiceInterface interface {
	LoadDTTOT() ([]models.MasterTeroris, error)
	LoadWMD() ([]models.MasterWMD, error)
	LoadLocalBlacklist() ([]models.MasterLocalBlacklist, error)
	LoadAllWatchlists() ([]models.MasterWatchlist, error)
}

type WatchlistServiceImpl struct {
	MasterDTTOT          repositories.MasterTerorisRepository
	MasterWMD            repositories.MasterWMDRepository
	MasterLocalblacklist repositories.MasterLocalBlacklistRepository
}

type WatchlistOption func(*WatchlistServiceImpl)

func WithMasterTeroris(r repositories.MasterTerorisRepository) WatchlistOption {
	return func(w *WatchlistServiceImpl) {
		w.MasterDTTOT = r
	}
}

func WithMasterWMD(r repositories.MasterWMDRepository) WatchlistOption {
	return func(w *WatchlistServiceImpl) {
		w.MasterWMD = r
	}
}

func WithMasterLocalBlacklist(r repositories.MasterLocalBlacklistRepository) WatchlistOption {
	return func(w *WatchlistServiceImpl) {
		w.MasterLocalblacklist = r
	}
}

func NewWatchlistService(opts ...WatchlistOption) WatchlistServiceInterface {
	svc := &WatchlistServiceImpl{}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func (w *WatchlistServiceImpl) LoadDTTOT() ([]models.MasterTeroris, error) {
	return w.MasterDTTOT.LoadMasterTeroris()
}

func (w *WatchlistServiceImpl) LoadWMD() ([]models.MasterWMD, error) {
	return w.MasterWMD.LoadMasterWMD()
}

func (w *WatchlistServiceImpl) LoadLocalBlacklist() ([]models.MasterLocalBlacklist, error) {
	return w.MasterLocalblacklist.LoadMasterLocalBlacklist()
}

func (w *WatchlistServiceImpl) LoadAllWatchlists() ([]models.MasterWatchlist, error) {
	var combined []models.MasterWatchlist

	if w.MasterDTTOT != nil {
		dttot, err := w.MasterDTTOT.LoadMasterTeroris()
		if err != nil {
			return nil, err
		}
		combined = append(combined, ToWatchlistSlice(dttot)...)
	}

	if w.MasterWMD != nil {
		wmd, err := w.MasterWMD.LoadMasterWMD()
		if err != nil {
			return nil, err
		}
		combined = append(combined, ToWatchlistSlice(wmd)...)
	}

	if w.MasterLocalblacklist != nil {
		local, err := w.MasterLocalblacklist.LoadMasterLocalBlacklist()
		if err != nil {
			return nil, err
		}
		combined = append(combined, ToWatchlistSlice(local)...)
	}

	return combined, nil
}

func ToWatchlist(item models.Watchlistable) models.MasterWatchlist {
	return models.MasterWatchlist{
		ID:           item.GetID(),
		Nama:         item.GetNama(),
		Aliases:      textutil.CombineAliases(item, "Alias"),
		TempatLahir:  item.GetTempatLahir(),
		TanggalLahir: item.GetTanggalLahir(),
		KTP:          item.GetKTP(),
		NPWP:         item.GetNPWP(),
		NoPaspor:     item.GetNoPaspor(),
		Source:       item.GetSource(),
		Type:         item.GetType(),
		CreatedAt:    item.GetCreatedAt(),
		IsActive:     item.GetIsActive(),
	}
}

func ToWatchlistSlice[T any](items []T) []models.MasterWatchlist {
	result := make([]models.MasterWatchlist, 0, len(items))
	for i := range items {
		if wl, ok := any(&items[i]).(models.Watchlistable); ok {
			result = append(result, ToWatchlist(wl))
		}
	}
	return result
}
