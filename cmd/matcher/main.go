package main

import (
	"BlackListWorker/internal/bootstrap"
	"log"
)

func main() {
	app := bootstrap.NewApp()
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

	/*
		service := monitor.NewMonitorService(
			2*time.Second, // interval monitoring
			monitor.WithCollector(monitor.NewThreadCollector()),
			monitor.WithCollector(monitor.NewCPUCollector(1*time.Second)),
			monitor.WithCollector(monitor.NewMemoryCollector()),
		)

		// jalanin monitoring di background
		go service.Start()

		// biar program gak langsung exit
		select {}
	*/
}
