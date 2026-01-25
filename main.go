package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

func marco(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Polo!")
}

func main() {
	http.HandleFunc("/", marco)

	slog.Info("Running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
