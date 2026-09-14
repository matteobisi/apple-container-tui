# Data Model: Kubernetes Workflow

## KubernetesCluster

Represents a local Apple Container Kubernetes cluster gathered from the list
output.

| Field | Description | Validation |
|-------|-------------|------------|
| `Name` | Cluster name used by target commands | Required, non-blank |
| `Nodes` | Nodes belonging to the cluster | May be empty when list output is partial |
| `State` | Safe aggregate display state | `running`, `stopped`, or `unknown` |
| `CPUs` | Displayed allocated CPU count | Optional display field |
| `Memory` | Displayed allocated memory | Optional display field |
| `Address` | Displayed API address | Optional display field |
| `Ports` | Displayed port mapping | Optional display field |

State is derived from reported node states: all running yields `running`, all
stopped yields `stopped`, and mixed, absent, or unrecognized input yields
`unknown`.

## KubernetesNode

Represents one row reported by `container k8s list`.

| Field | Description | Validation |
|-------|-------------|------------|
| `Name` | Node name | Required to retain a parsed row |
| `Role` | Node role text | Optional |
| `State` | Reported node state | Normalized to known state or `unknown` |
| `CPUs` | Allocated CPU count | Optional |
| `Memory` | Allocated memory | Optional |
| `Address` | Node/API address | Optional |
| `Ports` | Port mapping | Optional; may wrap in CLI output |

## KubernetesCreateInput

Captures user-entered values for `container k8s create`.

| Field | Required | Validation |
|-------|----------|------------|
| `Name` | No | Non-blank token when provided; CLI defaults to `k8s-dev` |
| `CPUs` | No | Positive integer when provided |
| `Memory` | No | CLI-supported size with optional K/M/G/T/P suffix |
| `RemoveOnStop` | No | Boolean |
| `NodeImage` | No | Non-blank image reference when provided |

`Scheme` and `MaxConcurrentDownloads` remain at Container CLI defaults in this
feature because the workflow does not need them for a usable local cluster
lifecycle.

## KubernetesLoadImageInput

| Field | Required | Validation |
|-------|----------|------------|
| `ClusterName` | No | Non-blank token when provided; CLI default applies |
| `Image` | Yes | Non-blank image reference |
| `Platform` | No | `os/arch[/variant]` when provided |

## KubernetesWriteConfigInput

| Field | Required | Validation |
|-------|----------|------------|
| `ClusterName` | No | Non-blank token when provided; CLI default applies |
| `KubeconfigPath` | No | Non-blank path when provided; CLI default applies |

## KubernetesOperation

Tracks an operation initiated from the workflow: action, optional target,
constructed command, approval state, and resulting stdout/stderr. Read-only
list actions execute directly. Create, start, load-image, and write-config need
a command preview. Delete requires a type-to-confirm state keyed to the target
cluster name. The existing dry-run executor returns a result without executing
the constructed command.

## BuildVersion

| Field | Description | Rule |
|-------|-------------|------|
| `Version` | Value printed by `actui --version` | Release builds use the GitHub tag without `v`; local builds use `dev` |
| `ReleaseTag` | CI-generated semantic tag | Persisted from build to publish and used unchanged |