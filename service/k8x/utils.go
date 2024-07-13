package k8x

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wizenheimer/cascade/internal/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

// Fetch Target specific Configurations from Echo's Context
func ParseTargetConfig(c echo.Context) (*TargetConfig, error) {
	namespacesStr := c.FormValue("namespaces")
	namespaces, err := parseNamespaces(namespacesStr)
	if err != nil {
		return nil, err
	}

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

// Parse strings into selectors
func parseNamespaces(str string) (labels.Selector, error) {
	selector, err := labels.Parse(str)
	if err != nil {
		return nil, err
	}
	return selector, nil
}

// =====================
// Filtering Functions
// =====================

func includePodsByNodeName(pods []v1.Pod, includedNodeNames string) (filteredPods []v1.Pod) {
	includedPodNamesList := strings.Split(includedNodeNames, ",")

	var resultingPods []v1.Pod
	for _, pod := range pods {
		for _, podNameToInclude := range includedPodNamesList {
			if strings.Contains(pod.Spec.NodeName, podNameToInclude) {
				resultingPods = append(resultingPods, pod)
			}
		}
	}

	return resultingPods
}

func includePodsByPodName(pods []v1.Pod, includedPodNames string) (filteredPods []v1.Pod) {
	includedPodNamesList := strings.Split(includedPodNames, ",")

	var resultingPods []v1.Pod
	for _, pod := range pods {
		for _, podNameToInclude := range includedPodNamesList {
			if strings.Contains(pod.ObjectMeta.Name, podNameToInclude) {
				resultingPods = append(resultingPods, pod)
			}
		}
	}

	return resultingPods
}

func excludePodsByPodName(pods []v1.Pod, excludedPodNames string) (filteredPods []v1.Pod) {
	excludedPodNamesList := strings.Split(excludedPodNames, ",")

	if len(excludedPodNamesList) == 1 && excludedPodNamesList[0] == "" {
		return pods
	}

	var resultingPods []v1.Pod
	for _, pod := range pods {
		for _, podNameToExclude := range excludedPodNamesList {
			if !strings.Contains(pod.ObjectMeta.Name, podNameToExclude) {
				resultingPods = append(resultingPods, pod)
			}
		}
	}

	return resultingPods
}

func filterByNamespaces(pods []v1.Pod, namespaces labels.Selector) ([]v1.Pod, error) {
	if namespaces.Empty() {
		return pods, nil
	}

	requirements, _ := namespaces.Requirements()
	var includeRequirements []labels.Requirement
	var excludeRequirements []labels.Requirement

	for _, req := range requirements {
		switch req.Operator() {
		case selection.Exists:
			includeRequirements = append(includeRequirements, req)
		case selection.DoesNotExist:
			excludeRequirements = append(excludeRequirements, req)
		default:
			return nil, fmt.Errorf("unsupported operator: %s", req.Operator())
		}
	}

	var filteredPods []v1.Pod

	for _, pod := range pods {
		included := len(includeRequirements) == 0

		selector := labels.Set{pod.Namespace: ""}

		for _, req := range includeRequirements {
			if req.Matches(selector) {
				included = true
				break
			}
		}

		for _, req := range excludeRequirements {
			if !req.Matches(selector) {
				included = false
				break
			}
		}

		if included {
			filteredPods = append(filteredPods, pod)
		}
	}

	return filteredPods, nil
}

func filterTerminatingPods(pods []v1.Pod) []v1.Pod {
	var filteredList []v1.Pod
	for _, pod := range pods {
		if pod.DeletionTimestamp != nil {
			continue
		}
		filteredList = append(filteredList, pod)
	}
	return filteredList
}

func RandomPodSlice(pods []v1.Pod, percentageToKill float64) []v1.Pod {
	count := int(float64(len(pods)) * percentageToKill)

	rand.Shuffle(len(pods), func(i, j int) { pods[i], pods[j] = pods[j], pods[i] })
	res := pods[0:count]
	return res
}
