package report

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/utils"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type MatchSectionData struct {
	SheetName       string
	title           string
	NasabahCount    int
	WatchlistCount  int
	SimilarityScore *float64
	Reason          string
}

func formatDate(t *time.Time) string {
	if t != nil {
		return t.Format("02/01/2006")
	}
	return ""
}

// ================= MAIN ==================

func (r *ReportGenerator) GenerateSection(
	sheetName string,
	nasabah []models.MasterNasabah,
	watchlist []models.MasterWatchlist,
	ms MatchSectionData,
) error {

	index, err := r.file.NewSheet(sheetName)
	if err != nil {
		return err
	}
	r.file.SetActiveSheet(index)

	_ = r.file.MergeCell(sheetName, "A1", "G1")
	_ = r.file.MergeCell(sheetName, "H1", "N1")
	_ = r.file.MergeCell(sheetName, "O1", "R1")

	r.file.SetCellValue(sheetName, "A1", "Data Nasabah Bank Index")
	r.file.SetCellValue(sheetName, "H1", "Data Watchlist")
	r.file.SetCellValue(sheetName, "O1", "Pemadanan")

	headers := []string{
		"Nomor CIF", "Nama", "KTP", "NPWP", "Tgl. Lahir", "Tempat Lahir", "Kewarganegaraan", // A–G
		"Nama", "KTP", "Passport", "NPWP", "Tgl. Lahir", "Tempat Lahir", "Kewarganegaraan", // H–N
		"Hasil Similarity", "Reason", "Tanggal", "Waktu", // O–R
	}
	for i, h := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		cell := fmt.Sprintf("%s2", col)
		r.file.SetCellValue(sheetName, cell, h)
	}

	row := 3
	rowCount := utils.Max(len(nasabah), len(watchlist))
	now := time.Now()

	for i := 0; i < rowCount; i++ {
		var n models.MasterNasabah
		if i < len(nasabah) {
			n = nasabah[i]
		}
		var w models.MasterWatchlist
		if i < len(watchlist) {
			w = watchlist[i]
		}

		score := "0%"
		if ms.SimilarityScore != nil {
			score = fmt.Sprintf("%.0f%%", *ms.SimilarityScore*100)
		}

		values := []interface{}{
			n.CifNumber, n.NamaNasabah, utils.ValStr(n.KTP), utils.ValStr(n.NPWP), formatDate(n.TanggalLahir), utils.ValStr(n.TempatLahir), utils.ValStr(n.StatusNasabah),
			w.Nama, utils.ValStr(w.KTP), utils.ValStr(w.NoPaspor), utils.ValStr(w.NPWP), formatDate(w.TanggalLahir), utils.ValStr(w.TempatLahir), "-", // kewarganegaraan kosong
			score, ms.Reason, now.Format("02/01/2006"), now.Format("15:04:05"),
		}

		for j, v := range values {
			col, _ := excelize.ColumnNumberToName(j + 1)
			cell := fmt.Sprintf("%s%d", col, row+i)
			r.file.SetCellValue(sheetName, cell, v)
		}
	}

	for i := 1; i <= len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i)
		_ = r.file.SetColWidth(sheetName, col, col, 18)
	}

	return nil
}

