package kubernetes

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func Scenario(ctx context.Context, logger *zap.Logger, requestID string) {
	i := 0

	for {
		select {
		case <-ctx.Done():
			logger.Info("Processing stopped due to context cancellation")
			return
		default:
			logger.Info(fmt.Sprintf("%s Processing item %d", requestID, i))
			logger.Warn(fmt.Sprintf("%s This is a warning for item %d", requestID, i))
			time.Sleep(time.Second) // Simulate some work
			i++
		}
	}
}

func Session(ctx context.Context, logger *zap.Logger, requestID string) {
	i := 0

	for {
		select {
		case <-ctx.Done():
			logger.Info("Processing stopped due to context cancellation")
			return
		default:
			logger.Info(fmt.Sprintf("%s Processing item %d", requestID, i))
			logger.Warn(fmt.Sprintf("%s This is a warning for item %d", requestID, i))
			time.Sleep(time.Second) // Simulate some work
			i++
		}
	}
}
