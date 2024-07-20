package models

import "time"

// Session represents a chaos engineering experiment during execution
type Session struct {
	ID         int       `gorm:"primaryKey;column:session_id" json:"id"`
	ScenarioID string    `gorm:"column:scenario_id;not null" json:"scenario_id"`
	UserID     string    `gorm:"column:user_id;not null" json:"user_id"`
	StartTime  time.Time `gorm:"column:start_time;not null;default:CURRENT_TIMESTAMP()" json:"start_time"`
	EndTime    time.Time `gorm:"column:end_time" json:"end_time"`
	Status     string    `gorm:"column:status;size:20;not null" json:"status"`
	Scenario   Scenario  `gorm:"foreignKey:ScenarioID" json:"scenario"`
	Version    int       `gorm:"column:version;not null;default:1"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
}
