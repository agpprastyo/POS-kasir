package websocket

import (
	"encoding/json"
	"sync"

	"github.com/sirupsen/logrus"
)

// Event types
const (
	EventOrderCreated = "ORDER_CREATED"
	EventOrderUpdated = "ORDER_UPDATED"
)

// Event represents a WebSocket message payload.
type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients (if we need to handle them).
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	mu sync.Mutex
}

var (
	// DefaultHub is the global hub instance
	DefaultHub *Hub
)

func InitHub() *Hub {
	DefaultHub = &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	go DefaultHub.run()
	return DefaultHub
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			logrus.Infof("WebSocket client connected. Total clients: %d", len(h.clients))
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				logrus.Infof("WebSocket client disconnected. Total clients: %d", len(h.clients))
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// Register adding a new client to the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister removing a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// BroadcastEvent is a helper to broadcast a JSON event
func (h *Hub) BroadcastEvent(eventType string, payload interface{}) {
	event := Event{
		Type:    eventType,
		Payload: payload,
	}
	
	bytes, err := json.Marshal(event)
	if err != nil {
		logrus.Errorf("Failed to marshal websocket event: %v", err)
		return
	}
	
	h.broadcast <- bytes
}
