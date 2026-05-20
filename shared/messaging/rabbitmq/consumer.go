package rabbitmq

import (
	"context"
	"fmt"

	"github.com/dinno7/ride-sharing/shared/logger"
)

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	DLXEnabled bool
}

type BindingConfig struct {
	Exchange    string
	RoutingKeys []string
}

type Consumer struct {
	conn   *RabbitMQConnection
	logger logger.Logger
}

func NewConsumer(conn *RabbitMQConnection, logger logger.Logger) *Consumer {
	return &Consumer{
		conn:   conn,
		logger: logger,
	}
}

// DeclareAndBind creates a queue and binds it — called by each service for ITS OWN queues
func (c *Consumer) DeclareAndBind(queue QueueConfig, bindings []BindingConfig) error {
	if queue.DLXEnabled {
		// TODO:
	}

	q, err := c.conn.Channel().QueueDeclare(
		queue.Name,
		queue.Durable,
		queue.AutoDelete,
		queue.Exclusive,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for _, bd := range bindings {
		for _, routeKey := range bd.RoutingKeys {
			err := c.conn.Channel().QueueBind(q.Name, routeKey, bd.Exchange, false, nil)
			if err != nil {
				return fmt.Errorf("failed to bind queue %s to %s with key %s: %w",
					q.Name, bd.Exchange, routeKey, err)
			}
			c.logger.Info(
				"bound queue",
				"queue", q.Name,
				"exchange", bd.Exchange,
				"routing_key", routeKey,
			)
		}
	}

	return nil
}

func (c *Consumer) Consume(queueName string, handler MessageHandler) error {
	err := c.conn.Channel().Qos(1, 0, false)
	if err != nil {
		return err
	}

	messages, err := c.conn.Channel().
		ConsumeWithContext(context.Background(),
			queueName,
			"", // consumer tag (auto-generated)
			false,
			false,
			false,
			false,
			nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			ctx := context.Background()
			if err := handler(ctx, &msg); err != nil {
				c.logger.Error(
					"failed to handle message",
					err,
					"queue", queueName,
					"message_id", msg.MessageId,
				)
				if err := msg.Nack(false, false); err != nil {
					c.logger.Error(
						"failed to nack message",
						err,
						"queue", queueName,
						"type", "NACK",
						"message_id", msg.MessageId,
					)
				}
			} else {
				if err := msg.Ack(false); err != nil {
					c.logger.Error(
						"failed to ack message",
						err,
						"queue", queueName,
						"type", "ACK",
						"message_id", msg.MessageId,
					)
				}
			}
		}
	}()
	return nil
}
