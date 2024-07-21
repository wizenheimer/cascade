package rest

import (
	"github.com/labstack/echo/v4"
)

// Inject routes onto the instance
func (rest *APIServer) injectRoutes(e *echo.Echo) {
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
