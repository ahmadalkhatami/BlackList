package config

import "os"

type DBConfig struct {
	DBServer   string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() DBConfig {
	return DBConfig{
		DBServer:   os.Getenv("DB_SERVER"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}
}
