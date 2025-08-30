package services

import (
	dbcon "BlackListWorker/internal/db"
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repositories"
	"BlackListWorker/internal/utils"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/agext/levenshtein"
)

const threshold = 0.85

type MatchingServiceInterface interface {
	MatchCIFWithTeroris() ([]models.MatchingResult, map[int64][]models.MatchingDetail, error)
}

func MatchCIFWithTeroris() ([]models.MatchingResult, map[int64][]models.MatchingDetail, error) {

	connector := dbcon.GetConnector()
	sqlDB, err := connector.Connect()
	if err != nil {
		return []models.MatchingResult{}, map[int64][]models.MatchingDetail{}, err
	}
	defer sqlDB.Close()

	var results []models.MatchingResult
	detailsMap := make(map[int64][]models.MatchingDetail)
	fieldConfig := make(map[string]models.JoinedMatchingConfig)

	masterMatchingRepo := repositories.NewSQLMasterMatchingRepository(sqlDB)
	masterMatchingConfigRepo := repositories.NewSQLMasterMatchingConfigRepository(sqlDB)
	systemConfigRepo := repositories.NewSQLSystemConfigRepository(sqlDB)

	configService := NewConfigService(
		WithMasterMatching(masterMatchingRepo),
		WithMasterMatchingConfig(masterMatchingConfigRepo),
		WithSystemConfig(systemConfigRepo),
	)

	masterNasabahRepo := repositories.NewSQLMasterNasabahRepository(sqlDB)
	masterNasabahService := NewMasterNasabah(masterNasabahRepo)

	masterDTTOTRepo := repositories.NewSQLMasterTerorisRepository(sqlDB)
	masterWMDRepo := repositories.NewSQLMasterWMDRepository(sqlDB)
	masterLocalBalcklistRepo := repositories.NewSQLMasterLocalBlacklistRepository(sqlDB)

	watchlistService := NewWatchlistService(
		WithMasterTeroris(masterDTTOTRepo),
		WithMasterWMD(masterWMDRepo),
		WithMasterLocalBlacklist(masterLocalBalcklistRepo),
	)

	allWatchlists, err := watchlistService.LoadAllWatchlists()
	if err != nil {
		panic(err)
	}

	fmt.Println("Total Watchlists:", len(allWatchlists))

	// Ambil threshold dari SYSTEM_CONFIG
	threshold, err := configService.GetThresholdFromConfig("MATCHING_THRESHOLD")
	if err != nil {
		fmt.Printf("❌ Error getting threshold: %v\n", err)
		return []models.MatchingResult{}, map[int64][]models.MatchingDetail{}, err
	}

	configs, err := configService.GetJoinedMatchingConfig()
	if err != nil {
		fmt.Printf("❌ Error getting matching config: %v\n", err)
		return []models.MatchingResult{}, map[int64][]models.MatchingDetail{}, err
	}

	// Ambil config untuk MASTER_TERORIS saja
	for _, cfg := range configs {
		if strings.ToUpper(cfg.WatchlistSource) == "MASTER_TERORIS" && cfg.IsActive {
			fieldConfig[strings.ToLower(cfg.FieldName)] = cfg
		}
	}

	// fmt.Printf("📌 Field config untuk MASTER_TERORIS: %+v\n", fieldConfig)
	cifs, err := masterNasabahService.Load()
	if err != nil {
		return []models.MatchingResult{}, map[int64][]models.MatchingDetail{}, err
	}

	terorisList, err := watchlistService.LoadDTTOT()
	if err != nil {
		return []models.MatchingResult{}, map[int64][]models.MatchingDetail{}, err
	}

	cifKeys := utils.GetStructKeys(&models.MasterNasabah{})
	terorisKeys := utils.GetStructKeys(&models.MasterTeroris{})
	fmt.Printf("📌 Field di MasterNasabah: %+v\n", cifKeys)
	fmt.Printf("📌 Field di MasterTeroris: %+v\n", terorisKeys)
	fmt.Printf("📌 Field di Config: %+v\n", configs)
	fmt.Printf("📌 Field di FieldConfig: %+v\n", fieldConfig)

	return nil, nil, nil

	for _, cif := range cifs {
		for _, wl := range terorisList {
			if !wl.IsActive {
				continue
			}

			// fmt.Printf("\n🚀 Proses CIF: %s | DTTOT: %s",
			// 	cif.NamaNasabah, wl.Nama)

			totalWeight := 0.0
			totalScore := 0.0
			var matchedFields []models.MatchingDetail

			// for field, cfg := range fieldConfig {

			// 	/*
			// 		custVal := getCIFValueByField(cif, field)
			// 		watchlistValues := getWatchlistValuesByField(wl, field)

			// 		maxScore := 0.0
			// 		bestMatchVal := ""
			// 		for _, wlVal := range watchlistValues {
			// 			score := matchScore(custVal, wlVal, cfg.MatchingAlgorithm)
			// 			if score > maxScore {
			// 				maxScore = score
			// 				bestMatchVal = wlVal
			// 			}
			// 		}

			// 		// Debug per field
			// 		fmt.Printf("   🔍 Field: %s | CIF: '%s' | Watchlist: '%s' | Score: %.2f | Algoritma: %s\n",
			// 			field, custVal, bestMatchVal, maxScore, cfg.MatchingAlgorithm)

			// 		totalScore += maxScore * cfg.FieldWeight
			// 		totalWeight += cfg.FieldWeight

			// 		matchedFields = append(matchedFields, models.MatchingDetail{
			// 			FieldName:      field,
			// 			CustomerValue:  custVal,
			// 			WatchlistValue: bestMatchVal,
			// 			FieldScore:     maxScore,
			// 			FieldWeight:    cfg.FieldWeight,
			// 			AlgorithmUsed:  cfg.MatchingAlgorithm,
			// 		}) */

			// 	fmt.Printf("   🔍 Field: %s | Algoritma: %s | Bobot: %.2f\n", field, cfg.MatchingAlgorithm, cfg.FieldWeight)

			// 	matchedFields = append(matchedFields, models.MatchingDetail{
			// 		FieldName:      "Not Implemented",
			// 		CustomerValue:  "Not Implemented",
			// 		WatchlistValue: "Not Implemented",
			// 		FieldScore:     0,
			// 		FieldWeight:    0.0,
			// 		AlgorithmUsed:  "Not Implemented",
			// 	})
			// }

			if totalWeight == 0 {
				fmt.Println("⚠️ Skip: totalWeight = 0 (tidak ada config aktif)")
				continue
			}

			finalScore := totalScore / totalWeight
			// fmt.Printf("➡️ FinalScore CIF %s vs Watchlist %s = %.2f (Threshold %.2f)\n",
			// 	cif.NamaNasabah, wl.Nama, finalScore, threshold)

			if finalScore >= threshold {
				result := models.MatchingResult{
					CIFNumber:    cif.CIFNumber,
					CustomerName: cif.NamaNasabah,
					// WatchlistID:     wl.ID,
					WatchlistID:     0,
					WatchlistSource: "MASTER_TERORIS",
					SimilarityScore: finalScore,
					Status:          "SUCCESS",
					ProcessDate:     time.Now(),
					ProcessTime:     time.Now(),
					CreatedAt:       time.Now(),
				}

				resultIndex := int64(len(results)) // index sebelum append
				results = append(results, result)
				detailsMap[resultIndex] = matchedFields

				fmt.Printf("✅ MATCH ditemukan! CIF %s cocok dengan Watchlist %s (Score: %.2f)\n",
					cif.NamaNasabah, wl.Nama, finalScore)
			} else {
				fmt.Printf("❌ Tidak match (FinalScore %.2f < Threshold %.2f)\n", finalScore, threshold)
			}
		}
	}

	return results, detailsMap, nil
}

