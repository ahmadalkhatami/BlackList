package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/similarity"
	"BlackListWorker/internal/utils"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

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
	var wg sync.WaitGroup
	resultsChan := make(chan *MatchResults, 2)
	errChan := make(chan error, 2)

	// Individu
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := matchMaster(ctx, m.Source(), m.nasabahSvc, m.watchlistSvc, m.configSvc, m.watchlistSvc.LoadDTTOTIndividu, false)
		if err != nil {
			errChan <- err
			return
		}
		resultsChan <- res
	}()

	// Corporate
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := matchMaster(ctx, m.Source(), m.nasabahSvc, m.watchlistSvc, m.configSvc, m.watchlistSvc.LoadDTTOTCorporate, false)
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
	var wg sync.WaitGroup
	resultsChan := make(chan *MatchResults, 2)
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := matchMaster(ctx, m.Source(), m.nasabahSvc, m.watchlistSvc, m.configSvc, m.watchlistSvc.LoadWMDIndividu, false)
		if err != nil {
			errChan <- err
			return
		}
		resultsChan <- res
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := matchMaster(ctx, m.Source(), m.nasabahSvc, m.watchlistSvc, m.configSvc, m.watchlistSvc.LoadWMDCorporate, false)
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
	var wg sync.WaitGroup
	resultsChan := make(chan *MatchResults, 2)
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := matchMaster(ctx, m.Source(), m.nasabahSvc, m.watchlistSvc, m.configSvc, m.watchlistSvc.LoadLocalBlacklistIndividu, false)
		if err != nil {
			errChan <- err
			return
		}
		resultsChan <- res
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := matchMaster(ctx, m.Source(), m.nasabahSvc, m.watchlistSvc, m.configSvc, m.watchlistSvc.LoadLocalBlacklistCorporate, false)
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

// ------------------ GENERIC MATCH FUNCTION ------------------
func matchMaster(
	ctx context.Context,
	source string,
	nasabahSvc MasterNasabahService,
	watchlistSvc WatchlistService,
	configSvc SystemConfigService,
	loadWatchlist func(context.Context) ([]models.MasterWatchlist, error),
	debug bool,
) (*MatchResults, error) {

	if nasabahSvc == nil || watchlistSvc == nil || configSvc == nil {
		return nil, errors.New("missing dependencies")
	}

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

	cifList, err := nasabahSvc.Load(ctx)
	if err != nil {
		return nil, err
	}

	cfgIndividu, err := configSvc.LoadIndividu(ctx)
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

			for _, cfg := range cfgIndividu {
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
					FieldName:        cfg.FieldName,
					CustomerValue:    custVal,
					WatchlistValue:   bestMatch,
					FieldScore:       maxScore,
					FieldWeight:      cfg.FieldWeight,
					AlgorithmUsed:    algorithm,
					MatchingResultID: wl.ID,
				})

				if debug {
					fmt.Printf("   🔍 Field: %s | CIF: '%s' | WL: '%s' | Score: %.2f | Weight: %.2f\n",
						cfg.FieldName, custVal, bestMatch, maxScore, cfg.FieldWeight)
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

			// if finalScore >= threshold {
			// 	result := models.MatchingResult{
			// 		CIFNumber:       cif.CIFNumber,
			// 		CustomerName:    cif.NamaNasabah,
			// 		WatchlistID:     wl.ID,
			// 		WatchlistSource: source,
			// 		SimilarityScore: finalScore,
			// 		Status:          "SUCCESS",
			// 		ProcessDate:     time.Now(),
			// 		ProcessTime:     time.Now(),
			// 		CreatedAt:       time.Now(),
			// 	}

			// 	matchResults = append(matchResults, result)

			// 	for i := range fieldMatches {
			// 		fieldMatches[i].MatchingResultID = result.ID
			// 		matchDetails = append(matchDetails, fieldMatches[i])
			// 	}

			// 	if debug {
			// 		fmt.Printf("✅ MATCH: CIF=%s vs Watchlist=%s Score=%.2f\n",
			// 			cif.NamaNasabah, wl.Nama, finalScore)
			// 	}
			// } else if debug {
			// 	fmt.Printf("❌ NO MATCH: CIF=%s vs Watchlist=%s Score=%.2f < Threshold %.2f\n",
			// 		cif.NamaNasabah, wl.Nama, finalScore, threshold)
			// }
		}
	}

	return &MatchResults{
		MatchResult: matchResults,
		MatchDetail: matchDetails,
	}, nil
}
