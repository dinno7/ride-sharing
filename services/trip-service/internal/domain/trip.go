package domain

import (
	"context"

	"github.com/dinno7/ride-sharing/shared/util"
)

type RideFare struct {
	ID                string  `json:"id"`
	UserID            string  `json:"user_id"`
	PackageSlug       string  `json:"package_slug"` // van, luxury, sedan
	TotalPriceInCents float64 `json:"total_price_in_cents"`
}

type Trip struct {
	ID       string      `json:"id"`
	UserID   string      `json:"user_id"`
	Status   *tripStatus `json:"status"`
	RideFare *RideFare   `json:"ride_fare"`
}

func NewTrip(userID string, status *tripStatus, fare *RideFare) *Trip {
	return &Trip{
		ID:       util.GenRandomID(),
		UserID:   userID,
		Status:   status,
		RideFare: fare,
	}
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *Trip) (*Trip, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFare) (*Trip, error)
}
