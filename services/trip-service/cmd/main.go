package main

import (
	"log"
	"net/http"

	infraHttp "github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/http"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/infrastructure/repository/inmem"
	"github.com/dinno7/ride-sharing/services/trip-service/internal/service"
	"github.com/dinno7/ride-sharing/shared/env"
)

var httpAddr = env.GetString("HTTP_ADDR", ":7000")

func main() {
	log.Println("Starting Trip Service")

	tripRepo := inmem.NewInMemTripRepository()
	tripService := service.NewTripService(tripRepo)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from Trip Service"))
	})

	tripHandler := infraHttp.NewTripHttpHandler(tripService)
	http.HandleFunc("POST /preview", tripHandler.CreateTrip)

	http.ListenAndServe(httpAddr, nil)
}
