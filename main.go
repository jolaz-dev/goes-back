package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/jolaz-dev/goes-back/src/goesback"
)

func main() {
	http.HandleFunc("/", goesback.GoesBack)

	slog.Info("Running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
