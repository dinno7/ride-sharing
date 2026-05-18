package http

import (
	"log"
	"net/http"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/shared/types"
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

type TripPreviewResponse struct {
	Route *types.Route `json:"route"`
}

func (h *TripHttpHandler) PreviewTrip(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pickup := types.Coordinate{Latitude: 12.12, Longitude: 12.12}
	destination := types.Coordinate{Latitude: 12.12, Longitude: 12.12}
	userID := util.GenRandomID()
	tripPreviewRoute, err := h.tripService.PreviewTrip(ctx, userID, &pickup, &destination)
	if err != nil {
		log.Fatalf("Err", err)
		return
	}
	util.HttpOkResponse(w, "Proccess successfull",
		&TripPreviewResponse{
			Route: tripPreviewRoute.Route,
		})
}
