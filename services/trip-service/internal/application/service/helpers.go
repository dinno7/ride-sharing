package service

import (
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/types"
)

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
