package config

import (
	"fmt"
	"os"
)

type Config struct {
	DB *DB
}

type DB struct {
	DSN        string
	DriverName string
}

func LoadConfig() *Config {
	driverName := os.Getenv("DATABASE_DRIVER_NAME")

	if driverName == "" {
		fmt.Println("DATABASE_DRIVER_NAME is not set")
	}

	dsn := os.Getenv("DATABASE_DSN")

	if dsn == "" {
		fmt.Println("DATABASE_DSN is not set")
	}

	cfg := Config{
		DB: &DB{
			DSN:        dsn,
			DriverName: driverName,
		},
	}

	return &cfg
}
