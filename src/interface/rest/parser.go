package rest

import (
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wizenheimer/cascade/internal/config"
	"github.com/wizenheimer/cascade/internal/models"
	k8x "github.com/wizenheimer/cascade/service/kubernetes"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/labels"
)

// Parse DB Scenario
func ParseDBScenario(scenario models.Scenario) (*k8x.TargetConfig, *k8x.RuntimeConfig, error) {
	// Convert it into an intermediate config
	cfg := Config{
		Scenario: Scenario{
			ID:          scenario.ID,
			Description: scenario.Description,
		},
		Target: Target{
			Namespaces:        scenario.Namespaces,
			IncludedPodNames:  scenario.IncludedPodNames,
			IncludedNodeNames: scenario.IncludedNodeNames,
			ExcludedPodNames:  scenario.ExcludedPodNames,
		},
		Runtime: Runtime{
			Interval: scenario.Interval,
			Grace:    scenario.Grace,
			Mode:     scenario.Mode,
			Ordering: scenario.Ordering,
		},
	}

	tc, err := ParseTargetConfig(&cfg)
	if err != nil {
		return nil, nil, err
	}

	rc, err := ParseRuntimeConfig(&cfg)
	if err != nil {
		return nil, nil, err
	}
	rc.Ratio = scenario.Ratio

	return tc, rc, nil
}

func ParseYAMLConfigToScenario(cfg *Config) (*models.Scenario, error) {
	scenario := models.Scenario{
		Description: cfg.Scenario.Description,
		Namespaces:  cfg.Target.Namespaces,
	}

	scenario.IncludedPodNames = cfg.Target.IncludedPodNames
	scenario.IncludedNodeNames = cfg.Target.IncludedNodeNames
	scenario.ExcludedPodNames = cfg.Target.ExcludedPodNames

	scenario.Interval = cfg.Runtime.Interval
	scenario.Grace = cfg.Runtime.Grace
	scenario.Mode = cfg.Runtime.Mode
	scenario.Ordering = cfg.Runtime.Ordering

	ratioStr := cfg.Runtime.Ratio
	ratio, err := strconv.ParseFloat(ratioStr, 64)
	if err != nil {
		return nil, err
	}

	scenario.Ratio = ratio

	return &scenario, nil
}

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
