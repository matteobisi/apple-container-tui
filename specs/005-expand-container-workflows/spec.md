# Feature Specification: Expanded Container Workflows

**Feature Branch**: `005-expand-container-workflows`  
**Created**: March 31, 2026  
**Status**: Draft  
**Input**: User description: "i want to add some new feature to actui to expand the features 1. container registry list -> new Registries screen 2. container export -> new action in container submenu to export containers 3. --pull flag -> checkbox in build form 4. system status --format -> more robust daemon status parsing"

## Clarifications

### Session 2026-03-31

- Q: Which registries should the new Registries screen show? -> A: Registries explicitly configured or exposed by the current container environment/runtime.
- Q: How should the container export destination be collected? -> A: The user selects a destination directory and the app generates the export filename.
- Q: What should the default value be for the build refresh checkbox? -> A: Default to enabled.
- Q: Which containers should be eligible for export? -> A: Only stopped containers.
- Q: How should daemon status parsing use formatted output? -> A: Use structured formatted output as the primary source, falling back to unknown if required fields are missing.

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Browse Registries in the TUI (Priority: P1)

Users who manage images need a dedicated Registries screen so they can review which container registries are available without leaving the TUI or guessing from image references.

**Why this priority**: Registry visibility fills a clear navigation gap and supports downstream workflows such as pulling, building, and troubleshooting authentication. It adds immediate operational value as a standalone screen.

**Independent Test**: Open the Registries screen from the TUI, confirm that available registries are listed with distinguishing details, and verify the user can leave the screen and return to their previous context without using external commands.

**Acceptance Scenarios**:

1. **Given** the user is in the main TUI, **When** they navigate to Registries, **Then** a dedicated Registries screen is shown instead of requiring shell commands outside the application
2. **Given** one or more registries are available, **When** the Registries screen loads, **Then** each registry is displayed with enough identifying information for the user to distinguish entries confidently
3. **Given** no registries are available or the list cannot be loaded, **When** the Registries screen is shown, **Then** the user sees a clear empty or error state with guidance on what to do next

---

### User Story 2 - Export a Container from Its Submenu (Priority: P2)

Users who want to archive, share, or move a container need an export action directly in the container submenu so they can complete that workflow from the same place they already manage container lifecycle actions.

**Why this priority**: Export is a high-value operational task and fits naturally beside existing container actions. Adding it to the submenu reduces workflow fragmentation and keeps advanced container management inside the TUI.

**Independent Test**: Open any container submenu, select Export, provide a destination, run the export, and confirm the workflow reports success or failure without altering the source container unexpectedly.

**Acceptance Scenarios**:

1. **Given** the user has opened a stopped container submenu, **When** they review available actions, **Then** Export container appears as a selectable action
2. **Given** the user chooses Export container, **When** they select a destination directory, **Then** the TUI generates the export filename, starts the export workflow, and shows progress or completion feedback
3. **Given** the export cannot be completed, **When** the workflow ends, **Then** the user receives a specific failure message and the original container remains available in the list

---

### User Story 3 - Choose Freshness During Image Build (Priority: P3)

Users building images need a visible checkbox in the build form that controls whether the build should fetch newer base layers before execution, so they can make that decision intentionally instead of relying on hidden defaults.

**Why this priority**: Build behavior affects correctness and repeatability. Exposing the choice in the form improves trust in the build workflow and reduces accidental use of stale or unexpectedly refreshed base images.

**Independent Test**: Open the build form, confirm the pull option is visible as a checkbox, toggle it on and off, and verify the selected value is reflected in the build confirmation or preview before execution.

**Acceptance Scenarios**:

1. **Given** the user opens the build workflow, **When** the form is displayed, **Then** a checkbox is available for choosing whether the build refreshes base image content before running and it is enabled by default
2. **Given** the user changes the checkbox value, **When** they review the pending build settings, **Then** the selected behavior is shown clearly before the build starts
3. **Given** the user submits the build form, **When** the build begins, **Then** the chosen refresh behavior is applied consistently to that build request

---

### User Story 4 - Trust Daemon Status Feedback (Priority: P4)

Users monitoring the container environment need the daemon status view to interpret formatted status output reliably, so the TUI does not misreport whether the system is running, stopped, or in an unknown state.

**Why this priority**: Incorrect status reporting undermines confidence in operational controls and leads users to take the wrong action. This change improves reliability in an existing screen rather than introducing a new workflow, which makes it lower priority than the first three stories.

**Independent Test**: Open the daemon control screen against representative status outputs, including formatted output variants, and verify the displayed state matches the real daemon condition or falls back to an explicit unknown state.

