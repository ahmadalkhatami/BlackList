package report

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/utils"
)

func (r *ReportGenerator) GenerateFromResultsCSVStream(filename string, results *models.MatchResults, matcherType string) error {
	// Buat file CSV
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("gagal membuat file csv: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header CSV
	header := []string{
		"CIF Number", "Nama Nasabah", "KTP", "NPWP", "No Paspor", "Tempat Lahir", "Tanggal Lahir", "Kewarganegaraan",
		"Watchlist Nama", "Watchlist KTP", "Watchlist Passport", "Watchlist NPWP", "Watchlist Tempat Lahir", "Watchlist Tanggal Lahir", "Watchlist Kewarganegaraan",
		"Similarity Score", "Reason", "Tanggal", "Waktu",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("gagal menulis header csv: %w", err)
	}

	now := time.Now()

	normalize := func(s string) string {
		return strings.ToLower(strings.TrimSpace(s))
	}

	// Fungsi untuk generate alias fields sesuai matcher
	aliasFieldsForMatcher := func(matcher string) []string {
		fields := []string{"nama"}
		aliasCount := 4
		if matcher == "MASTER_WMD" {
			aliasCount = 10
		}
		for i := 1; i <= aliasCount; i++ {
			fields = append(fields, fmt.Sprintf("alias%d", i))
		}
		return fields
	}

	aliasFields := aliasFieldsForMatcher(matcherType)

	// Ambil WatchlistValue case-insensitive dengan fallback alias
	getWatchlistValue := func(resultID int64, target string) string {
		target = normalize(target)
		var fallback string
		for _, d := range results.MatchDetail {
			if d.MatchingResultId != nil && utils.ValInt64(d.MatchingResultId) == resultID {
				field := normalize(utils.ValStr(d.FieldName))
				if field == target {
					return utils.ValStr(d.WatchlistValue)
				}
				for _, a := range aliasFields {
					if field == a && fallback == "" {
						fallback = utils.ValStr(d.WatchlistValue)
					}
				}
			}
		}
		return fallback
	}

	for _, result := range results.MatchResult {
		// log.Printf("DEBUG: Processing result.Id=%d, CustomerName=%s", result.Id, utils.ValStr(result.CustomerName))

		// Ambil data Nasabah
		n := models.MasterNasabah{
			CifNumber:   utils.ValStr(result.CifNumber),
			NamaNasabah: utils.ValStr(result.CustomerName),
		}

		for _, detail := range results.MatchDetail {
			if detail.MatchingResultId != nil && utils.ValInt64(detail.MatchingResultId) == result.Id {
				if detail.FieldName != nil {
					switch normalize(*detail.FieldName) {
					case "nama":
						n.NamaNasabah = utils.ValStr(detail.CustomerValue)
					case "ktp":
						n.KTP = detail.CustomerValue
					case "npwp":
						n.NPWP = detail.CustomerValue
					case "nopaspor":
						n.NoPaspor = detail.CustomerValue
					case "tempatlahir":
						n.TempatLahir = detail.CustomerValue
					case "tanggallahir":
						n.TanggalLahir = utils.ParseDatePtrAny(utils.ValStr(detail.CustomerValue))
					}
				}
			}
		}

		// Ambil data Watchlist dari MatchDetail (dengan alias fallback)
		w := models.MasterWatchlist{
			KTP:          utils.PtrOrEmpty(getWatchlistValue(result.Id, "ktp")),
			NoPaspor:     utils.PtrOrEmpty(getWatchlistValue(result.Id, "nopaspor")),
			NPWP:         utils.PtrOrEmpty(getWatchlistValue(result.Id, "npwp")),
			TempatLahir:  utils.PtrOrEmpty(getWatchlistValue(result.Id, "tempatlahir")),
			TanggalLahir: utils.ParseDatePtrAny(getWatchlistValue(result.Id, "tanggallahir")),
		}

		row := []string{
			n.CifNumber,
			n.NamaNasabah,
			utils.ValStr(n.KTP),
			utils.ValStr(n.NPWP),
			utils.ValStr(n.NoPaspor),
			utils.ValStr(n.TempatLahir),
			utils.ValTime(n.TanggalLahir),
			"-",

			utils.ValStr(result.CustomerName),
			utils.ValStr(w.KTP),
			utils.ValStr(w.NoPaspor),
			utils.ValStr(w.NPWP),
			utils.ValStr(w.TempatLahir),
			utils.ValTime(w.TanggalLahir),
			"-",

			fmt.Sprintf("%.0f%%", utils.ValFloat64(result.SimilarityScore)*100),
			utils.ValStr(result.Status),
			now.Format("02/01/2006"),
			now.Format("15:04:05"),
		}

		// log.Printf("DEBUG CSV ROW: result.Id=%d, row=%v", result.Id, row)

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("gagal menulis baris csv: %w", err)
		}
	}

	return nil
}
