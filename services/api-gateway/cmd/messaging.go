package main

import (
	"github.com/dinno7/ride-sharing/shared/contracts"
	"github.com/dinno7/ride-sharing/shared/logger"
	rmqMessaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
)

const (
	// INFO: For driver
	DriverCmdTripRequestQueue = "q.driver_cmd_trip_request"

	// INFO: For rider
	NotifyNoDriversFoundQueue        = "q.notify_no_driver_found"
	NotifyDriverAssignedQueue        = "q.notify_driver_assign"
	NotifyPaymentSessionCreatedQueue = "q.notify_payment_session_created"
)

func setupMessageConsumer(
	conn *rmqMessaging.RabbitMQConnection,
	logger logger.Logger,
) (*rmqMessaging.Consumer, error) {
	consumer := rmqMessaging.NewConsumer(conn, logger)

	// INFO: Drivers
	if err := consumer.DeclareAndBind(rmqMessaging.QueueConfig{
		Name:       DriverCmdTripRequestQueue,
		Durable:    true,
		DLXEnabled: true,
		Exclusive:  false,
		AutoDelete: false,
	}, []rmqMessaging.BindingConfig{
		{
			Exchange: rmqMessaging.ExchangeMain,
			RoutingKeys: []string{
				contracts.DriverCmdTripRequest,
			},
		},
	}); err != nil {
		return nil, err
	}

	// INFO: Riders
	if err := consumer.DeclareAndBind(rmqMessaging.QueueConfig{
		Name:       NotifyNoDriversFoundQueue,
		Durable:    true,
		DLXEnabled: true,
		Exclusive:  false,
		AutoDelete: false,
	}, []rmqMessaging.BindingConfig{
		{
			Exchange: rmqMessaging.ExchangeMain,
			RoutingKeys: []string{
				contracts.TripEventNoDriversFound,
			},
		},
	}); err != nil {
		return nil, err
	}

	if err := consumer.DeclareAndBind(rmqMessaging.QueueConfig{
		Name:       NotifyDriverAssignedQueue,
		Durable:    true,
		DLXEnabled: true,
		Exclusive:  false,
		AutoDelete: false,
	}, []rmqMessaging.BindingConfig{
		{
			Exchange: rmqMessaging.ExchangeMain,
			RoutingKeys: []string{
				contracts.TripEventDriverAssigned,
			},
		},
	}); err != nil {
		return nil, err
	}

	if err := consumer.DeclareAndBind(rmqMessaging.QueueConfig{
		Name:       NotifyPaymentSessionCreatedQueue,
		Durable:    true,
		DLXEnabled: true,
		Exclusive:  false,
		AutoDelete: false,
	}, []rmqMessaging.BindingConfig{
		{
			Exchange: rmqMessaging.ExchangeMain,
			RoutingKeys: []string{
				contracts.PaymentEventSessionCreated,
			},
		},
	}); err != nil {
		return nil, err
	}

	return consumer, nil
}

func setupMessagePublisher(
	conn *rmqMessaging.RabbitMQConnection,
) *rmqMessaging.Publisher {
	return rmqMessaging.NewPublisher(conn, "gateway")
}
