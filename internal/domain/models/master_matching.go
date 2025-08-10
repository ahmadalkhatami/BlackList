package models

type MasterMatching struct {
	Id string // use NEWSEQUENTIALID() in sqlserver for unique like uuid but more sequential for better indexing
	// Name string
	WatchlistSource string
	Type            bool // IsIndividual 1 for Individu, 0 for Corporate
	// Description string
	IsActive bool
}
