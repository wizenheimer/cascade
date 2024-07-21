package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (client *APIServer) CreateScenario(c echo.Context) error {
	// TODO: Create a Scenario using YAML
	return c.NoContent(http.StatusOK) // TODO
}

func (client *APIServer) ListScenario(c echo.Context) error {
	// TODO: List out Scenario for the given Team by means of Query Params
	return c.NoContent(http.StatusOK) // TODO
}

func (client *APIServer) DetailScenario(c echo.Context) error {
	// TODO: List out properties of the scenario
	return c.NoContent(http.StatusOK) // TODO
}

func (client *APIServer) UpdateScenario(c echo.Context) error {
	// TODO: Bump up the version + Update properties of the scenario
	return c.NoContent(http.StatusOK) // TODO
}
