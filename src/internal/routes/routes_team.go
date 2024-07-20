package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateTeam(c echo.Context) error {
	// TODO: Create a Team
	return c.NoContent(http.StatusOK) // TODO
}

func DeleteTeam(c echo.Context) error {
	// TODO: Delete a Team
	return c.NoContent(http.StatusOK) // TODO
}

func ListUsers(c echo.Context) error {
	// TODO: List Team Users
	return c.NoContent(http.StatusOK) // TODO
}

func ManageTeam(c echo.Context) error {
	// TODO: Implement Team Attribute Management
	return c.NoContent(http.StatusOK) // TODO
}

func ManageUsers(c echo.Context) error {
	// TODO: Implement User Management
	return c.NoContent(http.StatusOK) // TODO
}
