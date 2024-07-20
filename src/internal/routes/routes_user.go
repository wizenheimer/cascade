package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SignUp(c echo.Context) error {
	// TODO: Create a Team
	return c.NoContent(http.StatusOK) // TODO
}

func Login(c echo.Context) error {
	// TODO: Create a Team
	return c.NoContent(http.StatusOK) // TODO
}

func Logout(c echo.Context) error {
	// TODO: Create a Team
	return c.NoContent(http.StatusOK) // TODO
}

func Churn(c echo.Context) error {
	return c.NoContent(http.StatusOK) // TODO
}
