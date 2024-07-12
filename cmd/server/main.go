package main

import (
	"github.com/wizenheimer/cascade/interface/rest"
	"go.uber.org/zap"
)

func main() {
	// Create a logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// Create a New RESTful API Service
	server := rest.NewAPIServer(logger)
	server.Serve()
}
