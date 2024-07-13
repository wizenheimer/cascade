package k8x

import (
	"context"
	"os"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
)

// Holds the executor for running Chaos Experiments
type Executor struct {
	// Good'ol Kubernetes Clientset
	Client kubernetes.Interface
	// Publishes events using EventRecorder Controller
	EventRecorder record.EventRecorder
	// Determines the Target for Chaos Scenarios
	Target *TargetConfig
	// Determine the Runtime for Chaos Scenarios
	Runtime *RuntimeConfig
}

// Initializes an executor instance
func CreateExecutor(cc *ClusterConfig, tc *TargetConfig, rc *RuntimeConfig) (*Executor, error) {
	client, err := getK8Client(cc)
	if err != nil {
		return nil, err
	}

	recorder := getEventRecorder(client)

	return &Executor{
		Client:        client,   // Kubernetes Client Instance
		EventRecorder: recorder, // Event Recorder Instance
		Target:        tc,
		Runtime:       rc,
	}, nil
}

// Returns Kubernetes Client
func getK8Client(cc *ClusterConfig) (*kubernetes.Clientset, error) {
	// look for kubeconfig in home if not set
	if cc.Kubeconfig == "" {
		if _, err := os.Stat(clientcmd.RecommendedHomeFile); err == nil {
			cc.Kubeconfig = clientcmd.RecommendedHomeFile
		}
	}

	config, err := clientcmd.BuildConfigFromFlags(cc.Master, cc.Kubeconfig)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Returns an EventRecorder that can be used to log events via EventRecorder Controller
func getEventRecorder(client *kubernetes.Clientset) record.EventRecorderLogger {
	broadcaster := record.NewBroadcaster()
	broadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: client.CoreV1().Events(v1.NamespaceAll)})
	recorder := broadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "thanos"})
	return recorder
}

// Execute the chaos engineering scenario
// Return an error incase, pods deletion got interupped
func (executor *Executor) Execute(ctx context.Context, logger *zap.Logger) error {
	// Identify the pods to kill
	podsToKill, err := executor.SelectPodsToKill(ctx)
	if err != nil {
		return err
	}
	if err == errPodNotFound {
		logger.Debug(podNotFound)
		return nil
	}

	// Trigger deletion
	var result *multierror.Error
	for _, victim := range podsToKill {
		err = executor.DeletePod(victim, ctx)
		if err != nil {
			logger.Error("failed to delete pod", zap.Any("pod", victim.Name))
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}
