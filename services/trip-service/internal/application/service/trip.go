package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/types"
	"github.com/dinno7/ride-sharing/shared/util"
)

type tripService struct {
	repo             ports.TripRepository
	routeCalculator  ports.RouteCalculator
	tripEventHandler ports.TripEventHandler
}

func NewTripService(
	repo ports.TripRepository,
	routeCalculator ports.RouteCalculator,
	tripEventHandler ports.TripEventHandler,
) ports.TripService {
	return &tripService{
		repo:             repo,
		routeCalculator:  routeCalculator,
		tripEventHandler: tripEventHandler,
	}
}

func (s *tripService) StartTrip(ctx context.Context, fareID, userID string) (*domain.Trip, error) {
	// NOTE: Check trip fare & userid exists in db
	fetchedFare, err := s.repo.GetFareByID(fareID)
	if err != nil {
		return nil, err
	}
	if fetchedFare.UserID != userID {
		return nil, errors.New("fare not belongs you")
	}

	// NOTE: Create New trip
	trip := domain.NewTrip(fetchedFare.UserID, domain.TripStatusPending, fetchedFare)
	persistTrip, err := s.repo.CreateTrip(ctx, trip)
	if err != nil {
		return nil, fmt.Errorf("failed to %w", err)
	}

	if err := s.tripEventHandler.TripCreated(ctx, trip); err != nil {
		return nil, err
	}

	return persistTrip, nil
}

func (s *tripService) PreviewTrip(
	ctx context.Context,
	userID string,
	pickup, destination *types.Coordinate,
) (*ports.PreviewTripOutput, error) {
	route, err := s.routeCalculator.CalcRoutes(pickup, destination)
	if err != nil {
		return nil, fmt.Errorf("failed to %w", err)
	}
	// NOTE: Estimate pkg price with route
	estimatedPrices := estimatePackagePrice(route)
	// NOTE: Generate & Save to DB Trip fare
	rideFares := make([]*domain.RideFare, len(estimatedPrices))
	for i, rfp := range estimatedPrices {
		rideFares[i] = &domain.RideFare{
			ID:                util.GenRandomID(),
			UserID:            userID,
			PackageSlug:       rfp.PackageSlug,
			TotalPriceInCents: rfp.TotalPriceInCents,
			Route:             route,
		}
	}

	rideFares, err = s.repo.SaveRideFares(ctx, rideFares)
	if err != nil {
		return nil, err
	}

	return &ports.PreviewTripOutput{
		RideFares: rideFares,
		Route:     route,
	}, nil
}
