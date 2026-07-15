package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	echocustoms "github.com/dinno7/ride-sharing/services/api-gateway/cmd/echo-customs"
	"github.com/dinno7/ride-sharing/services/api-gateway/cmd/ws"
	"github.com/dinno7/ride-sharing/shared/env"
	"github.com/dinno7/ride-sharing/shared/logger"
	messaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
	"github.com/dinno7/ride-sharing/shared/tracing"
	wshub "github.com/dinno7/ride-sharing/shared/ws"
	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const serviceName = "api-gateway"

var (
	httpAddr        = env.GetString("HTTP_ADDR", ":7000")
	amqpURI         = env.GetString("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/")
	appEnv          = env.GetString("APP_ENV", "development")
	otelExporterURL = env.GetString("JAEGER_URL", "http://jaeger:4318")
	tracer          = tracing.GetTracer(serviceName)
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	e := echo.New()
	appLogger := logger.NewSlogLogger(e.Logger)

	tracingCfg := tracing.OTelConfig{
		ServiceName: serviceName,
		Environment: appEnv,
		ExporterURL: otelExporterURL,
	}

	tp, err := tracing.InitTracer(tracingCfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			appLogger.Error("Failed to shutdown tracer provider", err)
		}
	}()

	e.Validator = echocustoms.NewEchoValidator()
	e.HTTPErrorHandler = echocustoms.CustomHTTPErrorHandler
	e.Use(middleware.BodyLimit(2_097_152)) // 2MB
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS("*"))
	e.Use(middleware.Recover())
	e.Use(echootel.NewMiddlewareWithConfig(echootel.Config{
		ServerName:     serviceName,
		TracerProvider: tp,
	}))

	sc := echo.StartConfig{
		Address:         httpAddr,
		GracefulTimeout: time.Second * 3,
		BeforeServeFunc: func(s *http.Server) error {
			e.Logger.Info("Starting API Gateway")
			return nil
		},
	}

	wsHub := wshub.NewHub()
	go wsHub.Run()

	messagingConn, err := messaging.NewRabbitMQBroker(amqpURI, appLogger)
	if err != nil {
		panic(err)
	}
	consumer, err := setupMessageConsumer(messagingConn, appLogger)
	if err != nil {
		panic(err)
	}
	publisher := setupMessagePublisher(messagingConn)

	driverWsHandler := ws.NewDriverWSHandler(wsHub, publisher)
	riderWsHandler := ws.NewRiderWSHandler(wsHub)

	// INFO: Consuming queues
	messageConsumerHandler := NewMessageConsumerHandler(wsHub, appLogger)
	consumeAndForwardQueues := []string{
		DriverCmdTripRequestQueue,
		NotifyNoDriversFoundQueue,
		NotifyDriverAssignedQueue,
		NotifyPaymentSessionCreatedQueue,
	}
	for _, queueName := range consumeAndForwardQueues {
		if err = consumer.Consume(
			queueName,
			messageConsumerHandler.Forward,
		); err != nil {
			panic(err)
		}
	}

	e.POST("/trip/preview", handleTripPreview, echocustoms.TracingMiddlewareWithName("previewTrip"))
	e.POST("/trip/start", handleTripStart, echocustoms.TracingMiddlewareWithName("startTrip"))
	e.GET(
		"/ws/drivers",
		driverWsHandler.Handle,
		echocustoms.TracingMiddlewareWithName("driverWS"),
	)
	e.GET("/ws/riders", riderWsHandler.Handle, echocustoms.TracingMiddlewareWithName("riderdWS"))
	e.POST(
		"/webhook/stripe",
		handleStripeWebhook,
		echocustoms.TracingMiddlewareWithName("stripeWebhook"),
	)

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
