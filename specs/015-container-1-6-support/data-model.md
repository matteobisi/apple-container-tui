# Data Model: Container 1.6 Support

**Feature**: specs/015-container-1-6-support/spec.md | **Date**: 2026-11-02

No new persisted models. All entities are transient TUI state plus CLI command arguments. Existing `models` types are reused; only builder inputs change.

## Entities

### Local Kubernetes Cluster (existing `models.KubernetesCluster`, unchanged)

- Fields: `Name`, `State` (running/stopped/unknown), `CPUs`, `Memory`, `Address`, `Ports`, `Nodes[]` (`Name`, `Role`, `State`, `CPUs`, `Memory`, `Address`, `Ports`).
- Source: parsed from `container cluster list` whitespace table (same shape as before; parser unchanged).
- Validation: empty output -> empty list + "No Kubernetes clusters found."; unparseable rows are skipped, never crash.
- State transitions: list -> select -> submenu; delete success -> back to list. No start transition exists anymore.

### Cluster Operation (transient, no struct)

- list: no args. create: name?, cpus?, memory?, rm?, node-image? (existing) + workers? (positive int, `0`/empty omitted), cni? (manifest path or `NONE`, empty omitted).
- load-image: cluster name, image (required), platform? + nodes? (comma-separated input split into repeated `--node <name>` flags, empty omitted).
- write-config: cluster name, kubeconfig path? (unchanged).
- delete: cluster name (required), type-to-confirm guard (unchanged).
- Validation: workers must be a non-negative integer when provided; node names are non-empty trimmed tokens; empty CNI/workers/nodes fields are omitted from the command.

### Committed Image (transient, no new struct)

- Inputs: source container ID (from selected `models.Container`), image reference (required, non-empty token).
- Command: `container commit <id> <reference>`. Preview + explicit approval; honors dry-run; result reported via existing result rendering.

### Container Machine (existing `models.ContainerMachine`, unchanged)

- Create inputs gain: cpus? (positive int), memory? (non-empty token, e.g. `8G`), home-mount? (`ro|rw|none`).
- `MachineSetBuilder` already covers cpus/memory/home-mount, so edit-resources needs no model change.
- Validation: cpus positive integer; home-mount one of ro/rw/none; invalid values rejected before preview.

## Relationships

- Container (1) -> Committed Image (0..n) via commit action from container submenu.
- Cluster (1) -> Nodes (0..n) via list parsing; load-image optionally targets a subset of nodes.
- Machine (1) <-> edit/create inputs via shared cpus/memory/home-mount value set.
