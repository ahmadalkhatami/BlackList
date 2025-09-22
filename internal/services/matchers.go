package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/similarity"
	"BlackListWorker/internal/monitor"
	"BlackListWorker/internal/utils"
)

type GenericMatcher struct {
	source        string
	nasabahSvc    MasterNasabahService
	watchlistSvc  WatchlistService
	configSvc     SystemConfigService
	idSvc         *IDService
	loadIndividu  func(context.Context) ([]models.MasterWatchlist, error)
	loadCorporate func(context.Context) ([]models.MasterWatchlist, error)
}

func NewGenericMatcher(
	source string,
	nasabahSvc MasterNasabahService,
	watchlistSvc WatchlistService,
	configSvc SystemConfigService,
	idSvc *IDService,
	loadIndividu func(context.Context) ([]models.MasterWatchlist, error),
	loadCorporate func(context.Context) ([]models.MasterWatchlist, error),
) *GenericMatcher {
	return &GenericMatcher{
		source:        source,
		nasabahSvc:    nasabahSvc,
		watchlistSvc:  watchlistSvc,
		configSvc:     configSvc,
		idSvc:         idSvc,
		loadIndividu:  loadIndividu,
		loadCorporate: loadCorporate,
	}
}

func (m *GenericMatcher) NextResultID() int64 {
	return atomic.AddInt64(&m.idSvc.resultIDCounter, 1)
}

func (m *GenericMatcher) NextDetailID() int64 {
	return atomic.AddInt64(&m.idSvc.detailIDCounter, 1)
}

func (m *GenericMatcher) NextBatchID(batch *models.BatchProcessing) (int64, error) {
	return m.idSvc.GenerateBatchID(batch)
}

func (m *GenericMatcher) Source() string { return m.source }

type matcherIDWrapper struct {
	matcher *GenericMatcher
}

func (w *matcherIDWrapper) GenerateResultID() int64 {
	return w.matcher.NextResultID()
}

func (w *matcherIDWrapper) GenerateDetailID() int64 {
	return w.matcher.NextDetailID()
}

func (w *matcherIDWrapper) GenerateBatchID(batch *models.BatchProcessing) (int64, error) {
	return w.matcher.NextBatchID(batch)
}

func (m *GenericMatcher) Match(ctx context.Context) (*MatchResults, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("⏱️ %s Matcher selesai dalam %v\n", m.source, time.Since(start))
	}()
	return runAdaptiveMatch(ctx, m.source, m.nasabahSvc, m.configSvc,
		&matcherIDWrapper{matcher: m}, m.loadIndividu, m.loadCorporate)
}

func NewDTTOTMatcher(nasabah MasterNasabahService, wl WatchlistService, cfg SystemConfigService, idSvc *IDService) *GenericMatcher {
	return NewGenericMatcher("MASTER_TERORIS", nasabah, wl, cfg, idSvc, wl.LoadDTTOTIndividu, wl.LoadDTTOTCorporate)
}

func NewWMDMatcher(nasabah MasterNasabahService, wl WatchlistService, cfg SystemConfigService, idSvc *IDService) *GenericMatcher {
	return NewGenericMatcher("MASTER_WMD", nasabah, wl, cfg, idSvc, wl.LoadWMDIndividu, wl.LoadWMDCorporate)
}

func NewLocalBlacklistMatcher(nasabah MasterNasabahService, wl WatchlistService, cfg SystemConfigService, idSvc *IDService) *GenericMatcher {
	return NewGenericMatcher("MASTER_LOCAL_BLACKLIST", nasabah, wl, cfg, idSvc, wl.LoadLocalBlacklistIndividu, wl.LoadLocalBlacklistCorporate)
}

type MatchConfig struct {
	Threshold float64
	Algorithm string
	SimCalc   similarity.Calculator
	Individu  []models.JoinedMatchingConfig
	Corporate []models.JoinedMatchingConfig
}

func loadMatchConfig(ctx context.Context, cfgSvc SystemConfigService) (*MatchConfig, error) {
	threshold, err := cfgSvc.GetThreshold(ctx)
	if err != nil {
		return nil, err
	}
	algorithm, err := cfgSvc.GetMatchingAlgorithm(ctx)
	if err != nil {
		return nil, err
	}
	simCalc, err := similarity.NewCalculator(similarity.Algorithm(algorithm))
	if err != nil {
		return nil, err
	}
	cfgIndividu, err := cfgSvc.LoadIndividu(ctx)
	if err != nil {
		return nil, err
	}
	cfgCorporate, err := cfgSvc.LoadCorporate(ctx)
	if err != nil {
		return nil, err
	}
	return &MatchConfig{
		Threshold: threshold,
		Algorithm: algorithm,
		SimCalc:   simCalc,
		Individu:  cfgIndividu,
		Corporate: cfgCorporate,
	}, nil
}