func precomputeCIFValues(cifs []models.MasterNasabah, fieldConfig map[string]models.MatchingConfig) []map[string]string {
	result := make([]map[string]string, len(cifs))
	for i, cif := range cifs {
		m := make(map[string]string)
		for field := range fieldConfig {
			m[field] = getCIFValueByField(cif, field)
		}
		result[i] = m
	}
	return result
}

func precomputeWMDValues(wmdList []models.MasterWatchlist, fieldConfig map[string]models.MatchingConfig) []map[string][]string {
	result := make([]map[string][]string, len(wmdList))
	for i, wl := range wmdList {
		m := make(map[string][]string)
		for field := range fieldConfig {
			m[field] = getWatchlistValuesByField(wl, field)
		}
		result[i] = m
	}
	return result
}

func MatchCIFWithWMD(
	db *sql.DB,
	cifs []models.MasterNasabah,
	wmdList []models.MasterWatchlist,
	configs []models.MatchingConfig,
) ([]models.MatchingResult, map[int64][]models.MatchingDetail) {

	var (
		results     []models.MatchingResult
		detailsMap  = make(map[int64][]models.MatchingDetail)
		fieldConfig = make(map[string]models.MatchingConfig)
		mu          sync.Mutex
		wg          sync.WaitGroup
	)

	// threshold := GetThresholdFromConfig("MATCHING_THRESHOLD")

	// 1. Build config map
	for _, cfg := range configs {
		if strings.ToUpper(cfg.WatchlistSource) == "MASTER_WMD" && cfg.IsActive {
			fieldConfig[strings.ToLower(cfg.FieldName)] = cfg
		}
	}

	// 2. Filter aktif WMD
	activeWMD := make([]models.MasterWatchlist, 0, len(wmdList))
	for _, wl := range wmdList {
		if wl.IsActive && wl.Source == "MASTER_WMD" {
			activeWMD = append(activeWMD, wl)
		}
	}

	// 3. Precompute values
	cifValues := precomputeCIFValues(cifs, fieldConfig)
	wmdValues := precomputeWMDValues(activeWMD, fieldConfig)

	// 4. Parallel processing per CIF
	for i, cif := range cifs {
		wg.Add(1)
		go func(i int, cif models.MasterNasabah) {
			defer wg.Done()

			var localResults []models.MatchingResult
			localDetails := make(map[int64][]models.MatchingDetail)

			for j, wl := range activeWMD {
				totalWeight := 0.0
				totalScore := 0.0
				var matchedFields []models.MatchingDetail

				for field, cfg := range fieldConfig {
					custVal := cifValues[i][field]
					watchlistVals := wmdValues[j][field]

					maxScore := 0.0
					bestMatch := ""
					for _, wlVal := range watchlistVals {
						score := matchScore(custVal, wlVal, cfg.MatchingAlgorithm)
						if score > maxScore {
							maxScore = score
							bestMatch = wlVal
						}
					}

					totalScore += maxScore * cfg.FieldWeight
					totalWeight += cfg.FieldWeight

					matchedFields = append(matchedFields, models.MatchingDetail{
						FieldName:      field,
						CustomerValue:  custVal,
						WatchlistValue: bestMatch,
						FieldScore:     maxScore,
						FieldWeight:    cfg.FieldWeight,
						AlgorithmUsed:  cfg.MatchingAlgorithm,
					})
				}

				if totalWeight == 0 {
					continue
				}

				finalScore := totalScore / totalWeight
				if finalScore >= threshold {
					result := models.MatchingResult{
						CIFNumber:       cif.CIFNumber,
						CustomerName:    cif.NamaNasabah,
						WatchlistID:     wl.ID,
						WatchlistSource: "MASTER_WMD",
						SimilarityScore: finalScore,
						Status:          "SUCCESS",
						ProcessDate:     time.Now(),
						ProcessTime:     time.Now(),
						CreatedAt:       time.Now(),
					}
					idx := int64(len(localResults) + len(results)) // sementara
					localResults = append(localResults, result)
					localDetails[idx] = matchedFields
				}
			}

			// Gabungkan hasil ke shared slice/map
			mu.Lock()
			baseIdx := int64(len(results))
			results = append(results, localResults...)
			for k, v := range localDetails {
				detailsMap[baseIdx+k] = v
			}
			mu.Unlock()
		}(i, cif)
	}

	wg.Wait()
	return results, detailsMap
}

