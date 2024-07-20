package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/wizenheimer/cascade/internal/routes"
)

// Inject routes onto the instance
func injectRoutes(e *echo.Echo) {
	// =======================
	//       SCENARIO
	// =======================
	scenario := e.Group("/scenario")
	scenario.POST("", routes.CreateScenario)      // Create a scenario
	scenario.GET("", routes.ListScenario)         // List out all scenarios for the given team
	scenario.GET("/:id", routes.DetailScenario)   // List out properties of the scenario
	scenario.PATCH("/:id", routes.UpdateScenario) // Update properties of the scenario

	// =======================
	//       SESSION
	// =======================
	session := e.Group("/session")
	session.POST("/:scenario/:version", routes.CreateSession) // Trigger Chaos Experiment and Stream Logs via SSE

	// =======================
	//      METRIC
	// =======================
	metric := e.Group("/metric")
	metric.GET("", routes.GetMetrics) // Get Metrics for the given Scenario via Query Params

	// =======================
	//      TEAM
	// =======================
	team := e.Group("/team")
	team.POST("", routes.CreateTeam)            // Create a Team
	team.DELETE("/:id", routes.DeleteTeam)      // Delete a Team
	team.GET("/:id/users", routes.ListUsers)    // List Team Users
	team.PATCH("/:id", routes.ManageTeam)       // Implement Team Attribute Management
	team.POST("/:id/users", routes.ManageUsers) // Implement User Management

	// =======================
	//     USER
	// =======================
	auth := e.Group("/auth")
	auth.POST("", routes.SignUp)            // Sign Up a User
	auth.POST("/:id", routes.Login)         // Login a User
	auth.POST("/:id/logout", routes.Logout) // Logout a User
	auth.POST("/:id/delete", routes.Churn)  // Delete a User
}
