package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
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

func tripPreview(c *echo.Context) error {
	payload := new(HTTPTripPreviewRequestPayload)
	if err := c.Bind(payload); err != nil {
		return err
	}
	if err := c.Validate(payload); err != nil {
		return err
	}

	// TODO: Send to service
	client := new(http.Client)
	jsonToService, _ := json.Marshal(payload)
	reader := bytes.NewReader(jsonToService)
	req, err := http.NewRequest(http.MethodPost, "http://trip-service:7000/preview", reader)
	if err != nil {
		c.Logger().Error("failed to create request to trup service", "error", err)
		return echo.ErrInternalServerError.Wrap(
			errors.New("something went wrong, please try again"),
		)
	}

	resp, err := client.Do(req)
	if err != nil {
		c.Logger().Error("failed to send request to trup service", "error", err)
		return echo.ErrInternalServerError.Wrap(
			errors.New("something went wrong, please try again"),
		)
	}

	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	// TODO: Send success response
	return c.JSONBlob(http.StatusOK, b)
}
