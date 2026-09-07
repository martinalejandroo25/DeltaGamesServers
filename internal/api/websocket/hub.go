package websocket

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/deltagames/deltagamesservers/internal/gameserver"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow CORS in development
	},
}

type ClientConnection struct {
	conn     *websocket.Conn
	serverID string
	send     chan []byte
}

type Hub struct {
	connections map[*ClientConnection]bool
	register    chan *ClientConnection
	unregister  chan *ClientConnection
	manager     *gameserver.Manager
	mu          sync.Mutex
}

func NewHub(manager *gameserver.Manager) *Hub {
	return &Hub{
		connections: make(map[*ClientConnection]bool),
		register:    make(chan *ClientConnection),
		unregister:  make(chan *ClientConnection),
		manager:     manager,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.connections[client] = true
			h.mu.Unlock()
			log.Printf("[WebSocket] Client connected for server %s", client.serverID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.connections[client]; ok {
				delete(h.connections, client)
				close(client.send)
				log.Printf("[WebSocket] Client disconnected for server %s", client.serverID)
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) HandleLogsWS(c *gin.Context, serverID string) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WebSocket ERROR] Upgrade failed: %v", err)
		return
	}

	client := &ClientConnection{
		conn:     ws,
		serverID: serverID,
		send:     make(chan []byte, 256),
	}

	h.register <- client

	// Start reader & writer loops
	go h.writePump(client)
	go h.streamLogsPump(client)
}

func (h *Hub) writePump(c *ClientConnection) {
	defer func() {
		c.conn.Close()
	}()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}

func (h *Hub) streamLogsPump(c *ClientConnection) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	var lastLogs string

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			logs, err := h.manager.GetServerLogs(ctx, c.serverID, 50)
			cancel()

			if err == nil && logs != lastLogs {
				lastLogs = logs
				c.send <- []byte(logs)
			}
		}
	}
}
