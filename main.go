package main

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	"BlackListWorker/config"
// 	"BlackListWorker/services"

// 	"github.com/joho/godotenv"
// )

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Println("No .env file found, skipping...")
// 	}
// }

// func main() {
// 	fmt.Println("🚀 BlackListWorker started...")

// 	// Step 1: Koneksi database
// 	db, err := config.ConnectDB()
// 	if err != nil {
// 		log.Fatalf("❌ Gagal koneksi ke DB: %v", err)
// 	}
// 	defer db.Close()

// 	// Step 2: Ambil batch_id berikutnya
// 	batchID, err := services.GetNextBatchID(db)
// 	if err != nil {
// 		log.Fatalf("❌ Gagal ambil batch ID: %v", err)
// 	}
// 	fmt.Printf("📌 Menggunakan batch_id: %d\n", batchID)

// 	// Step 3: Load data dari DB
// 	startLoad := time.Now()
// 	fmt.Println("📦 Memuat data CIF dan watchlist...")

// 	cifList, err := services.LoadCIF(db)
// 	if err != nil {
// 		log.Fatalf("❌ Gagal load CIF: %v", err)
// 	}

// 	// ambil watchlist dari 3 tabel
// 	watchlistList, err := services.LoadAllWatchlists(db)
// 	if err != nil {
// 		log.Fatalf("❌ Gagal load Watchlist: %v", err)
// 	}

// 	configList, err := services.LoadMatchingConfig(db)
// 	if err != nil {
// 		log.Fatalf("❌ Gagal load MATCHING_CONFIG: %v", err)
// 	}
// 	fmt.Printf("⏱️ Load data selesai dalam %s\n", time.Since(startLoad))

// 	// Step 4: Proses Matching
// 	startMatch := time.Now()
// 	fmt.Printf("🔍 Memproses %d data CIF terhadap %d data watchlist...\n", len(cifList), len(watchlistList))

// 	// results, detailsMap := services.MatchCIFWithTeroris(db, cifList, terorisList, configList)

// 	results, detailsMap := services.MatchCIFAll(db, cifList, watchlistList, configList)

// 	// Set batch_id untuk semua hasil
// 	for i := range results {
// 		results[i].BatchID = batchID
// 	}
// 	fmt.Printf("✅ Ditemukan %d hasil match\n", len(results))
// 	fmt.Printf("⏱️ Matching selesai dalam %s\n", time.Since(startMatch))

// 	// Step 5: Simpan hasil ke database
// 	if len(results) > 0 {
// 		startInsert := time.Now()
// 		err = services.InsertMatchingResults(db, results, detailsMap)
// 		if err != nil {
// 			log.Fatalf("❌ Gagal insert MATCHING_RESULTS & MATCHING_DETAILS: %v", err)
// 		}
// 		fmt.Printf("💾 %d hasil match berhasil disimpan ke database (batch_id: %d)\n", len(results), batchID)
// 		fmt.Printf("⏱️ Insert DB selesai dalam %s\n", time.Since(startInsert))
// 	} else {
// 		fmt.Println("ℹ️ Tidak ada hasil match yang disimpan.")
// 	}

// 	// Step 6: Cetak hasil
// 	for _, result := range results {
// 		fmt.Printf("[MATCH] Batch: %d | CIF: %s | Score: %.2f | Watchlist ID: %d\n",
// 			result.BatchID, result.CIFNumber, result.SimilarityScore, result.WatchlistID)
// 	}
// }
