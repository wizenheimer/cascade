package k8x

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wizenheimer/cascade/internal/config"
)

// Fetch Target specific Configurations from Echo's Context
func ParseTargetConfig(c echo.Context) (*TargetConfig, error) {
	namespaces := c.FormValue("namespaces")
	includedPodNames := c.FormValue("includedPodNames")
	includedNodeNames := c.FormValue("includedNodeNames")
	excludedPodNames := c.FormValue("excludedPodNames")

	healthcheck := c.FormValue("healthcheck")
	// Incase it's empty default it
	if healthcheck == "" {
		healthcheck = config.HEALTH_CHECK_PORT
	}

	return &TargetConfig{
		Namespaces:        namespaces,
		IncludedPodNames:  includedPodNames,
		IncludedNodeNames: includedNodeNames,
		ExcludedPodNames:  excludedPodNames,
		Healthcheck:       healthcheck,
	}, nil
}

// Fetch Cluster specific Configurations from Echo's Context
func ParseClusterConfig(c echo.Context) (*ClusterConfig, error) {
	kubeconfig := c.FormValue("kubeconfig")
	master := c.FormValue("master")

	return &ClusterConfig{
		Kubeconfig: kubeconfig,
		Master:     master,
	}, nil
}

// Fetch Runtime Configurations from Echo's Context
func ParseRuntimeConfig(c echo.Context) (*RuntimeConfig, error) {
	intervalStr := c.FormValue("interval")
	if intervalStr == "" {
		intervalStr = config.RUNTIME_INTERVAL
	}

	ratioStr := c.FormValue("ratio")
	if ratioStr == "" {
		ratioStr = config.RATIO
	}

	modeStr := c.FormValue("mode")
	if modeStr == "" {
		modeStr = config.MODE
	}

	// Parse interval
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, err
	}

	// Parse ratio
	ratio, err := strconv.ParseFloat(ratioStr, 64)
	if err != nil {
		return nil, err
	}

	// Convert modeStr to ExecutionMode enum
	mode := ParseExecutionMode(modeStr)

	return &RuntimeConfig{
		Interval: interval,
		Ratio:    ratio,
		Mode:     mode,
	}, nil
}

// Fetch Interface related Configurations from Echo's Context
func ParseInterfaceConfig(c echo.Context) (*InterfaceConfig, error) {
	synchronousStr := c.FormValue("synchronous")
	if synchronousStr == "" {
		synchronousStr = config.SYNC
	}

	watchStr := c.FormValue("watch")
	if watchStr == "" {
		watchStr = config.WATCH
	}

	// Parse synchronicity
	synchronous, err := strconv.ParseBool(synchronousStr)
	if err != nil {
		return nil, err
	}

	// Parse watch mode
	watch, err := strconv.ParseBool(watchStr)
	if err != nil {
		return nil, err
	}

	return &InterfaceConfig{
		Synchronous: synchronous,
		Watch:       watch,
	}, nil
}

// Parse Chaos Engineering Configs from Echo's Context
func ParseConfigs(c echo.Context) (*ClusterConfig, *TargetConfig, *RuntimeConfig, *InterfaceConfig, error) {
	targetConfig, err := ParseTargetConfig(c)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	clusterConfig, err := ParseClusterConfig(c)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	runtimeConfig, err := ParseRuntimeConfig(c)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	interfaceConfig, err := ParseInterfaceConfig(c)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return clusterConfig, targetConfig, runtimeConfig, interfaceConfig, nil
}
