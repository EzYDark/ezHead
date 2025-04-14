package server

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug().Msg("AuthMiddleware called")

		apiKey := r.Header.Get("Authorization")
		if apiKey == "" || !strings.HasPrefix(apiKey, "Bearer ") {
			http.Error(w, "Unauthorized: Invalid API key", http.StatusUnauthorized)
			return
		}

		// Validate the API key here
		// Remove "Bearer " prefix and check against your stored keys
		// key := strings.TrimPrefix(apiKey, "Bearer ")
		if false {
			http.Error(w, "Unauthorized: Invalid API key", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
