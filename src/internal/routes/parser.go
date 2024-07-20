package routes

import (
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wizenheimer/cascade/internal/config"
	k8x "github.com/wizenheimer/cascade/service/kubernetes"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/labels"
)

// Parse Target Configurations
func ParseTargetConfig(cfg *Config) (*k8x.TargetConfig, error) {
	namespaces, err := parseNamespaces(cfg.Target.Namespaces)
	if err != nil {
		return nil, err
	}

	return &k8x.TargetConfig{
		Namespaces:        namespaces,
		IncludedPodNames:  cfg.Target.IncludedNodeNames,
		IncludedNodeNames: cfg.Target.IncludedPodNames,
		ExcludedPodNames:  cfg.Target.ExcludedPodNames,
	}, nil
}

// Fetch Cluster specific Configurations from Echo's Context
func ParseClusterConfigFromContext(c echo.Context) (*k8x.ClusterConfig, error) {
	kubeconfig := c.FormValue("kubeconfig")
	master := c.FormValue("master")
	healthcheck := c.FormValue("healthcheck")
	if healthcheck == "" {
		healthcheck = config.GetEnv("HEALTH_CHECK_PORT", config.HEALTH_CHECK_PORT)
	}

	return &k8x.ClusterConfig{
		Kubeconfig:  kubeconfig,
		Master:      master,
		Healthcheck: healthcheck,
	}, nil
}

// Fetch Runtime Configurations from Echo's Context
func ParseRuntimeConfig(cfg *Config) (*k8x.RuntimeConfig, error) {
	intervalStr := cfg.Runtime.Interval
	if intervalStr == "" {
		intervalStr = config.GetEnv("RUNTIME_INTERVAL", config.RUNTIME_INTERVAL)
	}

	ratioStr := cfg.Runtime.Ratio
	if ratioStr == "" {
		ratioStr = config.GetEnv("RATIO", config.RATIO)
	}

	modeStr := cfg.Runtime.Mode
	if modeStr == "" {
		modeStr = config.GetEnv("MODE", config.MODE)
	}

	graceStr := cfg.Runtime.Grace
	if graceStr == "" {
		graceStr = config.GetEnv("GRACE", config.GRACE)
	}

	orderStr := cfg.Runtime.Ordering
	if orderStr == "" {
		orderStr = config.GetEnv("ORDERING", config.ORDERING)
	}

	// Parse Ordering
	ordering := k8x.ParseOrderingStrategy(orderStr)

	// Parse interval
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, err
	}

	// Parse grace
	grace, err := strconv.ParseInt(graceStr, 10, 64)
	if err != nil {
		return nil, err
	}

	// Parse ratio
	ratio, err := strconv.ParseFloat(ratioStr, 64)
	if err != nil {
		return nil, err
	}

	// Convert modeStr to ExecutionMode enum
	mode := k8x.ParseExecutionMode(modeStr)

	return &k8x.RuntimeConfig{
		Interval: interval,
		Ratio:    ratio,
		Mode:     mode,
		Grace:    grace,
		Order:    ordering,
	}, nil
}

// Parse Chaos Engineering Configs from Echo's Context
func ParseConfig(c echo.Context) (*k8x.TargetConfig, *k8x.RuntimeConfig, error) {
	var config Config

	// Handle YAML based parameterization
	file, err := c.FormFile("config")
	if err != nil {
		return nil, nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, nil, err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, nil, err
	}

	targetConfig, err := ParseTargetConfig(&config)
	if err != nil {
		return nil, nil, err
	}

	runtimeConfig, err := ParseRuntimeConfig(&config)
	if err != nil {
		return nil, nil, err
	}

	return targetConfig, runtimeConfig, nil
}

// Parse strings into selectors
func parseNamespaces(str string) (labels.Selector, error) {
	selector, err := labels.Parse(str)
	if err != nil {
		return nil, err
	}
	return selector, nil
}