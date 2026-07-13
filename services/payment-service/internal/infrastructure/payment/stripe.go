package payment

import (
	"context"

	"github.com/dinno7/ride-sharing/services/payment-service/internal/domain"
	"github.com/stripe/stripe-go/v86"
)

type stripePaymentProcessor struct {
	cfg *domain.PaymentConfig
}

func NewStripePaymentProcessor(cfg *domain.PaymentConfig) domain.PaymentProcessor {
	return &stripePaymentProcessor{
		cfg: cfg,
	}
}

func (h *stripePaymentProcessor) CreatePaymentSession(
	ctx context.Context,
	amount int64,
	currency string,
	metadata map[string]string,
) (string, error) {
	sc := stripe.NewClient(h.cfg.StripeSecretKey)
	params := &stripe.CheckoutSessionCreateParams{
		SuccessURL: stripe.String(h.cfg.SuccessURL),
		CancelURL:  stripe.String(h.cfg.CancelURL),
		Metadata:   metadata,
		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionCreateLineItemPriceDataParams{
					Currency: stripe.String(currency),
					ProductData: &stripe.CheckoutSessionCreateLineItemPriceDataProductDataParams{
						Name: stripe.String("Ride Payment"),
					},
					UnitAmount: stripe.Int64(amount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(stripe.CheckoutSessionModePayment),
	}
	result, err := sc.V1CheckoutSessions.Create(context.TODO(), params)
	if err != nil {
		return "", err
	}

	return result.ID, nil
}

func (h *stripePaymentProcessor) GetSessionStatus(
	ctx context.Context,
	sessionID string,
) (domain.PaymentStatus, error) {
	panic("not implemented") // TODO: Implement
}
