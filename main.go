package main

import (
	"fmt"
	"log"

	"BlackListWorker/config"
	"BlackListWorker/services"
)

func main() {
	fmt.Println("🚀 BlackListWorker started...")

	// Step 1: Koneksi database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke DB: %v", err)
	}
	defer db.Close()

	// Step 2: Load data dari DB
	fmt.Println("📦 Memuat data CIF dan watchlist...")
	cifList, err := services.LoadCIF(db)
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

	// Step 3: Proses Matching
	fmt.Printf("🔍 Memproses %d data CIF terhadap %d data watchlist...\n", len(cifList), len(terorisList))
	results := services.MatchCIFWithTeroris(cifList, terorisList, configList)

	fmt.Printf("✅ Ditemukan %d hasil match\n", len(results))

	// Step 4: Cetak hasil (atau simpan ke DB/Excel)
	for _, result := range results {
		fmt.Printf("[MATCH] CIF: %s | Score: %.2f | Watchlist ID: %d\n",
			result.CIFNumber, result.SimilarityScore, result.WatchlistID)
	}

	// TODO: Simpan ke database atau ekspor ke Excel
}
