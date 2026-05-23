package events

import (
	"context"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/contracts"
	messaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
	pb "github.com/dinno7/ride-sharing/shared/proto/trip"
)

type tripEventHandler struct {
	publisher *messaging.Publisher
}

func NewTripEventHandler(publisher *messaging.Publisher) ports.TripEventHandler {
	return &tripEventHandler{
		publisher: publisher,
	}
}

func (eh *tripEventHandler) TripCreated(ctx context.Context, trip *domain.Trip) error {
	return eh.publisher.PublishEvent(
		ctx,
		contracts.TripEventCreated,
		trip.UserID,
		contracts.TripCreatedEventData{
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
		},
	)
}
