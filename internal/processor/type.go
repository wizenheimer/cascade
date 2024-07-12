package processor

import (
	"context"

	"go.uber.org/zap"
)

type kubernetesCallback func(context.Context, *zap.Logger, string)
