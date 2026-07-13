package main

import (
	"github.com/dinno7/ride-sharing/shared/contracts"
	"github.com/dinno7/ride-sharing/shared/logger"
	messaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
)

const DriverTripAnswerQueue = "q.driver_trip_answer"

func setupMessageConsumer(
	conn *messaging.RabbitMQConnection,
	logger logger.Logger,
) (*messaging.Consumer, error) {
	consumer := messaging.NewConsumer(conn, logger)

	if err := consumer.DeclareAndBind(
		messaging.QueueConfig{
			Name:       DriverTripAnswerQueue,
			Durable:    true,
			DLXEnabled: true,
		},
		[]messaging.BindingConfig{{
			Exchange: messaging.ExchangeMain,
			RoutingKeys: []string{
				contracts.DriverCmdTripAccept,
				contracts.DriverCmdTripDecline,
			},
		}},
	); err != nil {
		return nil, err
	}

	return consumer, nil
}
