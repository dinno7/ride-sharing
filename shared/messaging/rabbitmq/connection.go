package rabbitmq

import (
	"context"

	"github.com/dinno7/ride-sharing/shared/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	MessageHandler     func(ctx context.Context, message *amqp.Delivery) error
	RabbitMQConnection struct {
		conn    *amqp.Connection
		channel *amqp.Channel
		logger  logger.Logger
	}
)

func NewRabbitMQBroker(uri string, logger logger.Logger) (*RabbitMQConnection, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	rmqConnection := &RabbitMQConnection{conn: conn, channel: ch, logger: logger}

	if err := rmqConnection.setupExchanges(); err != nil {
		rmqConnection.Close()
		return nil, err
	}

	return rmqConnection, nil
}

func (c *RabbitMQConnection) setupExchanges() error {
	for _, ex := range exchanges {
		if err := c.channel.ExchangeDeclare(
			ex.Name,
			ex.Kind,
			ex.Durable,
			ex.AutoDelete,
			false, // internal
			false, // no-wait
			nil,
		); err != nil {
			return err
		}
		c.logger.Info("declared exchange", "name", ex.Name, "kind", ex.Kind)
	}
	return nil
}

func (c *RabbitMQConnection) Channel() *amqp.Channel {
	return c.channel
}

func (c *RabbitMQConnection) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
