package ports

import (
	"context"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *domain.Trip) (*domain.Trip, error)
	SaveRideFares(
		ctx context.Context,
		estimatedPackagePrices []*domain.RideFare,
	) ([]*domain.RideFare, error)
	GetFareByID(fareID string) (*domain.RideFare, error)
}
