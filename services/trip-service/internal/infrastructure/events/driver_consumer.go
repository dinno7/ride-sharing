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

type DriverConsumerHandler struct {
	tripService      ports.TripService
	messagePublisher *messaging.Publisher
}

func NewDriverConsumerHandler(
	tripService ports.TripService,
	messagePublisher *messaging.Publisher,
) *DriverConsumerHandler {
	return &DriverConsumerHandler{
		tripService:      tripService,
		messagePublisher: messagePublisher,
	}
}

func (h *DriverConsumerHandler) Handle(ctx context.Context, message *amqp091.Delivery) error {
	var payload messaging.MessageInfo[contracts.DriverResponseToTripData]
	if err := json.Unmarshal(message.Body, &payload); err != nil {
		return err
	}

	switch message.RoutingKey {
	case contracts.DriverCmdTripAccept:
		// INFO: Update the trip's status
		trip, err := h.tripService.UpdateTripStatus(
			ctx,
			payload.Data.TripID,
			domain.TripStatusAccepted,
		)
		if err != nil {
			return err
		}
		// INFO: Publish driver assigend event
		if err := h.messagePublisher.PublishEvent(
			ctx,
			contracts.TripEventDriverAssigned,
			trip.UserID,
			trip,
		); err != nil {
			return err
		}

		// TODO: Notify payment service to start payment link
	case contracts.DriverCmdTripDecline:
	}

	return nil
}
