package services

import (
	"strings"
	"time"

	"BlackListWorker/models"

	"github.com/agext/levenshtein"
)

// MatchCIFWithTeroris mencocokkan CIF dengan daftar teroris berdasarkan konfigurasi
func MatchCIFWithTeroris(
	cifs []models.MASTER_NASABAH,
	terorisList []models.MasterWatchlist,
	configs []models.MatchingConfig,
) []models.MatchingResult {

	var results []models.MatchingResult
	fieldConfig := map[string]models.MatchingConfig{}

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
				wlVal := getWatchlistValueByField(teroris, field)
				score := matchScore(custVal, wlVal, cfg.MatchingAlgorithm)

				// Tambah ke akumulasi total skor
				totalScore += score * cfg.FieldWeight
				totalWeight += cfg.FieldWeight

				matchedFields = append(matchedFields, models.MatchingDetail{
					FieldName:      field,
					CustomerValue:  custVal,
					WatchlistValue: wlVal,
					FieldScore:     score,
					FieldWeight:    cfg.FieldWeight,
					AlgorithmUsed:  cfg.MatchingAlgorithm,
				})
			}

			if totalWeight == 0 {
				continue
			}

			finalScore := totalScore / totalWeight

			// Threshold 85% match
			if finalScore >= 0.85 {
				result := models.MatchingResult{
					CIFNumber:       cif.CIFNumber,
					CustomerName:    cif.NamaNasabah,
					WatchlistID:     teroris.ID,
					WatchlistSource: "DTTOT",
					SimilarityScore: finalScore,
					Status:          "PENDING",
					ProcessDate:     time.Now(),
					ProcessTime:     time.Now(),
					CreatedAt:       time.Now(),
				}
				results = append(results, result)
				// matchedFields bisa dimasukkan juga ke MATCHING_DETAILS jika dibutuhkan
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

func getWatchlistValueByField(wl models.MasterWatchlist, field string) string {
	switch field {
	case "nama":
		return wl.Nama
	case "tempat_lahir":
		return wl.TempatLahir
	case "tanggal_lahir":
		return wl.TanggalLahir.Format("2006-01-02")
	case "ktp":
		return wl.KTP
	case "npwp":
		return wl.NPWP
	case "no_paspor":
		return wl.NoPaspor
	default:
		return ""
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
		// Default fallback ke Similarity
		return levenshtein.Similarity(a, b, nil)
	}
}
