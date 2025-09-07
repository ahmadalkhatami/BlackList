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

func calculateResult() error {
	return nil
}
