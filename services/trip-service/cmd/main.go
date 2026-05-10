package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	infraHttp "github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/http"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/osrm"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/repository/inmem"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/service"
	"github.com/dinno7/ride-sharing/shared/env"
)

var httpAddr = env.GetString("HTTP_ADDR", ":7000")

func main() {
	tripRepo := inmem.NewInMemTripRepository()
	osrmRouteCalculator := osrm.NewRouteCalculator()
	tripService := service.NewTripService(tripRepo, osrmRouteCalculator)
	tripHandler := infraHttp.NewTripHttpHandler(tripService)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		<-time.Tick(time.Second * 5)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from Trip Service"))
	})

	mux.HandleFunc("POST /preview", tripHandler.PreviewTrip)

	server := &http.Server{Addr: httpAddr, Handler: mux}

	stopServerCtx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	go func() {
		log.Println("Starting Trip Service")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	<-stopServerCtx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	log.Println("Shuting down Trip Service HTTP Server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln("Failed to shutdown Trip Service HTTP Server: %w", err)
	}
	log.Println("Trip Service server shutdown successfully")
}
