package services

import (
	"errors"
)

var ErrUnknownAlgorithm = errors.New("Unknown Algorithm")

type MatchService interface {
	RunMatch() error
	MatchCIFWithDTTOT() error
	MatchCIFWithWMD() error
	MatchCIFWithLocalBlacklist() error
}

type MatchServiceImpl struct {
	NumWorker int
}

type MatchServiceOption func(*MatchServiceImpl)

func NewMatchService(numWorker int) MatchService {

	if numWorker < 1 {
		numWorker = 1
	}

	return &MatchServiceImpl{
		NumWorker: numWorker,
	}
}

func (s *MatchServiceImpl) RunMatch() error {
	return nil
}

func (s *MatchServiceImpl) MatchCIFWithDTTOT() error {
	return nil
}

func (s *MatchServiceImpl) MatchCIFWithWMD() error {
	return nil
}

func (s *MatchServiceImpl) MatchCIFWithLocalBlacklist() error {
	return nil
}
