package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"BlackListWorker/models"

	"github.com/agext/levenshtein"
)

// ===== MATCHING UNTUK MASTER_TERORIS =====
func MatchCIFWithTeroris(
	db *sql.DB,
	cifs []models.MASTER_NASABAH,
	terorisList []models.MasterWatchlist,
	configs []models.MatchingConfig,
) ([]models.MatchingResult, map[int64][]models.MatchingDetail) {

	var results []models.MatchingResult
	detailsMap := make(map[int64][]models.MatchingDetail)
	fieldConfig := map[string]models.MatchingConfig{}

	// Ambil threshold dari SYSTEM_CONFIG
	threshold := GetThresholdFromConfig(db, "MATCHING_THRESHOLD")
	fmt.Printf("🔧 Threshold MATCHING_THRESHOLD: %.2f\n", threshold)

	// Ambil config untuk MASTER_TERORIS saja
	for _, cfg := range configs {
		if strings.ToUpper(cfg.WatchlistSource) == "MASTER_TERORIS" && cfg.IsActive {
			fieldConfig[strings.ToLower(cfg.FieldName)] = cfg
		}
	}
	// fmt.Printf("📌 Field config untuk MASTER_TERORIS: %+v\n", fieldConfig)

	for _, cif := range cifs {
		for _, wl := range terorisList {
			if !wl.IsActive || wl.Source != "MASTER_TERORIS" {
				continue
			}

			fmt.Printf("\n🚀 Proses CIF: %s | Watchlist: %s (Source: %s)\n",
				cif.NamaNasabah, wl.Nama, wl.Source)

			totalWeight := 0.0
			totalScore := 0.0
			var matchedFields []models.MatchingDetail

			for field, cfg := range fieldConfig {
				custVal := getCIFValueByField(cif, field)
				watchlistValues := getWatchlistValuesByField(wl, field)

				fmt.Print("custVal:", custVal, " | wlVals:", watchlistValues, " ")

				maxScore := 0.0
				bestMatchVal := ""
				for _, wlVal := range watchlistValues {
					score := matchScore(custVal, wlVal, cfg.MatchingAlgorithm)
					if score > maxScore {
						maxScore = score
						bestMatchVal = wlVal
					}
				}

				// Debug per field
				// fmt.Printf("   🔍 Field: %s | CIF: '%s' | Watchlist: '%s' | Score: %.2f | Algoritma: %s\n",
				// 	field, custVal, bestMatchVal, maxScore, cfg.MatchingAlgorithm)

				totalScore += maxScore * cfg.FieldWeight
				totalWeight += cfg.FieldWeight

				fmt.Printf("totalWeight:", totalWeight, " | totalScore:", totalScore, "\n")

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
				fmt.Println("⚠️ Skip: totalWeight = 0 (tidak ada config aktif)")
				continue
			}

			finalScore := totalScore / totalWeight
			fmt.Printf("➡️ FinalScore CIF %s vs Watchlist %s = %.2f (Threshold %.2f)\n",
				cif.NamaNasabah, wl.Nama, finalScore, threshold)

			if finalScore >= threshold {
				result := models.MatchingResult{
					CIFNumber:       cif.CIFNumber,
					CustomerName:    cif.NamaNasabah,
					WatchlistID:     wl.ID,
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

	return results, detailsMap
}


// ===== MATCHING UNTUK MASTER_WMD =====
func MatchCIFWithWMD(
	db *sql.DB,
	cifs []models.MASTER_NASABAH,
	wmdList []models.MasterWatchlist,
	configs []models.MatchingConfig,
) ([]models.MatchingResult, map[int64][]models.MatchingDetail) {

	var results []models.MatchingResult
	detailsMap := make(map[int64][]models.MatchingDetail)
	fieldConfig := map[string]models.MatchingConfig{}

	threshold := GetThresholdFromConfig(db, "MATCHING_THRESHOLD")

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
}

// ===== MATCHING UNTUK MASTER_LOCAL_BLACKLIST =====
func MatchCIFWithLocalBlacklist(
	db *sql.DB,
	cifs []models.MASTER_NASABAH,
	localList []models.MasterWatchlist,
	configs []models.MatchingConfig,
) ([]models.MatchingResult, map[int64][]models.MatchingDetail) {

	var results []models.MatchingResult
	detailsMap := make(map[int64][]models.MatchingDetail)
	fieldConfig := map[string]models.MatchingConfig{}

	// Ambil threshold
	threshold := GetThresholdFromConfig(db, "MATCHING_THRESHOLD")

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
	cifs []models.MASTER_NASABAH,
	watchlist []models.MasterWatchlist,
	configs []models.MatchingConfig,
) ([]models.MatchingResult, map[int64][]models.MatchingDetail) {

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
	}

	return allResults, allDetails
}

// ===== UTIL =====
func getCIFValueByField(cif models.MASTER_NASABAH, field string) string {
	switch field {
	case "nama", "namanasabah":
		return cif.NamaNasabah
	case "tempatlahir":
		return cif.TempatLahir
	case "tanggallahir":
		if !cif.TanggalLahir.IsZero() {
			return cif.TanggalLahir.Format("2006-01-02")
		}
		return ""
	case "ktp":
		return cif.KTP
	case "npwp":
		return cif.NPWP
	case "nopaspor":
		return cif.NoPaspor
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
		for _, alias := range wl.Alias { // alias sudah slice di struct
			if strings.TrimSpace(alias) != "" {
				values = append(values, alias)
			}
		}
		return values
	case "tempatlahir":
		if wl.TempatLahir != "" {
			return []string{wl.TempatLahir}
		}
	case "tanggallahir":
		if !wl.TanggalLahir.IsZero() {
			return []string{wl.TanggalLahir.Format("2006-01-02")}
		}
	case "ktp":
		if wl.KTP != "" {
			return []string{wl.KTP}
		}
	case "npwp":
		if wl.NPWP != "" {
			return []string{wl.NPWP}
		}
	case "nopaspor":
		if wl.NoPaspor != "" {
			return []string{wl.NoPaspor}
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

// 		fmt.Printf(`
// 💾 QUERY RESULT:
// %s
// VALUES (
//   BatchId=%v,
//   CIFNumber='%v',
//   CustomerName='%v',
//   WatchlistId=%v,
//   WatchlistSource='%v',
//   SimilarityScore=%.4f,
//   Status='%v',
//   ProcessDate='%v',
//   ProcessTime='%v',
//   CreatedAt='%v'
// )
// `,
//     queryResult,
//     r.BatchID,
//     r.CIFNumber,
//     r.CustomerName,
//     r.WatchlistID,
//     r.WatchlistSource,
//     r.SimilarityScore,
//     r.Status,
//     r.ProcessDate,
//     r.ProcessTime,
//     r.CreatedAt,
// )

		// fmt.Print("similarity score:", r.SimilarityScore)

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
