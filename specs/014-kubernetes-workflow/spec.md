# Feature Specification: Kubernetes Workflow

**Feature Branch**: `014-kubernetes-workflow`  
**Created**: 2026-09-14  
**Status**: Draft  
**Input**: User description: "Update the app for Container 1.4.1, add a `k` keybinding for a Kubernetes operations submenu, synchronize the displayed app version with GitHub releases, and update internal documentation starting with the README."

## Clarifications

### Session 2026-09-15

- Q: Should the approved analysis remediations be applied? -> A: Apply the constitution amendment, cluster-list navigation, corrected app-routing task order, and UI-level unavailable-command handling.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Open Kubernetes Operations (Priority: P1)

A developer using Apple Container 1.4.1 opens the Kubernetes workflow from the container list with the `k` key and can select a Kubernetes operation from a dedicated submenu.

**Why this priority**: The keyboard entry point and operation menu are the minimum usable Kubernetes experience and expose Container's Kubernetes capability without requiring users to remember command syntax.

**Independent Test**: From the container list, press `k`, confirm the Kubernetes submenu appears, move through each available operation, and return to the container list without executing an operation.

**Acceptance Scenarios**:

1. **Given** the container list is visible, **When** the developer presses lowercase `k`, **Then** the application opens the Kubernetes cluster list.
2. **Given** the Kubernetes cluster list is visible, **When** the developer selects a cluster with `enter`, **Then** the application opens its Kubernetes operations submenu.
3. **Given** the Kubernetes submenu is visible, **When** the developer chooses an available operation, **Then** the application presents the operation's next step or command preview.
4. **Given** the Kubernetes cluster list or submenu is visible, **When** the developer cancels or goes back, **Then** the application returns to its parent screen and ultimately to the container list.

---

### User Story 2 - Safely Manage a Kubernetes Cluster (Priority: P1)

A developer can use the Kubernetes submenu to view local Kubernetes clusters and perform the supported operations to create, start, remove, load an image into, or write configuration for one.

**Why this priority**: The workflow must support both discovery and the routine lifecycle work developers perform after enabling Kubernetes locally.

**Independent Test**: Using a supported local Container installation, complete each supported Kubernetes operation from the submenu and verify that the application reports the resulting state or a readable failure.

**Acceptance Scenarios**:

1. **Given** one or more Kubernetes clusters exist, **When** the developer selects the list operation, **Then** the application displays the available clusters, nodes, and states.
2. **Given** a cluster is selected, **When** the developer requests start, image loading, or configuration writing, **Then** the application shows the selected operation and its result.
3. **Given** a developer requests cluster creation, image loading, configuration writing, or removal, **When** the operation requires confirmation, **Then** the application displays the action for approval before execution.
4. **Given** a removal operation is pending, **When** the developer declines it, **Then** no cluster is removed.

---

### User Story 3 - Understand Unsupported or Failed Kubernetes Actions (Priority: P2)

A developer receives a clear, actionable result when the installed Container version does not support Kubernetes, no clusters exist, or an operation fails.

**Why this priority**: Kubernetes support depends on the local Container installation; clear feedback prevents ambiguous or unsafe actions.

**Independent Test**: Run the Kubernetes workflow against an unsupported installation, an installation with no clusters, and a simulated command failure; verify each case gives a readable outcome and leaves the user in a navigable state.

**Acceptance Scenarios**:

1. **Given** Kubernetes is unavailable in the local Container installation, **When** the developer opens or uses the Kubernetes workflow, **Then** the application explains that Kubernetes support is unavailable and preserves a route back to the container list.
2. **Given** no Kubernetes clusters exist, **When** the developer lists them, **Then** the application clearly reports the empty state and offers the create operation.
3. **Given** a Kubernetes operation fails, **When** its result is shown, **Then** the application identifies the failed operation and surfaces readable diagnostic output.

---

### User Story 4 - Identify the Installed Release (Priority: P1)

A developer running a downloaded GitHub release sees the same release version from `actui --version` that is shown for that release on GitHub.

**Why this priority**: Accurate version output is necessary for support, diagnostics, and verifying that a downloaded binary is current.

**Independent Test**: Download or build at least two distinct versioned release artifacts, run `actui --version` for each, and compare the output to its GitHub release label.

**Acceptance Scenarios**:

1. **Given** a developer runs a published release binary, **When** they execute `actui --version`, **Then** its version exactly matches the associated GitHub release label without the release prefix.
2. **Given** a developer builds the application outside the release process, **When** they execute `actui --version`, **Then** the output identifies the build in a documented, non-misleading way.

---

### User Story 5 - Find Current Kubernetes Guidance (Priority: P2)

A developer can find the Kubernetes workflow, supported operations, safety behavior, and release-version behavior in the README and affected internal guides.

**Why this priority**: The new workflow and versioning behavior must be discoverable and maintainable as the application evolves.

