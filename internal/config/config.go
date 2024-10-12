package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

func NewConfig() *Config {
	var flagRunAddress, flagDatabaseURI, flagAccrualSystemAddress string

	flag.StringVar(&flagRunAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagDatabaseURI, "d", "host=localhost port=5433 user=postgres password=password dbname=gophermart sslmode=disable", "database URI")
	flag.StringVar(&flagAccrualSystemAddress, "r", "", "accrual system address")

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		flagRunAddress = envRunAddress
	}
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		flagDatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		flagAccrualSystemAddress = envAccrualSystemAddress
	}

	return &Config{
		RunAddress:           flagRunAddress,
		DatabaseURI:          flagDatabaseURI,
		AccrualSystemAddress: flagAccrualSystemAddress,
	}
}
