package grpc

import (
	"context"
	"log"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/ports"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/domain"
	pb "github.com/dinno7/ride-sharing/shared/proto/trip"
	"github.com/dinno7/ride-sharing/shared/types"
	"google.golang.org/grpc"
)

type tripGrpcHandler struct {
	pb.UnimplementedTripServiceServer
	tripService ports.TripService
}

func NewTripGrpcHandler(server *grpc.Server, tripService ports.TripService) *tripGrpcHandler {
	handler := &tripGrpcHandler{
		tripService: tripService,
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
	routes, err := h.tripService.PreviewTrip(ctx,
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

	return &pb.PreviewTripResponse{
		Route:    routes.ToGrpc(),
		RideFare: []*pb.RideFare{},
	}, nil
}

func (h *tripGrpcHandler) StartTrip(
	ctx context.Context,
	req *pb.StartTripRequest,
) (*pb.StartTripResponse, error) {
	fare := domain.NewRideFare(req.UserID, "", 12.0)
	trip, err := h.tripService.CreateTrip(ctx, fare)
	if err != nil {
		return nil, err
	}

	return &pb.StartTripResponse{
		TripID: trip.ID,
	}, nil
}