**Independent Test**: Review the README and relevant documentation, then use only that guidance to locate the `k` entry point, describe the available operations and confirmations, and explain how release versions are reported.

**Acceptance Scenarios**:

1. **Given** a developer reads the README, **When** they look for supported workflows, **Then** they find Kubernetes support and its `k` entry point.
2. **Given** a contributor changes the Kubernetes flow or release process, **When** they consult internal documentation, **Then** they can identify the relevant ownership, command-safety expectations, and versioning policy.

### Edge Cases

- The `k` key must remain distinct from existing lowercase navigation keys and must not trigger an action while text input or confirmation is focused.
- Kubernetes commands may be unavailable, return no clusters, or return partially populated cluster data; the workflow must present a safe, readable outcome rather than an incorrect state.
- A lifecycle target may disappear or change state between selection, preview, and execution; the action must report the failure and return the developer to a usable screen.
- The release label may be unavailable during a local build; version output must not claim a published release version it cannot establish.
- A release artifact must not report an older fixed version after a newer GitHub release is published.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide a Kubernetes workflow from the container list when the developer presses lowercase `k`.
- **FR-002**: The Kubernetes workflow MUST present a dedicated submenu for Kubernetes operations and allow the developer to return to the container list without performing an operation.
- **FR-003**: The submenu MUST support listing, creating, starting, removing, loading an image into, and writing Kubernetes configuration for local Kubernetes clusters supported by the installed Container CLI.
- **FR-004**: Before executing any state-changing Kubernetes operation, the system MUST present a visible operation preview and require explicit approval.
- **FR-005**: Before removing a Kubernetes cluster, the system MUST require an explicit destructive-action confirmation and provide a cancellation path that makes no change.
- **FR-006**: The system MUST display the outcome of each Kubernetes operation, including readable standard output and error details when available.
- **FR-007**: When Kubernetes is unsupported, unavailable, or has no clusters, the system MUST show a specific, readable status and retain keyboard navigation.
- **FR-008**: The system MUST support a dry-run mode for all Kubernetes operations that displays the intended action without changing the local cluster.
- **FR-009**: Each published binary MUST report a version through `actui --version` that exactly matches its associated GitHub release version, apart from an optional `v` prefix.
- **FR-010**: Builds that are not associated with a published release MUST display a documented version identifier that is distinguishable from a release version.
- **FR-011**: The README MUST document Kubernetes availability, the `k` entry point, the supported Kubernetes operations, and destructive-action safeguards.
- **FR-012**: Internal documentation that maps navigation or release behavior MUST be updated to reflect the Kubernetes workflow and the version-reporting policy.
- **FR-013**: Automated checks MUST cover Kubernetes command composition, destructive-action guardrails, keyboard navigation to and from the submenu, and version reporting for a release-labeled build and a non-release build.

### Key Entities

- **Kubernetes Cluster**: A local Kubernetes cluster managed through Apple Container, identified by its displayed name, nodes, state, and connectivity details.
- **Kubernetes Operation**: A requested list, create, start, remove, load-image, or write-config action together with its target, preview, approval state, and outcome.
- **Build Version**: The version identifier displayed by an application binary, associated either with a published GitHub release or a clearly distinguished non-release build.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In keyboard-driven usability testing, 100% of participants can open the Kubernetes submenu from the container list using `k` and return without assistance within 15 seconds.
- **SC-002**: All six supported Kubernetes operations can be initiated from the submenu, and every state-changing operation displays an approval step before execution.
- **SC-003**: In 100% of tested unsupported, empty-state, and failed-operation cases, the application provides a readable explanation and a usable route back to the container list.
- **SC-004**: For 100% of tested published release artifacts, `actui --version` matches the corresponding GitHub release version.
- **SC-005**: Documentation review confirms the README and every affected internal guide describe the Kubernetes entry point, supported operations, safeguards, and release-version behavior without conflicting instructions.
- **SC-006**: At least 90% of developers following the updated documentation can locate and describe the Kubernetes workflow and version-reporting behavior on their first attempt.

## Assumptions

- The target environment is macOS on Apple Silicon with Apple Container 1.4.1 installed; Kubernetes support is available starting with Container 1.2.0.
- The Kubernetes workflow manages only local clusters exposed by the installed Apple Container CLI; remote clusters, cloud providers, workload deployment, networking configuration, and persistent-volume administration are out of scope.
- The installed Container 1.4.1 CLI supports list, create, start, delete, load-image, and write-config. It does not expose separate stop or inspect subcommands; optional capabilities outside these operations are deferred unless they are necessary to make one of these flows usable.
- Existing application conventions for keyboard navigation, command previews, dry-run behavior, confirmations, result reporting, and local-only operation apply to the Kubernetes workflow.
- Published GitHub releases remain the source of truth for release version labels, while local development builds are allowed a documented non-release identifier.
- Documentation updates include the README plus the AI menu map and binary-build/release guide when those documents describe affected behavior.
