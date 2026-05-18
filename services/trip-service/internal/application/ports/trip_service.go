package ports

import (
	"context"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/types"
)

type PreviewTripOutput struct {
	RideFares []*domain.RideFare
	Route     *types.Route
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *domain.RideFare) (*domain.Trip, error)
	PreviewTrip(
		ctx context.Context,
		userID string,
		pickup, destination *types.Coordinate,
	) (*PreviewTripOutput, error)
}
