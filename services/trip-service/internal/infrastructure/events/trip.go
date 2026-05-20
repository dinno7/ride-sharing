package events

import (
	"context"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/contracts"
	messaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
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
		contracts.TripCreatedEventData{
			TripID:      trip.ID,
			PackageSlug: trip.RideFare.PackageSlug,
		},
	)
}
