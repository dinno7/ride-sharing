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

func (p *Publisher) PublishEvent(ctx context.Context, routingKey, ownerID string, data any) error {
	messageInfo := MessageInfo[any]{
		ID:         uuid.New().String(),
		RoutingKey: routingKey,
		Kind:       MessageInfoKindEvent,
		OwnerID:    ownerID,
		Source:     p.sourceName,
		Data:       data,
	}

	return p.publish(ctx, messageInfo)
}

func (p *Publisher) PublishCommand(
	ctx context.Context,
	routingKey, ownerID string,
	data any,
) error {
	messageInfo := MessageInfo[any]{
		ID:         uuid.New().String(),
		RoutingKey: routingKey,
		Kind:       MessageInfoKindCommand,
		OwnerID:    ownerID,
		Source:     p.sourceName,
		Data:       data,
	}
	return p.publish(ctx, messageInfo)
}

func (p *Publisher) ReQueueWithNewID(
	ctx context.Context, msg *amqp.Delivery,
) error {
	var msgData MessageInfo[any]

	if err := json.Unmarshal(msg.Body, &msgData); err != nil {
		return err
	}

	return p.publish(ctx, msgData)
}

func (p *Publisher) publish(ctx context.Context, data MessageInfo[any]) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return p.conn.Channel().PublishWithContext(
		ctx,
		ExchangeMain,    // exchange
		data.RoutingKey, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    data.ID,
			Timestamp:    time.Now().UTC(),
			Body:         body,
		},
	)
}
