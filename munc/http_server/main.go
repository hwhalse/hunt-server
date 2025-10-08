package http_server

import (
	"eolian/munc/utils"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rs/cors"
)

func StartHttpServer() {
	port := utils.GetEnvString("MUNC_HTTP_PORT", "10000")

	mux := http.NewServeMux()

	mux.HandleFunc("POST /liveData", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "Hello, World!")
	})
	mux.HandleFunc("GET /config", configHandler)
	mux.HandleFunc("/", fallbackHandler)

	addr := fmt.Sprintf(":%s", port)
	slog.Info("HTTP Server Started", "port", port)
	c := cors.New(cors.Options{
		AllowedHeaders:   []string{"X-Api-Key"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET"},
	})
	err := http.ListenAndServe(addr, c.Handler(mux))
	if err != nil {
		panic(err)
	}
}

func configHandler(response http.ResponseWriter, request *http.Request) {
	// if request.Method != http.MethodGet {
    //     response.Header().Set("Allow", http.MethodGet)
    //     http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
    //     return
    // }

	fmt.Fprint(response, utils.GetEnvString("MUNC_HTTP_PORT", "10000"))
}

func fallbackHandler(response http.ResponseWriter, request *http.Request) {
    endpoint := request.URL.Path
    slog.Info("Received request", "endpoint", endpoint)
    response.WriteHeader(http.StatusOK)
}