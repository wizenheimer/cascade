package models

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;column:user_id" json:"id"`
	Email     string    `gorm:"column:email;size:100;not null;unique" json:"email"`
	Role      string    `gorm:"column:role;size:20;not null;default:'user'" json:"role"`
	IsActive  bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP()" json:"updated_at"`
	Teams     []Team    `gorm:"many2many:user_team;" json:"teams"`
}
