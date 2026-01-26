package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/jolaz-dev/goes-back/config"
	"github.com/jolaz-dev/goes-back/goesback"
	"github.com/jolaz-dev/goes-back/internal"
	"github.com/jolaz-dev/goes-back/utils/compression"
)

const (
	headerModifiersPrefix  = "x-goesback-header-"
	headerStatusModifier   = "x-goesback-status"
	headerBodyModifier     = "x-goesback-body"
	headerGoesBackResponse = "X-GoesBack-Response"
)

func GoesBack(config *config.Config) http.HandlerFunc {
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

		checkForHeaderModifiers(w, r)

		status := checkForStatusModifier(r)

		bodyModifier := r.Header.Get(headerBodyModifier)
		if bodyModifier != "" {
			w.Header().Add(headerGoesBackResponse, string(output))
			compression.Compress(w, r, []byte(bodyModifier), status)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		compression.Compress(w, r, output, status)
	}
}

func checkForHeaderModifiers(w http.ResponseWriter, r *http.Request) {
	for name, values := range r.Header {
		if len(values) == 0 {
			continue
		}

		lowerName := strings.ToLower(name)

		if strings.HasPrefix(lowerName, headerModifiersPrefix) {
			modifierName := name[len(headerModifiersPrefix):]
			if modifierName == "" {
				continue
			}

			for _, value := range values {
				w.Header().Add(modifierName, value)
			}
		}
	}
}

func checkForStatusModifier(r *http.Request) int {
	strStatus := r.Header.Get(headerStatusModifier)
	if strStatus == "" {
		return http.StatusOK
	}

	status, err := strconv.Atoi(strStatus)
	if err != nil {
		slog.Error("Error converting status modifier to integer", "error", err)
		return http.StatusOK
	}

	if status < http.StatusContinue || status > 599 {
		slog.Error("Status modifier out of valid HTTP status code range", "status", status)
		return http.StatusOK
	}

	return status
}
