package location

import (
	"github.com/rs/zerolog"
	"net/http"
)

func HandleLocation(log zerolog.Logger) http.Handler {
	locationMux := http.NewServeMux()
	locationMux.HandleFunc("GET /{uid}", func(w http.ResponseWriter, r *http.Request) {
		getByUid(w, r, log)
	})
	locationMux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		upsertLocation(w, r, log)
	})
	locationMux.HandleFunc("GET /all", func(w http.ResponseWriter, r *http.Request) {
		getAll(w, r, log)
	})
	locationMux.HandleFunc("PUT /readMany", func(w http.ResponseWriter, r *http.Request) {
		readMany(w, r, log)
	})
	return locationMux
}
