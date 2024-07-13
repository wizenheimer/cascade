package logger

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Instantiate a logger
func CreateLogger() *LoggerWithChannel {
	logChan := make(chan LogEntry, 100)
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}

	hook := func(entry zapcore.Entry) error {
		logEntry := LogEntry{
			Timestamp: entry.Time.UnixNano() / 1e6,
			Level:     entry.Level.String(),
			Message:   entry.Message,
		}
		select {
		case logChan <- logEntry:
		default:
			// Channel full, log dropped
		}
		return nil
	}

	logger, _ := config.Build(zap.Hooks(hook))
	return &LoggerWithChannel{
		Logger:  logger,
		LogChan: logChan,
	}
}

// Manually parse log for streaming it back to client
func ParseLog(level, msg string) ([]byte, error) {
	logEntry := LogEntry{
		Timestamp: time.Now().Unix(),
		Level:     level,
		Message:   msg,
	}
	// Serialize logEntry to JSON
	data, err := json.Marshal(logEntry)
	if err != nil {
		return nil, err
	}

	data = []byte(fmt.Sprintf("data: %s\n\n", data))

	return data, nil
}
