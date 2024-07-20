package k8x

import (
	"errors"
	"time"

	"k8s.io/apimachinery/pkg/labels"
)

// Determines which resources to target for chaos engineering scenarios
type TargetConfig struct {
	// A namespace or a set of namespaces to restrict thanoskube
	Namespaces labels.Selector `json:"namespaces" yaml:"namespaces"`
	// A string to select which pods to kill
	IncludedPodNames string `json:"includedPodNames" yaml:"includedPodNames"`
	// A string to select nodes, pods within the selected nodes will be killed
	IncludedNodeNames string `json:"includedNodeNames" yaml:"includedNodeNames"`
	// A string to exclude pods to kill
	ExcludedPodNames string `json:"excludedPodNames" yaml:"excludedPodNames"`
}

// Determines which clusters to target for chaos engineering scenarios
type ClusterConfig struct {
	// Path to a kubeconfig file
	Kubeconfig string `json:"kubeconfig"`
	// The address of the Kubernetes cluster to target, if none looks under $HOME/.kube L
	Master string `json:"master"`
	// Listens this endpoint for healtcheck
	Healthcheck string `json:"healthcheck"`
}

// Determines the Pod Termination Strategy
type ExecutionMode int

const (
	Delete ExecutionMode = iota // Triggers Deletion
	DryRun                      // Executes as a Dry Run
	Evict                       // Eviction API instead of Deletion
)

// ParseExecutionMode converts a string representation of ExecutionMode to its enum value.
func ParseExecutionMode(modeStr string) ExecutionMode {
	switch modeStr {
	case "delete":
		return Delete
	case "dry-run":
		return DryRun
	case "evict":
		return Evict
	default:
		// Default to Delete
		return Delete
	}
}

// Determines the Pod Ordering Strategy
type OrderingStrategy int

const (
	Random   OrderingStrategy = iota // Order the pods randomly
	Default                          // Avoid reordering the pods
	Cost                             // Lower cost pods are rank higher
	Youngest                         // Younger pods are ranked higher
	Oldest                           // Older pods are ranked higher
)

// ParseOrderingStrategy converts a string representation of OrderingStrategy to its enum value.
func ParseOrderingStrategy(orderingStr string) OrderingStrategy {
	switch orderingStr {
	case "random":
		return Random
	case "default":
		return Default
	case "cost":
		return Cost
	case "youngest":
		return Youngest
	case "oldest":
		return Oldest
	default:
		// Default to Random
		return Random
	}
}

// Determine the Runtime Configurations for chaos engineering scenarios
type RuntimeConfig struct {
	// Interval between killing pods
	Interval time.Duration `json:"interval" yaml:"interval"`
	// Grace Time after which pods are terminated
	Grace int64 `json:"grace" yaml:"grace"`
	// Ratio of pods to kill
	Ratio float64 `json:"ratio" yaml:"ratio"`
	// Pod termination strategy
	Mode ExecutionMode `json:"mode" yaml:"mode"`
	// Pod Ordering strategy
	Order OrderingStrategy `json:"ordering" yaml:"ordering"`
}

var podNotFound = "pod not found"
var errPodNotFound = errors.New(podNotFound)
