package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (client *APIServer) SignUp(c echo.Context) error {
	// TODO: Create a Team
	return c.NoContent(http.StatusOK) // TODO
}

func (client *APIServer) Login(c echo.Context) error {
	// TODO: Create a Team
	return c.NoContent(http.StatusOK) // TODO
}

func (client *APIServer) Logout(c echo.Context) error {
	// TODO: Create a Team
	return c.NoContent(http.StatusOK) // TODO
}

func (client *APIServer) Churn(c echo.Context) error {
	return c.NoContent(http.StatusOK) // TODO
}
