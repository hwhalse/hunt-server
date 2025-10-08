package socket

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"hunt/collections"
	"hunt/constants"
	"hunt/logging"
	"net/http"
	"sync"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex
	handlers map[int]EventHandler
	logger   zerolog.Logger
	ctx      context.Context
}

func NewManager(ctx context.Context) *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[int]EventHandler),
		logger:   logging.NewLogger(),
		ctx:      ctx,
	}
	m.setupHandlers()
	return m
}

// routeEvent allows us to pass errors out of our handlers and handle all errors at once
func (m *Manager) routeEvent(event Event, c *Client) error {
	m.logger.Info().Int("type", event.Type).Str("msg", event.Payload).Msg("incoming")
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			m.logger.Error().Err(err).Msg("unable to route to handler")
			return err
		}
		return nil
	} else {
		return constants.ErrUnknownEvent
	}
}

// setupHandlers adds our handler functions to the manager
func (m *Manager) setupHandlers() {
	m.handlers[Init] = handleInit
	m.handlers[LocationUpdate] = handleLocationUpdate
	m.handlers[CommandNodeUpdate] = handleCommandNodeUpdate
	m.handlers[TargetUpdate] = handleTargetUpdate
	m.handlers[NewCommandNode] = handleNewCommandNode
	m.handlers[NewGroup] = handleNewGroup
	m.handlers[UpdateCallsign] = handleUpdateCallsign
	m.handlers[CommandNodeStatusUpdate] = handleNodeStatusUpdate
}

func (m *Manager) BroadcastMessage(event Event) error {
	for broadcastClient, _ := range m.clients {
		broadcastClient.writeChan <- event
	}
	return nil
}

func (m *Manager) Start(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		m.logger.Error().Err(err).Msg("unable to upgrade websocket")
		return
	}
	client := NewClient(conn, m)
	m.addClient(client)
	go client.read()
	go client.write()
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}

func (m *Manager) removeClient(c *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[c]; ok {
		uid := connMap[c.conn]
		if uid != "" {
			err := collections.LocationCollection.SetActive(c.ctx, uid)
			if err != nil {
				m.logger.Error().Err(err).Msg("unable to set client as active")
			}
		}
		delete(m.clients, c)
		err := c.conn.Close()
		if err != nil {
			m.logger.Error().Err(err).Msg("unable to close client")
		}
	}
}
