![Cascade](media/banner.png)

# Cascade

Ever had that sinking feeling of your system going down at 3 AM? Yeah, me too. It sucks. With distributed systems, failures aren't an exception—they're the norm.

But what if failures could be a feature, not a bug? Enter Cascade. We break things on purpose so they don't break by accident. Simple as that.

Cascade pushes your deployments to the brink and helps you patch them up—all before your users even notice.

## Why Chaos Engineering?

Cascade doesn't just break things and leave you to pick up the pieces. We let you simulate failure scenarios in a controlled environment. Crash and burn all you want—no real users harmed in the process.

**Who uses chaos engineering?**
- **Netflix** ensures your binge-watching is never interrupted
- **Zomato** keeps your impulsive orders flowing 24/7
- **NASA** uses it (yes, actual rocket scientists)

If it's good enough for rocket scientists, it's probably good enough for your app too.

## Overview

Chaos Engineering involves inducing controlled failures in a system to observe how it responds under stress conditions. Cascade facilitates chaos experiments to test the resiliency and recovery capabilities of Kubernetes applications through controlled disruptions.

This project is a work in progress, but changes will always be backwards compatible.

## Core Concepts

### Scenario

A Scenario organizes chaos experiments with specific configurations:

- **Target**: Defines the resources (namespaces, pods, nodes) to target for chaos experiments
- **Runtime**: Specifies runtime parameters such as interval, grace period, ratio, and execution modes
- **Interface**: Determines how chaos events are handled, including synchronous/asynchronous execution and event logging options

#### Target Configuration

Configure specific targets for chaos experiments using namespaces, labels, selectors, and resource types.

| Field | Description | Example Value |
|-------|-------------|---------------|
| `Namespaces` | A selector for namespaces to restrict chaos experiments | `labels.Selector` instance |
| `IncludedPodNames` | Names of specific pods to include in chaos experiments | `"app-1-pod, app-2-pod"` |
| `IncludedNodeNames` | Names of specific nodes to include in pod chaos experiments | `"node-1, node-2"` |
| `ExcludedPodNames` | Names of specific pods to exclude from chaos experiments | `"app-3-pod"` |
| `Healthcheck` | Endpoint for health checks during chaos experiments | `"/healthcheck"` |

#### Runtime Parameters

Specifies runtime parameters for chaos scenarios.

| Field | Description | Example Value |
|-------|-------------|---------------|
| `Interval` | Interval between killing pods | `5` |
| `Grace` | Grace time after which pods are terminated | `30` |
| `Ratio` | Ratio of pods to kill | `0.2` |
| `Mode` | Pod termination strategy (`ExecutionMode` enum) | `ExecutionMode.Delete` |
| `Order` | Pod ordering strategy (`OrderingStrategy` enum) | `OrderingStrategy.Random` |

#### Interface Options

Determines how chaos events are handled.

| Field | Description | Example Value |
|-------|-------------|---------------|
| `Synchronous` | Determines if chaos experiments run synchronously or asynchronously | `true` |
| `Watch` | Enables logging of chaos events as server-side events | `true` |

### Session

A Session represents an active chaos experiment within a Scenario.

## Test Types

### Pod Termination

Pod termination tests are the foundation of chaos engineering in Kubernetes. They simulate real-world failure scenarios where pods unexpectedly crash, get evicted, or are forcibly terminated.

#### Termination Strategies

**Delete Mode**

Immediately removes pods from the cluster, simulating hard crashes or node failures. This is the most aggressive termination strategy and tests your application's ability to handle sudden, ungraceful shutdowns.

- Use case: Testing resilience to node failures, out-of-memory kills, or infrastructure issues
- What it validates: Pod restart policies, replica set recovery, service mesh behavior during abrupt failures

**Evict Mode**

Uses Kubernetes' eviction API to gracefully remove pods, respecting PodDisruptionBudgets (PDBs). This simulates planned maintenance scenarios or resource pressure situations.

- Use case: Testing rolling updates, node drains, and cluster autoscaling scenarios
- What it validates: PDB configurations, graceful shutdown handlers, connection draining

**Dry-Run Mode**

Identifies target pods without actually terminating them. Perfect for validating your chaos experiment configuration before executing it for real.

- Use case: Testing chaos scenario configurations, verifying target selection logic
- What it validates: Your understanding of which pods will be affected

#### What Pod Termination Tests Reveal

Pod termination experiments help you answer critical questions about your system:

- **Recovery Time**: How quickly does Kubernetes reschedule terminated pods?
- **Service Availability**: Do your services remain available during pod failures?
- **Data Integrity**: Are in-flight requests handled gracefully?
- **Dependencies**: How do downstream services react when upstream pods fail?
- **Monitoring & Alerting**: Do your observability tools correctly detect and alert on failures?

#### Ordering Strategies

Control which pods are terminated first based on different criteria:

- **Random**: Unpredictable termination pattern, simulating chaotic real-world failures
- **Oldest**: Targets the longest-running pods first, useful for testing rolling update scenarios
- **Youngest**: Targets recently started pods, testing startup and initialization resilience
- **Cost**: Terminates pods based on resource consumption metrics

### Additional Test Types *(Work in Progress)*

Cascade is expanding to include node disruptions (drain, restart, taint), container-level failures (Docker service kills), and cluster-wide disruptions (ETCD, API server, network partitions).

