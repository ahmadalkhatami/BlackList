package worker

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/similarity"
	"math"
	"sync"
	"time"
)

type WorkerManager struct {
	NumWorker int
	Threshold float64
	Calc      similarity.Calculator
}

func (wm *WorkerManager) chunkify(data []models.MasterNasabah, chunks int) [][]models.MasterNasabah {
	if chunks <= 0 {
		chunks = 1
	}
	size := int(
		math.Ceil(
			float64(len(data)) / float64(chunks)))
	var out [][]models.MasterNasabah
	for i := 0; i < len(data); i += size {
		end := i + size
		if end > len(data) {
			end = len(data)
		}
		out = append(out, data[i:end])
	}
	return out
}

// Start Worker
func (wm *WorkerManager) Run(tableA, tableB []models.MasterNasabah) []models.MatchingResult {

	jobs := make(chan []models.MasterNasabah, wm.NumWorker)
	results := make(chan []models.MatchingResult, wm.NumWorker)
	var wg sync.WaitGroup

	for i := 0; i < wm.NumWorker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range jobs {
				var out []models.MatchingResult
				for _, a := range chunk {
					for _, b := range tableB {
						score := wm.Calc.Calculate(a.NamaNasabah, b.NamaNasabah)
						//if score >= wm.Threshold {
						out = append(out, models.MatchingResult{
							ID:              a.Id,
							BatchID:         123, // make function for generate batching ID
							CIFNumber:       a.CIFNumber,
							CustomerName:    a.NamaNasabah,
							WatchlistID:     123, //will get from master_mathing and master_matching_config
							WatchlistSource: "",
							SimilarityScore: score,
							Status:          "",
							ProcessDate:     time.Now(),
							ProcessTime:     time.Now(),
							CreatedAt:       time.Now(),
						})
						//}
					}
				}

				results <- out
			}
		}()
	}

	// Feed Jobs
	go func() {
		for _, c := range wm.chunkify(tableA, wm.NumWorker) {
			jobs <- c
		}
		close(jobs)
	}()

	// Collect
	go func() {
		wg.Wait()
		close(results)
	}()

	var final []models.MatchingResult
	for r := range results {
		final = append(final, r...)
	}

	return final
}
