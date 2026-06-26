package ws

import (
	"errors"

	"github.com/gorilla/websocket"
)

type Client interface {
	ID() string
	GetType() string
	IsEqual(currentConn *websocket.Conn) bool
	SendMessage(msg string) error
	SendJSON(data any) error
	ReadMessage() ([]byte, error)
}

type (
	clientStorage = map[string]Client
	Hub           struct {
		clients    clientStorage
		register   chan Client
		unregister chan Client
	}
)

func NewHub() *Hub {
	return &Hub{
		clients:    make(clientStorage),
		register:   make(chan Client),
		unregister: make(chan Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case newClient := <-h.register:
			h.clients[newClient.ID()] = newClient
		case client := <-h.unregister:
			delete(h.clients, client.ID())
		}
	}
}

func (h *Hub) ClientOnline(client Client) {
	h.register <- client
}

func (h *Hub) ClientOffline(client Client) {
	h.unregister <- client
}

// Returning client by id, it can be nil if no client found via id
func (h *Hub) GetClientByID(clientID string) (Client, error) {
	if client, ok := h.clients[clientID]; ok {
		return client, nil
	}
	return nil, errors.New("client not found")
}

func (h *Hub) BroadcastTo(cType string, msg string) {
	h.broadcastTo(cType, msg)
}

func (h *Hub) Broadcast(msg string) {
	h.broadcast(msg)
}

func (h *Hub) broadcastTo(cType string, msg string) {
	for clientID := range h.clients {
		client := h.clients[clientID]
		if client.GetType() == cType {
			_ = client.SendMessage(msg)
		}
	}
}

func (h *Hub) broadcast(msg string) {
	for clientID := range h.clients {
		client := h.clients[clientID]
		_ = client.SendMessage(msg)
	}
}
