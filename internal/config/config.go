package config

import (
	"flag"
	"os"
)

type Config struct {
	runAddress string
}

func New() *Config {
	var flagRunAddress string

	flag.StringVar(&flagRunAddress, "a", ":8080", "address and port to run server")

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		flagRunAddress = envRunAddress
	}

	return &Config{
		runAddress: flagRunAddress,
	}
}
