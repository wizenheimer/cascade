package database

import (
	"context"
	"errors"

	"github.com/wizenheimer/cascade/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateUser creates a new user, duh
func (c Client) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	user.ID = uuid.NewString()
	result := c.DB.WithContext(ctx).Create(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, &ConflictError{}
		}
		return nil, result.Error
	}

	return user, nil
}

func (c Client) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := c.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &NotFoundError{}
		}
		return nil, result.Error
	}
	return &user, nil
}

func (c Client) UpdateUser(ctx context.Context, email string, updatedUser *models.User) (*models.User, error) {
	user, err := c.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if updatedUser.Email != "" {
		user.Email = updatedUser.Email
	}
	if updatedUser.Role != "" {
		user.Role = updatedUser.Role
	}
	result := c.DB.WithContext(ctx).Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (c Client) DeleteUser(ctx context.Context, email string) (*models.User, error) {
	user, err := c.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	result := c.DB.WithContext(ctx).Delete(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (c Client) DeactivateUser(ctx context.Context, email string) (*models.User, error) {
	user, err := c.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	user.IsActive = false
	result := c.DB.WithContext(ctx).Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
