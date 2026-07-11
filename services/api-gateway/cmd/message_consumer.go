package main

import (
	"context"
	"encoding/json"

	"github.com/dinno7/ride-sharing/shared/contracts"
	"github.com/dinno7/ride-sharing/shared/logger"
	messaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
	wshub "github.com/dinno7/ride-sharing/shared/ws"
	"github.com/rabbitmq/amqp091-go"
)

type MessageConsumerHandler struct {
	wsHub  *wshub.Hub
	logger logger.Logger
}

func NewMessageConsumerHandler(
	wsHub *wshub.Hub,
	logger logger.Logger,
) *MessageConsumerHandler {
	return &MessageConsumerHandler{
		wsHub:  wsHub,
		logger: logger,
	}
}

func (h *MessageConsumerHandler) Forward(
	ctx context.Context,
	message *amqp091.Delivery,
) error {
	var payload messaging.MessageInfo[any]
	if err := json.Unmarshal(message.Body, &payload); err != nil {
		h.logger.Error(
			"Faild to unmartial message body", err,
			"payload", payload,
		)
		return err
	}

	client, err := h.wsHub.GetClientByID(payload.OwnerID)
	if err != nil {
		h.logger.Error(
			"Getting ws client failed", err,
			"payload", payload,
		)
		return err
	}

	data := &contracts.WSMessage{
		Type: message.RoutingKey,
		Data: payload.Data,
	}
	if err := client.SendJSON(data); err != nil {
		h.logger.Error(
			"Sending ws message failed", err,
			"client_id", client.ID(),
			"client_type", client.GetType(),
			"data", data,
		)
		return err
	}
	return nil
}
