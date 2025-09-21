package main

import (
	"BlackListWorker/config"
	"BlackListWorker/internal/bootstrap"
	"log"
)

func main() {

	config.LoadEnv()

	app := bootstrap.NewApp()
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
