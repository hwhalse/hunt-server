package main

import (
	"context"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/rs/zerolog"
	"log"
	"net/http"
	"os"
)

//func printMemoryUsage() {
//	var memStats runtime.MemStats
//	runtime.ReadMemStats(&memStats)
//	fmt.Printf("Allocated memory: %v\n", memStats.Alloc)
//	fmt.Printf("Heap allocated memory: %v\n", memStats.HeapAlloc)
//	fmt.Printf("System memory obtained from OS: %v \n", memStats.Sys)
//	fmt.Printf("Number of garbage collection cycles: %v\n", memStats.NumGC)
//}

func main() {
	http.HandleFunc("/", homeHandler)
	err := http.ListenAndServe(":6666", nil)
	if err != nil {
		panic(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		http.Error(w, "Failed to establish WebSocket connection", http.StatusBadRequest)
		return
	}
	defer c.Close(websocket.StatusNormalClosure, "")

	ctx := context.Background()
	for {
		var v interface{}
		err := wsjson.Read(ctx, c, &v)
		if err != nil {
			log.Println(err)
			break
		}
		log.Printf("Received message: %+v\n", v)

		err = wsjson.Write(ctx, c, "Hello from server!")
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func newLogger() *zerolog.Logger {
	// Create a new logger
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        log.Writer(),
		TimeFormat: "2006-01-02 15:04:05",
		FormatLevel: func(i interface{}) string {
			level := i.(string)
			switch level {
			case "trace":
				return "[T] "
			case "debug":
				return "[D] "
			case "info":
				return "[I] "
			case "warn":
				return "[W] "
			case "error":
				return "[E] "
			case "fatal":
				return "[F] "
			case "panic":
				return "[P] "
			default:
				return level
			}
		},
	}).
		With().
		Timestamp().
		Str("app", "MUNC").
		Str("env", os.Getenv("ENV")).
		Str("host", os.Getenv("HOSTNAME")).
		Str("pid", os.Getenv("PID")).
		Stack().
		CallerWithSkipFrameCount(2).
		Logger()

	return &logger
}
