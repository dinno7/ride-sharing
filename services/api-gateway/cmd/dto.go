package main

import (
	pb "github.com/dinno7/ride-sharing/shared/proto/trip"
	"github.com/dinno7/ride-sharing/shared/types"
)

type HTTPTripPreviewRequestPayload struct {
	UserID      string           `json:"userID"      validate:"required"`
	Pickup      types.Coordinate `json:"pickup"      validate:"required"`
	Destination types.Coordinate `json:"destination" validate:"required"`
}

func (r *HTTPTripPreviewRequestPayload) ToGrpc() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserID: r.UserID,
		StartLocation: &pb.Cordinate{
			Latitude:  r.Pickup.Latitude,
			Longitude: r.Pickup.Longitude,
		},
		EndLocation: &pb.Cordinate{
			Latitude:  r.Destination.Latitude,
			Longitude: r.Destination.Longitude,
		},
	}
}