## CLI Usage

![Cascade](media/cli.gif)

Cascade CLI is a powerful tool for orchestrating chaos engineering experiments in Kubernetes clusters. It provides an intuitive interface for creating, managing, and executing chaos scenarios.

### Visualizing Chaos with Kite
<img width="3160" height="2334" alt="image" src="https://github.com/user-attachments/assets/14d95591-fc37-4673-b3f1-2953fdb34d4c" />

For a better chaos engineering experience, we recommend using [Kite](https://github.com/zxh326/kite) - a modern, lightweight Kubernetes dashboard that provides real-time visualization of your chaos experiments.

**Why use Kite with Cascade?**

- **Real-time Monitoring**: Watch pod terminations and recovery in real-time with live metrics
- **Resource Relationships**: Visualize how chaos experiments affect related resources (Deployments → ReplicaSets → Pods)
- **Live Logs**: Stream logs from affected pods during chaos experiments
- **Event Tracking**: See Kubernetes events triggered by chaos experiments as they happen
- **Multi-Cluster Support**: Monitor chaos experiments across different clusters simultaneously

**Quick Setup:**

```bash
# Install Kite via Helm
helm repo add kite https://zxh326.github.io/kite
helm repo update
helm install kite kite/kite -n kube-system

# Access Kite dashboard
kubectl port-forward -n kube-system svc/kite 8080:8080
```

Then navigate to `http://localhost:8080` to visualize your cluster while running Cascade experiments.

**Using Kite During Chaos Experiments:**

1. Start Kite dashboard and navigate to your target namespace
2. Run your Cascade chaos experiment: `cascade exec`
3. Watch in real-time as pods are terminated and Kubernetes recovers
4. Monitor metrics, logs, and events to validate resilience

This combination gives you both the chaos orchestration power of Cascade and the observability of Kite.

### Building the CLI

Clone the repository and build the CLI:

```bash
make cli
```

### Starting the API Server

Start the API server to trigger chaos engineering tests on your Kubernetes cluster:

```bash
cascade serve start
```

This initializes the server, allowing you to interact with it via API calls or other Cascade CLI commands.

### Stopping the API Server

Stop the API server and clean up resources:

```bash
cascade serve stop
```

This gracefully shuts down the server and removes any associated Docker containers.

### Creating a Chaos Scenario

Create a new chaos engineering scenario:

```bash
cascade create
```

This interactive command walks you through defining a chaos scenario, including target selection, fault injection parameters, and execution strategies.

### Executing a Chaos Scenario

Execute an existing chaos scenario:

```bash
cascade exec
```

This command guides you through selecting and triggering a predefined chaos engineering scenario on your target Kubernetes cluster.

## Configuration

Cascade CLI uses a YAML configuration file to store settings and scenario definitions. By default, it looks for a `config.yaml` file in the current directory.

### Example Configuration

```yaml
scenario:
  # Name of the chaos experiment
  id: scenario
  # Description of the chaos experiment
  description: this is a sample scenario

# Defines the targets for chaos experiment
target:
  # Namespace or set of namespaces to target
  namespaces: test
  # Pods containing the given string in their name become targets
  includedPodNames: chaos
  # Exclude pods on nodes containing the given string
  includedNodeNames: chaos
  # Pods containing the given string are spared from experiments
  excludedPodNames: chaos

# Defines the session attributes for the chaos experiment
runtime:
  # Interval at which chaos experiments are triggered (defaults to 10m)
  interval: 10m
  # Grace time before the chaos experiment starts (defaults to 1m)
  grace: 1m
  # Execution strategy: evict, delete, dry-run (defaults to delete)
  mode: dry-run
  # Pod ordering strategy: oldest, youngest, cost, random (defaults to random)
  ordering: default
  # Ratio of candidate pods to target (defaults to 0.5)
  ratio: 0.5

# Defines the cluster attributes for the chaos experiment
cluster:
  # Path to the kubeconfig file
  kubeconfig: "/path/to/kubeconfig"
  # Path to the master configuration file
  master: "https://master.example.com"
  # Type of origin: host or cluster (defaults to host)
  origin: host
  # Health check port for the pods
  healthcheck: ":8080"
```

## Docker Setup

### Multi-Stage Build

The Dockerfile uses a multi-stage build for both development and production environments. The development stage includes auto-reload functionality backed by Air, allowing for seamless code changes during development.

### Building Docker Images

**Development:**

```bash
docker build --target development -t cascade:dev .
```

**Production:**

```bash
docker build --target production -t cascade:prod .
```

### Docker Compose

The project includes a Docker Compose configuration with a PostgreSQL container for database management.

**Start services:**

```bash
docker compose up -d --build
```

**Teardown:**

```bash
docker compose down -v
```

## Contributing

Contributions are what make the open source community an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

Feel free to contribute by sending us your suggestions, bug reports, or cat videos.

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

Cascade is provided "as is" and comes with absolutely no guarantees. If it breaks your system, well, that's kind of the point, isn't it? Congratulations, you're now doing chaos engineering!

**Use at your own risk.** Side effects may include improved system resilience, fewer 3 AM panic attacks, and an irresistible urge to push big red buttons.

## Credits

Created by an engineer who loves to break things for a living and sleep soundly at night.

Special thanks to Murphy's Law for the inspiration.
