package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dinno7/ride-sharing/shared/contracts"
	"github.com/dinno7/ride-sharing/shared/logger"
	messaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
)

var ErrNoDriverAvailable = errors.New("no any driver available")

type tripConsumer struct {
	logger        logger.Logger
	driverService *DriverService
	publisher     *messaging.Publisher
}

func NewTripConsumer(
	driverService *DriverService,
	publisher *messaging.Publisher,
	logger logger.Logger,
) *tripConsumer {
	return &tripConsumer{logger: logger, driverService: driverService, publisher: publisher}
}

func (c *tripConsumer) HandleTripCreatedEvent(
	ctx context.Context,
	message *amqp091.Delivery,
) error {
	var payload messaging.MessageInfo[contracts.TripCreatedEventData]
	if err := json.Unmarshal(message.Body, &payload); err != nil {
		return err
	}

	drivers := c.driverService.GetAvailableDriverIDs(payload.Data.Trip.SelectedRideFare.PackageSlug)
	if len(drivers) == 0 {
		c.logger.Info("no any driver found")
		if err := c.publisher.PublishEvent(
			ctx,
			contracts.TripEventNoDriversFound,
			payload.Data.Trip.UserId,
			nil,
		); err != nil {
			c.logger.Error("failed to publish message", err)
			return err
		}
		return ErrNoDriverAvailable
	}

	selectedDriver := drivers[0]
	if err := c.publisher.PublishCommand(
		ctx,
		contracts.DriverCmdTripRequest,
		selectedDriver,
		payload.Data,
	); err != nil {
		c.logger.Error(
			"Failed to publish driver trip request command",
			err,
			"data", payload,
		)
		return err
	}
	c.logger.Info(
		"New driver found for trip",
		"owner_id", payload.OwnerID,
	)
	return nil
}
