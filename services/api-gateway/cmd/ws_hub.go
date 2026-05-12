package main

import (
	"github.com/gorilla/websocket"
)

// type Client struct {
// 	ID         string
// 	conn       *websocket.Conn
// 	hub        *Hub
// 	clientType ClientType
// }

type Client interface {
	GetType() string
	IsEqual(currentConn *websocket.Conn) bool
	SendMessage(msg string) error
}

type (
	clientStorage = map[Client]struct{}
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
			h.clients[newClient] = struct{}{}
		case c := <-h.unregister:
			delete(h.clients, c)
		}
	}
}

func (h *Hub) ClientOnline(client Client) {
	h.register <- client
}

func (h *Hub) ClientOffline(client Client) {
	h.unregister <- client
}

func (h *Hub) BroadcastTo(cType string, msg string) {
	h.broadcastTo(cType, msg)
}

func (h *Hub) Broadcast(msg string) {
	h.broadcast(msg)
}

func (h *Hub) broadcastTo(cType string, msg string) {
	for client := range h.clients {
		if client.GetType() == cType {
			_ = client.SendMessage(msg)
		}
	}
}

func (h *Hub) broadcast(msg string) {
	for client := range h.clients {
		_ = client.SendMessage(msg)
	}
}
