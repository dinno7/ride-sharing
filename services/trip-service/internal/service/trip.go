package service

import (
	"context"
	"fmt"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/types"
)

type tripService struct {
	repo            ports.TripRepository
	routeCalculator ports.RouteCalculator
}

func NewTripService(
	repo ports.TripRepository,
	routeCalculator ports.RouteCalculator,
) ports.TripService {
	return &tripService{repo: repo, routeCalculator: routeCalculator}
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

func (s *tripService) PreviewTrip(
	ctx context.Context,
	pickup, destination *types.Coordinate,
) (*types.Route, error) {
	routes, err := s.routeCalculator.CalcRoutes(pickup, destination)
	if err != nil {
		return nil, fmt.Errorf("failed to %w", err)
	}
	return routes, nil
}
