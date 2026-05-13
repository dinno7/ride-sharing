package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/osrm"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/repository/inmem"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/service"
	"github.com/dinno7/ride-sharing/shared/env"
	googlegrpc "google.golang.org/grpc"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":7000")
	grpcAddr = env.GetString("GRPC_ADDR", ":9000")
)

func main() {
	tripRepo := inmem.NewInMemTripRepository()
	osrmRouteCalculator := osrm.NewRouteCalculator()
	tripService := service.NewTripService(tripRepo, osrmRouteCalculator)

	addr, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	server := googlegrpc.NewServer()
	grpc.NewTripGrpcHandler(server, tripService)
	go func() {
		log.Println("Starting Trip Service")
		if err := server.Serve(addr); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	server.GracefulStop()
}
