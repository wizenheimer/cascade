package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wizenheimer/cascade/internal/config"
	"github.com/wizenheimer/cascade/service/database"
	"go.uber.org/zap"
)

// Initialize API Instance
func NewAPIServer(logger *zap.Logger) APIServer {
	db, err := initalizeDBClient(logger)
	if err != nil {
		logger.Fatal("failed to initialize Database Client", zap.Any("error", err))
	}

	api := APIServer{
		// Inject Logger
		Logger: logger,
		// Inject Database Client
		DB: db,
	}

	// Create Echo
	e := echo.New()

	// Add request ID middleware
	e.Use(middleware.RequestID())

	// Create a server
	s := http.Server{
		Addr:    config.SERVER_PORT,
		Handler: e,
	}

	// Inject Routes
	api.injectRoutes(e)

	// Inject Server
	api.server = &s

	// Return Instance
	return api
}

// Initalize Database Client using environment variables
func initalizeDBClient(logger *zap.Logger) (database.DatabaseClient, error) {
	// Load .env file incase of local development
	if os.Getenv("ENVIRONMENT") != "docker" {
		err := godotenv.Load("../.env")
		if err != nil {
			logger.Fatal("couldn't load .env file", zap.Any("error", err))
		}
	}
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	sslMode := os.Getenv("SSL_MODE")
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		logger.Fatal("failed to convert port to int", zap.Any("error", err))
	}

	// Create Database Client
	db, err := database.NewDatabaseClient(host, user, password, dbname, int32(port), sslMode)
	if err != nil {
		logger.Fatal("failed to initialize Database Client", zap.Any("error", err))
	} else {
		logger.Info("Connected to Database", zap.Any("host", host), zap.Any("user", user), zap.Any("dbname", dbname), zap.Any("port", port), zap.Any("sslMode", sslMode))
	}

	return db, nil
}

// Trigger Serving
func (api *APIServer) Serve() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := api.server.ListenAndServe(); err != http.ErrServerClosed {
			api.Logger.Fatal(err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := api.server.Shutdown(ctx); err != nil {
		api.Logger.Error(err.Error())
	}
}
