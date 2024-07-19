package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wizenheimer/cascade/internal/config"
	"go.uber.org/zap"
)

// Initialize API Instance
func NewAPIServer(logger *zap.Logger) APIServer {
	// Create Echo
	e := echo.New()

	// Add request ID middleware
	e.Use(middleware.RequestID())

	// Inject Routes
	injectRoutes(e)

	// Create a server
	s := http.Server{
		Addr:    config.SERVER_PORT,
		Handler: e,
	}

	api := APIServer{
		// Inject Server
		server: &s,
		// Inject Logger
		logger: logger,
	}

	// Return Instance
	return api
}

// Trigger Serving
func (api *APIServer) Serve() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := api.server.ListenAndServe(); err != http.ErrServerClosed {
			api.logger.Fatal(err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := api.server.Shutdown(ctx); err != nil {
		api.logger.Error(err.Error())
	}
}
