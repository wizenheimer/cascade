package models

// SessionMetrics represents aggregated metrics for a given scenario
type SessionMetrics struct {
	ScenarioID      string `json:"scenario_id"`
	Version         int    `json:"version"`
	CompletedCount  int    `json:"completed_count"`
	TerminatedCount int    `json:"terminated_count"`
	InProgressCount int    `json:"in_progress_count"`
	QueuedCount     int    `json:"queued_count"`
}