**Acceptance Scenarios**:

1. **Given** the daemon returns structured formatted status output indicating it is running, **When** the daemon screen refreshes, **Then** the TUI shows the daemon as running
2. **Given** the daemon returns structured formatted status output indicating it is not running, **When** the daemon screen refreshes, **Then** the TUI shows the daemon as stopped
3. **Given** the structured status response is incomplete, unexpected, or missing required fields, **When** the daemon screen refreshes, **Then** the TUI shows an unknown state with guidance instead of a misleading running or stopped result

---

### Edge Cases

- What happens when the user opens the Registries screen but the environment has no configured or discoverable registries? The screen should remain navigable and explain that no registry entries are currently available.
- What happens when registry details are partially available, such as a host without authentication or status metadata? Each entry should still render with the information that is known rather than being dropped entirely.
- What happens when the user starts a container export and the generated output file already exists or the chosen directory is not writable? The export must stop safely and present a clear corrective message.
- What happens when the user opens the submenu for a running container? Export should not be offered as an available action for that container state.
- What happens when the user opens the build form from a previously used session? The pull checkbox should still default to enabled so users do not have to guess the current setting.
- What happens when daemon status output changes format again in a future release or omits required structured fields? The application should prefer an explicit unknown state over incorrectly labeling the daemon as running or stopped.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a dedicated Registries screen within the TUI navigation model.
- **FR-002**: The Registries screen MUST list registries explicitly configured or exposed by the current container environment/runtime at the time of the request.
- **FR-003**: Each registry entry MUST display enough identifying information for the user to distinguish one registry from another.
- **FR-004**: The Registries screen MUST provide a clear empty or error state when registry data cannot be shown.
- **FR-005**: The container submenu MUST include an Export container action only for stopped containers.
- **FR-006**: The export workflow MUST collect a destination directory before execution begins and generate the export filename automatically.
- **FR-007**: The export workflow MUST report success or failure to the user when the export finishes.
- **FR-008**: A failed export MUST not remove, rename, or otherwise alter the source container unexpectedly.
- **FR-009**: The build form MUST expose a user-editable checkbox that controls whether the build refreshes base image content before execution.
- **FR-010**: The build refresh checkbox MUST default to enabled whenever the build form is opened.
- **FR-011**: The selected build refresh option MUST be applied consistently to the submitted build request.
- **FR-012**: The daemon status workflow MUST use structured formatted status output as the primary source for classification when that output is available.
- **FR-013**: The daemon screen MUST distinguish at least three outcomes: running, stopped, and unknown.
- **FR-014**: When required structured status fields are missing or daemon status cannot be classified confidently, the system MUST show an unknown state and guidance instead of an incorrect definitive state.
- **FR-015**: The selected build refresh option MUST be visible to the user before the build is submitted.

### Key Entities *(include if feature involves data)*

- **Registry Entry**: Represents one container registry explicitly configured or exposed by the current container environment/runtime. Includes the identifying details needed to recognize the registry and any user-relevant availability or authentication context that can be shown.
- **Container Export Request**: Represents the user’s intent to export a specific stopped container to a chosen destination directory, along with the generated output filename and the final outcome reported back to the user.
- **Build Options**: Represents the configurable choices attached to an image build, including whether the build should refresh base image content before starting.
- **Daemon Status Result**: Represents the interpreted state of the container environment as shown to the user, including running, stopped, or unknown, derived primarily from structured formatted status output when available.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can open the Registries screen and determine whether the environment has available registries within 30 seconds.
- **SC-002**: Users can start a container export from the container submenu without leaving the TUI and complete the setup steps in under 1 minute.
- **SC-003**: In usability checks, 100% of build attempts show the refresh choice before submission with the default state set to enabled.
- **SC-004**: Daemon state is classified correctly for all representative running and stopped structured formatted outputs used during acceptance testing.
- **SC-005**: When daemon output is unrecognized, 100% of tested cases display an explicit unknown state instead of a false running or stopped result.

## Assumptions

- The feature is limited to expanding existing TUI workflows and does not introduce separate background services or non-interactive commands.
- Registry information comes only from registries explicitly configured or exposed by the current container environment/runtime; managing registry credentials is out of scope for this feature.
- Container export is available only for stopped containers, targets a user-selected destination directory, uses an application-generated export filename, and does not require additional post-export import or distribution workflows in this feature.
- The build workflow already has a concept of refreshing base image content, but that choice is not yet visible in the form and should default to enabled when first shown.
- The daemon status command can provide structured formatted output suitable for reliable classification, and the screen should continue to support the current running and stopped states while adding a safe fallback to unknown when required fields are missing.
