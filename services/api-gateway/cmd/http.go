package main

import (
	"errors"
	"fmt"
	"net/http"

	grpcclients "github.com/dinno7/ride-sharing/services/api-gateway/cmd/grpc_clients"
	"github.com/labstack/echo/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	if closeErr := tripService.Close(); closeErr == nil {
		c.Logger().Info("Connection closed")
	}

	if err != nil {
		c.Logger().Error("failed to preview trip service", "error", err)
		return fmt.Errorf("failed get trip preview from trip service: %w", err)
	}
	c.Logger().Info("Done")

	// NOTE Send success response
	return c.JSON(http.StatusOK, tripServicePayload)
}

func handleTripStart(c *echo.Context) error {
	payload := new(HTTPTripStartRequestPayload)
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
		return status.Errorf(codes.Internal, "something went wrong, please try again: %v", err)
	}

	// NOTE Send req via grpc
	c.Logger().Info("Sending grpc start trip request")
	tripServicePayload, err := tripService.Client.StartTrip(c.Request().Context(), payload.ToGrpc())
	c.Logger().Info("Closing connection")

	if err := tripService.Close(); err == nil {
		c.Logger().Info("Connection closed")
	}
	if err != nil {
		return fmt.Errorf("failed start trip from trip service: %w", err)
	}
	c.Logger().Info("Done")

	// NOTE Send success response
	return c.JSON(http.StatusOK, payload.ToHttp(tripServicePayload))
}
