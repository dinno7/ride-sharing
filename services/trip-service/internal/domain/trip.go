package domain

import (
	"context"

	"github.com/dinno7/ride-sharing/shared/util"
)

type RideFare struct {
	ID                string
	UserID            string
	PackageSlug       string // van, luxury, sedan
	TotalPriceInCents float64
}

type Trip struct {
	ID       string
	UserID   string
	Status   *tripStatus
	RideFare *RideFare
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
