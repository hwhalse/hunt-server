// Package socket handles the websocket server
package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"hunt/logging"
	"time"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
	connMap      = make(ConnMap)
)

type Client struct {
	conn      *websocket.Conn
	manager   *Manager
	writeChan chan Event
	ctx       context.Context
	logger    zerolog.Logger
}

type ClientList map[*Client]bool
type ConnMap map[*websocket.Conn]string

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		conn:      conn,
		manager:   manager,
		writeChan: make(chan Event),
		ctx:       context.Background(),
		logger:    logging.NewLogger(),
	}
}

// read leave the set read readline commented out for now, it was causing the server to drop clients
func (c *Client) read() {
	defer c.manager.removeClient(c)
	c.conn.SetReadLimit(1500)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.logger.Error().Err(err).Msg("error reading message")
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error().Err(err).Msg("unexpected close")
			}
			c.handleDisconnect()
			break
		}
		var request Event
		err = json.Unmarshal(message, &request)
		if err != nil {
			c.logger.Error().Err(err).Msg("unable to unmarshal request")
			c.handleDisconnect()
			break
		}
		if err := c.manager.routeEvent(request, c); err != nil {
			c.logger.Error().Err(err).Msg("unable to route event")
		}
		b := make([]byte, len(message))
		copy(b, message)
	}
}

func (c *Client) handleDisconnect() {
	c.manager.removeClient(c)
}

func (c *Client) write() {
	defer func() {
		c.manager.removeClient(c)
	}()
	ticker := time.NewTicker(pingInterval)
	for {
		select {
		case message, ok := <-c.writeChan:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					c.logger.Error().Err(err).Msg("unable to write msg")
				}
				c.handleDisconnect()
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				fmt.Println("Failed to send message to client: ", err)
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				c.logger.Error().Err(err).Msg("unable to write ping")
				return
			}
		}
	}
}
