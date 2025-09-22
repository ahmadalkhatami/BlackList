package services

import (
	"BlackListWorker/internal/domain/models"
	"context"
)

type MatchResults struct {
	MatchResult []models.MatchingResult
	MatchDetail []models.MatchingDetail
	BatchID     int64
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
	matchers []Matcher
	results  *MatchResults
}

func NewMatchService(matchers ...Matcher) *MatchServiceImpl {
	return &MatchServiceImpl{matchers: matchers}
}

func (s *MatchServiceImpl) RunMatch(ctx context.Context) error {
	var allResults MatchResults

	for _, m := range s.matchers {
		res, err := m.Match(ctx)
		if err != nil {
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
