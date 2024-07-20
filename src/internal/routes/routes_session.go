package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateSession(c echo.Context) error {
	// TODO: Stream Logs via SSE
	return c.NoContent(http.StatusOK) // TODO
}
