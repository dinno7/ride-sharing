package ws

import (
	"fmt"

	grpcclients "github.com/dinno7/ride-sharing/services/api-gateway/cmd/grpc_clients"
	"github.com/dinno7/ride-sharing/shared/contracts"
	"github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
	driverPb "github.com/dinno7/ride-sharing/shared/proto/driver"
	"github.com/dinno7/ride-sharing/shared/ws"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

type DriverWSHandler struct {
	hub               *ws.Hub
	messengerConsumer *rabbitmq.Consumer
}

func NewDriverWSHandler(hub *ws.Hub, messengerConsumer *rabbitmq.Consumer) *DriverWSHandler {
	return &DriverWSHandler{
		hub,
		messengerConsumer,
	}
}

type driverWSInput struct {
	UserID      string `query:"userID"      validate:"required"`
	PackageSlug string `query:"packageSlug" validate:"required,oneof=sedan suv van luxury"`
}

func (h *DriverWSHandler) Handle(c *echo.Context) error {
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

	client := newDriverClient(ws, payload.UserID)
	h.hub.ClientOnline(client)
	defer h.hub.ClientOffline(client)

	driverService, err := grpcclients.NewDriverServiceClient()
	if err != nil {
		return err
	}
	defer func() {
		driverService.Client.UnregisterDriver(
			ctx,
			&driverPb.DriverRequest{DriverId: payload.UserID, PackageSlug: payload.PackageSlug},
		)
		driverService.Close()
	}()

	driver, err := driverService.Client.RegisterDriver(
		ctx,
		&driverPb.DriverRequest{DriverId: payload.UserID, PackageSlug: payload.PackageSlug},
	)
	if err != nil {
		return err
	}

	if err := client.SendJSON(contracts.WSMessage{
		Type: contracts.DriverCmdRegister,
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
		fmt.Println("💀 New ws message from Driver -> ", msg)
	}
}