// ===== MATCHING UNTUK MASTER_WMD =====
/*
func MatchCIFWithWMD(
	db *sql.DB,
	cifs []models.MasterNasabah,
	wmdList []models.MasterWatchlist,
	configs []models.MatchingConfig,
) ([]models.MatchingResult, map[int64][]models.MatchingDetail) {

	var results []models.MatchingResult
	detailsMap := make(map[int64][]models.MatchingDetail)
	fieldConfig := map[string]models.MatchingConfig{}

	// threshold := GetThresholdFromConfig("MATCHING_THRESHOLD")

	for _, cfg := range configs {
		if strings.ToUpper(cfg.WatchlistSource) == "MASTER_WMD" && cfg.IsActive {
			fieldConfig[strings.ToLower(cfg.FieldName)] = cfg
		}
	}

	for _, cif := range cifs {
		for _, wl := range wmdList {
			if !wl.IsActive || wl.Source != "MASTER_WMD" {
				continue
			}

			totalWeight := 0.0
			totalScore := 0.0
			var matchedFields []models.MatchingDetail

			for field, cfg := range fieldConfig {
				custVal := getCIFValueByField(cif, field)
				watchlistValues := getWatchlistValuesByField(wl, field)

				maxScore := 0.0
				bestMatchVal := ""
				for _, wlVal := range watchlistValues {
					score := matchScore(custVal, wlVal, cfg.MatchingAlgorithm)
					if score > maxScore {
						maxScore = score
						bestMatchVal = wlVal
					}
				}

				totalScore += maxScore * cfg.FieldWeight
				totalWeight += cfg.FieldWeight

				matchedFields = append(matchedFields, models.MatchingDetail{
					FieldName:      field,
					CustomerValue:  custVal,
					WatchlistValue: bestMatchVal,
					FieldScore:     maxScore,
					FieldWeight:    cfg.FieldWeight,
					AlgorithmUsed:  cfg.MatchingAlgorithm,
				})
			}

			if totalWeight == 0 {
				continue
			}

			finalScore := totalScore / totalWeight

			if finalScore >= threshold {
				result := models.MatchingResult{
					CIFNumber:       cif.CIFNumber,
					CustomerName:    cif.NamaNasabah,
					WatchlistID:     wl.ID,
					WatchlistSource: "MASTER_WMD",
					SimilarityScore: finalScore,
					Status:          "SUCCESS",
					ProcessDate:     time.Now(),
					ProcessTime:     time.Now(),
					CreatedAt:       time.Now(),
				}

				resultIndex := int64(len(results))
				results = append(results, result)
				detailsMap[resultIndex] = matchedFields
			}
		}
	}
	return results, detailsMap
}*/

