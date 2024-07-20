package database

import (
	"context"
	"errors"

	"github.com/wizenheimer/cascade/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateScenario creates and persists a new scenario
func (c Client) CreateScenario(ctx context.Context, scenario *models.Scenario) (*models.Scenario, error) {
	// Generate a new UUID for the scenario
	scenario.ID = uuid.NewString()
	scenario.Version = 1

	// Create the scenario
	result := c.DB.Create(scenario)
	if result.Error != nil {
		return nil, result.Error
	}

	return scenario, nil
}

// UpdateScenario updates an existing scenario or attempts to create one, incase it does not exist
func (c Client) UpdateScenario(ctx context.Context, scenarioID string, updatedScenario *models.Scenario) (*models.Scenario, error) {
	// Find the highest versioned scenario
	var latestScenario models.Scenario
	result := c.DB.Where("scenario_id = ?", scenarioID).Order("version DESC").First(&latestScenario)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.CreateScenario(ctx, updatedScenario)
		}
		return nil, result.Error
	}

	// Create a new scenario based on the latest one
	newScenario := latestScenario

	// Update fields from updatedScenario
	newScenario.Description = updatedScenario.Description
	if updatedScenario.TeamID != "" {
		newScenario.TeamID = updatedScenario.TeamID
	}

	// Increment version and update timestamps
	newScenario.Version = latestScenario.Version + 1

	// Save the new version
	result = c.DB.Create(&newScenario)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newScenario, nil
}

// DeleteScenario deletes a scenario and all its associated sessions
func (c Client) DeleteScenario(ctx context.Context, scenarioID string) error {
	tx := c.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&models.Scenario{}, scenarioID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("scenario_id = ?", scenarioID).Delete(&models.Session{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// ListScenario returns a list of scenarios inside an Organization
func (c Client) ListScenarios(ctx context.Context) ([]models.Scenario, error) {
	var scenarios []models.Scenario
	result := c.DB.Select("scenario_id", "version", "created_at").Find(&scenarios)
	return scenarios, result.Error
}

// ListScenarioVersion returns a list of versions of a given scenario
func (c Client) ListScenarioVersion(ctx context.Context, scenarioID string) ([]models.Scenario, error) {
	var scenarios []models.Scenario
	result := c.DB.Select("scenario_id", "version", "created_at").Where("scenario_id = ?", scenarioID).Find(&scenarios)
	return scenarios, result.Error
}

// ListScenarioByTeamID returns a list of scenarios inside a Team
func (c Client) ListScenariosByTeamID(ctx context.Context, teamID string) ([]models.Scenario, error) {
	var scenarios []models.Scenario
	result := c.DB.Select("scenario_id").Where("team_id = ?", teamID).Find(&scenarios)
	return scenarios, result.Error
}

// GetScenarioByID returns a scenario by its ID
func (c Client) GetScenarioByID(ctx context.Context, scenarioID string) ([]models.Scenario, error) {
	var scenarios []models.Scenario
	result := c.DB.Select("scenario_name", "version", "description", "created_at").Where("scenario_id = ?", scenarioID).Find(&scenarios).Order("version DESC")
	return scenarios, result.Error
}