func systemHealthy() bool {
	cpuCollector := monitor.NewCPUCollector(0)
	memCollector := monitor.NewMemoryCollector()
	threadCollector := monitor.NewThreadCollector()

	_, cpuVal := cpuCollector.Collect()
	_, memVal := memCollector.Collect()
	_, threads := threadCollector.Collect()

	cpu := cpuVal.(float64)
	mem := memVal.(float64)
	goroutines := threads.(int)

	if cpu > 80.0 {
		fmt.Printf("⚠️ CPU tinggi: %.2f%%\n", cpu)
		return false
	}
	if mem > 80.0 {
		fmt.Printf("⚠️ Memory tinggi: %.2f%%\n", mem)
		return false
	}
	if goroutines > 500 {
		fmt.Printf("⚠️ Goroutine terlalu banyak: %d\n", goroutines)
		return false
	}
	return true
}

func runAdaptiveMatch(
	ctx context.Context,
	source string,
	nasabahSvc MasterNasabahService,
	configSvc SystemConfigService,
	idGen interface {
		GenerateResultID() int64
		GenerateDetailID() int64
		GenerateBatchID(batch *models.BatchProcessing) (int64, error)
	},
	loadIndividu func(context.Context) ([]models.MasterWatchlist, error),
	loadCorporate func(context.Context) ([]models.MasterWatchlist, error),
) (*MatchResults, error) {
	if systemHealthy() {
		fmt.Println("🚀 Resource sehat → gunakan parallel matching")
		return runParallelMatch(ctx, source, nasabahSvc, configSvc, idGen, loadIndividu, loadCorporate)
	}
	fmt.Println("🐢 Resource terbatas → gunakan serial matching")
	return runSerialMatch(ctx, source, nasabahSvc, configSvc, idGen, loadIndividu, loadCorporate)
}

func runParallelMatch(
	ctx context.Context,
	source string,
	nasabahSvc MasterNasabahService,
	configSvc SystemConfigService,
	idGen interface {
		GenerateResultID() int64
		GenerateDetailID() int64
		GenerateBatchID(batch *models.BatchProcessing) (int64, error)
	},
	loadIndividu func(context.Context) ([]models.MasterWatchlist, error),
	loadCorporate func(context.Context) ([]models.MasterWatchlist, error),
) (*MatchResults, error) {

	debug := utils.IsDebugMode()
	var wg sync.WaitGroup
	resultsChan := make(chan *MatchResults, 2)
	errChan := make(chan error, 2)

	cfg, err := loadMatchConfig(ctx, configSvc)
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := matchMaster(ctx, source, nasabahSvc, cfg.Individu, idGen, loadIndividu, cfg.SimCalc, cfg.Threshold, cfg.Algorithm, debug)
		if err != nil {
			errChan <- err
			return
		}
		resultsChan <- res
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := matchMaster(ctx, source, nasabahSvc, cfg.Corporate, idGen, loadCorporate, cfg.SimCalc, cfg.Threshold, cfg.Algorithm, debug)
		if err != nil {
			errChan <- err
			return
		}
		resultsChan <- res
	}()

	wg.Wait()
	close(resultsChan)
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	var allResults MatchResults
	for r := range resultsChan {
		allResults.MatchResult = append(allResults.MatchResult, r.MatchResult...)
		allResults.MatchDetail = append(allResults.MatchDetail, r.MatchDetail...)
	}
	return &allResults, nil
}

func runSerialMatch(
	ctx context.Context,
	source string,
	nasabahSvc MasterNasabahService,
	configSvc SystemConfigService,
	idGen interface {
		GenerateResultID() int64
		GenerateDetailID() int64
		GenerateBatchID(batch *models.BatchProcessing) (int64, error)
	},
	loadIndividu func(context.Context) ([]models.MasterWatchlist, error),
	loadCorporate func(context.Context) ([]models.MasterWatchlist, error),
) (*MatchResults, error) {
	debug := utils.IsDebugMode()

	cfg, err := loadMatchConfig(ctx, configSvc)
	if err != nil {
		return nil, err
	}

	allResults := &MatchResults{}

	resInd, err := matchMaster(ctx, source, nasabahSvc, cfg.Individu, idGen, loadIndividu, cfg.SimCalc, cfg.Threshold, cfg.Algorithm, debug)
	if err != nil {
		return nil, err
	}
	allResults.MatchResult = append(allResults.MatchResult, resInd.MatchResult...)
	allResults.MatchDetail = append(allResults.MatchDetail, resInd.MatchDetail...)

	resCorp, err := matchMaster(ctx, source, nasabahSvc, cfg.Corporate, idGen, loadCorporate, cfg.SimCalc, cfg.Threshold, cfg.Algorithm, debug)
	if err != nil {
		return nil, err
	}
	allResults.MatchResult = append(allResults.MatchResult, resCorp.MatchResult...)
	allResults.MatchDetail = append(allResults.MatchDetail, resCorp.MatchDetail...)

	return allResults, nil
}

