package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dinno7/ride-sharing/services/payment-service/internal/application"
	"github.com/dinno7/ride-sharing/services/payment-service/internal/domain"
	"github.com/dinno7/ride-sharing/services/payment-service/internal/infrastructure/events"
	"github.com/dinno7/ride-sharing/services/payment-service/internal/infrastructure/payment"
	"github.com/dinno7/ride-sharing/shared/env"
	"github.com/dinno7/ride-sharing/shared/logger"
	messaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
)

var (
	GrpcAddr    = env.GetString("GRPC_ADDR", ":9004")
	rabbitMqURI = env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")
	appURL      = env.GetString("APP_URL", "http://localhost:3000")
)

func main() {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	appLogger := logger.NewSlogLogger(slogger)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stripe config
	stripeCfg := &domain.PaymentConfig{
		StripeSecretKey: env.GetString("STRIPE_SECRET_KEY", ""),
		SuccessURL:      env.GetString("STRIPE_SUCCESS_URL", appURL+"?payment=success"),
		CancelURL:       env.GetString("STRIPE_CANCEL_URL", appURL+"?payment=cancel"),
	}

	if stripeCfg.StripeSecretKey == "" {
		panic("STRIPE_SECRET_KEY is not set")
	}

	// RabbitMQ connection
	appLogger.Info("Starting RabbitMQ connection")
	rabbitmq := mustWithValue(messaging.NewRabbitMQBroker(rabbitMqURI, appLogger))
	defer rabbitmq.Close()

	consumer := mustWithValue(setupConsumers(rabbitmq, appLogger))
	publisher := messaging.NewPublisher(rabbitmq, "payment")

	paymentProcessor := payment.NewStripePaymentProcessor(stripeCfg)
	paymentService := application.NewPaymentService(paymentProcessor)

	paymentSessionCreatorConsumer := events.NewPaymentConsumer(publisher, paymentService)

	consumer.Consume(
		PaymentCreateSessionQueue,
		paymentSessionCreatorConsumer.HandleCreatePaymentSession,
	)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	appLogger.Info("Shutting down payment service...")
}

func mustWithValue[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
