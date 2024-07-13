package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wizenheimer/cascade/internal/processor"
)

func createScenario(c echo.Context) error {
	return c.NoContent(http.StatusOK) // TODO
}

func listScenario(c echo.Context) error {
	return c.NoContent(http.StatusOK) // TODO
}

func detailScenatio(c echo.Context) error {
	return c.NoContent(http.StatusOK) // TODO
}

func patchScenario(c echo.Context) error {
	return c.NoContent(http.StatusOK) // TODO
}

func createSession(c echo.Context) error {
	return processor.Process(c)
}

func listSession(c echo.Context) error {
	return c.NoContent(http.StatusOK) // TODO
}

// Inject routes onto the instance
func injectRoutes(e *echo.Echo) {
	// =======================
	//       SCENARIO
	// =======================
	scenario := e.Group("/scenario")
	scenario.POST("", createScenario)     // Create a scenario
	scenario.GET("", listScenario)        // List out all scenarios
	scenario.GET("/:id", detailScenatio)  // List out properties of the scenario
	scenario.PATCH("/:id", patchScenario) // Update properties of the scenario

	// =======================
	//       SESSION
	// =======================
	session := e.Group("/session")
	session.POST("/:scenario", createSession) // Stream Logs via SSE
	session.GET("/:scenario", listSession)    // List out sessions
}
