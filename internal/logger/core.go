package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Instantiate a logger
func CreateLogger(logChan chan<- LogEntry) *zap.Logger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}

	// Attach a hook that parses the logs into server side events
	hook := func(entry zapcore.Entry) error {
		logEntry := LogEntry{
			Timestamp: entry.Time.UnixNano() / 1e6,
			Level:     entry.Level.String(),
			Message:   entry.Message,
		}
		select {
		case logChan <- logEntry:
			// Relay the logs back to the client
		default:
			// Channel full, log dropped
		}
		return nil
	}

	logger, err := config.Build(zap.Hooks(hook))
	if err != nil {
		panic(err.Error())
	}

	return logger
}