/* INI LEBIH OKE
func (r *ReportGenerator) GenerateSectionStream(
	sheetName string,
	nasabah []models.MasterNasabah,
	watchlist []models.MasterWatchlist,
	ms MatchSectionData,
) error {
	// Buat sheet baru
	_, err := r.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// Pakai stream writer
	sw, err := r.file.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}

	// Pakai Styler
	st := NewStyler(r.file)

	// ========================
	// 1. HEADER (ROW 1 & 2)
	// ========================
	row1 := []interface{}{
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Data Nasabah Bank Index"},
		"", "", "", "", "", "",
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Data Watchlist"},
		"", "", "", "", "", "",
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Pemadanan"},
		"", "", "", "",
	}
	if err := sw.SetRow("A1", row1); err != nil {
		return err
	}

	headers2 := []interface{}{
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Nomor CIF"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Nama"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "KTP"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "NPWP"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Tgl. Lahir"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Tempat Lahir"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Kewarganegaraan"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Nama"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "KTP"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Passport"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "NPWP"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Tgl. Lahir"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Tempat Lahir"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Kewarganegaraan"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Hasil Similarity"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Reason"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Tanggal"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Waktu"},
	}
	if err := sw.SetRow("A2", headers2); err != nil {
		return err
	}

	// ========================
	// 2. DATA (ROW 3++)
	// ========================
	rowIndex := 3
	now := time.Now()
	for i := 0; i < utils.Max(len(nasabah), len(watchlist)); i++ {
		row := []interface{}{}

		// Data nasabah
		if i < len(nasabah) {
			n := nasabah[i]
			row = append(row,
				excelize.Cell{StyleID: st.RowStyle, Value: n.CifNumber},
				excelize.Cell{StyleID: st.RowStyle, Value: n.NamaNasabah},
				excelize.Cell{StyleID: st.RowStyle, Value: n.KTP},
				excelize.Cell{StyleID: st.RowStyle, Value: n.NPWP},
				excelize.Cell{StyleID: st.RowStyle, Value: utils.ValTime(n.TanggalLahir)},
				excelize.Cell{StyleID: st.RowStyle, Value: n.TempatLahir},
				excelize.Cell{StyleID: st.RowStyle, Value: "-"}, // kewarganegaraan kosong
			)
		} else {
			for j := 0; j < 7; j++ {
				row = append(row, excelize.Cell{StyleID: st.RowStyle, Value: "-"})
			}
		}

		// Data watchlist
		if i < len(watchlist) {
			w := watchlist[i]
			row = append(row,
				excelize.Cell{StyleID: st.RowStyle, Value: w.Nama},
				excelize.Cell{StyleID: st.RowStyle, Value: w.KTP},
				excelize.Cell{StyleID: st.RowStyle, Value: w.NoPaspor},
				excelize.Cell{StyleID: st.RowStyle, Value: w.NPWP},
				excelize.Cell{StyleID: st.RowStyle, Value: utils.ValTime(w.TanggalLahir)},
				excelize.Cell{StyleID: st.RowStyle, Value: w.TempatLahir},
				excelize.Cell{StyleID: st.RowStyle, Value: "-"}, // kewarganegaraan kosong
			)
		} else {
			for j := 0; j < 7; j++ {
				row = append(row, excelize.Cell{StyleID: st.RowStyle, Value: "-"})
			}
		}

		// Data pemadanan
		score := "0%"
		if ms.SimilarityScore != nil {
			score = fmt.Sprintf("%.0f%%", *ms.SimilarityScore*100)
		}
		row = append(row,
			excelize.Cell{StyleID: st.RowStyle, Value: score},
			excelize.Cell{StyleID: st.RowStyle, Value: ms.Reason},
			excelize.Cell{StyleID: st.RowStyle, Value: now.Format("02/01/2006")},
			excelize.Cell{StyleID: st.RowStyle, Value: now.Format("15:04:05")},
		)

		cell, _ := excelize.CoordinatesToCellName(1, rowIndex)
		if err := sw.SetRow(cell, row); err != nil {
			return err
		}
		rowIndex++
	}

	// ========================
	// 3. Flush stream
	// ========================
	if err := sw.Flush(); err != nil {
		return err
	}

	// ========================
	// 4. Merge Cells (Setelah Flush)
	// ========================
	// Merge untuk "Data Nasabah Bank Index" dari A1 ke G1 (7 kolom)
	if err := r.file.MergeCell(sheetName, "A1", "G1"); err != nil {
		return err
	}
	// Merge untuk "Data Watchlist" dari H1 ke N1 (7 kolom)
	if err := r.file.MergeCell(sheetName, "H1", "N1"); err != nil {
		return err
	}
	// Merge untuk "Pemadanan" dari O1 ke R1 (4 kolom)
	if err := r.file.MergeCell(sheetName, "O1", "R1"); err != nil {
		return err
	}

	return nil
} */

