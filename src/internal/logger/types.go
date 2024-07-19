package logger

import (
	"sync"

	"go.uber.org/zap"
)

type LogEntry struct {
	Timestamp int64  `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type LoggerWithChannel struct {
	Logger  *zap.Logger
	LogChan chan LogEntry
}

// Pool for verbose loggers
var LoggerPool = sync.Pool{
	New: func() interface{} {
		return CreateLogger()
	},
}