// ===== MATCHING UNTUK MASTER_LOCAL_BLACKLIST =====
func MatchCIFWithLocalBlacklist(
	db *sql.DB,
	cifs []models.MasterNasabah,
	localList []models.MasterWatchlist,
	configs []models.MatchingConfig,
) ([]models.MatchingResult, map[int64][]models.MatchingDetail) {

	var results []models.MatchingResult
	detailsMap := make(map[int64][]models.MatchingDetail)
	fieldConfig := map[string]models.MatchingConfig{}

	// Ambil threshold
	// threshold := GetThresholdFromConfig("MATCHING_THRESHOLD")

	// Ambil config utk MASTER_LOCAL_BLACKLIST
	for _, cfg := range configs {
		if strings.ToUpper(cfg.WatchlistSource) == "MASTER_LOCAL_BLACKLIST" && cfg.IsActive {
			fieldConfig[strings.ToLower(cfg.FieldName)] = cfg
		}
	}

	for _, cif := range cifs {
		for _, wl := range localList {
			if !wl.IsActive || wl.Source != "MASTER_LOCAL_BLACKLIST" {
				continue
			}

			totalWeight := 0.0
			totalScore := 0.0
			var matchedFields []models.MatchingDetail

			for field, cfg := range fieldConfig {
				custVal := getCIFValueByField(cif, field)
				watchlistValues := getWatchlistValuesByField(wl, field)

				maxScore := 0.0
				bestMatchVal := ""
				for _, wlVal := range watchlistValues {
					score := matchScore(custVal, wlVal, cfg.MatchingAlgorithm)
					if score > maxScore {
						maxScore = score
						bestMatchVal = wlVal
					}
				}

				totalScore += maxScore * cfg.FieldWeight
				totalWeight += cfg.FieldWeight

				matchedFields = append(matchedFields, models.MatchingDetail{
					FieldName:      field,
					CustomerValue:  custVal,
					WatchlistValue: bestMatchVal,
					FieldScore:     maxScore,
					FieldWeight:    cfg.FieldWeight,
					AlgorithmUsed:  cfg.MatchingAlgorithm,
				})
			}

			if totalWeight == 0 {
				continue
			}

			finalScore := totalScore / totalWeight

			if finalScore >= threshold {
				result := models.MatchingResult{
					CIFNumber:       cif.CIFNumber,
					CustomerName:    cif.NamaNasabah,
					WatchlistID:     wl.ID,
					WatchlistSource: "MASTER_LOCAL_BLACKLIST",
					SimilarityScore: finalScore,
					Status:          "SUCCESS",
					ProcessDate:     time.Now(),
					ProcessTime:     time.Now(),
					CreatedAt:       time.Now(),
				}

				resultIndex := int64(len(results))
				results = append(results, result)
				detailsMap[resultIndex] = matchedFields
			}
		}
	}
	return results, detailsMap
}

