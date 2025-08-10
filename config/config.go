package config

import "os"

type AppConfig struct {
	DBServer   string
	DBUser     string
	DBPassword string
	DBName     string
	NumWorker  int
}

func Load() AppConfig {
	return AppConfig{
		DBServer:   os.Getenv("DB_SERVER"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		NumWorker:  0, // DEFAULT
	}
}
