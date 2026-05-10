package service

import (
	"context"
	"fmt"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
)

type tripService struct {
	repo ports.TripRepository
}

func NewTripService(repo ports.TripRepository) ports.TripService {
	return &tripService{repo: repo}
}

func (s *tripService) CreateTrip(
	ctx context.Context,
	fare *domain.RideFare,
) (*domain.Trip, error) {
	trip := domain.NewTrip(fare.UserID, domain.TripStatusPending, fare)
	persistTrip, err := s.repo.CreateTrip(ctx, trip)
	if err != nil {
		return nil, fmt.Errorf("failed to %w", err)
	}
	return persistTrip, nil
}
