package main

import (
	"fmt"
	"net/http"

	grpcclients "github.com/dinno7/ride-sharing/services/api-gateway/cmd/grpc_clients"
	"github.com/dinno7/ride-sharing/shared/contracts"
	driverPb "github.com/dinno7/ride-sharing/shared/proto/driver"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleDriverWS(hub *Hub) echo.HandlerFunc {
	type driverWSInput struct {
		UserID      string `query:"userID"      validate:"required"`
		PackageSlug string `query:"packageSlug" validate:"required,oneof=sedan suv van luxury"`
	}
	return func(c *echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			c.Logger().Error("Failed to connect ws", "error", err)
			return err
		}
		defer ws.Close()

		// NOTE: Validation
		payload := new(driverWSInput)
		err = c.Bind(payload)
		if err != nil {
			c.Logger().Error("Failed to bind payload", "payload", payload, "error", err)
			return websocket.ErrBadHandshake
		}
		if err := c.Validate(payload); err != nil {
			c.Logger().Error("Failed to validate payload", "payload", payload, "error", err)
			return websocket.ErrBadHandshake
		}

		ctx := c.Request().Context()

		client := newDriverClient(ws)
		hub.ClientOnline(client)
		defer hub.ClientOffline(client)

		driverService, err := grpcclients.NewDriverServiceClient()
		if err != nil {
			return err
		}

		defer func() {
			driverService.Client.UnregisterDriver(
				ctx,
				&driverPb.DriverRequest{DriverID: payload.UserID, PackageSlug: payload.PackageSlug},
			)
			driverService.Close()
		}()

		driver, err := driverService.Client.RegisterDriver(
			ctx,
			&driverPb.DriverRequest{DriverID: payload.UserID, PackageSlug: payload.PackageSlug},
		)
		if err != nil {
			return err
		}

		if err := client.conn.WriteJSON(contracts.WSMessage{
			Type: "driver.cmd.register",
			Data: driver.Driver,
		}); err != nil {
			return err
		}

		// NOTE: New Connection established
		c.Logger().Info("New connection")
		for {
			msg, err := client.ReadMessage()
			if err != nil {
				c.Logger().Error("failed read msg", "error", err)
				return err
			}
			fmt.Println("💀 New message from Driver -> ", msg)
		}
	}
}

func handleRiderWS(hub *Hub) echo.HandlerFunc {
	return func(c *echo.Context) error {
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

		client := newRiderClient(ws)
		hub.ClientOnline(client)
		defer hub.ClientOffline(client)

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

			fmt.Println("Got message", msg)
		}
	}
}
