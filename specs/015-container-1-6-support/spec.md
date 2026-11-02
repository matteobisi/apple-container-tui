# Feature Specification: Container 1.6 Support

**Feature Branch**: `015-container-1-6-support`
**Created**: 2026-11-02
**Status**: Draft
**Input**: User description: "i want to update the actui code, respecting the what's new included in last 2 releases. Current release was build on container 1.5.0 now we are 1.6.1. you could use the official apple container help command and see the GitHub documentation. installed on this mac i've already updated container to 1.6.1. also update all internal documentation starting from readme.md"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Keep managing local Kubernetes clusters after the rename (Priority: P1)

A developer who used the `k` Kubernetes workflow with Apple Container 1.5.0 upgrades the installed CLI to 1.6.1 and can still list, create, delete, load images into, and write kubeconfig for local clusters from actui without learning new syntax.

**Why this priority**: The upstream 1.6.1 release renames `container k8s` to `container cluster`, and the old name no longer works (verified: `container k8s` returns "unknown command"). Without this change every Kubernetes action in actui is broken on 1.6.1, so this is the minimum viable update.

**Independent Test**: With Apple Container 1.6.1 installed, press `k` from the container list, list clusters, create a cluster, load an image, write kubeconfig, and delete the cluster — each showing a command preview using `container cluster ...` and reporting its result.

**Acceptance Scenarios**:

1. **Given** Apple Container 1.6.1 is installed, **When** the developer opens the Kubernetes workflow with `k` and lists clusters, **Then** the application queries local clusters through the supported 1.6.1 command and displays clusters, nodes, and states.
2. **Given** a cluster is selected, **When** the developer requests delete, load-image, or write-config, **Then** the application shows a preview of the supported 1.6.1 command and reports the outcome after approval.
3. **Given** the developer used the old workflow that offered a cluster start action, **When** they open the cluster submenu on 1.6.1, **Then** no start action is offered (upstream removed it; stopped nodes are recreated by delete plus create) and the remaining actions all work.

---

### User Story 2 - Create and use multi-worker clusters (Priority: P1)

A developer creates a local Kubernetes cluster with worker nodes and a chosen network plugin setup, then loads an image optionally targeting specific nodes.

**Why this priority**: Multi-worker clusters (`--workers`, `--cni`, per-node `--node` image loading) are the headline Kubernetes feature of the 1.6.0 release; without it actui can only create single-node clusters and cannot use the new capability.

**Independent Test**: Create a cluster specifying worker count and CNI choice, verify the previewed command carries those choices, then load an image with an optional node target and verify the preview and result.

**Acceptance Scenarios**:

1. **Given** the cluster creation form is open, **When** the developer enters a worker count and a CNI choice (manifest path or `NONE`), **Then** the previewed creation command includes those choices and creation reports success or a readable failure.
2. **Given** a multi-node cluster exists, **When** the developer loads an image and optionally names target nodes, **Then** the preview reflects the image, platform default, and node targets, and the result is reported.
3. **Given** the developer leaves worker/CNI/node fields empty, **When** they confirm, **Then** the application behaves exactly as a single-node default cluster creation or cluster-wide image load.

---

### User Story 3 - Commit a container filesystem to an image (Priority: P1)

A developer selects a container and creates a new image from its filesystem through a previewed, guarded action.

**Why this priority**: `container commit <id> <reference>` is the headline new core command of 1.6.0 and has no equivalent in actui today; it is the only new container-lifecycle capability in the last two releases.

**Independent Test**: Select a container, invoke the commit action, enter an image reference, approve the previewed command, and verify the new image appears in the image list (or a readable error is shown).

**Acceptance Scenarios**:

1. **Given** a container is selected, **When** the developer chooses the commit action and provides an image reference, **Then** the application previews the exact commit command before executing it.
2. **Given** the commit preview is shown, **When** the developer approves it, **Then** the application executes it and reports success or a readable failure without leaving the user on a dead screen.
3. **Given** dry-run mode is active, **When** the developer approves the commit preview, **Then** no image is created and the intended command is displayed.

---

### User Story 4 - Create machines with current 1.6 options (Priority: P2)

A developer creates a container machine with resource and boot options supported by Apple Container 1.6.1, and adjusts machine resources including the new settings.

**Why this priority**: Machine creation in actui only supports image plus name, while 1.6.1 supports CPUs, memory, home-mount mode, nested virtualization, custom kernel, and boot behavior; closing the most-used gaps keeps machine management usable without exposing every exotic flag.

**Independent Test**: Create a machine specifying CPUs, memory, and home-mount mode; verify the preview carries them; edit an existing machine's resources and verify the preview and result.

**Acceptance Scenarios**:

1. **Given** the machine creation form is open, **When** the developer fills in CPUs, memory, and home-mount mode, **Then** the previewed creation command includes those values.
2. **Given** a machine is selected, **When** the developer edits its resources, **Then** the application supports CPU, memory, and home-mount changes consistent with the installed CLI.
3. **Given** an unsupported or invalid resource value is entered, **When** the developer confirms, **Then** the application shows a readable validation message before executing anything.

---

### User Story 5 - Find current 1.6 guidance in README and docs (Priority: P2)

A developer or contributor reads the README and internal guides and learns the supported Apple Container version, the renamed Kubernetes commands, the removed start action, the new commit and multi-worker capabilities, and the machine options.

**Why this priority**: The user explicitly requested documentation updates starting with the README; without them users will follow stale `k8s`/start instructions and miss new capabilities.

