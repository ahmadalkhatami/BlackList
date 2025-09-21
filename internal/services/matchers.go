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
	loadIndividu  func(context.Context) ([]models.MasterWatchlist, error)
	loadCorporate func(context.Context) ([]models.MasterWatchlist, error)
}

func NewGenericMatcher(
	source string,
	nasabahSvc MasterNasabahService,
	watchlistSvc WatchlistService,
	configSvc SystemConfigService,
	loadIndividu func(context.Context) ([]models.MasterWatchlist, error),
	loadCorporate func(context.Context) ([]models.MasterWatchlist, error),
) *GenericMatcher {
	return &GenericMatcher{
		source: source, nasabahSvc: nasabahSvc,
		watchlistSvc: watchlistSvc, configSvc: configSvc,
		loadIndividu: loadIndividu, loadCorporate: loadCorporate,
	}
}

func (m *GenericMatcher) Source() string { return m.source }

func (m *GenericMatcher) Match(ctx context.Context) (*MatchResults, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("⏱️ %sMatcher selesai dalam %v\n", m.Source(), time.Since(start))
	}()
	return runAdaptiveMatch(ctx, m.Source(), m.nasabahSvc, m.configSvc, m.loadIndividu, m.loadCorporate)
}

func NewDTTOTMatcher(nasabah MasterNasabahService, wl WatchlistService, cfg SystemConfigService) *GenericMatcher {
	return NewGenericMatcher("MASTER_TERORIS", nasabah, wl, cfg, wl.LoadDTTOTIndividu, wl.LoadDTTOTCorporate)
}

func NewWMDMatcher(nasabah MasterNasabahService, wl WatchlistService, cfg SystemConfigService) *GenericMatcher {
	return NewGenericMatcher("MASTER_WMD", nasabah, wl, cfg, wl.LoadWMDIndividu, wl.LoadWMDCorporate)
}

func NewLocalBlacklistMatcher(nasabah MasterNasabahService, wl WatchlistService, cfg SystemConfigService) *GenericMatcher {
	return NewGenericMatcher("MASTER_LOCAL_BLACKLIST", nasabah, wl, cfg, wl.LoadLocalBlacklistIndividu, wl.LoadLocalBlacklistCorporate)
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
	cpuCollector := monitor.NewCPUCollector(0) // interval = 0 → snapshot instan
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
	loadIndividu func(context.Context) ([]models.MasterWatchlist, error),
	loadCorporate func(context.Context) ([]models.MasterWatchlist, error),
) (*MatchResults, error) {

	if systemHealthy() {
		fmt.Println("🚀 Resource sehat → gunakan parallel matching")
		return runParallelMatch(ctx, source, nasabahSvc, configSvc, loadIndividu, loadCorporate)
	}

	fmt.Println("🐢 Resource terbatas → gunakan serial matching")
	return runSerialMatch(ctx, source, nasabahSvc, configSvc, loadIndividu, loadCorporate)
}

var resultIDCounter int64 = 0
var detailIDCounter int64 = 0

func generateResultID() int64 {
	return atomic.AddInt64(&resultIDCounter, 1)
}

func generateDetailID() int64 {
	return atomic.AddInt64(&detailIDCounter, 1)
}

func runParallelMatch(
	ctx context.Context,
	source string,
	nasabahSvc MasterNasabahService,
	configSvc SystemConfigService,
	loadIndividu func(context.Context) ([]models.MasterWatchlist, error),
	loadCorporate func(context.Context) ([]models.MasterWatchlist, error),
	// debug bool,
) (*MatchResults, error) {

	debug := utils.IsDebugMode()
	// fmt.Printf("Debug Mode : %t\n", debug)

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
		res, err := matchMaster(ctx, source, nasabahSvc, cfg.Individu, loadIndividu, cfg.SimCalc, cfg.Threshold, cfg.Algorithm, debug)
		if err != nil {
			errChan <- err
			return
		}
		resultsChan <- res
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := matchMaster(ctx, source, nasabahSvc, cfg.Corporate, loadCorporate, cfg.SimCalc, cfg.Threshold, cfg.Algorithm, debug)
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
	loadIndividu func(context.Context) ([]models.MasterWatchlist, error),
	loadCorporate func(context.Context) ([]models.MasterWatchlist, error),
) (*MatchResults, error) {

	debug := utils.IsDebugMode()

	cfg, err := loadMatchConfig(ctx, configSvc)
	if err != nil {
		return nil, err
	}

	allResults := &MatchResults{}

	resInd, err := matchMaster(ctx, source, nasabahSvc, cfg.Individu, loadIndividu, cfg.SimCalc, cfg.Threshold, cfg.Algorithm, debug)
	if err != nil {
		return nil, err
	}
	allResults.MatchResult = append(allResults.MatchResult, resInd.MatchResult...)
	allResults.MatchDetail = append(allResults.MatchDetail, resInd.MatchDetail...)

	resCorp, err := matchMaster(ctx, source, nasabahSvc, cfg.Corporate, loadCorporate, cfg.SimCalc, cfg.Threshold, cfg.Algorithm, debug)
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

				// fieldLower := strings.ToLower(cfg.FieldName)
				// var custVal string
				// if strings.HasPrefix(fieldLower, "alias") {
				// 	custVal = cif.NamaNasabah
				// } else {
				// 	custVal = utils.GetCIFValueByField(cif, cfg.FieldName)
				// }

				// var wlValues []string
				// switch fieldLower {
				// case "nama", "namanasabah":
				// 	wlValues = append([]string{wl.Nama}, wl.Aliases...)
				// default:
				// 	wlValues = utils.GetWatchlistValuesByField(wl, cfg.FieldName)
				// }

				custVal := utils.GetCIFValueByField(cif, cfg.FieldName)
				wlValues := utils.GetWatchlistValuesByField(wl, cfg.FieldName)

				if debug {
					fmt.Printf("\n🔧 DEBUG Field=%s\n", cfg.FieldName)
					fmt.Printf("   ↳ CustomerValue = %q\n", custVal)
					fmt.Printf("   ↳ WatchlistValues = %v\n", wlValues)
					fmt.Printf("   ↳ FieldWeight = %.2f | Active=%v\n", cfg.FieldWeight, cfg.IsActive)
				}

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

				if debug {
					fmt.Printf("   🔍 Watchlist Id: %d | Field: %s | CIF: '%s' | WL: '%s' | Score: %.2f | Weight: %.2f\n",
						wl.ID, cfg.FieldName, custVal, bestMatch, maxScore, cfg.FieldWeight)
				}
			}

			if totalWeight == 0 {
				continue
			}

			finalScore := totalScore / totalWeight

			result := models.MatchingResult{
				Id:              generateResultID(),
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
				fieldMatches[i].Id = generateDetailID()
				matchDetails = append(matchDetails, fieldMatches[i])
			}

			matchResults = append(matchResults, result)

			if debug {
				fmt.Printf("✅ MATCH: CIF=%s vs Watchlist=%s Score=%.2f\n | Threshold : %d",
					cif.NamaNasabah, wl.Nama, finalScore, threshold)
			}
		}
	}

	return &MatchResults{
		MatchResult: matchResults,
		MatchDetail: matchDetails,
	}, nil
}
