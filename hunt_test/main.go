package main

import (
	"fmt"
	"time"
)

const (
	PingInterval = 50 * time.Second
	PingWait     = 100 * time.Second
	UpdatePeriod = 200 * time.Millisecond
	ClientCount  = 200

	// message protocol
	Init           = 1
	InitResponse   = 2
	LocationUpdate = 6
	NewUid         = 14
)

func main() {
	go startMetricsServer()
	for i := 0; i < ClientCount; i++ {
		time.Sleep(100 * time.Millisecond)
		go func(i int) {
			client := NewClient(i)
			if err := client.Connect(); err != nil {
				fmt.Printf("Client %d failed to connect: %v", i, err)
			}
		}(i)
	}
	select {}
}
