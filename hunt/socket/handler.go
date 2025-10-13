package socket

import (
	"context"
	"encoding/json"
	"hunt/state"
	"hunt/collections"
	"hunt/models"
	"log"
	"time"

	"github.com/lxzan/gws"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	PingInterval = 50 * time.Second
	PingWait     = 100 * time.Second
)

// Change to true to use state object instead of writing to DB
var useState = false

func NewHandler() *Handler {
	return &Handler{}
}

type Handler struct{}

var huntState = state.NewStateObject()
var locations = state.NewLocationsObject()

var connMap = gws.NewConcurrentMap[*gws.Conn, models.HuntUser](16, 256)

func (c *Handler) OnOpen(socket *gws.Conn) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
}

func (c *Handler) OnClose(socket *gws.Conn, err error) {
	if client, ok := connMap.Load(socket); ok {
		collections.UsersCollectionManager.Collection.DeleteOne(context.Background(), bson.M{
			"uid": client.Uid,
		})
		connMap.Delete(socket)
	}
}

func (c *Handler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
	_ = socket.WritePong(nil)
}

func (c *Handler) OnPong(socket *gws.Conn, payload []byte) {}

func (c *Handler) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	e := &Event{}
	err := json.Unmarshal(message.Bytes(), e)
	if err != nil {
		log.Fatal(err)
	}
	switch e.Type {
	case Init:
		handleInit(e, socket, useState)
	case LocationUpdate:
		handleLocationUpdate(e, message.Bytes(), useState)
	}
}

func Broadcast(conns []*gws.Conn, opcode gws.Opcode, payload []byte) {
    var b = gws.NewBroadcaster(opcode, payload)
    defer b.Close()
    for _, item := range conns {
        _ = b.Broadcast(item)
    }
}