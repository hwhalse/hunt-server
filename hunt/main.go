package main

import (
	"hunt/logging"
	huntSocket "hunt/socket"
	"net/http"
	"github.com/lxzan/gws"
)

func main() {
	port := "8080"
	mux := http.NewServeMux()
	upgrader := gws.NewUpgrader(huntSocket.NewHandler(), &gws.ServerOption{
		ParallelEnabled:  true,                               
		Recovery:          gws.Recovery,                    
		PermessageDeflate: gws.PermessageDeflate{Enabled: true},
	})
	mux.HandleFunc("/connect", func(writer http.ResponseWriter, request *http.Request) {
		socket, err := upgrader.Upgrade(writer, request)
		if err != nil {
			return
		}
		go func() {
			socket.ReadLoop()
		}()
	})
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)
	logger := logging.NewLogger()

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
	logger.Info().Msgf("HUNT HTTP server listening on port %s", port)
}