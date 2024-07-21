package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/wizenheimer/cascade/internal/logger"
	k8x "github.com/wizenheimer/cascade/service/kubernetes"
	"go.uber.org/zap"
)

func (client *APIServer) CreateSession(c echo.Context) error {
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

	// Create a new context that will be cancelled when the request is done
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	scenarioStr := c.Param("scenario")
	scenarioVersionStr := c.Param("version")
	version, err := strconv.Atoi(scenarioVersionStr)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// Fetch the scenario
	scenario, err := client.DB.GetScenarioByIDByVersion(c.Request().Context(), scenarioStr, version)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	// Trigger a session
	session, err := client.DB.CreateSession(c.Request().Context(), scenarioStr, version)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Parse cluster configs
	cc, err := ParseClusterConfigFromContext(c)
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

	// Parse the target and runtime config
	tc, rc, err := ParseDBScenario(scenario)
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

	// Create executor
	executor, err := k8x.CreateExecutor(cc, tc, rc, logger)
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

	// Create ticker
	ticker := time.NewTicker(rc.Interval)
	defer ticker.Stop()

	// Start processing in a goroutine
	go func(executor *k8x.Executor, sessionID string, ctx context.Context, next <-chan time.Time) {
		for {
			executor.Logger.Info("Chaos Session Triggered", zap.Any("Session", sessionID))
			executor.Logger.Info("Chaos Scenario", zap.Any("Scenario", scenario))

			// Trigger Execution
			if err = executor.Execute(ctx); err != nil {
				executor.Logger.Error(err.Error())
			}

			select {
			case <-next:
				// trigger next session
			case <-ctx.Done():
				// skip subsequent execution
				return
			}
		}
	}(executor, string(session.ID), ctx, ticker.C)

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
			client.DB.GracefullyEndSession(c.Request().Context(), string(session.ID))
			return nil
		}
	}
}
