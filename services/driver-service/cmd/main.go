package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dinno7/ride-sharing/shared/env"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var grpcAddr = env.GetString("GRPC_ADDR", ":9000")

func main() {
	addr, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	driverService := NewDriverService()

	server := googlegrpc.NewServer()

	reflection.Register(server)

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