func matchMaster(
	ctx context.Context,
	source string,
	nasabahSvc MasterNasabahService,
	cfgList []models.JoinedMatchingConfig,
	idGen interface {
		GenerateResultID() int64
		GenerateDetailID() int64
		GenerateBatchID(batch *models.BatchProcessing) (int64, error)
	},
	loadWatchlist func(context.Context) ([]models.MasterWatchlist, error),
	simCalc similarity.Calculator,
	threshold float64,
	algorithm string,
	debug bool,
) (*MatchResults, error) {

	if nasabahSvc == nil {
		return nil, errors.New("missing nasabahSvc dependency")
	}

	cifList, err := nasabahSvc.Load(ctx)
	if err != nil {
		return nil, err
	}
	if debug {
		fmt.Printf("DEBUG: Loaded CIF total = %d\n", len(cifList))
	}

	watchlist, err := loadWatchlist(ctx)
	if err != nil {
		return nil, err
	}
	if debug {
		fmt.Printf("DEBUG: Loaded Watchlist total = %d\n", len(watchlist))
	}

	var matchResults []models.MatchingResult
	var matchDetails []models.MatchingDetail

	batch := &models.BatchProcessing{
		Id:          0,
		ProcessType: utils.Ptr(source),
		/* Status: running | completed | failed */
		Status:       utils.Ptr("running"),
		TotalRecords: utils.Ptr(0),
		InitiatedBy:  utils.TrigeredBy(),
	}

	batchID, err := idGen.GenerateBatchID(batch)
	if err != nil {
		fmt.Printf("⚠️ warning: gagal create batch in DB: %v. using fallback id %d\n", err, batchID)
	}

	for _, cif := range cifList {
		for _, wl := range watchlist {
			if !*wl.GetIsActive() || *wl.GetSource() != source {
				continue
			}

			totalScore := 0.0
			totalWeight := 0.0
			var fieldMatches []models.MatchingDetail

			for _, cfg := range cfgList {
				if !cfg.IsActive || cfg.WatchlistSource != source {
					continue
				}

				custVal := utils.GetCIFValueByField(cif, cfg.FieldName)
				wlValues := utils.GetWatchlistValuesByField(wl, cfg.FieldName)

				maxScore := 0.0
				bestMatch := ""
				for _, wVal := range wlValues {
					score := simCalc.Calculate(custVal, wVal)
					if score > maxScore {
						maxScore = score
						bestMatch = wVal
					}
				}

				totalScore += maxScore * cfg.FieldWeight
				totalWeight += cfg.FieldWeight

				fieldMatches = append(fieldMatches, models.MatchingDetail{
					FieldName:      utils.Ptr(cfg.FieldName),
					CustomerValue:  utils.Ptr(custVal),
					WatchlistValue: utils.Ptr(bestMatch),
					FieldScore:     utils.Ptr(maxScore),
					FieldWeight:    utils.Ptr(cfg.FieldWeight),
					AlgorithmUsed:  utils.Ptr(algorithm),
				})
			}

			if totalWeight == 0 {
				continue
			}

			finalScore := totalScore / totalWeight

			if debug && finalScore < threshold {
				fmt.Printf(`\n| finalScore:%.2f, threshold:%.2f |\n`, finalScore, threshold)
			}

			result := models.MatchingResult{
				Id:              idGen.GenerateResultID(),
				BatchId:         utils.Ptr(batchID),
				CifNumber:       utils.Ptr(cif.CifNumber),
				CustomerName:    utils.Ptr(cif.NamaNasabah),
				WatchlistId:     utils.Ptr(wl.ID),
				WatchlistSource: utils.Ptr(source),
				SimilarityScore: utils.Ptr(finalScore),
				Status:          utils.Ptr("SUCCESS"),
				ProcessDate:     utils.Ptr(time.Now()),
				ProcessTime:     utils.Ptr(time.Now()),
				CreatedAt:       utils.Ptr(time.Now()),
			}

			for i := range fieldMatches {
				fieldMatches[i].MatchingResultId = utils.Ptr(result.Id)
				fieldMatches[i].Id = idGen.GenerateDetailID()
				matchDetails = append(matchDetails, fieldMatches[i])
			}

			matchResults = append(matchResults, result)
		}
	}

	return &MatchResults{
		MatchResult: matchResults,
		MatchDetail: matchDetails,
		BatchID:     batchID,
	}, nil
}
