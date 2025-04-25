package handles

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Error().Msgf("Failed to write health response:\n%v", err)
	}
}
