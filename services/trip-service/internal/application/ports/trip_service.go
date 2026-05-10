package ports

import (
	"context"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
)

type TripService interface {
	CreateTrip(ctx context.Context, fare *domain.RideFare) (*domain.Trip, error)
}
