package main

import (
	"errors"
	"net/http"

	grpcclients "github.com/dinno7/ride-sharing/services/api-gateway/cmd/grpc_clients"
	pb "github.com/dinno7/ride-sharing/shared/proto/trip"
	"github.com/dinno7/ride-sharing/shared/types"
	"github.com/labstack/echo/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func handleTripPreview(c *echo.Context) error {
	payload := new(HTTPTripPreviewRequestPayload)
	if err := c.Bind(payload); err != nil {
		return err
	}
	if err := c.Validate(payload); err != nil {
		return err
	}

	// NOTE Send to service
	c.Logger().Info("Initiate grpc connection with trip service")
	tripService, err := grpcclients.NewTripServiceClient()
	if err != nil {
		c.Logger().Error("failed to connect trip service", "error", err)
		return echo.ErrInternalServerError.Wrap(
			errors.New("something went wrong, please try again"),
		)
	}

	// NOTE Send req via grpc
	c.Logger().Info("Sending grpc preview trip request")
	tripServicePayload, err := tripService.Client.PreviewTrip(
		c.Request().Context(),
		payload.ToGrpc(),
	)
	c.Logger().Info("Closing connection")

	if err := tripService.Close(); err == nil {
		c.Logger().Info("Connection closed")
	}
	if err != nil {
		return status.Errorf(codes.Internal, "failed get trip preview from trip service: %v", err)
	}
	c.Logger().Info("Done")

	// NOTE Send success response
	return c.JSON(http.StatusOK, tripServicePayload)
}
