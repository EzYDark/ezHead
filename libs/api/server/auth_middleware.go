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
			log.Error().Msg("Invalid API key. No API key or no Bearer prefix.")
			http.Error(w, "Unauthorized: Invalid API key", http.StatusUnauthorized)
			return
		}

		// key := strings.TrimPrefix(apiKey, "Bearer ")
		if false {
			log.Error().Msg("Invalid API key")
			http.Error(w, "Unauthorized: Invalid API key", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
