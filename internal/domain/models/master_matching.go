package models

type MasterMatching struct {
	// use NEWSEQUENTIALID() in sqlserver for unique like uuid but more sequential for better indexing
	Id              string
	WatchlistSource string
	//Type <-> IsIndividual 1 for INDIVIDU, 0 for CORPORATE
	Type     bool
	IsActive bool
	// Name string
	// Description string
}