**Independent Test**: Using only the updated README plus internal guides, describe the supported Container version, the Kubernetes entry point and operations, and the new commit action.

**Acceptance Scenarios**:

1. **Given** a developer reads the README, **When** they look for supported versions and workflows, **Then** they find Apple Container 1.6.1 as the supported backend, the `k` entry point, and the current operation set with no start action.
2. **Given** a contributor changes a Kubernetes, commit, or machine flow, **When** they consult internal documentation, **Then** they can identify screen ownership, command-safety expectations, and the version-reporting policy without conflicting instructions.

### Edge Cases

- The installed CLI may still be 1.5.0 or older: Kubernetes actions using the new `cluster` name will fail there; the failure output must be shown readably and the user must retain navigation (no hard requirement for dual-version support).
- A cluster or container may disappear or change state between selection, preview, and execution; the action must report the failure and return the developer to a usable screen.
- Cluster creation with workers or a custom CNI may partially fail (e.g. host port conflicts noted upstream); the application must surface the CLI's error rather than report success.
- An image reference for commit may be invalid or already exist; validation and CLI errors must be shown before or after preview without crashing.
- Empty cluster list, missing node targets, or unavailable Kubernetes support must produce a readable empty-state or unavailable message with a route back to the container list.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The Kubernetes workflow MUST invoke local cluster operations through the Apple Container 1.6.1 `container cluster` command family (list, create, delete, load-image, write-config).
- **FR-002**: The system MUST NOT offer a cluster start action, because the installed 1.6.1 CLI exposes no cluster start subcommand (upstream removed `k8s start` in 1.5.0 and 1.6.x `cluster` has no start).
- **FR-003**: Cluster creation MUST support worker count and CNI choice (manifest path or `NONE`) in addition to the existing name, CPUs, memory, remove-on-stop, and node-image fields.
- **FR-004**: Cluster image loading MUST support an optional target-node selection for multi-node clusters, in addition to image reference and platform.
- **FR-005**: The system MUST provide a container commit action that builds a `container commit <id> <reference>` command with a visible preview and explicit approval before execution.
- **FR-006**: Machine creation MUST support CPUs, memory, and home-mount mode in addition to image and name.
- **FR-007**: Machine resource editing MUST remain consistent with the installed CLI's settable values (CPUs, memory, home-mount).
- **FR-008**: Before executing any state-changing Kubernetes, commit, or machine operation, the system MUST present a visible command preview and require explicit approval; cluster/machine removal MUST keep the type-to-confirm destructive-action guard.
- **FR-009**: All new and changed operations MUST honor dry-run mode (show the intended command without executing it) and MUST surface success, failure, stdout, and stderr readably.
- **FR-010**: Automated checks MUST cover the renamed cluster command composition, removal of the start action, new create/load-image options, commit command composition and guardrails, and extended machine creation options.
- **FR-011**: The README MUST document Apple Container 1.6.1 as the supported backend, the `container cluster` command family, the absence of a start action, the commit action, multi-worker creation, and machine options.
- **FR-012**: Internal guides that describe navigation, command ownership, or release behavior (user guide, AI menu map, and any workflow doc covering Kubernetes/machines) MUST be updated to match the new behavior without conflicting instructions.

### Key Entities

- **Local Kubernetes Cluster**: A cluster managed by the installed Apple Container CLI, described by name, nodes (role, state, resources), worker topology, CNI choice, and connectivity details.
- **Cluster Operation**: A list, create, delete, load-image, or write-config action with its target, options (workers, CNI, node targets), preview, approval state, and outcome. Start is explicitly not an operation on 1.6.1.
- **Committed Image**: A new image reference produced from an existing container's filesystem via the commit action, with its source container, preview, approval, and outcome.
- **Container Machine**: A machine managed by the installed CLI, described by image, name, CPUs, memory, home-mount mode, and state.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: With Apple Container 1.6.1 installed, 100% of Kubernetes operations offered by the application (list, create, delete, load-image, write-config) execute through the supported command family and report a readable outcome.
- **SC-002**: In keyboard-driven testing, participants can open the Kubernetes workflow with `k`, create a multi-worker cluster, load an image, and delete the cluster, with every state-changing step showing an approval before execution.
- **SC-003**: Participants can commit a container to a new image reference in under 2 minutes from discovering the action, with a preview shown in 100% of attempts.
- **SC-004**: Documentation review confirms the README and every affected internal guide describe the 1.6.1 backend, the renamed cluster commands, the removed start action, the commit action, and machine options with no conflicting instructions.
- **SC-005**: 100% of tested unsupported-version, empty-state, and failed-operation cases produce a readable explanation and preserve a route back to a usable screen.

## Assumptions

- Target environment is macOS on Apple Silicon with Apple Container 1.6.1 installed (verified on the author's machine via `container --version`); 1.5.0 compatibility shims (e.g. trying both `k8s` and `cluster`) are out of scope.
- Upstream `container build`/`container builder` becoming CLI plugins does not change their command syntax, so no actui build-flow change is required for that packaging change alone.
- DNS alias options on `run`/`create --network` (1.6.0 networking) do not require actui changes because actui has no container run/create screen.
- Network and volume management screens remain out of scope per the constitution (container networking configuration and volume management beyond CLI defaults are excluded).
- Existing conventions apply unchanged: keyboard navigation, command preview before execution, type-to-confirm for destructive actions, dry-run support, JSONL logging, and local-only operation.
- Documentation updates cover README.md plus the user guide, AI menu map, and any workflow doc that describes Kubernetes, commit, machine, or release behavior.
