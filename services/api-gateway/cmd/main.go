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
	wshub "github.com/dinno7/ride-sharing/shared/ws"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":7000")
	amqpURI  = env.GetString("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/")
)

func main() {
	e := echo.New()
	appLogger := logger.NewSlogLogger(e.Logger)

	e.Validator = echocustoms.NewEchoValidator()
	e.HTTPErrorHandler = echocustoms.CustomHTTPErrorHandler
	e.Use(middleware.BodyLimit(2_097_152)) // 2MB
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS("*"))
	e.Use(middleware.Recover())

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
	}
	for _, queueName := range consumeAndForwardQueues {
		if err = consumer.Consume(
			queueName,
			messageConsumerHandler.Forward,
		); err != nil {
			panic(err)
		}
	}

	e.POST("/trip/preview", handleTripPreview)
	e.POST("/trip/start", handleTripStart)
	e.GET("/ws/drivers", driverWsHandler.Handle)
	e.GET("/ws/riders", riderWsHandler.Handle)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
