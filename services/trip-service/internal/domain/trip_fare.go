package domain

import (
	"github.com/dinno7/ride-sharing/shared/types"
	"github.com/dinno7/ride-sharing/shared/util"
)

type RideFare struct {
	ID                string       `json:"id"`
	UserID            string       `json:"user_id"`
	PackageSlug       string       `json:"package_slug"` // van, luxury, sedan
	TotalPriceInCents float64      `json:"total_price_in_cents"`
	Route             *types.Route `json:"route"`
}

func NewRideFare(userID, packageSlug string, priceCents float64) *RideFare {
	return &RideFare{
		ID:                util.GenRandomID(),
		UserID:            userID,
		PackageSlug:       packageSlug,
		TotalPriceInCents: priceCents,
	}
}
