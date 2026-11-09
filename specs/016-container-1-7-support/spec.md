# Feature Specification: Container 1.7 Support

**Feature Branch**: `016-container-1-7-support`
**Created**: 2026-11-09
**Status**: Draft
**Input**: User description: "i want to update the actui code, respecting the what's new Apple Continer. Current release was build on container 1.6.1 now we are 1.7.1. you could use the official apple container help command and see the GitHub documentation. installed on this mac i've already updated container to 1.7.1. also update all internal documentation starting from readme.md. I saw readme say as pre-req macOS 26, update it to 27, Apple Silicon"

Verified baseline on author machine: `container CLI version 1.7.1 (build: release, commit: a15b460)`, `sw_vers ProductVersion 27.2`. Upstream highlights source: Apple Container 1.7.1 release notes (`container cluster load-image` removed in favor of `--shared-image-store` on create; `cluster list` shows one row per cluster with `--nodes` for per-node view; container lifecycle commands no longer operate on cluster nodes; worker CPU/memory via `--workers count=<n>,cpus=<n>,memory=<size>`; `stats` reports `cpuUtilization` in machine-readable output).

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Manage Kubernetes clusters without the removed load-image command (Priority: P1)

A developer who used the `k` Kubernetes workflow with Apple Container 1.6.1 upgrades the installed CLI to 1.7.1 and can still manage local clusters from actui without hitting broken commands.

On 1.7.1 `container cluster load-image` no longer exists (verified: `container cluster load-image --help` falls back to cluster help; subcommands are only `create`, `delete`, `list`, `write-config`). Upstream replaces it with `--shared-image-store` on cluster creation, so locally built images are usable from cluster nodes without an explicit load step.

**Why this priority**: Every load-image action offered by actui today is broken on 1.7.1. Removing the dead workflow and exposing its replacement is the minimum viable update.

**Independent Test**: With Apple Container 1.7.1 installed, press `k`, list clusters, create a cluster with the shared image store option, write kubeconfig, and delete the cluster — each showing a command preview using only supported 1.7.1 subcommands, with no load-image action offered anywhere.

**Acceptance Scenarios**:

1. **Given** Apple Container 1.7.1 is installed, **When** the developer opens the Kubernetes workflow with `k` and selects a cluster, **Then** no load-image action is offered and the remaining actions (list, create, delete, write-config) all work.
2. **Given** the cluster creation form is open, **When** the developer enables the shared image store option, **Then** the previewed creation command carries that choice and creation reports success or a readable failure.
3. **Given** dry-run mode is active, **When** the developer approves any cluster operation, **Then** no command is executed and the intended 1.7.1 command is displayed.

---

### User Story 2 - Read cluster topology in the one-row-per-cluster list (Priority: P1)

A developer lists local Kubernetes clusters and understands at a glance which clusters exist and, when needed, the state of their individual nodes.

Upstream changed `container cluster list` to show one row per cluster by default, with `--nodes` restoring the per-node view.

**Why this priority**: actui parses cluster list output to render the cluster browser; if parsing still assumes the old per-node table shape, the cluster list misleads or breaks on 1.7.1.

**Independent Test**: With one multi-node cluster present, list clusters from actui and confirm clusters appear once each with correct state, then drill into or toggle the per-node detail view and confirm node rows match the CLI's `--nodes` output.

**Acceptance Scenarios**:

1. **Given** clusters exist on 1.7.1, **When** the developer lists clusters, **Then** each cluster appears once with its name, state, and resources matching the supported list output.
2. **Given** a multi-node cluster exists, **When** the developer requests node detail, **Then** per-node rows (role, state, resources) are shown consistent with the CLI's per-node view.
3. **Given** no clusters exist, **When** the developer lists clusters, **Then** a readable empty-state message with a route back to the container list is shown.

---

### User Story 3 - Size worker nodes at creation time (Priority: P2)

A developer creates a local Kubernetes cluster with worker nodes sized for the workload, choosing worker count plus per-group CPUs and memory.

Upstream 1.7.1 documents `--workers count=<n>,cpus=<n>,memory=<size>` (repeatable for differently sized groups); control-plane sizing via `--cpus`/`--memory` and `--cni`, `--node-image`, `--rm` remain as before.

**Why this priority**: Worker sizing is the headline new cluster capability in 1.7.1. Without it actui can only create default-sized workers and cannot use the new functionality.

**Independent Test**: Create a cluster specifying worker count with CPUs and memory, verify the previewed command carries those choices, and confirm creation reports success or a readable failure.

**Acceptance Scenarios**:

1. **Given** the cluster creation form is open, **When** the developer enters a worker specification with count, CPUs, and memory, **Then** the previewed creation command includes that specification.
2. **Given** the developer leaves worker sizing empty, **When** they confirm, **Then** creation behaves exactly as a default single-node cluster creation.
3. **Given** an invalid worker specification is entered, **When** the developer confirms, **Then** a readable validation message appears before anything is executed.

---

### User Story 4 - Find current 1.7.1 guidance in README and docs (Priority: P2)

A developer or contributor reads the README and internal guides and learns the supported Apple Container version (1.7.1), the supported macOS version (27.x on Apple Silicon), the removal of the load-image workflow, the shared-image-store replacement, the one-row-per-cluster list behavior, and the new worker sizing options.

**Why this priority**: The user explicitly requested documentation updates starting with the README, including the macOS 26 to 27 prerequisite change. Without them users will follow stale load-image instructions and install on an unsupported OS.

**Independent Test**: Using only the updated README plus internal guides, describe the supported Container and macOS versions, the Kubernetes entry point and operation set, and what replaced image loading.

**Acceptance Scenarios**:

