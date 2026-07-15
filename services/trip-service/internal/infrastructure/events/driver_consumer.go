package events

import (
	"context"
	"encoding/json"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/contracts"
	messaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
	pb "github.com/dinno7/ride-sharing/shared/proto/trip"
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

		// INFO: Notify payment service to start payment link
		if err := h.messagePublisher.PublishCommand(
			ctx,
			contracts.PaymentCmdCreateSession,
			trip.UserID,
			contracts.PaymentCmdCreateSessionData{
				TripID:   trip.ID,
				UserID:   trip.UserID,
				DriverID: payload.Data.Driver.Id,
				Amount:   trip.RideFare.TotalPriceInCents,
				Currency: "USD",
			},
		); err != nil {
			return err
		}
	case contracts.DriverCmdTripDecline:
		trip, err := h.tripService.GetTripByID(ctx, payload.Data.TripID)
		if err != nil {
			return err
		}

		data := contracts.TripCreatedEventData{
			Trip: &pb.Trip{
				Id:     trip.ID,
				Route:  trip.RideFare.Route.ToGrpc(),
				Status: trip.Status.String(),
				UserId: trip.UserID,
				Driver: nil,
				SelectedRideFare: &pb.RideFare{
					Id:                trip.RideFare.ID,
					UserId:            trip.RideFare.UserID,
					PackageSlug:       trip.RideFare.PackageSlug,
					TotalPriceInCents: trip.RideFare.TotalPriceInCents,
				},
			},
		}
		return h.messagePublisher.PublishEvent(
			ctx,
			contracts.TripEventDriverNotInterested,
			payload.OwnerID,
			data,
		)
	}

	return nil
}
