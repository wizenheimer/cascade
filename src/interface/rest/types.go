package rest

import (
	"net/http"

	"github.com/wizenheimer/cascade/service/database"
	"go.uber.org/zap"
)

type APIServer struct {
	server *http.Server
	Logger *zap.Logger
	DB     database.DatabaseClient
}

// Config represents the configuration for chaos engineering scenarios
type Config struct {
	Scenario Scenario `yaml:"scenario"`
	Target   Target   `yaml:"target"`
	Runtime  Runtime  `yaml:"runtime"`
}

// Scenario represents the chaos engineering scenario
type Scenario struct {
	ID          string `yaml:"id"`
	Description string `yaml:"description"`
}

// Target represents the resources to target for chaos engineering scenarios
type Target struct {
	Namespaces        string `yaml:"namespaces"`
	IncludedPodNames  string `yaml:"includedPodNames"`
	IncludedNodeNames string `yaml:"includedNodeNames"`
	ExcludedPodNames  string `yaml:"excludedPodNames"`
}

// Runtime represents the runtime arguments for executing the scenario
type Runtime struct {
	Interval string `yaml:"interval"`
	Grace    string `yaml:"grace"`
	Mode     string `yaml:"mode"`
	Ordering string `yaml:"ordering"`
	Ratio    string `yaml:"ratio"`
}
