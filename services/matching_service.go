package services

import (
	"database/sql"
	"log"
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
) []models.MatchingResult {

	var results []models.MatchingResult
	fieldConfig := map[string]models.MatchingConfig{}

	// Ambil threshold dari SYSTEM_CONFIG (fallback defaultThreshold kalau tidak ada)
	threshold := GetThresholdFromConfig(db, "MATCHING_THRESHOLD")

	// Ambil config untuk DTTOT saja
	for _, cfg := range configs {
		if strings.ToUpper(cfg.WatchlistSource) == "DTTOT" && cfg.IsActive {
			fieldConfig[strings.ToLower(cfg.FieldName)] = cfg
		}
	}

	// Loop setiap CIF dan cocokan ke data teroris
	for _, cif := range cifs {
		for _, teroris := range terorisList {
			if !teroris.IsActive || teroris.Source != "DTTOT" {
				continue
			}

			totalWeight := 0.0
			totalScore := 0.0
			var matchedFields []models.MatchingDetail

			for field, cfg := range fieldConfig {
				custVal := getCIFValueByField(cif, field)

				// Ambil semua nilai watchlist (nama utama + alias jika field "nama")
				watchlistValues := getWatchlistValuesByField(teroris, field)

				// Cari skor tertinggi
				maxScore := 0.0
				bestMatchVal := ""
				for _, wlVal := range watchlistValues {
					score := matchScore(custVal, wlVal, cfg.MatchingAlgorithm)
					if score > maxScore {
						maxScore = score
						bestMatchVal = wlVal
					}
				}

				// Tambah ke akumulasi total skor
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

			// Bandingkan dengan threshold dari config
			if finalScore >= threshold {
				result := models.MatchingResult{
					CIFNumber:       cif.CIFNumber,
					CustomerName:    cif.NamaNasabah,
					WatchlistID:     teroris.ID,
					WatchlistSource: "DTTOT",
					SimilarityScore: finalScore,
					Status:          "SUCCESS",
					ProcessDate:     time.Now(),
					ProcessTime:     time.Now(),
					CreatedAt:       time.Now(),
				}
				results = append(results, result)
				// matchedFields bisa dimasukkan ke MATCHING_DETAILS kalau dibutuhkan
			}
		}
	}

	return results
}

func getCIFValueByField(cif models.MASTER_NASABAH, field string) string {
	switch field {
	case "nama", "nama_nasabah":
		return cif.NamaNasabah
	case "tempat_lahir":
		return cif.TempatLahir
	case "tanggal_lahir":
		return cif.TanggalLahir.Format("2006-01-02")
	case "ktp":
		return cif.KTP
	case "npwp":
		return cif.NPWP
	case "no_paspor":
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
	case "tempat_lahir":
		return []string{wl.TempatLahir}
	case "tanggal_lahir":
		return []string{wl.TanggalLahir.Format("2006-01-02")}
	case "ktp":
		return []string{wl.KTP}
	case "npwp":
		return []string{wl.NPWP}
	case "no_paspor":
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

func InsertMatchingResults(db *sql.DB, results []models.MatchingResult) error {
	if len(results) == 0 {
		return nil
	}

	query := `
		INSERT INTO MATCHING_RESULTS 
		(batch_id, cif_number, customer_name, watchlist_id, watchlist_source, similarity_score, status, process_date, process_time, created_at)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10)
	`

	for _, r := range results {
		_, err := db.Exec(query,
			r.BatchID,         // @p1
			r.CIFNumber,       // @p2
			r.CustomerName,    // @p3
			r.WatchlistID,     // @p4
			r.WatchlistSource, // @p5
			r.SimilarityScore, // @p6
			r.Status,          // @p7
			r.ProcessDate,     // @p8
			r.ProcessTime,     // @p9
			r.CreatedAt,       // @p10
		)
		if err != nil {
			log.Println("❌ Gagal insert MATCHING_RESULTS:", err)
			return err
		}
	}

	return nil
}

func GetNextBatchID(db *sql.DB) (int64, error) {
    var lastBatchID sql.NullInt64
    err := db.QueryRow(`SELECT ISNULL(MAX(batch_id), 0) FROM MATCHING_RESULTS`).Scan(&lastBatchID)
    if err != nil {
        return 0, err
    }
    return lastBatchID.Int64 + 1, nil
}

