package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/xuri/excelize/v2"
)

func main() {
	connString := "server=100.126.31.12;user id=sa;password=yourStrong(!)Password;database=belajar-migration"

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error opening connection: ", err.Error())
	}
	defer db.Close()

	query := `
SELECT 
    mr.WatchlistSource,
    mn.CifNumber     AS NomorCIFNasabah,
    mn.NamaNasabah   AS NamaNasabah,
    mn.Ktp           AS KTPNasabah,
    mn.NoPaspor      AS PassportNasabah,
    mn.Npwp          AS NPWPNasabah,
    mn.TanggalLahir  AS TanggalLahirNasabah,
    mn.TempatLahir   AS TempatLahirNasabah,

    mt.Nama          AS NamaTeroris,
    mt.Ktp           AS KTPTeroris,
    mt.NoPaspor      AS PassportTeroris,
    mt.Npwp          AS NPWPTeroris,
    mt.TanggalLahir  AS TanggalLahirTeroris,
    mt.TempatLahir   AS TempatLahirTeroris,

    mr.SimilarityScore,
    r.Reason,
    mr.ProcessDate AS Tanggal,
    CONVERT(VARCHAR(8), mr.ProcessTime, 108) AS Waktu
FROM MATCHING_RESULTS mr WITH (NOLOCK)
INNER JOIN MASTER_NASABAH mn WITH (NOLOCK) 
    ON mr.CifNumber = mn.CifNumber
LEFT JOIN MASTER_TERORIS mt WITH (NOLOCK) 
    ON mr.WatchlistId = mt.ID
LEFT JOIN (
    SELECT 
        md.MatchingResultId,
        STRING_AGG(CONCAT(md.FieldName, ' ', CAST(md.FieldScore * 100 AS VARCHAR(10)), '%'), '; ') AS Reason
    FROM MATCHING_DETAILS md
    WHERE md.FieldScore > 0
    GROUP BY md.MatchingResultId
) r ON r.MatchingResultId = mr.Id
WHERE mr.WatchlistId = 1051
`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Error running query: ", err)
	}
	defer rows.Close()

	// Buat Excel baru
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.NewSheet(sheet)

	// Header
	headers := []string{
		"WatchlistSource", "NomorCIFNasabah", "NamaNasabah", "KTPNasabah", "PassportNasabah", "NPWPNasabah",
		"TanggalLahirNasabah", "TempatLahirNasabah",
		"NamaTeroris", "KTPTeroris", "PassportTeroris", "NPWPTeroris",
		"TanggalLahirTeroris", "TempatLahirTeroris",
		"SimilarityScore", "Reason", "Tanggal", "Waktu",
	}
	for i, h := range headers {
		col := string(rune('A' + i))
		f.SetCellValue(sheet, fmt.Sprintf("%s1", col), h)
	}

	// Data
	rowIndex := 2
	cols := len(headers)
	for rows.Next() {
		values := make([]interface{}, cols)
		pointers := make([]interface{}, cols)
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			log.Fatal("Error scanning row: ", err)
		}

		for i, val := range values {
			col := string(rune('A' + i))
			f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, rowIndex), val)
		}
		rowIndex++
	}

	// Simpan ke file lokal
	if err := f.SaveAs("output.xlsx"); err != nil {
		log.Fatal("Error saving file: ", err)
	}

	fmt.Println("✅ File berhasil dibuat: output.xlsx")
}
