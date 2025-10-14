package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lxzan/gws"
)

type Client struct {
	socket   *gws.Conn
	uid      string
	callsign string
	lat      float64
	lon      float64
	cancel   context.CancelFunc
}

// i represents the 0-based index of the client as it's spawned in a loop
func NewClient(i int) *Client {
	callsign := fmt.Sprintf("test user %d", i)
	lat, lon := randomUSCoordinate()

	return &Client{
		callsign: callsign,
		lat:      lat,
		lon:      lon,
	}
}

func (c *Client) Connect() error {
	handler := &Handler{client: c}
	socket, _, err := gws.NewClient(handler, &gws.ClientOption{
		Addr: "ws://localhost:8080/connect",
		PermessageDeflate: gws.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
	})
	if err != nil {
		return err
	}
	c.socket = socket

	initMsg := InitMsg{Callsign: c.callsign}
	data, _ := json.Marshal(initMsg)
	event := Event{Type: Init, Payload: string(data), Time: time.Now()}
	msg, _ := json.Marshal(event)
	_ = socket.WriteMessage(gws.OpcodeBinary, msg)

	go socket.ReadLoop()
	return nil
}

func (c *Client) startSending(ctx context.Context) {
	ticker := time.NewTicker(UpdatePeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			lat, lon := randomNearbyCoordinate(c.lat, c.lon)
			c.lat, c.lon = lat, lon
			update := LocationUpdatePayload{
				Callsign: c.callsign,
				Uid:      c.uid,
				Location: Location{Lat: lat, Lon: lon, Alt: 0},
			}
			data, _ := json.Marshal(update)
			event := Event{Type: LocationUpdate, Payload: string(data), Time: time.Now()}
			msg, _ := json.Marshal(event)
			c.socket.WriteAsync(gws.OpcodeBinary, msg, func(err error) {
				if err != nil {
					log.Printf("%s send failed: %v", c.callsign, err)
				} else {
					recordSend()
				}
			})
		}
	}
}