package main

import (
	"BlackListWorker/config"
	"BlackListWorker/internal/bootstrap"
	"log"
	"os"
)

func main() {

	config.LoadEnv()

	app := bootstrap.NewApp()
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

	// Exit Code 0 untuk menandakan proses berhasil
	// Exit Code 1 untuk menandakan proses gagal
	// Exit Code 2 (golang memberikan exit code untuk panic)
	os.Exit(0)
}
