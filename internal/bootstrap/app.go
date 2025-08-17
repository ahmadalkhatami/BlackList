package bootstrap

import (
	"BlackListWorker/config"
	"BlackListWorker/internal/db"
	"BlackListWorker/internal/domain/repository"
	"BlackListWorker/internal/services"
	"fmt"
)

type app struct{}

func NewApp() *app {
	return &app{}
}

func (a *app) Start() error {

	config.LoadEnv()
	dbConfig := config.Load()

	connector := db.NewSQLServerConnector(dbConfig.DBServer, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName)
	sqlDB, err := connector.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer sqlDB.Close()

	rows, err := sqlDB.Query("SELECT DB_NAME() AS CurrentDB")
	var currentDB string
	for rows.Next() {
		rows.Scan(&currentDB)
	}
	fmt.Println("Current Database:", currentDB)

	/*
		masterMatchingRepo := repository.NewSQLMasterMatchingRepository(sqlDB)
		masterMatchingConfigRepo := repository.NewSQLMasterMatchingConfigRepository(sqlDB)
		systemConfigRepo := repository.NewSQLSystemConfigRepository(sqlDB)

		configService := services.NewConfigService(masterMatchingRepo, masterMatchingConfigRepo, systemConfigRepo)

		matchConfigs, err := configService.GetJoinedMatchingConfig()
		if err != nil {
			return fmt.Errorf("error while getting configs: %w", err)
		}

		for _, j := range matchConfigs {
			fmt.Printf("ID=%s, MatchingID=%s, Source=%s, Field=%s, Weight=%.2f, Type=%t, Algorithm=%s\n",
				j.Id, j.MatchingId, j.WatchlistSource, j.FieldName, j.FieldWeight, j.Type, j.MatchingAlgorithm)
		}
	*/

	// masterNasabahRepo := repository.NewSQLMasterNasabahRepository(sqlDB)
	// masterNasabahService := services.NewMasterNasabah(masterNasabahRepo)

	// masterNasabah, err := masterNasabahService.Load()

	// for _, j := range masterNasabah {
	// 	fmt.Printf("ID: %s, CIF: %s, Nama: %s, Tempat Lahir: %s, Tanggal Lahir: %s, KTP: %s, NPWP: %s, No Paspor: %s, Status: %s, Created Date: %s",
	// 		j.Id, j.CIFNumber, j.NamaNasabah, j.TanggalLahir, j.TanggalLahir, j.KTP, j.NPWP, j.NoPaspor, j.StatusNasabah, j.CreatedAt)
	// 	fmt.Println()
	// }

	masterDTTOTRepo := repository.NewSQLMasterTerorisRepository(sqlDB)
	masterWMDRepo := repository.NewSQLMasterWMDRepository(sqlDB)
	masterLocalBalcklistRepo := repository.NewSQLMasterLocalBlacklistRepository(sqlDB)

	watchlistService := services.NewWatchlistService(masterDTTOTRepo, masterWMDRepo, masterLocalBalcklistRepo)

	masterDTTOT, err := watchlistService.LoadDTTOT()
	masterWMD, err := watchlistService.LoadWMD()
	masterLocalBlacklist, err := watchlistService.LoadLocalBlacklist()

	for _, j := range masterDTTOT {
		fmt.Printf("ID: %s, Nama: %s, Alias1: %s, Alias2: %s, Alias3: %s, Alias4: %s, Type: %s, Tempat Lahir: %s, Tanggal Lahir: %s, KTP: %s, NPWP: %s, NoPaspor: %s, Created Date: %s, Status: %s",
			j.Id, j.Nama, j.Alias1, j.Alias2, j.Alias3, j.Alias4, j.Type, j.TempatLahir, j.TanggalLahir, j.KTP, j.NPWP, j.NoPaspor, j.CreatedAt, j.IsActive)
		fmt.Println()
	}

	for _, j := range masterWMD {
		fmt.Printf("ID: %s, Nama: %s, Alias1: %s, Alias2: %s, Alias3: %s, Alias4: %s, Alias5: %s, Alias6: %s, Alias7: %s, Alias8: %s, Alias9: %s, Alias10: %s, Type: %s, Tempat Lahir: %s, Tanggal Lahir: %s, KTP: %s, NPWP: %s, NoPaspor: %s, Created Date: %s, Status: %s",
			j.Id, j.Nama, j.Alias1, j.Alias2, j.Alias3, j.Alias4, j.Alias5, j.Alias6, j.Alias7, j.Alias8, j.Alias9, j.Alias10, j.Type, j.TempatLahir, j.TanggalLahir, j.KTP, j.NPWP, j.NoPaspor, j.CreatedAt, j.IsActive)
		fmt.Println()
	}

	for _, j := range masterLocalBlacklist {
		fmt.Printf("ID: %s, Nama: %s, Alias1: %s, Alias2: %s, Alias3: %s, Alias4: %s, Type: %s, Tempat Lahir: %s, Tanggal Lahir: %s, KTP: %s, NPWP: %s, NoPaspor: %s, Created Date: %s, Status: %s",
			j.Id, j.Nama, j.Alias1, j.Alias2, j.Alias3, j.Alias4, j.Type, j.TempatLahir, j.TanggalLahir, j.KTP, j.NPWP, j.NoPaspor, j.CreatedAt, j.IsActive)
		fmt.Println()
	}

	return nil
}
