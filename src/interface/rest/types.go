package rest

import (
	"net/http"

	"github.com/wizenheimer/cascade/service/database"
	"go.uber.org/zap"
)

type APIServer struct {
	server *http.Server
	Logger *zap.Logger
	DB     *database.DatabaseClient
}
