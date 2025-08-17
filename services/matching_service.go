package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"BlackListWorker/models"

	"github.com/agext/levenshtein"
)

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

	// fmt.Print("🔍 Menggunakan Config: %.2f\n", configs)
	// return results, detailsMap
	// Ambil config untuk MASTER_TERORIS saja
	for _, cfg := range configs {
		if strings.ToUpper(cfg.WatchlistSource) == "MASTER_TERORIS" && cfg.IsActive {
			fieldConfig[strings.ToLower(cfg.FieldName)] = cfg
		}
	}

	// Loop setiap CIF dan cocokan ke data teroris
	for _, cif := range cifs {
		for _, teroris := range terorisList {
			if !teroris.IsActive || teroris.Source != "MASTER_TERORIS" {
				continue
			}

			totalWeight := 0.0
			totalScore := 0.0
			var matchedFields []models.MatchingDetail

			for field, cfg := range fieldConfig {
				custVal := getCIFValueByField(cif, field)
				fmt.Print("🔍 Matching field: %s, Customer Value: %s\n", field, custVal)
				watchlistValues := getWatchlistValuesByField(teroris, field)

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
					WatchlistID:     teroris.ID,
					WatchlistSource: "MASTER_TERORIS",
					SimilarityScore: finalScore,
					Status:          "SUCCESS",
					ProcessDate:     time.Now(),
					ProcessTime:     time.Now(),
					CreatedAt:       time.Now(),
				}

				resultIndex := int64(len(results)) // pakai index array sebelum append
				results = append(results, result)
				detailsMap[resultIndex] = matchedFields
			}
		}
	}

	return results, detailsMap
}


func getCIFValueByField(cif models.MASTER_NASABAH, field string) string {
	switch field {
	case "nama", "namanasabah":
		return cif.NamaNasabah
	case "tempatlahir":
		return cif.TempatLahir
	case "tanggallahir":
		return cif.TanggalLahir.Format("2006-01-02")
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
		values := []string{wl.Nama}
		values = append(values, wl.Alias...) // wl.Alias []string
		return values
	case "tempatlahir":
		return []string{wl.TempatLahir}
	case "tanggallahir":
		return []string{wl.TanggalLahir.Format("2006-01-02")}
	case "ktp":
		return []string{wl.KTP}
	case "npwp":
		return []string{wl.NPWP}
	case "nopaspor":
		return []string{wl.NoPaspor}
	default:
		return []string{""}
	}
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

    for _, r := range results {
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

        // Masukkan details untuk result ini
        if detailList, ok := detailsMap[r.WatchlistID]; ok {
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

