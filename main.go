package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/jolaz-dev/goes-back/config"
	"github.com/jolaz-dev/goes-back/handlers"
)

func main() {
	config := config.New()
	slog.Info("Starting application", "AppName", config.AppName, "Version", config.Version)

	http.HandleFunc("/", handlers.GoesBack(config))

	slog.Info("Running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
