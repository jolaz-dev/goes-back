package main

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/jolaz-dev/goes-back/config"
	"github.com/jolaz-dev/goes-back/internal"
	"github.com/jolaz-dev/goes-back/src/goesback"
)

func handleGoesBackRequest(config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := goesback.GoesBack(r, config)
		if errors.Is(err, internal.ErrUnprocessableEntity) {
			http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
			return
		} else if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		output, err := json.Marshal(response)
		if err != nil {
			slog.Error("Error marshaling JSON:", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(output)
	}
}

func main() {
	config := config.New()
	slog.Info("Starting application", "AppName", config.AppName, "Version", config.Version)

	http.HandleFunc("/", handleGoesBackRequest(config))

	slog.Info("Running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
