package ports

import (
	"context"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/types"
)

type TripService interface {
	CreateTrip(ctx context.Context, fare *domain.RideFare) (*domain.Trip, error)
	PreviewTrip(ctx context.Context, pickup, destination *types.Coordinate) (*types.Route, error)
}
