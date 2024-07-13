package processor

import (
	"context"

	"go.uber.org/zap"
)

// Kubernetes Callback
type kubernetesCallback func(context.Context, *zap.Logger, string)
