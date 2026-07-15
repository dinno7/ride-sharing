package events

import (
	"context"
	"encoding/json"

	"github.com/dinno7/ride-sharing/services/payment-service/internal/domain"
	"github.com/dinno7/ride-sharing/shared/contracts"
	messaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
)

type paymentConsumer struct {
	paymentService domain.PaymentService
	publisher      *messaging.Publisher
}

func NewPaymentConsumer(
	publisher *messaging.Publisher,
	paymentService domain.PaymentService,
) *paymentConsumer {
	return &paymentConsumer{
		publisher:      publisher,
		paymentService: paymentService,
	}
}

func (c *paymentConsumer) HandleCreatePaymentSession(
	ctx context.Context,
	message *amqp091.Delivery,
) error {
	if message.RoutingKey != contracts.PaymentCmdCreateSession {
		return nil
	}

	var payload messaging.MessageInfo[contracts.PaymentCmdCreateSessionData]
	if err := json.Unmarshal(message.Body, &payload); err != nil {
		return err
	}

	payment, err := c.paymentService.CreatePaymentSession(
		ctx,
		payload.Data.TripID,
		payload.Data.UserID,
		payload.Data.DriverID,
		int64(payload.Data.Amount),
		payload.Data.Currency,
	)
	if err != nil {
		return err
	}
	if err := c.publisher.PublishEvent(
		ctx,
		contracts.PaymentEventSessionCreated,
		payload.OwnerID,
		contracts.PaymentEventSessionCreatedData{
			TripID:    payment.TripID,
			SessionID: payment.StripeSessionID,
			Amount:    float64(payment.Amount),
			Currency:  payment.Currency,
		},
	); err != nil {
		return err
	}

	return nil
}
