package config

import (
	"os"
)

const DefaultPort = "80"
const DefaultBackendPort = "6342"

type Config struct {
	Port        string
	BackendPort string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	backendPort := os.Getenv("BACKENDPORT")
	if backendPort == "" {
		backendPort = DefaultBackendPort
	}

	return &Config{
		Port:        port,
		BackendPort: backendPort,
	}
}
