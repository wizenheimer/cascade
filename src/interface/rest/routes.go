package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/wizenheimer/cascade/internal/logger"
	k8x "github.com/wizenheimer/cascade/service/kubernetes"
	"go.uber.org/zap"
)

// Inject routes onto the instance
func (rest *APIServer) injectRoutes(e *echo.Echo) {
	// =======================
	//       QuickStart
	// =======================
	e.POST("/quickstart", QuickStart) // QuickStart Endpoint for Stateless Runs
	// =======================
	//       SCENARIO
	// =======================
	scenario := e.Group("/scenario")
	scenario.POST("", rest.CreateScenario)      // Create a scenario
	scenario.GET("", rest.ListScenario)         // List out all scenarios for the given team
	scenario.GET("/:id", rest.DetailScenario)   // List out properties of the scenario
	scenario.PATCH("/:id", rest.UpdateScenario) // Update properties of the scenario

	// =======================
	//       SESSION
	// =======================
	session := e.Group("/session")
	session.POST("/:scenario/:version", rest.CreateSession) // Trigger Chaos Experiment and Stream Logs via SSE

	// =======================
	//      METRIC
	// =======================
	metric := e.Group("/metric")
	metric.GET("", rest.GetMetrics) // Get Metrics for the given Scenario via Query Params

	// =======================
	//      TEAM
	// =======================
	team := e.Group("/team")
	team.POST("", rest.CreateTeam)            // Create a Team
	team.DELETE("/:id", rest.DeleteTeam)      // Delete a Team
	team.GET("/:id/users", rest.ListUsers)    // List Team Users
	team.PATCH("/:id", rest.ManageTeam)       // Implement Team Attribute Management
	team.POST("/:id/users", rest.ManageUsers) // Implement User Management

	// =======================
	//     USER
	// =======================
	auth := e.Group("/auth")
	auth.POST("", rest.SignUp)            // Sign Up a User
	auth.POST("/:id", rest.Login)         // Login a User
	auth.POST("/:id/logout", rest.Logout) // Logout a User
	auth.POST("/:id/delete", rest.Churn)  // Delete a User
}

// QuickStart is the handler for the QuickStart endpoint
func QuickStart(c echo.Context) error {
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

	sessionID := c.Request().Header.Get("X-Request-ID")
	if sessionID == "" {
		sessionID = c.Response().Header().Get(echo.HeaderXRequestID)
	}

	scenario := c.Param("scenario")
	if scenario == "" {
		scenario = "undefined"
	}

	// Create a new context that will be cancelled when the request is done
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()

	// Parse Configs
	cc, tc, rc, err := ParseConfigsFromContext(c)
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
	}(executor, sessionID, ctx, ticker.C)

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