1. **Given** a developer reads the README, **When** they look for supported versions and workflows, **Then** they find Apple Container 1.7.1, macOS 27.x on Apple Silicon, the `k` entry point, and the current operation set with no load-image action.
2. **Given** a contributor changes a Kubernetes, commit, or machine flow, **When** they consult internal documentation, **Then** they can identify screen ownership, command-safety expectations, and the version-reporting policy without conflicting instructions.

### Edge Cases

- The installed CLI may still be 1.6.1 or older: cluster actions using removed or renamed options will fail there; the failure output must be shown readably and the user must retain navigation (no hard requirement for dual-version support).
- A cluster may disappear or change state between selection, preview, and execution; the action must report the failure and return the developer to a usable screen.
- Cluster creation with workers, a custom CNI, or the shared image store may partially fail (e.g. host port conflicts noted upstream); the application must surface the CLI's error rather than report success.
- Container lifecycle commands no longer operate on cluster nodes upstream; any attempt against a node identifier must surface the CLI's readable refusal rather than a success.
- Empty cluster list or unavailable Kubernetes support must produce a readable empty-state or unavailable message with a route back to the container list.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The Kubernetes workflow MUST NOT offer a load-image action, because the installed 1.7.1 CLI exposes no `container cluster load-image` subcommand.
- **FR-002**: Cluster creation MUST support the shared image store option in addition to the existing name, CPUs, memory, remove-on-stop, node-image, worker, and CNI fields, with the choice carried into the previewed creation command.
- **FR-003**: Cluster creation MUST support worker resource sizing (count with CPUs and memory, repeatable for multiple groups) matching the installed CLI's worker specification format.
- **FR-004**: The cluster list MUST correctly render the 1.7.1 one-row-per-cluster default output and MUST provide access to per-node detail consistent with the CLI's per-node view.
- **FR-005**: The Kubernetes workflow MUST continue to invoke local cluster operations through the Apple Container 1.7.1 `container cluster` command family (list, create, delete, write-config) with no start and no load-image operation.
- **FR-006**: Before executing any state-changing cluster operation, the system MUST present a visible command preview and require explicit approval; cluster removal MUST keep the type-to-confirm destructive-action guard.
- **FR-007**: All changed operations MUST honor dry-run mode (show the intended command without executing it) and MUST surface success, failure, stdout, and stderr readably.
- **FR-008**: Automated checks MUST cover the removal of the load-image path, the shared-image-store and worker-sizing command composition, the updated list parsing for both default and per-node shapes, and the absence of start/load-image actions.
- **FR-009**: The README MUST document Apple Container 1.7.1 as the supported backend, macOS 27.x on Apple Silicon as the prerequisite, the `container cluster` command family, the removed load-image action and its shared-image-store replacement, the one-row-per-cluster list behavior, and worker sizing.
- **FR-010**: Internal guides that describe navigation, command ownership, or release behavior (user guide, AI menu map, and any workflow doc covering Kubernetes) MUST be updated to match the new behavior without conflicting instructions.
- **FR-011**: The project constitution's platform target (currently macOS 26.x) MUST be amended to macOS 27.x on Apple Silicon as part of this change, recording rationale and migration notes per the governance rules.

### Key Entities

- **Local Kubernetes Cluster**: A cluster managed by the installed Apple Container CLI, described by name, nodes (role, state, resources), worker topology and sizing, shared-image-store choice, CNI choice, and connectivity details.
- **Cluster Operation**: A list, create, delete, or write-config action with its target, options (workers with CPU/memory sizing, shared image store, CNI, node targets for display), preview, approval state, and outcome. Start and load-image are explicitly not operations on 1.7.1.
- **Container Machine**: A machine managed by the installed CLI, described by image, name, CPUs, memory, home-mount mode, and state (unchanged from 1.6.1 support; no new machine flags were found in 1.7.1 help output).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: With Apple Container 1.7.1 installed, 100% of Kubernetes operations offered by the application (list, create, delete, write-config) execute through the supported command family and report a readable outcome, with zero load-image or start actions offered.
- **SC-002**: In keyboard-driven testing, participants can open the Kubernetes workflow with `k`, create a cluster with worker sizing and the shared image store, view per-node detail, write kubeconfig, and delete the cluster, with every state-changing step showing an approval before execution.
- **SC-003**: Documentation review confirms the README and every affected internal guide describe the 1.7.1 backend, macOS 27.x on Apple Silicon, the current cluster operation set, the load-image removal and its replacement, and worker sizing, with no conflicting instructions.
- **SC-004**: 100% of tested unsupported-version, empty-state, and failed-operation cases produce a readable explanation and preserve a route back to a usable screen.

## Assumptions

- Target environment is macOS 27.x on Apple Silicon with Apple Container 1.7.1 installed (verified on the author's machine via `container --version` and `sw_vers`: macOS 27.2); 1.6.1 compatibility shims are out of scope.
- Evaluated 1.7.1 highlights requiring no actui change: `container stats` `cpuUtilization` reporting (actui has no stats screen), Time Machine exclusions, Kata kernel and Containerization version bumps, default DNS `config.toml` settings (daemon-side configuration, no TUI surface), and BuildKit build fixes (no build syntax change found in `container build --help`).
- Container lifecycle, image, machine, registry, export, commit, and daemon flows are unchanged unless 1.7.1 help output proves otherwise; machine `create`/`set` help shows the same options as the 1.6.1 support baseline (CPUs, memory, home-mount, virtualization, kernel).
- Existing conventions apply unchanged: keyboard navigation, command preview before execution, type-to-confirm for destructive actions, dry-run support, JSONL logging, and local-only operation.
- Documentation updates cover README.md plus the user guide, AI menu map, and any workflow doc that describes Kubernetes behavior or release/version prerequisites.
