package services

import (
	"errors"
)

var ErrUnknownAlgorithm = errors.New("error unknow algorithm")

type MatchService interface {
	RunMatch() error
	MatchCIFWithDTTOT() error
	MatchCIFWithWMD() error
	MatchCIFWithLocalBlacklist() error
}

type MatchServiceImpl struct {
	NumWorker   int
	FieldWeight float64
}

type MatchServiceOption func(*MatchServiceImpl)

func WithNumOfWorker(nw int) MatchServiceOption {
	return func(msi *MatchServiceImpl) {
		msi.NumWorker = nw
	}
}

func WithFieldWeight(fw float64) MatchServiceOption {
	return func(msi *MatchServiceImpl) {
		msi.FieldWeight = fw
	}
}

func NewMatchService(opts ...MatchServiceOption) MatchService {
	svc := &MatchServiceImpl{}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func (s *MatchServiceImpl) RunMatch() error {
	return nil
}

func (s *MatchServiceImpl) MatchCIFWithDTTOT() error {
	return nil
}

func (s *MatchServiceImpl) MatchCIFWithDTTOTIndividu() error {
	return nil
}

func (s *MatchServiceImpl) MatchCIFWithDTTOTCorporate() error {
	return nil
}

func (s *MatchServiceImpl) MatchCIFWithWMD() error {
	return nil
}

func (s *MatchServiceImpl) MatchCIFWithLocalBlacklist() error {
	return nil
}

func (s *MatchServiceImpl) CalculateResult() (*float64, error) {
	return nil, nil
}
