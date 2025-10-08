package main

import (
	"context"
	"github.com/rs/cors"
	"hunt/handlers"
	"hunt/logging"
	"hunt/metrics"
	"hunt/socket"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	mux := http.NewServeMux()
	ctx := context.Background()
	manager := socket.NewManager(ctx)
	logger := logging.NewLogger()
	logger.Info().Msg("Hunt listening on port: " + port)
	wss := http.HandlerFunc(manager.Start)
	mux.Handle("/ws", wss)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleTestPing(w, r, logger)
	})
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics.HandleSendingMetrics(w, r, logger)
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	})
	err := http.ListenAndServe(":"+port, c.Handler(metrics.PerformanceLoggingMiddleware(mux, logger)))
	if err != nil {
		panic(err)
	}
}
