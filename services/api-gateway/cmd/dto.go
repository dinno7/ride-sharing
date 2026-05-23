package main

import (
	pb "github.com/dinno7/ride-sharing/shared/proto/trip"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type HTTPTripPreviewRequestPayload struct {
	UserID      string     `json:"userID"      validate:"required"`
	Pickup      Coordinate `json:"pickup"      validate:"required"`
	Destination Coordinate `json:"destination" validate:"required"`
}

func (r *HTTPTripPreviewRequestPayload) ToGrpc() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserId: r.UserID,
		StartLocation: &pb.Coordinate{
			Latitude:  r.Pickup.Latitude,
			Longitude: r.Pickup.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  r.Destination.Latitude,
			Longitude: r.Destination.Longitude,
		},
	}
}

type HTTPTripStartRequestPayload struct {
	RideFareID string `json:"rideFareID" validate:"required"`
	UserID     string `json:"userID"     validate:"required"`
}

func (r *HTTPTripStartRequestPayload) ToGrpc() *pb.StartTripRequest {
	return &pb.StartTripRequest{
		RideFareId: r.RideFareID,
		UserId:     r.UserID,
	}
}

type HTTPTripStartResponse struct {
	TripID string `json:"tripID"`
}

func (r *HTTPTripStartRequestPayload) ToHttp(grpcRes *pb.StartTripResponse) *HTTPTripStartResponse {
	return &HTTPTripStartResponse{
		TripID: grpcRes.TripId,
	}
}
