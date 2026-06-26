package ws

import (
	"fmt"

	"github.com/dinno7/ride-sharing/shared/ws"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

type RiderWSHandler struct {
	hub *ws.Hub
}

func NewRiderWSHandler(hub *ws.Hub) *RiderWSHandler {
	return &RiderWSHandler{
		hub,
	}
}

func (h *RiderWSHandler) Handle(c *echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		c.Logger().Error("Failed to connect ws", "error", err)
		return err
	}
	defer ws.Close()

	userID := c.QueryParam("userID")
	if len(userID) == 0 {
		c.Logger().Error("Failed get userID", "userID", userID)
		return websocket.ErrBadHandshake
	}

	client := newRiderClient(ws, userID)
	h.hub.ClientOnline(client)
	defer h.hub.ClientOffline(client)

	// NOTE: New Connection established
	c.Logger().Info("New connection")

	for {
		msg, err := client.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				c.Logger().Error("Failed to read ws message", "error", err)
			}
			return err
		}

		fmt.Println("💀 New ws message from rider", msg)
	}
}
