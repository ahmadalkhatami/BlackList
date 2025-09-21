package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repositories"
	"BlackListWorker/pkg/textutil"
	"context"
)

type WatchlistService interface {
	LoadDTTOT(ctx context.Context) ([]models.MasterTeroris, error)
	LoadWMD(ctx context.Context) ([]models.MasterWMD, error)
	LoadLocalBlacklist(ctx context.Context) ([]models.MasterLocalBlacklist, error)

	LoadDTTOTIndividu(ctx context.Context) ([]models.MasterWatchlist, error)
	LoadWMDIndividu(ctx context.Context) ([]models.MasterWatchlist, error)
	LoadLocalBlacklistIndividu(ctx context.Context) ([]models.MasterWatchlist, error)

	LoadDTTOTCorporate(ctx context.Context) ([]models.MasterWatchlist, error)
	LoadWMDCorporate(ctx context.Context) ([]models.MasterWatchlist, error)
	LoadLocalBlacklistCorporate(ctx context.Context) ([]models.MasterWatchlist, error)

	LoadAllWatchlists(ctx context.Context) ([]models.MasterWatchlist, error)
}

type WatchlistServiceImpl struct {
	MasterDTTOT          repositories.MasterTerorisRepository
	MasterWMD            repositories.MasterWMDRepository
	MasterLocalblacklist repositories.MasterLocalBlacklistRepository
}

type WatchlistOption func(*WatchlistServiceImpl)

func WithMasterTeroris(r repositories.MasterTerorisRepository) WatchlistOption {
	return func(w *WatchlistServiceImpl) { w.MasterDTTOT = r }
}

func WithMasterWMD(r repositories.MasterWMDRepository) WatchlistOption {
	return func(w *WatchlistServiceImpl) { w.MasterWMD = r }
}

func WithMasterLocalBlacklist(r repositories.MasterLocalBlacklistRepository) WatchlistOption {
	return func(w *WatchlistServiceImpl) { w.MasterLocalblacklist = r }
}

func NewWatchlistService(opts ...WatchlistOption) WatchlistService {
	svc := &WatchlistServiceImpl{}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// ==================== Master Loader ====================
func (w *WatchlistServiceImpl) LoadDTTOT(ctx context.Context) ([]models.MasterTeroris, error) {
	if w.MasterDTTOT == nil {
		return []models.MasterTeroris{}, nil
	}
	return w.MasterDTTOT.Load(ctx)
}

func (w *WatchlistServiceImpl) LoadWMD(ctx context.Context) ([]models.MasterWMD, error) {
	if w.MasterWMD == nil {
		return []models.MasterWMD{}, nil
	}
	return w.MasterWMD.Load(ctx)
}

func (w *WatchlistServiceImpl) LoadLocalBlacklist(ctx context.Context) ([]models.MasterLocalBlacklist, error) {
	if w.MasterLocalblacklist == nil {
		return []models.MasterLocalBlacklist{}, nil
	}
	return w.MasterLocalblacklist.Load(ctx)
}

// ==================== Watchlist Loader Individu ====================
func (w *WatchlistServiceImpl) LoadDTTOTIndividu(ctx context.Context) ([]models.MasterWatchlist, error) {
	if w.MasterDTTOT == nil {
		return []models.MasterWatchlist{}, nil
	}

	list, err := w.MasterDTTOT.Load(ctx,
		repositories.WithTerorisActive(true),
		repositories.WithTerorisType("individu"),
	)
	if err != nil {
		return nil, err
	}
	return ToWatchlistSlice(list), nil
}

func (w *WatchlistServiceImpl) LoadWMDIndividu(ctx context.Context) ([]models.MasterWatchlist, error) {
	if w.MasterWMD == nil {
		return []models.MasterWatchlist{}, nil
	}

	list, err := w.MasterWMD.Load(ctx,
		repositories.WithWMDActive(true),
		repositories.WithWMDType("individu"),
	)
	if err != nil {
		return nil, err
	}
	return ToWatchlistSlice(list), nil
}

func (w *WatchlistServiceImpl) LoadLocalBlacklistIndividu(ctx context.Context) ([]models.MasterWatchlist, error) {
	if w.MasterLocalblacklist == nil {
		return []models.MasterWatchlist{}, nil
	}

	list, err := w.MasterLocalblacklist.Load(ctx,
		repositories.WithLBActive(true),
		repositories.WithLBType("individu"),
	)
	if err != nil {
		return nil, err
	}
	return ToWatchlistSlice(list), nil
}

// ==================== Watchlist Loader Corporate ====================
func (w *WatchlistServiceImpl) LoadDTTOTCorporate(ctx context.Context) ([]models.MasterWatchlist, error) {
	if w.MasterDTTOT == nil {
		return []models.MasterWatchlist{}, nil
	}

	list, err := w.MasterDTTOT.Load(ctx,
		repositories.WithTerorisActive(true),
		repositories.WithTerorisType("korporasi"),
	)
	if err != nil {
		return nil, err
	}
	return ToWatchlistSlice(list), nil
}

func (w *WatchlistServiceImpl) LoadWMDCorporate(ctx context.Context) ([]models.MasterWatchlist, error) {
	if w.MasterWMD == nil {
		return []models.MasterWatchlist{}, nil
	}

	list, err := w.MasterWMD.Load(ctx,
		repositories.WithWMDActive(true),
		repositories.WithWMDType("korporasi"),
	)
	if err != nil {
		return nil, err
	}
	return ToWatchlistSlice(list), nil
}

func (w *WatchlistServiceImpl) LoadLocalBlacklistCorporate(ctx context.Context) ([]models.MasterWatchlist, error) {
	if w.MasterLocalblacklist == nil {
		return []models.MasterWatchlist{}, nil
	}

	list, err := w.MasterLocalblacklist.Load(ctx,
		repositories.WithLBActive(true),
		repositories.WithLBType("korporasi"),
	)
	if err != nil {
		return nil, err
	}
	return ToWatchlistSlice(list), nil
}

// ==================== Combine All ====================
func (w *WatchlistServiceImpl) LoadAllWatchlists(ctx context.Context) ([]models.MasterWatchlist, error) {
	var combined []models.MasterWatchlist

	loaders := []func(context.Context) ([]models.MasterWatchlist, error){
		w.LoadDTTOTIndividu,
		w.LoadWMDIndividu,
		w.LoadLocalBlacklistIndividu,
		w.LoadDTTOTCorporate,
		w.LoadWMDCorporate,
		w.LoadLocalBlacklistCorporate,
	}

	for _, loader := range loaders {
		list, err := loader(ctx)
		if err != nil {
			return nil, err
		}
		combined = append(combined, list...)
	}

	return combined, nil
}

// ==================== Mapper ====================
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
		var wl models.Watchlistable
		if v, ok := any(items[i]).(models.Watchlistable); ok { // value
			wl = v
		} else if v, ok := any(&items[i]).(models.Watchlistable); ok { // pointer
			wl = v
		}
		if wl != nil {
			result = append(result, ToWatchlist(wl))
		}
	}
	return result
}

// ==================== Loader Generic ====================
// func loadWatchlist[T any](ctx context.Context, loader func(context.Context) ([]T, error)) ([]models.MasterWatchlist, error) {
// 	list, err := loader(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ToWatchlistSlice(list), nil
// }
