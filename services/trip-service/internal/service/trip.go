package service

import (
	"context"
	"fmt"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/types"
	"github.com/dinno7/ride-sharing/shared/util"
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
	userID string,
	pickup, destination *types.Coordinate,
) (*ports.PreviewTripOutput, error) {
	route, err := s.routeCalculator.CalcRoutes(pickup, destination)
	if err != nil {
		return nil, fmt.Errorf("failed to %w", err)
	}
	// TODO: Estimate pkg price with route
	estimatedPrices := estimatePackagePrice(route)
	// TODO: Generate & Save to DB Trip fare
	rideFares := make([]*domain.RideFare, len(estimatedPrices))
	for i, rfp := range estimatedPrices {
		rideFares[i] = &domain.RideFare{
			ID:                util.GenRandomID(),
			UserID:            userID,
			PackageSlug:       rfp.PackageSlug,
			TotalPriceInCents: rfp.TotalPriceInCents,
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

func estimatePackagePrice(route *types.Route) []*domain.RideFare {
	baseFares := getBaseFares()
	estimatedFares := make([]*domain.RideFare, len(baseFares))
	for i, bf := range baseFares {
		estimatedFares[i] = estimateFareForRoute(bf, route.Distance, route.Duration)
	}
	return estimatedFares
}

func estimateFareForRoute(baseFare *domain.RideFare, distance, duration float64) *domain.RideFare {
	pricePerKM := 1.5
	pricePerMinute := 0.25

	distancePrice := pricePerKM * distance
	durationPrice := pricePerMinute * duration

	totalPrice := distancePrice + durationPrice + baseFare.TotalPriceInCents

	return &domain.RideFare{
		PackageSlug:       baseFare.PackageSlug,
		TotalPriceInCents: totalPrice,
	}
}

func getBaseFares() []*domain.RideFare {
	return []*domain.RideFare{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 400,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
	}
}
