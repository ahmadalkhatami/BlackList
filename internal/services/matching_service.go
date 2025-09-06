package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repositories"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/agext/levenshtein"
)

// MatchingService mengoordinasikan proses matching CIF dengan berbagai sumber watchlist.
type MatchingService struct {
	DB               *sql.DB
	configService    SystemConfigInterface
	masterNasabah    MasterNasabahInterface
	watchlistService WatchlistServiceInterface
}

// --- Contract minimal untuk service lain agar file ini self-contained ---
// (Sesuaikan dengan interface yang sudah kamu punya; ini kompatibel dengan kode sebelumnya.)

// NewMatchingService membangun semua dependency dari *sql.DB.
func NewMatchingService(db *sql.DB) (*MatchingService, error) {
	// Repositories
	masterMatchingRepo := repositories.NewSQLMasterMatchingRepository(db)
	masterMatchingConfigRepo := repositories.NewSQLMasterMatchingConfigRepository(db)
	systemConfigRepo := repositories.NewSQLSystemConfigRepository(db)

	// Services pendukung (mengikuti pola yang sudah ada di project-mu)
	configService := NewConfigService(
		WithMasterMatching(masterMatchingRepo),
		WithMasterMatchingConfig(masterMatchingConfigRepo),
		WithSystemConfig(systemConfigRepo),
	)

	masterNasabahRepo := repositories.NewSQLMasterNasabahRepository(db)
	masterNasabahService := NewMasterNasabah(masterNasabahRepo)

	masterDTTOTRepo := repositories.NewSQLMasterTerorisRepository(db)
	masterWMDRepo := repositories.NewSQLMasterWMDRepository(db)
	masterLocalBlacklistRepo := repositories.NewSQLMasterLocalBlacklistRepository(db)

	watchlistService := NewWatchlistService(
		WithMasterTeroris(masterDTTOTRepo),
		WithMasterWMD(masterWMDRepo),
		WithMasterLocalBlacklist(masterLocalBlacklistRepo),
	)

	return &MatchingService{
		DB:               db,
		configService:    configService,
		masterNasabah:    masterNasabahService,
		watchlistService: watchlistService,
	}, nil
}

// RunAll melakukan matching untuk semua sumber (TERORIS, WMD, LOCAL_BLACKLIST)
func (s *MatchingService) RunAll() ([]models.MatchingResult, map[int64][]models.MatchingDetail, error) {
	threshold, err := s.configService.GetThresholdFromConfig("MATCHING_THRESHOLD")
	if err != nil {
		return nil, nil, fmt.Errorf("get threshold: %w", err)
	}

	fmt.Printf("threshold: %f", threshold)

	configsJoined, err := s.configService.GetJoinedMatchingConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("get joined matching config: %w", err)
	}

	// fmt.Printf("threshold: %s", err)

	cifs, err := s.masterNasabah.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("load CIF: %w", err)
	}

	allWatchlists, err := s.watchlistService.LoadAllWatchlists()
	if err != nil {
		return nil, nil, fmt.Errorf("load watchlists: %w", err)
	}

	// Split watchlists by source (aktif saja)
	var terorisList, wmdList, localList []models.MasterWatchlist
	for _, wl := range allWatchlists {
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

	// Jalankan generic matching per sumber
	res1, det1 := genericMatch(cifs, terorisList, configsJoined, "MASTER_TERORIS", threshold)
	res2, det2 := genericMatch(cifs, wmdList, configsJoined, "MASTER_WMD", threshold)
	res3, det3 := genericMatch(cifs, localList, configsJoined, "MASTER_LOCAL_BLACKLIST", threshold)

	// Gabungkan sekaligus reindex details agar konsisten dengan urutan results akhir
	allResults := make([]models.MatchingResult, 0, len(res1)+len(res2)+len(res3))
	allDetails := make(map[int64][]models.MatchingDetail)

	appendWithRemap := func(res []models.MatchingResult, det map[int64][]models.MatchingDetail) {
		for i, r := range res {
			newIdx := int64(len(allResults))
			allResults = append(allResults, r)
			if d, ok := det[int64(i)]; ok {
				allDetails[newIdx] = d
			}
		}
	}

	appendWithRemap(res1, det1)
	appendWithRemap(res2, det2)
	appendWithRemap(res3, det3)

	return allResults, allDetails, nil
}

