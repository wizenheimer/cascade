package k8x

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/reference"
)

// Selects the Pods to Kill
func (executor *Executor) SelectPodsToKill(ctx context.Context) ([]v1.Pod, error) {
	// Figure out the Candidate Pods
	executor.Logger.Info("Preparing Candidate Pods")
	pods, err := executor.SelectCandidatePods(ctx)
	if err != nil {
		executor.Logger.Error("Error occured while preparing Candidate Pods")
		return []v1.Pod{}, err
	}
	if len(pods) == 0 {
		executor.Logger.Warn(errPodNotFound.Error())
		return []v1.Pod{}, errPodNotFound
	}

	// Prepare a Random Pod Slice
	executor.Logger.Info("Sampling from a list of candidate pods")
	pods = RandomPodSlice(pods, executor.Runtime.Ratio)

	// Reorder the Pods
	executor.Logger.Info("Reordering the Pods")
	reorderPod(pods, executor.Runtime.Order)

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

	executor.Logger.Info(fmt.Sprintf("Filtering down to %d Candidates", len(filteredPods)))
	return filteredPods, nil
}

// Trigger Pod Deletion based on Termination Strategy
func (executor *Executor) DeletePod(pod v1.Pod, ctx context.Context) error {

	opts := metav1.DeleteOptions{
		GracePeriodSeconds: &executor.Runtime.Grace,
	}
	var err error

	switch executor.Runtime.Mode {
	case DryRun:
		executor.Logger.Info("Terminating as per Dry Run Strategy")
		time.Sleep(time.Duration(executor.Runtime.Grace)) // To mock grace period
		err = nil
	case Evict:
		executor.Logger.Info("Terminating as per Eviction Strategy")
		err = executor.Client.CoreV1().Pods(pod.Namespace).Evict(ctx, &policyv1.Eviction{
			ObjectMeta:    metav1.ObjectMeta{Namespace: pod.Namespace, Name: pod.Name},
			DeleteOptions: &opts,
		})
	default:
		executor.Logger.Info("Terminating as per Deletion Strategy")
		err = executor.Client.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, opts)
	}

	if err != nil {
		return err
	}

	ref, err := reference.GetReference(scheme.Scheme, &pod)
	if err != nil {
		return err
	}

	executor.EventRecorder.Event(ref, v1.EventTypeNormal, "killing", "pod was killed by cascade.")

	return nil
}
