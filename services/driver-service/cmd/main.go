package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dinno7/ride-sharing/shared/contracts"
	"github.com/dinno7/ride-sharing/shared/env"
	"github.com/dinno7/ride-sharing/shared/logger"
	rmqMessaging "github.com/dinno7/ride-sharing/shared/messaging/rabbitmq"
	"github.com/dinno7/ride-sharing/shared/tracing"
	googlegrpc "google.golang.org/grpc"
)

const serviceName = "driver-service"

var (
	grpcAddr        = env.GetString("GRPC_ADDR", ":9000")
	amqpURI         = env.GetString("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/")
	appEnv          = env.GetString("APP_ENV", "development")
	otelExporterURL = env.GetString("JAEGER_URL", "http://jaeger:4318")
	tracer          = tracing.GetTracer(serviceName)
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	appLogger := logger.NewSlogLogger(slogger)

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

	addr, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	rabbitmq, err := rmqMessaging.NewRabbitMQBroker(amqpURI, appLogger)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()
	rmqConsumer, err := setupMessageConsumer(rabbitmq, appLogger)
	if err != nil {
		appLogger.Fatal("failed to declare message consumer", err)
	}
	rmqPublisher := setupMessagePublisher(rabbitmq)

	driverService := NewDriverService()

	tripConsumer := NewTripConsumer(driverService, rmqPublisher, appLogger)
	if err := rmqConsumer.Consume(FindAvailableDriverQueue, tripConsumer.Handle); err != nil {
		panic(err)
	}

	server := googlegrpc.NewServer(tracing.GRPCServerTracingOpts()...)

	NewDriverGrpcHandler(server, driverService)

	go func() {
		log.Println("Starting Driver Service")
		if err := server.Serve(addr); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	server.GracefulStop()
}

func setupMessagePublisher(conn *rmqMessaging.RabbitMQConnection) *rmqMessaging.Publisher {
	return rmqMessaging.NewPublisher(conn, "drivers")
}

func setupMessageConsumer(
	conn *rmqMessaging.RabbitMQConnection,
	logger logger.Logger,
) (*rmqMessaging.Consumer, error) {
	consumer := rmqMessaging.NewConsumer(conn, logger)

	err := consumer.DeclareAndBind(rmqMessaging.QueueConfig{
		Name:       FindAvailableDriverQueue,
		Durable:    true,
		DLXEnabled: true,
		Exclusive:  false,
		AutoDelete: false,
	}, []rmqMessaging.BindingConfig{
		{
			Exchange: rmqMessaging.ExchangeMain,
			RoutingKeys: []string{
				contracts.TripEventCreated,
				contracts.TripEventDriverNotInterested,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
