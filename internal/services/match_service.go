package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repositories"
	"BlackListWorker/internal/utils"
	"context"
	"fmt"
)

type MatchResults struct {
	MatchResult []models.MatchingResult
	MatchDetail []models.MatchingDetail
	BatchID     *int64
}

type MatchService interface {
	RunMatch(ctx context.Context) error
	GetResults() *MatchResults
}

type Matcher interface {
	Source() string
	Match(ctx context.Context) (*MatchResults, error)
}

type MatchServiceImpl struct {
	matchers      []Matcher
	results       *MatchResults
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
	var allResults MatchResults
	fmt.Printf("Matchers Count Queue : %v \n", len(s.matchers))

	var batchProcessing models.BatchProcessing

	for _, m := range s.matchers {
		res, err := m.Match(ctx)
		if err != nil {
			return err
		}

		if len(res.MatchResult) == 0 || res.MatchResult[0].BatchId == nil {
			fmt.Printf("⚠️  Matcher %s mengembalikan BatchId nil, dilewati\n", m.Source())
			continue
		}

		countRes := len(res.MatchResult)
		countDet := len(res.MatchDetail)
		totalCount := countRes + countDet
		status := "complete"

		batchProcessing.Id = *res.MatchResult[0].BatchId
		batchProcessing.Status = utils.Ptr(status)
		batchProcessing.ProcessedRecords = utils.Ptr(countRes)
		batchProcessing.MatchedRecords = utils.Ptr(countDet)
		batchProcessing.TotalRecords = utils.Ptr(totalCount)
		batchProcessing.ErrorMessage = nil

		if err := s.batchRepo.UpdateStatus(&batchProcessing); err != nil {
			return err
		}

		newResultService := s.resultService
		newResultService.SetResults(res.MatchResult)
		newResultService.SetDetails(res.MatchDetail)
		if err := newResultService.Save(ctx); err != nil {
			return err
		}

		allResults.MatchResult = append(allResults.MatchResult, res.MatchResult...)
		allResults.MatchDetail = append(allResults.MatchDetail, res.MatchDetail...)

	}

	s.results = &allResults
	return nil
}

func (s *MatchServiceImpl) GetResults() *MatchResults {
	return s.results
}
