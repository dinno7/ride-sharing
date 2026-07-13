package inmem

import (
	"context"
	"errors"
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

func (repo *tripRepositoryInMem) SaveRideFares(
	ctx context.Context,
	rideFares []*domain.RideFare,
) ([]*domain.RideFare, error) {
	repo.Lock()
	defer repo.Unlock()
	for _, rideFare := range rideFares {
		repo.rideFares[rideFare.ID] = rideFare
	}

	return rideFares, nil
}

func (repo *tripRepositoryInMem) GetFareByID(
	ctx context.Context,
	fareID string,
) (*domain.RideFare, error) {
	repo.RLock()
	defer repo.RUnlock()

	for id, fare := range repo.rideFares {
		if id == fareID {
			return fare, nil
		}
	}
	return nil, domain.ErrFareNotFound
}

func (repo *tripRepositoryInMem) GetTripByID(
	ctx context.Context,
	tripID string,
) (*domain.Trip, error) {
	repo.RLock()
	var trip *domain.Trip
	for i := range repo.trips {
		currentTrip := repo.trips[i]
		if currentTrip.ID == tripID {
			trip = currentTrip
		}
	}
	repo.RUnlock()

	if trip == nil {
		return nil, errors.New("trip not found")
	}

	return trip, nil
}

func (repo *tripRepositoryInMem) UpdateStatus(
	ctx context.Context,
	tripID string,
	newStatus *domain.TripStatus,
) (*domain.Trip, error) {
	repo.Lock()
	defer repo.Unlock()

	trip := repo.trips[tripID]
	trip.Status = newStatus

	return trip, nil
}
