package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/lxzan/gws"
)

type Event struct {
	Type    int    `json:"type"`
	Payload string `json:"payload"`
	Time time.Time `json:"time"`
}

type InitMsg struct {
	Uid      string `json:"uid,omitempty"`
	Callsign string `json:"callsign"`
}

type LocationUpdatePayload struct {
	Callsign string   `json:"callsign"`
	Uid      string   `json:"uid"`
	Location Location `json:"location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Alt float64 `json:"alt"`
}

type Handler struct {
	client *Client
}

func (h *Handler) OnOpen(socket *gws.Conn) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	log.Printf("[%s] connected", h.client.callsign)
}

func (h *Handler) OnClose(socket *gws.Conn, err error) {
	log.Printf("[%s] connection closed: %v", h.client.callsign, err)
	if h.client.cancel != nil {
		h.client.cancel()
	}
	socket.NetConn().Close()
}

func (h *Handler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.WritePong(nil)
}

func (h *Handler) OnPong(socket *gws.Conn, payload []byte) {}

func (h *Handler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()

	var event Event
	if err := json.Unmarshal(message.Bytes(), &event); err != nil {
		log.Printf("[%s] failed to parse message: %v", h.client.callsign, err)
		return
	}

	switch event.Type {
	case InitResponse:
		log.Printf("[%s] init response:", h.client.callsign)

	case NewUid:
		h.client.uid = event.Payload
		ctx, cancel := context.WithCancel(context.Background())
		h.client.cancel = cancel
		go h.client.startSending(ctx)

	case LocationUpdate:
		recordLatency(time.Since(event.Time))

	default:
		log.Printf("[%s] unhandled event type: %d", h.client.callsign, event.Type)
	}
}