package config

import (
	"os"
)

const DefaultPort = "6342"

type Config struct {
	Port string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	return &Config{
		Port: port,
	}
}
