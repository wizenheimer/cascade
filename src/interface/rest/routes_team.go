package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (client *APIServer) CreateTeam(c echo.Context) error {
	// Create a Team
	return c.NoContent(http.StatusOK) // TODO
}

func (client *APIServer) DeleteTeam(c echo.Context) error {
	// TODO: Delete a Team
	return c.NoContent(http.StatusOK) // TODO
}

func (client *APIServer) ListUsers(c echo.Context) error {
	// TODO: List Team Users
	return c.NoContent(http.StatusOK) // TODO
}

func (client *APIServer) ManageTeam(c echo.Context) error {
	// TODO: Implement Team Attribute Management
	return c.NoContent(http.StatusOK) // TODO
}

func (client *APIServer) ManageUsers(c echo.Context) error {
	// TODO: Implement User Management
	return c.NoContent(http.StatusOK) // TODO
}
