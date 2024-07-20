package k8x

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

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

// Calculates the Pod Deletion Cost
// Reference: https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/#pod-deletion-cost
func getPodDeletionCost(pod v1.Pod) int32 {
	costString, present := pod.ObjectMeta.Annotations["controller.kubernetes.io/pod-deletion-cost"]
	if !present {
		return 0
	}
	// per k8s doc: invalid values should be rejected by the API server
	cost, _ := strconv.ParseInt(costString, 10, 32)
	return int32(cost)
}

// ========================
//
//	Pod Sorting Strategy
//
// ========================

// Reorder Pods based on the ordering strategy
func reorderPod(pods []v1.Pod, strategy OrderingStrategy) {
	switch strategy {
	case Random:
		randomOrdering(pods)
	case Default:
		defaultOrdering(pods)
	case Cost:
		costBasedOrdering(pods)
	case Youngest:
		youngestFirstOrdering(pods)
	case Oldest:
		oldestFirstOrdering(pods)
	default:
		randomOrdering(pods)
	}
}

// Doesn't order the pod
func defaultOrdering([]v1.Pod) {}

// Randomizes the ordering
func randomOrdering(pods []v1.Pod) {
	rand.Shuffle(len(pods), func(i, j int) { pods[i], pods[j] = pods[j], pods[i] })
}

// Older pods rank higher
func oldestFirstOrdering(pods []v1.Pod) {
	sort.Slice(pods, func(i, j int) bool {
		if pods[i].Status.StartTime == nil {
			return false
		}
		if pods[j].Status.StartTime == nil {
			return true
		}
		return pods[i].Status.StartTime.Unix() < pods[j].Status.StartTime.Unix()
	})
}

// Younger pods rank higher
func youngestFirstOrdering(pods []v1.Pod) {
	sort.Slice(pods, func(i, j int) bool {
		if pods[i].Status.StartTime == nil {
			return false
		}
		if pods[j].Status.StartTime == nil {
			return true
		}
		return pods[j].Status.StartTime.Unix() < pods[i].Status.StartTime.Unix()
	})
}

// Pods with lower deletion cost are ranked higher
func costBasedOrdering(pods []v1.Pod) {
	sort.Slice(pods, func(i, j int) bool {
		return getPodDeletionCost(pods[i]) < getPodDeletionCost(pods[j])
	})
}
