package database

import (
	"context"
	"errors"
	"github.com/wizenheimer/cascade/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateTeam creates a new team, duh
func (c Client) CreateTeam(ctx context.Context, name string, description string, creator *models.User) (*models.Team, error) {
	if creator == nil {
		return nil, &ConflictError{}
	}

	team := &models.Team{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		IsActive:    true,
	}

	// Create the team
	result := c.DB.WithContext(ctx).Create(&team)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &ConflictError{}
		}
		return nil, result.Error
	}

	// Add the creator to the team
	if _, err := c.AddUserToTeam(ctx, creator, team); err != nil {
		return nil, err
	}

	return team, nil
}

func (c Client) GetTeamByID(ctx context.Context, teamID string) (*models.Team, error) {
	var team models.Team

	// Fetch the team by ID
	result := c.DB.WithContext(ctx).Where("id = ?", teamID).First(&team)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &NotFoundError{}
		}
		return nil, result.Error
	}
	return &team, nil
}

func (c Client) AddUserToTeam(ctx context.Context, user *models.User, team *models.Team) (*models.User, error) {
	user.Teams = append(user.Teams, *team)
	if _, err := c.UpdateUser(ctx, user.Email, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c Client) RemoveUserFromTeam(ctx context.Context, user *models.User, team *models.Team) (*models.User, error) {
	// Remove the team from the user
	for i, t := range user.Teams {
		if t.ID == team.ID {
			user.Teams = append(user.Teams[:i], user.Teams[i+1:]...)
			break
		}
	}

	// Update the user
	if _, err := c.UpdateUser(ctx, user.Email, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c Client) ListUsersByTeam(ctx context.Context, teamID string) ([]models.User, error) {
	var users []models.User
	result := c.DB.
		WithContext(ctx).
		Joins("JOIN cascade.user_team ON cascade.users.user_id = cascade.user_team.user_id").
		Where("cascade.user_team.team_id = ?", teamID).
		Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (c Client) UpdateTeambyTeamID(ctx context.Context, teamID string, updatedTeam *models.Team) (*models.Team, error) {
	team, err := c.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	if updatedTeam.Name != "" {
		team.Name = updatedTeam.Name
	}

	if updatedTeam.Description != "" {
		team.Description = updatedTeam.Description
	}

	result := c.DB.WithContext(ctx).Save(&team)
	if result.Error != nil {
		return nil, result.Error
	}

	return team, nil
}

func (c Client) DeleteTeam(ctx context.Context, teamID string) (*models.Team, error) {
	team, err := c.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	result := c.DB.WithContext(ctx).Delete(&team)
	if result.Error != nil {
		return nil, result.Error
	}

	return team, nil
}

func (c Client) DeactivateTeam(ctx context.Context, teamID string) (*models.Team, error) {
	team, err := c.GetTeamByID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	team.IsActive = false
	result := c.DB.WithContext(ctx).Save(&team)
	if result.Error != nil {
		return nil, result.Error
	}

	return team, nil
}
