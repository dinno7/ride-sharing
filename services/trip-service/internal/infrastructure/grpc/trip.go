package grpc

import (
	"context"
	"log"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	pb "github.com/dinno7/ride-sharing/shared/proto/trip"
	"github.com/dinno7/ride-sharing/shared/types"
	"github.com/dinno7/ride-sharing/shared/util"
	"google.golang.org/grpc"
)

type tripGrpcHandler struct {
	pb.UnimplementedTripServiceServer
	tripService      ports.TripService
	tripEventHandler ports.TripEventHandler
}

func NewTripGrpcHandler(
	server *grpc.Server,
	tripService ports.TripService,
	tripEventHandler ports.TripEventHandler,
) *tripGrpcHandler {
	handler := &tripGrpcHandler{
		tripService:      tripService,
		tripEventHandler: tripEventHandler,
	}
	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (h *tripGrpcHandler) PreviewTrip(
	ctx context.Context,
	req *pb.PreviewTripRequest,
) (*pb.PreviewTripResponse, error) {
	// userID := req.GetUserID()
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()
	log.Println("Calling preview trip application service")
	tripPreview, err := h.tripService.PreviewTrip(ctx,
		req.GetUserID(),
		&types.Coordinate{
			Latitude:  pickup.Latitude,
			Longitude: pickup.Longitude,
		},
		&types.Coordinate{
			Latitude:  destination.Latitude,
			Longitude: destination.Longitude,
		})
	if err != nil {
		return nil, err
	}

	return tripPreviewToGrpc(tripPreview), nil
}

func (h *tripGrpcHandler) StartTrip(
	ctx context.Context,
	req *pb.StartTripRequest,
) (*pb.StartTripResponse, error) {
	fareID := req.GetRideFareID()
	userID := req.GetUserID()
	trip, err := h.tripService.StartTrip(ctx, fareID, userID)
	if err != nil {
		return nil, err
	}

	err = h.tripEventHandler.TripCreated(ctx, trip)
	if err != nil {
		log.Println(err)
		// TODO: ...
	}

	return &pb.StartTripResponse{
		TripID: trip.ID,
		Trip: &pb.Trip{
			ID:               trip.ID,
			SelectedRideFare: rideFareToGrpc(trip.RideFare),
			Status:           trip.Status.String(),
			UserID:           trip.UserID,
			Driver: &pb.TripDriver{
				Id:             util.GenRandomID(),
				Name:           "Taha",
				ProfilePicture: "",
				CarPlate:       "123",
			},
		},
	}, nil
}
