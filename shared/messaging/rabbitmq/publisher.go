package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn       *RabbitMQConnection
	sourceName string // identifies which service is publishing
}

func NewPublisher(conn *RabbitMQConnection, serviceName string) *Publisher {
	return &Publisher{
		conn:       conn,
		sourceName: serviceName,
	}
}

func (p *Publisher) PublishEvent(ctx context.Context, eventName string, data any) error {
	event := Event[any]{
		ID:        uuid.New().String(),
		Type:      eventName,
		Source:    p.sourceName,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return p.conn.Channel().PublishWithContext(
		ctx,
		ExchangeMain, // exchange
		eventName,    // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    event.ID,
			Timestamp:    event.Timestamp,
			Body:         body,
		},
	)
}

func (p *Publisher) PublishCommand(ctx context.Context, commandName string, data any) error {
	event := Event[any]{
		ID:        uuid.New().String(),
		Type:      commandName,
		Source:    p.sourceName,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	return p.conn.Channel().PublishWithContext(
		ctx,
		ExchangeMain,
		commandName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    event.ID,
			Timestamp:    event.Timestamp,
			Body:         body,
		},
	)
}
