package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	echocustoms "github.com/dinno7/ride-sharing/services/api-gateway/cmd/echo-customs"
	"github.com/dinno7/ride-sharing/shared/env"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

var httpAddr = env.GetString("HTTP_ADDR", ":7000")

func main() {
	e := echo.New()
	e.Validator = echocustoms.NewEchoValidator()
	e.HTTPErrorHandler = echocustoms.CustomHTTPErrorHandler
	e.Use(middleware.BodyLimit(2_097_152)) // 2MB
	e.Pre(middleware.RemoveTrailingSlash())

	sc := echo.StartConfig{
		Address:         httpAddr,
		GracefulTimeout: time.Second * 3,
		BeforeServeFunc: func(s *http.Server) error {
			e.Logger.Info("Starting API Gateway")
			return nil
		},
	}

	e.POST("/trip/preview", handleTripPreview)

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
