package ports

import (
	"context"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
)

type TripEventHandler interface {
	TripCreated(ctx context.Context, trip *domain.Trip) error
}