func (r *ReportGenerator) GenerateSectionStream(
	sheetName string,
	nasabah []models.MasterNasabah,
	watchlist []models.MasterWatchlist,
	ms MatchSectionData,
) error {

	st := NewStyler(r.file)

	_, err := r.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	if err := r.file.SetCellValue(sheetName, "A1", "Data Nasabah Bank Index"); err != nil {
		return err
	}
	if err := r.file.SetCellValue(sheetName, "H1", "Data Watchlist"); err != nil {
		return err
	}
	if err := r.file.SetCellValue(sheetName, "O1", "Pemadanan"); err != nil {
		return err
	}

	if err := r.file.SetCellStyle(sheetName, "A1", "A1", st.HeaderStyle); err != nil {
		return err
	}
	if err := r.file.SetCellStyle(sheetName, "H1", "H1", st.HeaderStyle); err != nil {
		return err
	}
	if err := r.file.SetCellStyle(sheetName, "O1", "O1", st.HeaderStyle); err != nil {
		return err
	}

	if err := r.file.MergeCell(sheetName, "A1", "G1"); err != nil {
		fmt.Printf("ERROR: Merge A1:G1 Gagal: %v\n", err)
		return err
	}
	if err := r.file.MergeCell(sheetName, "H1", "N1"); err != nil {
		fmt.Printf("ERROR: Merge H1:N1 Gagal: %v\n", err)
		return err
	}
	if err := r.file.MergeCell(sheetName, "O1", "R1"); err != nil {
		fmt.Printf("ERROR: Merge O1:R1 Gagal: %v\n", err)
		return err
	}

	if err := r.file.SetRowHeight(sheetName, 1, 30); err != nil {
		return err
	}

	sw, err := r.file.NewStreamWriter(sheetName)
	if err != nil {
		return err
	}

	headers2 := []interface{}{
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Nomor CIF"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Nama"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "KTP"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "NPWP"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Tgl. Lahir"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Tempat Lahir"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Kewarganegaraan"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Nama"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "KTP"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Passport"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "NPWP"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Tgl. Lahir"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Tempat Lahir"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Kewarganegaraan"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Hasil Similarity"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Reason"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Tanggal"},
		excelize.Cell{StyleID: st.HeaderStyle, Value: "Waktu"},
	}

	if err := sw.SetRow("A2", headers2); err != nil {
		return err
	}

	if err := r.file.SetRowHeight(sheetName, 2, 35); err != nil {
		return err
	}

	totalRows := utils.Max(len(nasabah), len(watchlist))
	rowIndex := 3
	now := time.Now()

	for i := 0; i < totalRows; i++ {
		row := []interface{}{}

		if i < len(nasabah) {
			n := nasabah[i]
			row = append(row,
				excelize.Cell{StyleID: st.RowStyle, Value: n.CifNumber},
				excelize.Cell{StyleID: st.RowStyle, Value: n.NamaNasabah},
				excelize.Cell{StyleID: st.RowStyle, Value: n.KTP},
				excelize.Cell{StyleID: st.RowStyle, Value: n.NPWP},
				excelize.Cell{StyleID: st.RowStyle, Value: utils.ValTime(n.TanggalLahir)},
				excelize.Cell{StyleID: st.RowStyle, Value: n.TempatLahir},
				excelize.Cell{StyleID: st.RowStyle, Value: "-"},
			)
		} else {
			for j := 0; j < 7; j++ {
				row = append(row, excelize.Cell{StyleID: st.RowStyle, Value: "-"})
			}
		}

		if i < len(watchlist) {
			w := watchlist[i]
			row = append(row,
				excelize.Cell{StyleID: st.RowStyle, Value: w.Nama},
				excelize.Cell{StyleID: st.RowStyle, Value: w.KTP},
				excelize.Cell{StyleID: st.RowStyle, Value: w.NoPaspor},
				excelize.Cell{StyleID: st.RowStyle, Value: w.NPWP},
				excelize.Cell{StyleID: st.RowStyle, Value: utils.ValTime(w.TanggalLahir)},
				excelize.Cell{StyleID: st.RowStyle, Value: w.TempatLahir},
				excelize.Cell{StyleID: st.RowStyle, Value: "-"},
			)
		} else {
			for j := 0; j < 7; j++ {
				row = append(row, excelize.Cell{StyleID: st.RowStyle, Value: "-"})
			}
		}

		score := "0%"
		if ms.SimilarityScore != nil {
			score = fmt.Sprintf("%.0f%%", *ms.SimilarityScore*100)
		}
		row = append(row,
			excelize.Cell{StyleID: st.RowStyle, Value: score},
			excelize.Cell{StyleID: st.RowStyle, Value: ms.Reason},
			excelize.Cell{StyleID: st.RowStyle, Value: now.Format("02/01/2006")},
			excelize.Cell{StyleID: st.RowStyle, Value: now.Format("15:04:05")},
		)

		cell, _ := excelize.CoordinatesToCellName(1, rowIndex)
		if err := sw.SetRow(cell, row); err != nil {
			return err
		}
		rowIndex++
	}

	if err := sw.Flush(); err != nil {
		fmt.Printf("ERROR: Flush Stream Gagal: %v\n", err)
		return err
	}

	return nil
}

