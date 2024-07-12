package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wizenheimer/cascade/internal/logger"
)

func Process(c echo.Context, callback kubernetesCallback) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	logChan := make(chan logger.LogEntry, 100)
	logger := logger.CreateLogger(logChan)
	defer logger.Sync()

	requestID := c.Request().Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = c.Response().Header().Get(echo.HeaderXRequestID)
	}

	// Create a new context that will be cancelled when the request is done
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// Start processing in a goroutine
	go callback(ctx, logger, requestID)

	// Stream logs back to the client
	for {
		select {
		case logEntry := <-logChan:
			data, err := json.Marshal(logEntry)
			if err != nil {
				continue
			}
			_, err = c.Response().Write([]byte(fmt.Sprintf("data: %s\n\n", data)))
			if err != nil {
				return err
			}
			c.Response().Flush()
		case <-ctx.Done():
			logger.Info("Client disconnected, stopping log stream")
			return nil
		}
	}
}
