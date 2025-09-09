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
	LoadDTTOTIndividu() ([]models.MasterWatchlist, error)
	LoadWMDIndividu() ([]models.MasterWatchlist, error)
	LoadLocalBlacklistIndividu() ([]models.MasterWatchlist, error)
	LoadDTTOTCorporate() ([]models.MasterWatchlist, error)
	LoadWMDCorporate() ([]models.MasterWatchlist, error)
	LoadLocalBlacklistCorporate() ([]models.MasterWatchlist, error)
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
	return w.MasterDTTOT.Load()
}

func (w *WatchlistServiceImpl) LoadWMD() ([]models.MasterWMD, error) {
	return w.MasterWMD.Load()
}

func (w *WatchlistServiceImpl) LoadLocalBlacklist() ([]models.MasterLocalBlacklist, error) {
	return w.MasterLocalblacklist.Load()
}

func (w *WatchlistServiceImpl) LoadDTTOTIndividu() ([]models.MasterWatchlist, error) {
	var results []models.MasterWatchlist
	if w.MasterDTTOT != nil {
		dttot, err := w.MasterDTTOT.LoadIndividu()
		if err != nil {
			return []models.MasterWatchlist{}, err
		}
		results = append(results, ToWatchlistSlice(dttot)...)
	}
	return results, nil
}

func (w *WatchlistServiceImpl) LoadWMDIndividu() ([]models.MasterWatchlist, error) {
	var results []models.MasterWatchlist
	if w.MasterWMD != nil {
		list, err := w.MasterWMD.LoadIndividu()
		if err != nil {
			return []models.MasterWatchlist{}, err
		}
		results = append(results, ToWatchlistSlice(list)...)
	}
	return results, nil
}

func (w *WatchlistServiceImpl) LoadLocalBlacklistIndividu() ([]models.MasterWatchlist, error) {
	var results []models.MasterWatchlist
	if w.MasterLocalblacklist != nil {
		list, err := w.MasterLocalblacklist.LoadIndividu()
		if err != nil {
			return []models.MasterWatchlist{}, err
		}
		results = append(results, ToWatchlistSlice(list)...)
	}
	return results, nil
}

func (w *WatchlistServiceImpl) LoadDTTOTCorporate() ([]models.MasterWatchlist, error) {
	var results []models.MasterWatchlist
	if w.MasterDTTOT != nil {
		list, err := w.MasterDTTOT.LoadCorporate()
		if err != nil {
			return []models.MasterWatchlist{}, err
		}
		results = append(results, ToWatchlistSlice(list)...)
	}
	return results, nil
}

func (w *WatchlistServiceImpl) LoadWMDCorporate() ([]models.MasterWatchlist, error) {
	var results []models.MasterWatchlist
	if w.MasterWMD != nil {
		list, err := w.MasterWMD.LoadCorporate()
		if err != nil {
			return []models.MasterWatchlist{}, err
		}
		results = append(results, ToWatchlistSlice(list)...)
	}
	return results, nil
}

func (w *WatchlistServiceImpl) LoadLocalBlacklistCorporate() ([]models.MasterWatchlist, error) {
	var results []models.MasterWatchlist
	if w.MasterLocalblacklist != nil {
		list, err := w.MasterLocalblacklist.LoadCorporate()
		if err != nil {
			return []models.MasterWatchlist{}, err
		}
		results = append(results, ToWatchlistSlice(list)...)
	}
	return results, nil
}

func (w *WatchlistServiceImpl) LoadAllWatchlists() ([]models.MasterWatchlist, error) {
	var combined []models.MasterWatchlist

	if w.MasterDTTOT != nil {
		dttot, err := w.MasterDTTOT.Load()
		if err != nil {
			return nil, err
		}
		combined = append(combined, ToWatchlistSlice(dttot)...)
	}

	if w.MasterWMD != nil {
		wmd, err := w.MasterWMD.Load()
		if err != nil {
			return nil, err
		}
		combined = append(combined, ToWatchlistSlice(wmd)...)
	}

	if w.MasterLocalblacklist != nil {
		local, err := w.MasterLocalblacklist.Load()
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
