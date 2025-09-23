package services

import (
	"context"
	"fmt"
	"sync/atomic"

	"BlackListWorker/internal/domain/models"
	"BlackListWorker/internal/domain/repositories"
)

type IDService struct {
	resultRepo repositories.MatchingResultRepository
	detailRepo repositories.MatchingDetailsRepository
	batchRepo  repositories.BatchProcessingRepository

	resultIDCounter int64
	detailIDCounter int64
	batchIDCounter  int64
}

func NewIDService(
	resultRepo repositories.MatchingResultRepository,
	detailRepo repositories.MatchingDetailsRepository,
	batchRepo repositories.BatchProcessingRepository,
) *IDService {
	return &IDService{
		resultRepo: resultRepo,
		detailRepo: detailRepo,
		batchRepo:  batchRepo,
	}
}

func (s *IDService) InitAtomicIDs(ctx context.Context) error {
	lastResultID, err := s.resultRepo.GetLastId(ctx)
	if err != nil {
		return err
	}

	lastDetailID, err := s.detailRepo.GetLastId(ctx)
	if err != nil {
		return err
	}

	// lastBatchID, err := s.batchRepo.GetLastId(ctx)
	// if err != nil {
	// 	return err
	// }

	atomic.StoreInt64(&s.resultIDCounter, lastResultID)
	atomic.StoreInt64(&s.detailIDCounter, lastDetailID)
	// atomic.StoreInt64(&s.batchIDCounter, lastBatchID)

	fmt.Printf("Initialized atomic IDs: result=%d, detail=%d, batch=%d\n",
		s.resultIDCounter, s.detailIDCounter, s.batchIDCounter)
	return nil
}

func (s *IDService) GenerateResultID() int64 {
	return atomic.AddInt64(&s.resultIDCounter, 1)
}

func (s *IDService) GenerateDetailID() int64 {
	return atomic.AddInt64(&s.detailIDCounter, 1)
}

// func (s *IDService) GenerateBatchID(batch *models.BatchProcessing) (int64, error) {
// 	batchID, err := s.batchRepo.Create(batch)
// 	if err != nil {
// 		fallback := atomic.AddInt64(&s.batchIDCounter, 1)
// 		return fallback, err
// 	}
// 	atomic.StoreInt64(&s.batchIDCounter, batchID)
// 	return batchID, nil
// }

func (s *IDService) GenerateBatchID(batch *models.BatchProcessing) (int64, error) {
	batchID, err := s.batchRepo.Create(batch)
	if err != nil {
		panic(fmt.Sprintf("failed to create batch: %v", err))
	}
	atomic.StoreInt64(&s.batchIDCounter, batchID)
	return batchID, nil
}
