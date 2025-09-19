package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/similarity"
	"BlackListWorker/internal/utils"
)

// ===== DTTOT Matcher =====
type DTTOTMatcher struct {
	nasabahSvc   MasterNasabahService
	watchlistSvc WatchlistService
	configSvc    SystemConfigService
}

func NewDTTOTMatcher(nasabah MasterNasabahService, wl WatchlistService, cfg SystemConfigService) *DTTOTMatcher {
	return &DTTOTMatcher{nasabahSvc: nasabah, watchlistSvc: wl, configSvc: cfg}
}

func (m *DTTOTMatcher) Source() string { return "MASTER_TERORIS" }

func (m *DTTOTMatcher) Match(ctx context.Context) (*MatchResults, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("⏱️ %sMatcher selesai dalam %v\n", m.Source(), time.Since(start))
	}()

	return runParallelMatch(ctx, m.Source(), m.nasabahSvc, m.watchlistSvc,
		m.configSvc, m.watchlistSvc.LoadDTTOTIndividu, m.watchlistSvc.LoadDTTOTCorporate)
}

// ===== WMD Matcher =====
type WMDMatcher struct {
	nasabahSvc   MasterNasabahService
	watchlistSvc WatchlistService
	configSvc    SystemConfigService
}

func NewWMDMatcher(nasabah MasterNasabahService, wl WatchlistService, cfg SystemConfigService) *WMDMatcher {
	return &WMDMatcher{nasabahSvc: nasabah, watchlistSvc: wl, configSvc: cfg}
}

func (m *WMDMatcher) Source() string { return "MASTER_WMD" }

func (m *WMDMatcher) Match(ctx context.Context) (*MatchResults, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("⏱️ %sMatcher selesai dalam %v\n", m.Source(), time.Since(start))
	}()

	return runParallelMatch(ctx, m.Source(), m.nasabahSvc, m.watchlistSvc,
		m.configSvc, m.watchlistSvc.LoadWMDIndividu, m.watchlistSvc.LoadWMDCorporate)
}

// ===== Local Blacklist Matcher =====
type LocalBlacklistMatcher struct {
	nasabahSvc   MasterNasabahService
	watchlistSvc WatchlistService
	configSvc    SystemConfigService
}

func NewLocalBlacklistMatcher(nasabah MasterNasabahService, wl WatchlistService, cfg SystemConfigService) *LocalBlacklistMatcher {
	return &LocalBlacklistMatcher{nasabahSvc: nasabah, watchlistSvc: wl, configSvc: cfg}
}

func (m *LocalBlacklistMatcher) Source() string { return "MASTER_LOCAL_BLACKLIST" }

func (m *LocalBlacklistMatcher) Match(ctx context.Context) (*MatchResults, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("⏱️ %sMatcher selesai dalam %v\n", m.Source(), time.Since(start))
	}()

	return runParallelMatch(ctx, m.Source(), m.nasabahSvc, m.watchlistSvc,
		m.configSvc, m.watchlistSvc.LoadLocalBlacklistIndividu, m.watchlistSvc.LoadLocalBlacklistCorporate)
}

