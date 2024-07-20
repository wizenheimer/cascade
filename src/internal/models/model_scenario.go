package models

import "time"

// Scenario represents a chaos engineering experiment
type Scenario struct {
	ID                string    `gorm:"primaryKey;column:scenario_id" json:"id"`
	Version           int       `gorm:"column:version;not null;default:1" json:"version"`
	Description       string    `gorm:"column:description;type:text" json:"description"`
	Namespaces        string    `gorm:"column:namespaces;type:text" json:"namespaces"`
	IncludedPodNames  string    `gorm:"column:includedPodNames;type:text" json:"includedPodNames"`
	IncludedNodeNames string    `gorm:"column:includedNodeNames;type:text" json:"includedNodeNames"`
	ExcludedPodNames  string    `gorm:"column:excludedPodNames;type:text" json:"excludedPodNames"`
	Interval          string    `gorm:"column:interval;type:text" json:"interval"`
	Grace             string    `gorm:"column:grace;type:text" json:"grace"`
	Mode              string    `gorm:"column:mode;type:text" json:"mode"`
	Ordering          string    `gorm:"column:ordering;type:text" json:"ordering"`
	Ratio             float64   `gorm:"column:ratio;type:text" json:"ratio"`
	TeamID            string    `gorm:"column:team_id;not null" json:"team_id"`
	CreatedAt         time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP()" json:"created_at"`
	Team              Team      `gorm:"foreignKey:TeamID" json:"team"`
}
