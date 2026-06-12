# Feature Specification: Apple Container 1.0 Support with Container Machine

**Feature Branch**: `013-apple-container-1-0-support`  
**Created**: 2026-06-12  
**Status**: Draft  
**Input**: User description: "apple just released apple container 1.0 that supports also the brand new feature 'container machine'. I want to add the support for container machine on actui. This means that a new menu must be added on the main menu and a new sotto menu dedicated must be created for container machine with the related important commands (like list, edit resource machine and whatever could be relevant for machine management). Right now actui --version reply with actui version 0.4, but the release published on github is 0.1.12"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Browse and Manage Container Machines (Priority: P1)

A developer using actui wants to see all their container machines at a glance from the main menu, then take actions on a selected machine — view its details, start/stop it, or delete it.

**Why this priority**: Container machine is the headline feature of Apple Container 1.0 and the primary reason for this upgrade. Without a machine list and management submenu, the feature is entirely absent.

**Independent Test**: Navigate to the new "Machines" menu from the actui main screen. Verify a table of container machines is shown (ID, state, image, default marker). Select a machine with `enter` and verify a submenu with contextual actions appears.

**Acceptance Scenarios**:

1. **Given** the user is at the actui main container list, **When** they press the hotkey for "Machines", **Then** a machine list screen opens showing all container machines with columns: Name, State, Image, Default.
2. **Given** the machine list is open, **When** the user presses `enter` on a machine row, **Then** a machine submenu opens with available actions based on machine state (running or stopped).
3. **Given** the machine submenu is open on a running machine, **When** the user selects "Stop", **Then** a command preview is shown and the machine is stopped after confirmation.
4. **Given** the machine submenu is open on a stopped machine, **When** the user selects "Delete", **Then** a destructive-action confirmation prompt is shown before deletion proceeds.
5. **Given** the machine list is open, **When** no machines exist, **Then** an empty-state message is displayed rather than an empty or broken table.

---

### User Story 2 - View Machine Details and Logs (Priority: P2)

A developer wants to inspect detailed configuration of a container machine or follow its logs to diagnose issues, without leaving the TUI.

**Why this priority**: Inspect and logs are diagnostic workflows needed for day-to-day machine operation. They complement the management submenu and are independently deliverable.

**Independent Test**: From the machine submenu, select "Inspect" and verify a formatted JSON/detail view is shown. Select "Logs" and verify log output is displayed in the existing log viewer pattern.

**Acceptance Scenarios**:

1. **Given** the machine submenu is open, **When** the user selects "Inspect", **Then** a detail view shows machine configuration (name, image, CPUs, memory, home-mount, state).
2. **Given** the machine submenu is open, **When** the user selects "Logs", **Then** log output for that machine is displayed in a scrollable view.
3. **Given** the machine detail view is open, **When** the user presses `esc` or `q`, **Then** navigation returns to the machine submenu.

---

### User Story 3 - Edit Machine Resources (Priority: P2)

A developer wants to change the resource allocation (CPUs, memory) or home-mount mode of an existing machine without running raw CLI commands.

**Why this priority**: Resource editing (`container machine set`) is a key operational need, especially for developers who resize machines as workloads change.

**Independent Test**: From the machine submenu, select "Edit Resources". Verify an input form pre-populated with current values allows the user to update CPUs, memory, and home-mount. Submit and verify a `container machine set` command is previewed and executed.

**Acceptance Scenarios**:

1. **Given** the machine submenu is open, **When** the user selects "Edit Resources", **Then** a form displays the current cpus, memory, and home-mount values for that machine.
2. **Given** the resource edit form is open, **When** the user changes values and confirms, **Then** a `container machine set` command preview is shown before execution.
3. **Given** the edit is submitted, **When** the user returns to the machine list, **Then** a note indicates changes take effect after the next stop and restart.

---

### User Story 4 - Create a Container Machine (Priority: P3)

A developer wants to create a new container machine from an image reference directly from actui.

**Why this priority**: Creation is the onboarding flow for new machines. Useful but lower priority than management of existing machines; it can be delivered independently as a follow-on.

**Independent Test**: From the machine list screen, press the "Create" hotkey, enter an image reference (e.g. `alpine:latest`), optionally set a name, and confirm. Verify `container machine create` is executed with correct arguments.

**Acceptance Scenarios**:

1. **Given** the machine list is open, **When** the user presses the "Create" hotkey, **Then** a form is shown requesting at minimum an image reference.
2. **Given** the create form is filled, **When** the user submits, **Then** a command preview shows the `container machine create` invocation before it runs.
3. **Given** the create command completes successfully, **When** the user is returned to the machine list, **Then** the new machine appears in the list.

---

### User Story 5 - Version String Matches Published Release (Priority: P1)

A user running `actui --version` expects to see the version that matches the published GitHub release.

**Why this priority**: A mismatched version string is misleading for diagnostics and user trust. It is a one-line fix with no risk.

**Independent Test**: Build actui from source, run `actui --version`, and verify it reports `actui version 0.1.12`.

**Acceptance Scenarios**:

1. **Given** actui is installed, **When** the user runs `actui --version`, **Then** the output shows `actui version 0.1.12`.

---

### User Story 6 - Apple Container 1.0 Compatibility Review (Priority: P2)

