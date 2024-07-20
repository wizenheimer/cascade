package models

import "time"

type Team struct {
	ID          string    `gorm:"primaryKey;column:team_id" json:"id"`
	Name        string    `gorm:"column:team_name;size:100;not null" json:"name"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	IsActive    bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP()" json:"updated_at"`
}

type UserTeam struct {
	UserID int `gorm:"primaryKey;column:user_id" json:"user_id"`
	TeamID int `gorm:"primaryKey;column:team_id" json:"team_id"`
}
