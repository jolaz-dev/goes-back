package config

import (
	"log/slog"
	"os"

	"github.com/google/uuid"
)

type Config struct {
	AppName  string
	Hostname string
	Version  string
}

func New() *Config {
	config := Config{
		AppName: "GoesBack",
		Version: "0.2.0", // x-release-please-version
	}

	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("Error getting hostname", "error", err)
		hostname = uuid.New().String()
	}
	config.Hostname = hostname

	return &config
}
