package k8x

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Selects the Pods to Kill
func (executor *Executor) SelectPodsToKill(ctx context.Context) ([]v1.Pod, error) {
	// Figure out the Candidate Pods
	pods, err := executor.SelectCandidatePods(ctx)
	if err != nil {
		return []v1.Pod{}, err
	}
	if len(pods) == 0 {
		return []v1.Pod{}, errPodNotFound
	}

	// Prepare a Random Pod Slice
	pods = RandomPodSlice(pods, executor.Runtime.Ratio)

	// TODO: Reorder the Pods based on executor.Runtime.Strategy

	return pods, nil
}

// Returns the list of pods which qualify the targeting critera.
// Excludes terminating pods from Candidate List
func (executor *Executor) SelectCandidatePods(ctx context.Context) ([]v1.Pod, error) {
	listOptions := metav1.ListOptions{LabelSelector: ""} // get all labels

	allPods, err := executor.Client.CoreV1().Pods(executor.Target.Namespaces.String()).List(ctx, listOptions)
	if err != nil {
		return nil, err
	}
	filteredPods, err := filterByNamespaces(allPods.Items, executor.Target.Namespaces)
	if err != nil {
		return nil, err
	}

	filteredPods = includePodsByNodeName(filteredPods, executor.Target.IncludedNodeNames)
	filteredPods = includePodsByPodName(filteredPods, executor.Target.IncludedPodNames)
	filteredPods = excludePodsByPodName(filteredPods, executor.Target.ExcludedPodNames)
	filteredPods = filterTerminatingPods(filteredPods)

	return filteredPods, nil
}

// Trigger Pod Deletion
func (executor *Executor) DeletePod(pod v1.Pod, ctx context.Context) error {
	if executor.Runtime.Mode == DryRun {
		return nil
	}

	// TODO: Add Eviction based APIs
	opts := metav1.DeleteOptions{}
	err := executor.Client.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, opts)
	if err != nil {
		return err
	}

	return nil
}