// genericMatch mengerjakan matching untuk 1 sumber watchlist tertentu.
// - configsJoined: gunakan yang sudah include FieldName, FieldWeight, WatchlistSource, MatchingAlgorithm, IsActive.
// - threshold: ambang final score.
func genericMatch(
	cifs []models.MasterNasabah,
	watchlists []models.MasterWatchlist,
	configsJoined []models.JoinedMatchingConfig,
	source string,
	threshold float64,
) ([]models.MatchingResult, map[int64][]models.MatchingDetail) {
	results := make([]models.MatchingResult, 0)
	detailsMap := make(map[int64][]models.MatchingDetail)

	// 1) Build fieldConfig untuk source yang aktif
	fieldConfig := make(map[string]models.JoinedMatchingConfig)
	for _, cfg := range configsJoined {
		if strings.EqualFold(cfg.WatchlistSource, source) && cfg.IsActive {
			field := strings.ToLower(cfg.FieldName)
			fieldConfig[field] = cfg
		}
	}
	if len(fieldConfig) == 0 || len(cifs) == 0 || len(watchlists) == 0 {
		return results, detailsMap
	}

	// 2) Precompute values
	cifValues := precomputeCIFValuesJoined(cifs, fieldConfig)
	wlValues := precomputeWatchlistValuesJoined(watchlists, fieldConfig)

	// 3) Kerjakan paralel per CIF → hasil dikirim lewat channel agar aman dari race
	type resWithDetail struct {
		Result models.MatchingResult
		Detail []models.MatchingDetail
	}

	out := make(chan resWithDetail, 1024)
	var wg sync.WaitGroup

	for i, cif := range cifs {
		ci := i
		cifItem := cif
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j, wl := range watchlists {
				totalWeight := 0.0
				totalScore := 0.0
				matched := make([]models.MatchingDetail, 0, len(fieldConfig))

				for field, cfg := range fieldConfig {
					custVal := cifValues[ci][field]
					watchVals := wlValues[j][field]

					maxScore := 0.0
					best := ""
					for _, wv := range watchVals {
						s := matchScore(custVal, wv, cfg.MatchingAlgorithm)
						if s > maxScore {
							maxScore = s
							best = wv
						}
					}

					totalScore += maxScore * cfg.FieldWeight
					totalWeight += cfg.FieldWeight

					matched = append(matched, models.MatchingDetail{
						FieldName:      field,
						CustomerValue:  custVal,
						WatchlistValue: best,
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
					out <- resWithDetail{
						Result: models.MatchingResult{
							CIFNumber:       cifItem.CIFNumber,
							CustomerName:    cifItem.NamaNasabah,
							WatchlistID:     wl.ID,
							WatchlistSource: source,
							SimilarityScore: finalScore,
							Status:          "SUCCESS",
							ProcessDate:     time.Now(),
							ProcessTime:     time.Now(),
							CreatedAt:       time.Now(),
						},
						Detail: matched,
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	// Kumpulkan dan reindex lokal (0..n-1) untuk detailsMap
	idx := int64(0)
	for item := range out {
		results = append(results, item.Result)
		detailsMap[idx] = item.Detail
		idx++
	}

	return results, detailsMap
}

// ===== PRECOMPUTE (versi JoinedMatchingConfig) =====
func precomputeCIFValuesJoined(cifs []models.MasterNasabah, fieldCfg map[string]models.JoinedMatchingConfig) []map[string]string {
	res := make([]map[string]string, len(cifs))
	for i, cif := range cifs {
		m := make(map[string]string, len(fieldCfg))
		for field := range fieldCfg {
			m[field] = getCIFValueByField(cif, field)
		}
		res[i] = m
	}
	return res
}

func precomputeWatchlistValuesJoined(wls []models.MasterWatchlist, fieldCfg map[string]models.JoinedMatchingConfig) []map[string][]string {
	res := make([]map[string][]string, len(wls))
	for i, wl := range wls {
		m := make(map[string][]string, len(fieldCfg))
		for field := range fieldCfg {
			m[field] = getWatchlistValuesByField(wl, field)
		}
		res[i] = m
	}
	return res
}

// ===== UTIL FIELD EXTRACTORS (lebih aman pointer-nil) =====
func getCIFValueByField(cif models.MasterNasabah, field string) string {
	switch strings.ToLower(field) {
	case "nama", "namanasabah":
		return strings.TrimSpace(cif.NamaNasabah)
	case "tempatlahir":
		if cif.TempatLahir != nil {
			return strings.TrimSpace(*cif.TempatLahir)
		}
		return ""
	case "tanggallahir":
		if !cif.TanggalLahir.IsZero() {
			return cif.TanggalLahir.Format("2006-01-02")
		}
		return ""
	case "ktp":
		if cif.KTP != nil {
			return strings.TrimSpace(*cif.KTP)
		}
		return ""
	case "npwp":
		if cif.NPWP != nil {
			return strings.TrimSpace(*cif.NPWP)
		}
		return ""
	case "nopaspor":
		if cif.NoPaspor != nil {
			return strings.TrimSpace(*cif.NoPaspor)
		}
		return ""
	default:
		return ""
	}
}

func getWatchlistValuesByField(wl models.MasterWatchlist, field string) []string {
	switch strings.ToLower(field) {
	case "nama":
		values := []string{}
		if s := strings.TrimSpace(wl.Nama); s != "" {
			values = append(values, s)
		}
		for _, alias := range wl.Aliases {
			if s := strings.TrimSpace(alias); s != "" {
				values = append(values, s)
			}
		}
		return values
	case "tempatlahir":
		if wl.TempatLahir != nil {
			if s := strings.TrimSpace(*wl.TempatLahir); s != "" {
				return []string{s}
			}
		}
	case "tanggallahir":
		if wl.TanggalLahir != nil && !wl.TanggalLahir.IsZero() {
			return []string{wl.TanggalLahir.Format("2006-01-02")}
		}
	case "ktp":
		if wl.KTP != nil {
			if s := strings.TrimSpace(*wl.KTP); s != "" {
				return []string{s}
			}
		}
	case "npwp":
		if wl.NPWP != nil {
			if s := strings.TrimSpace(*wl.NPWP); s != "" {
				return []string{s}
			}
		}
	case "nopaspor":
		if wl.NoPaspor != nil {
			if s := strings.TrimSpace(*wl.NoPaspor); s != "" {
				return []string{s}
			}
		}
	}
	return []string{}
}

// Scoring function
func matchScore(a, b, algorithm string) float64 {
	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))
	if a == "" || b == "" {
		return 0
	}
	switch strings.ToLower(strings.TrimSpace(algorithm)) {
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

// ===== INSERT RESULT & DETAIL (dipertahankan, dengan minor perapian) =====
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

// ===== Helper untuk quick-run satu sumber tertentu jika dibutuhkan =====
func (s *MatchingService) RunForSource(source string) ([]models.MatchingResult, map[int64][]models.MatchingDetail, error) {
	threshold, err := s.configService.GetThresholdFromConfig("MATCHING_THRESHOLD")
	if err != nil {
		return nil, nil, err
	}
	configsJoined, err := s.configService.GetJoinedMatchingConfig()
	if err != nil {
		return nil, nil, err
	}
	cifs, err := s.masterNasabah.Load()
	if err != nil {
		return nil, nil, err
	}

	allWatchlists, err := s.watchlistService.LoadAllWatchlists()
	if err != nil {
		return nil, nil, err
	}

	filtered := make([]models.MasterWatchlist, 0)
	for _, wl := range allWatchlists {
		if wl.IsActive && strings.EqualFold(wl.Source, source) {
			filtered = append(filtered, wl)
		}
	}
	res, det := genericMatch(cifs, filtered, configsJoined, strings.ToUpper(source), threshold)
	return res, det, nil
}

// ===== Debug helpers (opsional) =====
// func DebugPrintFields(cifs []models.MasterNasabah, watchlists []models.MasterWatchlist, configs []models.JoinedMatchingConfig) {
// 	cifKeys := utils.GetStructKeys(&models.MasterNasabah{})
// 	fmt.Printf("\n📌 Field di MasterNasabah: %+v\n", cifKeys)

// 	if len(watchlists) > 0 {
// 		// Asumsikan salah satu tipe watchlist untuk melihat fields-nya
// 		fmt.Printf("📌 Contoh data Watchlist (first): %+v\n", watchlists[0])
// 	}
// 	fmt.Printf("📌 Jumlah Config Joined: %d\n", len(configs))
// }
