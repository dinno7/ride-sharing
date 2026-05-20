package main

import (
	"context"
	"encoding/json"

	"github.com/dinno7/ride-sharing/shared/contracts"
	"github.com/dinno7/ride-sharing/shared/logger"
	"github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
)

type tripConsumer struct {
	logger logger.Logger
}

func NewTripConsumer(logger logger.Logger) *tripConsumer {
	return &tripConsumer{logger: logger}
}

func (c *tripConsumer) HandleTripCreatedEvent(
	ctx context.Context,
	message *amqp091.Delivery,
) error {
	var payload rabbitmq.Event[contracts.TripCreatedEventData]
	if err := json.Unmarshal(message.Body, &payload); err != nil {
		return err
	}
	c.logger.Info(
		"💀 New trip created",
		"trip_id", payload.Data.TripID,
		"package_slug", payload.Data.PackageSlug,
	)
	return nil
}
