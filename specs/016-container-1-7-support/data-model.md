# Data Model: Container 1.7 Support

**Feature**: specs/016-container-1-7-support/spec.md | **Date**: 2026-11-09

No new persisted models. All entities are transient TUI state plus CLI command arguments. Existing `models` types are reused; only Kubernetes builder inputs and list parsing change.

## Entities

### Local Kubernetes Cluster (existing `models.KubernetesCluster`, parser extended)

- Fields (unchanged struct): `Name`, `State` (running/stopped/unknown), `CPUs`, `Memory`, `Address`, `Ports`, `Nodes[]` (`Name`, `Role`, `State`, `CPUs`, `Memory`, `Address`, `Ports`).
- Source: `container cluster list` whitespace table, now in TWO shapes:
  - Default (1.7.1): one row per cluster (`NAME NODES K8S VERSION STATE`) — cluster-level aggregate row, no per-node detail.
  - Per-node (`--nodes`): old shape (`CLUSTER NODE ROLE STATE CPUS MEMORY ADDR PORTS`) — existing parsing logic reused verbatim.
- Detection: by header tokens (`NAME` as first header = aggregate shape; `CLUSTER` = per-node shape). Unparseable rows are skipped, never crash; empty output -> empty list + "No Kubernetes clusters found."
- State transitions: list -> select -> submenu; delete success -> back to list. No start, no load-image transitions exist anymore.

### Cluster Operation (transient, no struct)

- list: `container cluster list` (default aggregate); per-node detail via `container cluster list --nodes`.
- create: name?, cpus?, memory?, rm?, node-image? (unchanged) + workers? (extended: bare `<n>` or `count=<n>,cpus=<n>,memory=<size>`, repeat groups passed as repeated `--workers` flags) + cni? (manifest path or `NONE`) + shared-image-store? (boolean flag `--shared-image-store`, no value).
- write-config: cluster name, kubeconfig path? (unchanged).
- delete: cluster name (required), type-to-confirm guard (unchanged).
- load-image: REMOVED — `models.KubernetesLoadImageInput` and `KubernetesLoadImageBuilder` deleted.
- Validation: workers spec must be a non-empty token with parseable non-negative count when provided; empty workers/CNI fields omitted; shared-image-store appends the bare flag only when enabled.

### Container Machine (existing `models.ContainerMachine`, unchanged)

- No 1.7.1 flag changes found in `machine create`/`machine set` help; create/edit flows untouched.

## Relationships

- Cluster (1) -> Nodes (0..n) via list parsing; aggregate rows carry node count only, per-node rows carry full topology. Submenu detail renders whichever the stored cluster holds.
- Cluster (1) -> shared-image-store choice (0..1) captured at create time; no post-create mutation (upstream provides no toggle subcommand).
- Screen removal: `ScreenKubernetesLoadImage` constant, `KubernetesLoadImageScreen` field/state, and its app.go wiring are deleted; no other screen gains or loses fields.
