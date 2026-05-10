package inmem

import (
	"context"
	"sync"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
)

type tripRepositoryInMem struct {
	trips     map[string]*domain.Trip
	rideFares map[string]*domain.RideFare
	sync.RWMutex
}

func NewInMemTripRepository() ports.TripRepository {
	return &tripRepositoryInMem{
		trips:     make(map[string]*domain.Trip),
		rideFares: make(map[string]*domain.RideFare),
	}
}

func (repo *tripRepositoryInMem) CreateTrip(
	ctx context.Context,
	trip *domain.Trip,
) (*domain.Trip, error) {
	repo.Lock()
	defer repo.Unlock()

	repo.trips[trip.ID] = trip

	return trip, nil
}
