package rest

import (
	"net/http"

	"go.uber.org/zap"
)

type APIServer struct {
	server *http.Server
	logger *zap.Logger
}