// ===== Helper: Parallel run for Individu & Corporate =====
func runParallelMatch(
	ctx context.Context,
	source string,
	nasabahSvc MasterNasabahService,
	watchlistSvc WatchlistService,
	configSvc SystemConfigService,
	loadIndividu func(context.Context) ([]models.MasterWatchlist, error),
	loadCorporate func(context.Context) ([]models.MasterWatchlist, error),
) (*MatchResults, error) {

	debug := utils.IsDebugMode()

	var wg sync.WaitGroup
	resultsChan := make(chan *MatchResults, 2)
	errChan := make(chan error, 2)

	// threshold & algorithm
	threshold, err := configSvc.GetThreshold(ctx)
	if err != nil {
		return nil, err
	}
	algorithm, err := configSvc.GetMatchingAlgorithm(ctx)
	if err != nil {
		return nil, err
	}
	simCalc, err := similarity.NewCalculator(similarity.Algorithm(algorithm))
	if err != nil {
		return nil, err
	}

	// config individu & corporate
	cfgIndividu, err := configSvc.LoadIndividu(ctx)
	if err != nil {
		return nil, err
	}
	cfgCorporate, err := configSvc.LoadCorporate(ctx)
	if err != nil {
		return nil, err
	}

	// run individu
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := matchMaster(ctx, source, nasabahSvc, cfgIndividu, loadIndividu, simCalc, threshold, algorithm, debug)
		if err != nil {
			errChan <- err
			return
		}
		resultsChan <- res
	}()

	// run corporate
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := matchMaster(ctx, source, nasabahSvc, cfgCorporate, loadCorporate, simCalc, threshold, algorithm, debug)
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

// ===== Core Matching Logic =====
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
	watchlist, err := loadWatchlist(ctx)
	if err != nil {
		return nil, err
	}

	var matchResults []models.MatchingResult
	var matchDetails []models.MatchingDetail

	for _, cif := range cifList {
		for _, wl := range watchlist {
			if !wl.IsActive || wl.Source != source {
				continue
			}

			if debug {
				fmt.Printf("\n🚀 Processing CIF: %s | Watchlist: %s (Source: %s)\n",
					cif.NamaNasabah, wl.Nama, wl.Source)
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

					if debug {
						fmt.Printf("      🔍 Compare: '%s' vs '%s' => Score=%.4f\n", custVal, wVal, score)
					}
				}

				totalScore += maxScore * cfg.FieldWeight
				totalWeight += cfg.FieldWeight

				fieldMatches = append(fieldMatches, models.MatchingDetail{
					FieldName:        cfg.FieldName,
					CustomerValue:    custVal,
					WatchlistValue:   bestMatch,
					FieldScore:       maxScore,
					FieldWeight:      cfg.FieldWeight,
					AlgorithmUsed:    algorithm,
					MatchingResultID: wl.ID,
				})

				// if debug {
				// 	fmt.Printf("   🔍 Field: %s | CIF: '%s' | WL: '%s' | Score: %.2f | Weight: %.2f\n",
				// 		cfg.FieldName, custVal, bestMatch, maxScore, cfg.FieldWeight)
				// }
				if debug {
					fmt.Printf("   ✅ BestMatch=%q | MaxScore=%.4f (Weighted=%.4f)\n",
						bestMatch, maxScore, maxScore*cfg.FieldWeight)
				}
			}

			if totalWeight == 0 {
				if debug {
					fmt.Println("⚠️ Skip: totalWeight = 0 (no active config)")
				}
				continue
			}

			finalScore := totalScore / totalWeight

			if debug {
				fmt.Printf("➡️ FinalScore: %.2f / totalWeight: %.2f (Threshold: %.2f)\n",
					totalScore, totalWeight, threshold)
			}

			result := models.MatchingResult{
				CIFNumber:       cif.CIFNumber,
				CustomerName:    cif.NamaNasabah,
				WatchlistID:     wl.ID,
				WatchlistSource: source,
				SimilarityScore: finalScore,
				Status:          "SUCCESS",
				ProcessDate:     time.Now(),
				ProcessTime:     time.Now(),
				CreatedAt:       time.Now(),
			}

			matchResults = append(matchResults, result)
			for i := range fieldMatches {
				fieldMatches[i].MatchingResultID = result.ID
				matchDetails = append(matchDetails, fieldMatches[i])
			}

			if debug {
				fmt.Printf("✅ MATCH: CIF=%s vs Watchlist=%s Score=%.2f\n",
					cif.NamaNasabah, wl.Nama, finalScore)
			}
		}
	}

	return &MatchResults{
		MatchResult: matchResults,
		MatchDetail: matchDetails,
	}, nil
}
