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
}