All existing actui workflows (containers, images, builder, daemon, registries) continue to function correctly against Apple Container 1.0 without broken output parsing.

**Why this priority**: Apple Container 1.0 introduced breaking changes to structured output shapes for container, image, network, and volume `ls` and `inspect`. Any parser in actui relying on the old JSON/TOML shapes must be validated and updated.

**Independent Test**: Run the full actui TUI against a system with Apple Container 1.0 installed. Verify each existing screen renders data correctly: ContainerList, ImageList, Registries, DaemonControl. Specific focus on JSON parsing in services that read `ls` and `inspect` output.

**Acceptance Scenarios**:

1. **Given** Apple Container 1.0 is installed, **When** the user opens actui, **Then** the ContainerList screen correctly lists containers with all columns populated.
2. **Given** Apple Container 1.0 is installed, **When** the user opens the ImageList screen, **Then** images are listed without parse errors.
3. **Given** Apple Container 1.0 is installed, **When** the user opens DaemonControl, **Then** system status and version information are displayed correctly.
4. **Given** Apple Container 1.0 is installed, **When** the user navigates to Registries, **Then** registry logins are listed correctly.

---

### Edge Cases

- What happens when the `container machine` CLI subcommand is absent (pre-1.0 installation)? The Machines menu entry should be hidden or show a clear "requires Apple Container 1.0+" message.
- What happens when the machine list returns an empty result? An empty-state message is shown.
- What happens when a machine is in an unknown or transitional state? The submenu should show only safe, applicable actions.
- What happens if `container machine set` changes fail (e.g. invalid memory value)? The error is surfaced to the user in the existing error display pattern.
- How does actui handle the removed `container system property get/set` subcommands from Apple Container 1.0? If any actui code calls these, they must be removed or replaced.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: actui MUST display a "Machines" entry in the main navigation, accessible via a keyboard shortcut from the ContainerList screen.
- **FR-002**: The machine list screen MUST display container machines in a table with at minimum these columns: Name, State, Image, Default marker.
- **FR-003**: Selecting a machine row MUST open a machine submenu with contextual actions based on machine state.
- **FR-004**: The machine submenu MUST offer: Inspect, Logs, Stop (if running), Start (if stopped), Edit Resources, Set as Default, Delete.
- **FR-005**: Delete action MUST display a destructive-action confirmation before executing `container machine delete`.
- **FR-006**: Stop action MUST display a command preview before executing `container machine stop`.
- **FR-007**: Edit Resources MUST allow editing cpus, memory, and home-mount and preview the `container machine set` command before execution.
- **FR-008**: The machine list screen MUST show an empty-state message when no machines exist.
- **FR-009**: All machine management screens MUST follow the existing navigation conventions: `esc`/`q` returns to the previous screen; `BackToListMsg` returns to root.
- **FR-010**: The version string in `cmd/actui/main.go` MUST be updated from `0.4` to `0.1.12`.
- **FR-011**: All existing service parsers that consume `container ls` or `container inspect` JSON/TOML output MUST be validated against Apple Container 1.0 output shapes and updated if necessary.
- **FR-012**: Any actui code calling removed Apple Container 1.0 subcommands (`container system property get/set`) MUST be removed or replaced.
- **FR-013**: New machine screens MUST be registered in `src/ui/messages.go` (screen IDs) and `src/ui/app.go` (AppModel switch logic) per project conventions.
- **FR-014**: New service builders for machine commands MUST reside in `src/services/` and follow the existing builder naming convention.
- **FR-015**: The machine list and submenu MUST be navigable with the keyboard only; no mouse input required.

### Key Entities

- **Container Machine**: A persistent Linux VM managed by Apple Container. Key attributes: name (ID), state (running/stopped), base image reference, resource config (cpus, memory, home-mount), default flag.
- **Machine Action**: A command that can be performed on a machine (inspect, logs, stop, start, set, set-default, delete). Availability depends on machine state.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can navigate from the actui main screen to the machine list and back in under 5 keystrokes.
- **SC-002**: All container machine management actions (list, inspect, stop, start, delete, set resources) are reachable from within actui without exiting to the CLI.
- **SC-003**: `actui --version` outputs `actui version 0.1.12`.
- **SC-004**: All existing actui screens render correct data when running against Apple Container 1.0 with no visible parse errors or empty fields where data is expected.
- **SC-005**: Destructive actions (delete, stop) require at least one confirmation step before executing the underlying CLI command.
- **SC-006**: The machine submenu shows only actions relevant to the current machine state (no "Start" on a running machine, no "Stop" on a stopped machine).

## Assumptions

- Apple Container 1.0 is installed on the target machine; actui does not manage the installation of the CLI.
- The `container machine` alias `m` is available but actui will use the full `container machine` command form for clarity and consistency with other service builders.
- Container machine create is lower priority and may be deferred if time-boxed; the spec includes it as P3.
- The home-mount edit option exposes `rw`, `ro`, and `none` as selectable values rather than a free-text field.
- Machine state detection relies on the `container machine list` output; no separate status polling is introduced.
- Breaking output shape changes in Apple Container 1.0 affect primarily the JSON `ls` and `inspect` responses; TOML config changes (`system property` removal) do not affect the existing actui UI flows since actui did not surface system property commands.
- The existing `DaemonControl` screen, which calls `container system status`, continues to work in Apple Container 1.0 as `container system status` is unchanged.
