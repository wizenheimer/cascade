package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/wizenheimer/cascade/internal/processor"
	"github.com/wizenheimer/cascade/service/k8x"
)

func createScenario(c echo.Context) error {
	return processor.Process(c, k8x.Scenario)
}

func listScenario(c echo.Context) error {
	return processor.Process(c, k8x.Scenario)
}

func detailScenatio(c echo.Context) error {
	return processor.Process(c, k8x.Scenario)
}

func patchScenario(c echo.Context) error {
	return processor.Process(c, k8x.Scenario)
}

func createSession(c echo.Context) error {
	return processor.Process(c, k8x.Session)
}

func listSession(c echo.Context) error {
	return processor.Process(c, k8x.Session)
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
