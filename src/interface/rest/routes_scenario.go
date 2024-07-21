package rest

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wizenheimer/cascade/internal/config"
	"github.com/wizenheimer/cascade/internal/parser"
	"gopkg.in/yaml.v2"
)

func (client *APIServer) CreateScenario(c echo.Context) error {
	// Create a Scenario using YAML
	var config config.Config

	// Handle file upload
	file, err := c.FormFile("config")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Read and parse YAML file
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	// Parse the Config
	scenario, err := parser.ParseYAMLConfigToScenario(&config)
	if err != nil {
		return err
	}

	// Persis the Scenario
	_, err = client.DB.CreateScenario(c.Request().Context(), scenario)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK) // TODO
}

func (client *APIServer) ListScenario(c echo.Context) error {
	// List out Scenario for the given Team by means of Query Params
	team := c.QueryParam("team")

	scenarios, err := client.DB.ListScenariosByTeamID(c.Request().Context(), team)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, scenarios)
}

func (client *APIServer) DetailScenario(c echo.Context) error {
	// Get the scenario ID from the request parameters
	scenarioID := c.Param("id")

	// Retrieve the scenario from the database
	scenario, err := client.DB.GetScenarioByID(c.Request().Context(), scenarioID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, scenario)
}

func (client *APIServer) UpdateScenario(c echo.Context) error {
	// Handle file upload
	file, err := c.FormFile("config")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Read and parse YAML file
	data, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	var config config.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	// Parse the Config
	updatedScenario, err := parser.ParseYAMLConfigToScenario(&config)
	if err != nil {
		return err
	}

	// Persist the updated scenario
	_, err = client.DB.UpdateScenario(c.Request().Context(), updatedScenario.ID, updatedScenario)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
