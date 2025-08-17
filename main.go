package main

import (
	"fmt"
	"log"

	"BlackListWorker/config"
	"BlackListWorker/services"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, skipping...")
	}
}

func main() {
	fmt.Println("🚀 BlackListWorker started...")

	// Step 1: Koneksi database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke DB: %v", err)
	}
	defer db.Close()

	// Step 2: Ambil batch_id berikutnya
	batchID, err := services.GetNextBatchID(db)
	if err != nil {
		log.Fatalf("❌ Gagal ambil batch ID: %v", err)
	}
	fmt.Printf("📌 Menggunakan batch_id: %d\n", batchID)

	// Step 3: Load data dari DB
	fmt.Println("📦 Memuat data CIF dan watchlist...")
	cifList, err := services.LoadCIF(db)
	// fmt.Println(cifList);
	if err != nil {
		log.Fatalf("❌ Gagal load CIF: %v", err)
	}

	terorisList, err := services.LoadMasterTeroris(db)
	if err != nil {
		log.Fatalf("❌ Gagal load MASTER_TERORIS: %v", err)
	}

	configList, err := services.LoadMatchingConfig(db)
	if err != nil {
		log.Fatalf("❌ Gagal load MATCHING_CONFIG: %v", err)
	}

	// Step 4: Proses Matching
	fmt.Printf("🔍 Memproses %d data CIF terhadap %d data watchlist...\n", len(cifList), len(terorisList))
	results, detailsMap := services.MatchCIFWithTeroris(db, cifList, terorisList, configList)
	// fmt.Print("config :", configList)
	// fmt.Print("cif :", cifList)
	// fmt.Print("terorisList :", terorisList)
	// Set batch_id untuk semua hasil
	for i := range results {
		results[i].BatchID = batchID
	}

	fmt.Printf("✅ Ditemukan %d hasil match\n", len(results))

	// Step 5: Simpan hasil ke database
	if len(results) > 0 {
		err = services.InsertMatchingResults(db, results, detailsMap)
		if err != nil {
			log.Fatalf("❌ Gagal insert MATCHING_RESULTS & MATCHING_DETAILS: %v", err)
		}
		fmt.Printf("💾 %d hasil match berhasil disimpan ke database (batch_id: %d)\n", len(results), batchID)
	} else {
		fmt.Println("ℹ️ Tidak ada hasil match yang disimpan.")
	}

	// Step 6: Cetak hasil
	for _, result := range results {
		fmt.Printf("[MATCH] Batch: %d | CIF: %s | Score: %.2f | Watchlist ID: %d\n",
			result.BatchID, result.CIFNumber, result.SimilarityScore, result.WatchlistID)
	}
}
