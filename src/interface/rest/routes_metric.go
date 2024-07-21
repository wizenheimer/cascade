package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (client *APIServer) GetMetrics(c echo.Context) error {
	scenarioStr := c.Param("scenario")
	metrics, err := client.DB.GetSessionMetrics(c.Request().Context(), scenarioStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, metrics)
}
