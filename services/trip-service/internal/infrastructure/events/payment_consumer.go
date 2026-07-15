package events

import (
	"context"
	"encoding/json"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/contracts"
	messaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
)

type PaymentConsumerHandler struct {
	tripService      ports.TripService
	messagePublisher *messaging.Publisher
}

func NewPaymentConsumerHandler(
	tripService ports.TripService,
	messagePublisher *messaging.Publisher,
) *PaymentConsumerHandler {
	return &PaymentConsumerHandler{
		tripService:      tripService,
		messagePublisher: messagePublisher,
	}
}

func (h *PaymentConsumerHandler) Handle(ctx context.Context, message *amqp091.Delivery) error {
	if message.RoutingKey != contracts.PaymentEventSuccess {
		return nil
	}

	var payload messaging.MessageInfo[contracts.PaymentEventPayed]
	if err := json.Unmarshal(message.Body, &payload); err != nil {
		return err
	}

	_, err := h.tripService.UpdateTripStatus(ctx, payload.Data.TripID, domain.TripStatusComplete)
	return err
}
