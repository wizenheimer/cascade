# Chaos Engineering
Cascade facilitates chaos experiments to test the resiliency and recovery capabilities of Kubernetes applications through controlled disruptions.

## Introduction
Chaos Engineering involves inducing controlled failures in a system to observe how it responds under stress conditions. This README outlines the configurations, scenarios, and executions available in Cascade. The project is a WIP but changes would always be backwards compatible.

## Scenario
A Scenario organizes chaos experiments with specific configurations:

- Target: Defines the resources (namespaces, pods, nodes) to target for chaos experiments.
- Runtime: Specifies runtime parameters such as interval, grace period, ratio, and execution modes.
- Interface: Determines how chaos events are handled, including synchronous/asynchronous execution and event logging options.


### Target
Configure specific targets for chaos experiments using a combination of namespaces, labels, selectors, and resource types.

| Field             | Description                                              | Example Value                |
|-------------------|----------------------------------------------------------|------------------------------|
| `Namespaces`      | A selector for namespaces to restrict chaos experiments. | `labels.Selector` instance   |
| `IncludedPodNames`| Names of specific pods to include in chaos experiments.  | `"app-1-pod, app-2-pod"`     |
| `IncludedNodeNames`| Names of specific nodes to include in pod chaos experiments. | `"node-1, node-2"`       |
| `ExcludedPodNames`| Names of specific pods to exclude from chaos experiments.  | `"app-3-pod"`                |
| `Healthcheck`     | Endpoint for health checks during chaos experiments.     | `"/healthcheck"`             |

### Runtime
Specifies runtime parameters for chaos scenarios.

| Field      | Description                                         | Example Value        |
|------------|-----------------------------------------------------|----------------------|
| `Interval` | Interval between killing pods                       | `5`    |
| `Grace`    | Grace time after which pods are terminated          | `30`                 |
| `Ratio`    | Ratio of pods to kill                               | `0.2`                |
| `Mode`     | Pod termination strategy (`ExecutionMode` enum)     | `ExecutionMode.Delete` |
| `Order`    | Pod ordering strategy (`OrderingStrategy` enum)     | `OrderingStrategy.Random` |

### Interface
Determines how chaos events are handled, including synchronous/asynchronous execution and event logging options.

| Field         | Description                                              | Example Value |
|---------------|----------------------------------------------------------|---------------|
| `Synchronous` | Determines if chaos experiments run synchronously or asynchronously. | `true`        |
| `Watch`       | Enables logging of chaos events as server-side events.   | `true`        |

## Session
A Session represents an active chaos experiment within a Scenario.

## Tests
Various tests can be conducted using Cascade to validate system resilience and recovery capabilities:

Pod Termination: Simulates failure scenarios by terminating pods (e.g., Delete Pods, Evict Pods).

- Significance: Tests application resilience to sudden pod failures and Kubernetes' ability to recover.

Node Disruption*:

- Node Drain: Evicts all pods from a node to test workload migration and cluster stability.
- Node Restart: Restarts a node to simulate node failures and observe application behavior.
- Node Taint: Applies taints to nodes to simulate degraded node conditions.

Container Failure*:

- Docker Service Kill: Simulates failure scenarios by killing Docker services or individual containers.

Cluster-wide Disruptions*:

- ETCD Disruption: Disrupts the ETCD cluster to test Kubernetes control plane resilience.
- API Server Breakdown: Induces failures in the Kubernetes API server to validate control plane recovery.
- Network Partition: Simulates network failures within the cluster to observe network resilience and service discovery.

*WIP
