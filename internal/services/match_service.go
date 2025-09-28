package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repositories"
	"BlackListWorker/internal/report"
	"BlackListWorker/internal/utils"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type MatchService interface {
	RunMatch(ctx context.Context) error
	GetResults() *models.MatchResults
}

type Matcher interface {
	Source() string
	Match(ctx context.Context) (*models.MatchResults, error)
	MatchType(ctx context.Context, t string) (*models.MatchResults, error)
}

type matchTask struct {
	matcher    Matcher
	masterType string
}

type MatchServiceImpl struct {
	matchers      []Matcher
	results       *models.MatchResults
	batchRepo     repositories.BatchProcessingRepository
	resultService MatchingResultService
}

func NewMatchService(batchRepo repositories.BatchProcessingRepository, resultService MatchingResultService, matchers ...Matcher) *MatchServiceImpl {
	return &MatchServiceImpl{
		matchers:      matchers,
		batchRepo:     batchRepo,
		resultService: resultService,
	}
}

func (s *MatchServiceImpl) RunMatch(ctx context.Context) error {
	startTotal := time.Now()
	log.Println("🚀 RunMatch dimulai")

	var allResults models.MatchResults
	var batchProcessing models.BatchProcessing

	// --- buat tasks ---
	var tasks []matchTask
	for _, m := range s.matchers {
		for _, t := range []string{"Individu", "Corporate"} {
			tasks = append(tasks, matchTask{matcher: m, masterType: t})
		}
	}

	resultsChan := make(chan matchResultWithSource, len(tasks))
	errChan := make(chan error, len(tasks))

	// --- cek healthy ---
	healthy := systemHealthy()
	maxWorkers := 4 // maksimal goroutine paralel, bisa set via env/config
	if !healthy {
		log.Println("⚠️ Resource terbatas → fallback ke serial matching")
		maxWorkers = 1
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxWorkers) // semaphore untuk bounded concurrency

	// --- jalankan matcher dengan bounded concurrency ---
	for _, task := range tasks {
		wg.Add(1)
		sem <- struct{}{}
		go func(m Matcher, t string) {
			defer wg.Done()
			defer func() { <-sem }()

			start := time.Now()
			res, err := m.MatchType(ctx, t)
			elapsed := time.Since(start)
			log.Printf("⏱️ Matcher %s type %s selesai dalam %v", m.Source(), t, elapsed)
			if err != nil {
				errChan <- err
				return
			}
			resultsChan <- matchResultWithSource{
				source: m.Source(),
				res:    res,
			}
		}(task.matcher, task.masterType)
	}

	// --- tunggu semua selesai ---
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	// --- CSV semaphore ---
	csvSem := make(chan struct{}, 2) // max 2 CSV simultan

	// --- proses hasil ---
	for item := range resultsChan {
		res := item.res
		source := item.source

		if len(res.MatchResult) == 0 || res.MatchResult[0].BatchId == nil {
			log.Printf("⚠️ Matcher %s type %s mengembalikan BatchId nil, dilewati\n",
				source, utils.ValStr(res.MasterType))
			continue
		}

		// Update hasil global
		allResults.MasterType = res.MasterType
		allResults.BatchID = res.BatchID
		allResults.MatchResult = append(allResults.MatchResult, res.MatchResult...)
		allResults.MatchDetail = append(allResults.MatchDetail, res.MatchDetail...)

		// Simpan hasil
		newResultService := s.resultService
		newResultService.SetResults(res.MatchResult)
		newResultService.SetDetails(res.MatchDetail)
		if err := newResultService.Save(ctx); err != nil {
			batchProcessing.Id = utils.ValInt64(res.MatchResult[0].BatchId)
			batchProcessing.Status = utils.Ptr("failed")
			batchProcessing.ErrorMessage = utils.Ptr(err.Error())
			_ = s.batchRepo.UpdateStatus(&batchProcessing)
			return err
		}

		// Generate CSV dengan semaphore
		filepath := fmt.Sprintf(
			"output/matcher/report_%s_%s_%s.csv",
			source,
			utils.ValStr(res.MasterType),
			time.Now().Format("20060102_150405"),
		)
		csvSem <- struct{}{}
		go func(res *models.MatchResults, path, matcherType string) {
			defer func() { <-csvSem }()
			reportGen := report.NewReportGenerator()
			if err := reportGen.GenerateFromResultsCSVStream(path, res, matcherType); err != nil {
				log.Printf("❌ Gagal generate CSV %s: %v", path, err)
			} else {
				log.Printf("✅ CSV generated: %s", path)
			}
		}(res, filepath, source)

		// Update batchProcessing
		countRes := len(res.MatchResult)
		countDet := len(res.MatchDetail)
		batchProcessing.Id = utils.ValInt64(res.MatchResult[0].BatchId)
		batchProcessing.Status = utils.Ptr("complete")
		batchProcessing.ProcessedRecords = utils.Ptr(countRes)
		batchProcessing.MatchedRecords = utils.Ptr(countDet)
		batchProcessing.TotalRecords = utils.Ptr(countRes + countDet)
		batchProcessing.ErrorMessage = nil
		batchProcessing.FilePath = utils.Ptr(filepath)

		if err := s.batchRepo.UpdateStatus(&batchProcessing); err != nil {
			batchProcessing.Status = utils.Ptr("failed")
			batchProcessing.ErrorMessage = utils.Ptr(err.Error())
			_ = s.batchRepo.UpdateStatus(&batchProcessing)
			return err
		}
		// log.Printf("✅ Report generated for matcher=%s, masterType=%v, batchId=%v",
		// 	source, utils.ValStr(res.MasterType), utils.ValInt64(res.BatchID))
	}

	// --- check error ---
	if len(errChan) > 0 {
		for err := range errChan {
			return err
		}
	}

	log.Printf("⏱️ Total RunMatch finished in %v", time.Since(startTotal))
	s.results = &allResults
	return nil
}

func (s *MatchServiceImpl) GetResults() *models.MatchResults {
	return s.results
}
