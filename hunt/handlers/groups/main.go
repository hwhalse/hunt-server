package groups

import (
	"github.com/rs/zerolog"
	"net/http"
)

func HandleUnits(log zerolog.Logger) http.Handler {
	unitsMux := http.NewServeMux()
	unitsMux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		handleGetAllGroups(w, r, log)
	})
	unitsMux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		handleAddGroup(w, r, log)
	})
	return unitsMux
}
