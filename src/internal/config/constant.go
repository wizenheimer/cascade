package config

// API Server Defaults
const (
	SERVER_PORT = ":8080"
)

// Executor Defaults
const (
	HEALTH_CHECK_PORT = ":8080"

	RUNTIME_INTERVAL = "10m"

	RATIO = "0.5"

	MODE = "dry-run"

	GRACE = "1m"

	ORDERING = "oldest"

	ORIGIN = "host"
)

// CLI Defaults
const (
	FILE_NAME = "scenario.yaml" // Default file name

	OUTPUT_FOLDER = "." // Current directory
)
