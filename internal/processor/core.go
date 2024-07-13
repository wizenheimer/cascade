package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/wizenheimer/cascade/internal/logger"
	"github.com/wizenheimer/cascade/service/k8x"
)

// Stream the response via Server Sent Event
func Process(c echo.Context, callback kubernetesCallback) error {
	// Set Headers
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	// Inject Loggers with Hook
	loggerWithChan := log.LoggerPool.Get().(*log.LoggerWithChannel)
	defer func() {
		// Clear the channel before putting the logger back in the pool
		for len(loggerWithChan.LogChan) > 0 {
			<-loggerWithChan.LogChan
		}
		log.LoggerPool.Put(loggerWithChan)
	}()

	logger := loggerWithChan.Logger
	logChan := loggerWithChan.LogChan

	requestID := c.Request().Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = c.Response().Header().Get(echo.HeaderXRequestID)
	}

	// Create a new context that will be cancelled when the request is done
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// Parse Configs
	_, _, _, _, err := k8x.ParseConfigs(c)
	if err != nil {
		// Parse the log
		data, err := log.ParseLog("error", err.Error())
		if err != nil {
			return err
		}

		// Send response back to client
		_, err = c.Response().Write(data)
		if err != nil {
			return err
		}

		c.Response().Flush()
		return nil
	}

	// Start processing in a goroutine
	go callback(ctx, logger, requestID)

	// Stream logs back to the client
	for {
		select {
		case logEntry := <-logChan:
			// Serialize logEntry to JSON
			data, err := json.Marshal(logEntry)
			if err != nil {
				continue
			}
			// Send log entry to client
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