// ===== MATCHING UNTUK TERORIS + WMD + LOCAL_BLACKLIST SEKALIGUS =====
func MatchCIFAll(
	db *sql.DB,
	cifs []models.MasterNasabah,
	watchlist []models.MasterWatchlist,
	configs []models.MatchingConfig,
) ([]models.MatchingResult, map[int64][]models.MatchingDetail, error) {

	var allResults []models.MatchingResult
	allDetails := make(map[int64][]models.MatchingDetail)

	// Pisahkan data watchlist berdasarkan source
	var terorisList []models.MasterWatchlist
	var wmdList []models.MasterWatchlist
	var localList []models.MasterWatchlist

	for _, wl := range watchlist {
		if !wl.IsActive {
			continue
		}
		switch strings.ToUpper(wl.Source) {
		case "MASTER_TERORIS":
			terorisList = append(terorisList, wl)
		case "MASTER_WMD":
			wmdList = append(wmdList, wl)
		case "MASTER_LOCAL_BLACKLIST":
			localList = append(localList, wl)
		}
	}

	/*
		// Matching TERORIS
		resultsTeroris, detailsTeroris := MatchCIFWithTeroris(db, cifs, terorisList, configs)
		for i, r := range resultsTeroris {
			idx := int64(len(allResults))
			allResults = append(allResults, r)
			allDetails[idx] = detailsTeroris[int64(i)]
		}

		// Matching WMD
		resultsWMD, detailsWMD := MatchCIFWithWMD(db, cifs, wmdList, configs)
		for i, r := range resultsWMD {
			idx := int64(len(allResults))
			allResults = append(allResults, r)
			allDetails[idx] = detailsWMD[int64(i)]
		}

		// Matching LOCAL BLACKLIST
		resultsLocal, detailsLocal := MatchCIFWithLocalBlacklist(db, cifs, localList, configs)
		for i, r := range resultsLocal {
			idx := int64(len(allResults))
			allResults = append(allResults, r)
			allDetails[idx] = detailsLocal[int64(i)]
		} */

	return allResults, allDetails, nil
}

// ===== UTIL =====
func getCIFValueByField(cif models.MasterNasabah, field string) string {
	switch field {
	case "nama", "namanasabah":
		return cif.NamaNasabah
	case "tempatlahir":
		return *cif.TempatLahir
	case "tanggallahir":
		if !cif.TanggalLahir.IsZero() {
			return cif.TanggalLahir.Format("2006-01-02")
		}
		return ""
	case "ktp":
		return *cif.KTP
	case "npwp":
		return *cif.NPWP
	case "nopaspor":
		return *cif.NoPaspor
	default:
		return ""
	}
}

