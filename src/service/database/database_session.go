package database

import (
	"context"
	"github.com/wizenheimer/cascade/internal/models"
	"time"
)

// CreateSession creates a new session for a given scenario
func (c Client) CreateSession(ctx context.Context, scenarioID string, version int) (*models.Session, error) {
	var scenario models.Scenario
	if err := c.DB.Where("scenario_id = ?", scenarioID).Order("version DESC").First(&scenario).Error; err != nil {
		return nil, err
	}

	if version <= 0 || version > scenario.Version {
		version = scenario.Version
	}

	session := &models.Session{
		ScenarioID: scenarioID,
		Version:    version,
		StartTime:  time.Now(),
		Status:     "queued",
	}

	result := c.DB.Create(session)
	if result.Error != nil {
		return nil, result.Error
	}

	return session, nil
}

// StartSession marks session as in-progress
func (c Client) StartSession(ctx context.Context, sessionID string) (*models.Session, error) {
	var session models.Session
	if err := c.DB.First(&session, sessionID).Error; err != nil {
		return nil, err
	}

	session.StartTime = time.Now()
	session.Status = "running"

	result := c.DB.Save(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	return &session, nil
}

// GracefullyEndSession marks session as a graceful exit
func (c Client) GracefullyEndSession(ctx context.Context, sessionID string) (*models.Session, error) {
	var session models.Session
	if err := c.DB.First(&session, sessionID).Error; err != nil {
		return nil, err
	}

	session.EndTime = time.Now()
	session.Status = "completed"

	result := c.DB.Save(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	return &session, nil
}

// TerminateSession marks session as terminated exit
func (c Client) TerminateSession(ctx context.Context, sessionID string) (*models.Session, error) {
	var session models.Session
	if err := c.DB.First(&session, sessionID).Error; err != nil {
		return nil, err
	}

	session.EndTime = time.Now()
	session.Status = "failed"

	result := c.DB.Save(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	return &session, nil
}

func (c Client) ListSessionByScenarioID(ctx context.Context, scenarioID string, version int) ([]models.Session, error) {
	var sessions []models.Session
	query := c.DB.Where("scenario_id = ?", scenarioID)

	if version > 0 {
		query = query.Where("version = ?", version)
	}

	result := query.Find(&sessions)
	return sessions, result.Error
}

func (c Client) GetSessionMetrics(ctx context.Context, scenarioID string) ([]models.SessionMetrics, error) {
	var metrics []models.SessionMetrics

	result := c.DB.Model(&models.Session{}).
		Select("scenario_id, version, "+
			"SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_count, "+
			"SUM(CASE WHEN status = 'terminated' THEN 1 ELSE 0 END) as terminated_count, "+
			"SUM(CASE WHEN status = 'in-progress' THEN 1 ELSE 0 END) as in_progress_count, "+
			"SUM(CASE WHEN status = 'queued' THEN 1 ELSE 0 END) as queued_count").
		Where("scenario_id = ?", scenarioID).
		Group("scenario_id, version").
		Scan(&metrics)

	if result.Error != nil {
		return nil, result.Error
	}

	return metrics, nil
}
