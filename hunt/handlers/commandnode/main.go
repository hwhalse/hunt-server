package commandnode

import (
	"github.com/rs/zerolog"
	"hunt/socket"
	"net/http"
)

func HandleCommandNodes(manager *socket.Manager, log zerolog.Logger) http.Handler {
	commandNodeMux := http.NewServeMux()
	commandNodeMux.HandleFunc("GET /{uid}", func(w http.ResponseWriter, r *http.Request) {
		handleGetByUid(w, r, log)
	})
	commandNodeMux.HandleFunc("PUT /add", func(w http.ResponseWriter, r *http.Request) {
		handleUpsert(w, r, manager, log)
	})
	commandNodeMux.HandleFunc("GET /all", func(w http.ResponseWriter, r *http.Request) {
		getAllNodes(w, r, log)
	})
	commandNodeMux.HandleFunc("PUT /readMany", func(w http.ResponseWriter, r *http.Request) {
		handleReadMany(w, r, log)
	})
	commandNodeMux.HandleFunc("PUT /delete", func(w http.ResponseWriter, r *http.Request) {
		handleDelete(w, r, manager, log)
	})
	return commandNodeMux
}