func getWatchlistValuesByField(wl models.MasterWatchlist, field string) []string {
	switch field {
	case "nama":
		values := []string{}
		if strings.TrimSpace(wl.Nama) != "" {
			values = append(values, wl.Nama)
		}
		for _, alias := range wl.Aliases {
			if strings.TrimSpace(alias) != "" {
				values = append(values, alias)
			}
		}
		return values
	case "tempatlahir":
		if wl.TempatLahir != nil && strings.TrimSpace(*wl.TempatLahir) != "" {
			return []string{*wl.TempatLahir}
		}
	case "tanggallahir":
		if wl.TanggalLahir != nil && !wl.TanggalLahir.IsZero() {
			return []string{wl.TanggalLahir.Format("2006-01-02")}
		}
	case "ktp":
		if wl.KTP != nil && strings.TrimSpace(*wl.KTP) != "" {
			return []string{*wl.KTP}
		}
	case "npwp":
		if wl.NPWP != nil && strings.TrimSpace(*wl.NPWP) != "" {
			return []string{*wl.NPWP}
		}
	case "nopaspor":
		if wl.NoPaspor != nil && strings.TrimSpace(*wl.NoPaspor) != "" {
			return []string{*wl.NoPaspor}
		}
	}
	return []string{}
}

func matchScore(a, b, algorithm string) float64 {
	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))

	if a == "" || b == "" {
		return 0
	}

	switch strings.ToLower(algorithm) {
	case "levenshtein", "fuzzywuzzy", "ratio":
		return levenshtein.Similarity(a, b, nil)
	case "exact":
		if a == b {
			return 1.0
		}
		return 0.0
	default:
		return levenshtein.Similarity(a, b, nil)
	}
}

// ===== INSERT RESULT & DETAIL =====
func InsertMatchingResults(db *sql.DB, results []models.MatchingResult, detailsMap map[int64][]models.MatchingDetail) error {
	if len(results) == 0 {
		return nil
	}

	queryResult := `
        INSERT INTO MATCHING_RESULTS 
        (BatchId, CifNumber, CustomerName, WatchlistId, WatchlistSource, SimilarityScore, Status, ProcessDate, ProcessTime, CreatedAt)
        OUTPUT INSERTED.Id
        VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10)
    `
	queryDetail := `
        INSERT INTO MATCHING_DETAILS 
        (MatchingResultId, FieldName, CustomerValue, WatchlistValue, FieldScore, FieldWeight, AlgorithmUsed)
        VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)
    `

	for idx, r := range results {
		var insertedID int64
		err := db.QueryRow(queryResult,
			r.BatchID,
			r.CIFNumber,
			r.CustomerName,
			r.WatchlistID,
			r.WatchlistSource,
			r.SimilarityScore,
			r.Status,
			r.ProcessDate,
			r.ProcessTime,
			r.CreatedAt,
		).Scan(&insertedID)
		if err != nil {
			return fmt.Errorf("insert MATCHING_RESULTS gagal: %w", err)
		}

		// pake index, bukan WatchlistID
		if detailList, ok := detailsMap[int64(idx)]; ok {
			for _, d := range detailList {
				_, err := db.Exec(queryDetail,
					insertedID,
					d.FieldName,
					d.CustomerValue,
					d.WatchlistValue,
					d.FieldScore,
					d.FieldWeight,
					d.AlgorithmUsed,
				)
				if err != nil {
					return fmt.Errorf("insert MATCHING_DETAILS gagal: %w", err)
				}
			}
		}
	}
	return nil
}

func GetNextBatchID(db *sql.DB) (int64, error) {
	var lastBatchID sql.NullInt64
	err := db.QueryRow(`SELECT ISNULL(MAX(BatchId), 0) FROM MATCHING_RESULTS`).Scan(&lastBatchID)
	if err != nil {
		return 0, err
	}
	return lastBatchID.Int64 + 1, nil
}
