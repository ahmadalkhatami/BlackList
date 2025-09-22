package services

import (
	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repositories"
	"context"
)

type MatchingResultService interface {
	SetBatchID(int64)
	SetResults([]models.MatchingResult)
	SetDetails([]models.MatchingDetail)
	Save(ctx context.Context) error
}

type MatchingResultImpl struct {
	MatchingResultsRepo        repositories.MatchingResultRepository
	MatchingResultsDetailsRepo repositories.MatchingDetailsRepository
	MatchingResults            []models.MatchingResult
	MatchingResultsDetails     []models.MatchingDetail
	BatchID                    *int64
}

/* //With no options
func MewMatchingResultService(
	matchingResultsRepo repositories.MatchingResultRepository,
	matchingResultsDetailsRepo repositories.MatchingDetailsRepository,
	matchingResults []models.MatchingResult,
	matchingResultsDetails []models.MatchingDetail,
	batchID *int64,
) *MatchingResultImpl {
	return &MatchingResultImpl{
		MatchingResultsRepo:        matchingResultsRepo,
		MatchingResultsDetailsRepo: matchingResultsDetailsRepo,
		MatchingResults:            matchingResults,
		MatchingResultsDetails:     matchingResultsDetails,
		BatchID:                    batchID,
	}
} */

type MatchingResultOption func(*MatchingResultImpl)

func WithMatchingResults(results []models.MatchingResult) MatchingResultOption {
	return func(m *MatchingResultImpl) {
		m.MatchingResults = results
	}
}

func WithMatchingDetails(details []models.MatchingDetail) MatchingResultOption {
	return func(m *MatchingResultImpl) {
		m.MatchingResultsDetails = details
	}
}

func WithBatchID(batchID int64) MatchingResultOption {
	return func(m *MatchingResultImpl) {
		m.BatchID = &batchID
	}
}

func NewMatchingResultService(
	resultsRepo repositories.MatchingResultRepository,
	detailsRepo repositories.MatchingDetailsRepository,
	opts ...MatchingResultOption,
) *MatchingResultImpl {
	svc := &MatchingResultImpl{
		MatchingResultsRepo:        resultsRepo,
		MatchingResultsDetailsRepo: detailsRepo,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (m *MatchingResultImpl) SetBatchID(batchID int64) {
	m.BatchID = &batchID
}

func (m *MatchingResultImpl) SetResults(results []models.MatchingResult) {
	m.MatchingResults = results
}

func (m *MatchingResultImpl) SetDetails(details []models.MatchingDetail) {
	m.MatchingResultsDetails = details
}

func (m *MatchingResultImpl) Save(ctx context.Context) error {

	if m.MatchingResults == nil || m.MatchingResultsDetails == nil {
		return nil
	}

	err1 := m.MatchingResultsRepo.SaveBatch(ctx, *m.BatchID, m.MatchingResults)
	if err1 != nil {
		return err1
	}

	err2 := m.MatchingResultsDetailsRepo.SaveBatch(ctx, *m.BatchID, m.MatchingResultsDetails)
	if err2 != nil {
		return err2
	}

	return nil
}

/*
// Pemakaian dengan WithXxx
svc := services.NewMatchingResultService(
    resultRepo,
    detailRepo,
    services.WithBatchID(123),
    services.WithMatchingResults(results),
    services.WithMatchingDetails(details),
)
err := svc.Save(ctx)

//Pemakaian dengan SetXxx
svc := services.NewMatchingResultService(resultRepo, detailRepo)
svc.SetBatchID(123)
svc.SetResults(results)
svc.SetDetails(details)
err := svc.Save(ctx)

*/