func (r *ReportGenerator) GenerateFromResults(results *models.MatchResults) error {
	nasabahMap := make(map[string][]models.MasterNasabah)
	watchlistMap := make(map[string][]models.MasterWatchlist)

	for _, result := range results.MatchResult {

		n := models.MasterNasabah{
			CifNumber:   utils.ValStr(result.CifNumber),
			NamaNasabah: utils.ValStr(result.CustomerName),
		}

		for _, detail := range results.MatchDetail {
			if detail.MatchingResultId != nil && *detail.MatchingResultId == result.Id {
				if detail.FieldName != nil && detail.CustomerValue != nil {
					switch *detail.FieldName {
					case "Nama":
						n.NamaNasabah = *detail.CustomerValue
					case "Ktp":
						n.KTP = detail.CustomerValue
					case "Npwp":
						n.NPWP = detail.CustomerValue
					case "NoPaspor":
						n.NoPaspor = detail.CustomerValue
					case "TempatLahir":
						n.TempatLahir = detail.CustomerValue
					case "TanggalLahir":
						n.TanggalLahir = utils.ParseDatePtrAny(*detail.CustomerValue)
					}
				}
			}
		}

		w := models.MasterWatchlist{
			ID:     utils.ValInt64(result.WatchlistId),
			Source: result.WatchlistSource,
		}

		sheetName := func() string {
			if result.WatchlistSource != nil {
				switch *result.WatchlistSource {
				case "MASTER_TERORIS":
					return "DTTOT"
				case "MASTER_WMD":
					return "WMD"
				case "MASTER_LOCAL_BLACKLIST":
					return "LOCAL_BLACKLIST"
				}
			}
			return "WATCHLIST"
		}()

		nasabahMap[sheetName] = append(nasabahMap[sheetName], n)
		watchlistMap[sheetName] = append(watchlistMap[sheetName], w)

		ms := MatchSectionData{
			SheetName:       sheetName,
			title:           "",
			NasabahCount:    len(nasabahMap[sheetName]),
			WatchlistCount:  len(watchlistMap[sheetName]),
			SimilarityScore: result.SimilarityScore,
			Reason:          utils.ValStr(result.Status),
		}

		if err := r.GenerateSectionStream(sheetName, nasabahMap[sheetName], watchlistMap[sheetName], ms); err != nil {
			return err
		}
	}

	return nil
}
