package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/application/service"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/events"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/osrm"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/repository/inmem"
	"github.com/dinno7/ride-sharing/shared/env"
	"github.com/dinno7/ride-sharing/shared/logger"
	rmqMessaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
	googlegrpc "google.golang.org/grpc"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":7000")
	grpcAddr = env.GetString("GRPC_ADDR", ":9000")
	amqpURI  = env.GetString("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/")
)

func main() {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	appLogger := logger.NewSlogLogger(slogger)

	rmqConnection, err := rmqMessaging.NewRabbitMQBroker(amqpURI, appLogger)
	if err != nil {
		log.Fatal(err)
	}
	defer rmqConnection.Close()

	rmqPublisher := SetupMessagePublisher(rmqConnection, appLogger)
	tripEventHandler := events.NewTripEventHandler(rmqPublisher)

	tripRepo := inmem.NewInMemTripRepository()
	osrmRouteCalculator := osrm.NewRouteCalculator()
	tripService := service.NewTripService(tripRepo, osrmRouteCalculator, tripEventHandler)

	addr, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
		return
	}
	server := googlegrpc.NewServer()
	grpc.NewTripGrpcHandler(server, tripService)

	// INFO: Consuming messaging broker
	rmqConsumer, err := setupMessageConsumer(rmqConnection, appLogger)
	if err != nil {
		panic(err)
	}
	driverConsumerHandler := events.NewDriverConsumerHandler(tripService, rmqPublisher)

	if err := rmqConsumer.Consume(DriverTripAnswerQueue, driverConsumerHandler.Handle); err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Println("Starting Trip Service")
		if err := server.Serve(addr); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	server.GracefulStop()
}

func SetupMessagePublisher(
	conn *rmqMessaging.RabbitMQConnection,
	logger logger.Logger,
) *rmqMessaging.Publisher {
	return rmqMessaging.NewPublisher(conn, "trips")
}
