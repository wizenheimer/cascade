package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetMetrics(c echo.Context) error {
	// TODO: Get Metrics for the given Scenario via Query Params
	return c.NoContent(http.StatusOK) // TODO
}
