package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	wd, _ := os.Getwd()
	envPath := filepath.Join(wd, ".env")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file at %s: %v", envPath, err)
	}
}
