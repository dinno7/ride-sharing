package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/util"
)

type TripHttpHandler struct {
	tripService ports.TripService
}

func NewTripHttpHandler(tripService ports.TripService) *TripHttpHandler {
	return &TripHttpHandler{
		tripService: tripService,
	}
}

func (h *TripHttpHandler) CreateTrip(w http.ResponseWriter, r *http.Request) {
	newTrip, err := h.tripService.CreateTrip(
		context.Background(),
		&domain.RideFare{
			ID:                util.GenRandomID(),
			UserID:            util.GenRandomID(),
			PackageSlug:       "luxury",
			TotalPriceInCents: 100.2,
		},
	)
	if err != nil {
		log.Fatalf("Err", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	newTripJSON, err := json.Marshal(newTrip)
	if err != nil {
		util.ErrorResponse(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// TODO: Publish Event
	log.Printf("trip.created: %s", newTrip.ID)
	w.Write(newTripJSON)
}
