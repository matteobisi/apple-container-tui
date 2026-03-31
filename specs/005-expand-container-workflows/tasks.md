# Tasks: Expanded Container Workflows

**Input**: Design documents from `/specs/005-expand-container-workflows/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/command-workflows.md, quickstart.md

**Tests**: Automated test coverage is required for command builders, JSON parsers, and UI routing/state changes called out in the plan and contracts.

**Organization**: Tasks are grouped by user story so each story can be implemented, verified, and demonstrated independently.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Prepare shared screen routing and navigation surfaces used by the new workflows.

- [X] T001 Add new active screen constants and screen-change payloads in src/ui/messages.go
- [X] T002 Register registries and container export screens in src/ui/app.go
- [X] T003 [P] Update shared navigation/help copy for new image and container actions in src/ui/help.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Add shared domain and state plumbing that the feature stories build on.

**⚠️ CRITICAL**: No user story work should start before this phase is complete.

- [X] T004 Create the registry domain model in src/models/registry_login.go
- [X] T005 [P] Extend shared daemon status state for running/stopped/unknown classification in src/models/daemon.go
- [X] T006 [P] Add generated export reference helpers in src/models/image_reference.go

**Checkpoint**: Shared routing and model foundations are ready for story implementation.

---

## Phase 3: User Story 1 - Browse Registries in the TUI (Priority: P1) 🎯 MVP

**Goal**: Let users open a dedicated Registries screen and inspect runtime-managed registry entries without leaving the TUI.

**Independent Test**: From the Images view, open Registries, confirm hostname/username rows load from runtime data, and verify empty/error/refresh paths remain navigable.

### Tests for User Story 1

- [X] T007 [P] [US1] Add registry list builder and parser coverage in src/services/services_test.go
- [X] T008 [P] [US1] Add Registries screen navigation and state coverage in src/ui/ui_additional_test.go

### Implementation for User Story 1

- [X] T009 [P] [US1] Implement the registry list command builder in src/services/registry_list_builder.go
- [X] T010 [P] [US1] Implement the registry list JSON parser in src/services/registry_parser.go
- [X] T011 [US1] Add the Registries entry point to image management in src/ui/image_list.go
- [X] T012 [US1] Implement registry loading, refresh, empty, and error states in src/ui/registries.go
- [X] T013 [US1] Wire Registries screen transitions and refresh behavior in src/ui/app.go and src/ui/messages.go

**Checkpoint**: User Story 1 is fully functional and independently testable.

---

## Phase 4: User Story 2 - Export a Container from Its Submenu (Priority: P2)

**Goal**: Let users export stopped containers from the existing container submenu using a previewed export/save/cleanup command sequence.

**Independent Test**: Open a stopped container submenu, choose Export, select a destination directory, preview the generated commands, execute the workflow, and confirm the OCI tar is created without altering the source container.

### Tests for User Story 2

- [X] T014 [P] [US2] Add export command-sequence and stopped-container guard coverage in src/services/services_test.go
- [X] T015 [P] [US2] Add container export submenu and screen-flow coverage in src/ui/ui_additional_test.go

### Implementation for User Story 2

- [X] T016 [P] [US2] Implement the container export command builder in src/services/export_container_builder.go
- [X] T017 [P] [US2] Implement the image save command builder in src/services/image_save_builder.go
- [X] T018 [US2] Implement export planning, sequential execution, and cleanup warning handling in src/services/export_workflow.go
- [X] T019 [US2] Add the stopped-only Export action to the container submenu in src/ui/container_submenu.go
- [X] T020 [US2] Implement destination-directory input, generated archive naming, and command preview in src/ui/container_export.go
- [X] T021 [US2] Wire container export screen transitions and result handling in src/ui/app.go and src/ui/messages.go

**Checkpoint**: User Story 2 is fully functional and independently testable.

---

## Phase 5: User Story 3 - Choose Freshness During Image Build (Priority: P3)

**Goal**: Make the build `--pull` behavior visible and user-controlled in the existing build form.

**Independent Test**: Open the build form, verify the pull checkbox defaults to enabled, toggle it, and confirm the previewed build command includes or omits `--pull` accordingly.

### Tests for User Story 3

- [X] T022 [P] [US3] Add build pull-flag command coverage in src/services/services_test.go
- [X] T023 [P] [US3] Add build form checkbox state coverage in src/ui/ui_additional_test.go

### Implementation for User Story 3

- [X] T024 [US3] Extend build command generation with PullLatest support in src/services/build_image_builder.go
- [X] T025 [US3] Add a default-enabled pull checkbox and preview state handling in src/ui/build.go
- [X] T026 [US3] Update build workflow help text for the pull toggle in src/ui/help.go

**Checkpoint**: User Story 3 is fully functional and independently testable.

---

## Phase 6: User Story 4 - Trust Daemon Status Feedback (Priority: P4)

**Goal**: Classify daemon state from structured JSON output and show an explicit unknown state when status cannot be trusted.

**Independent Test**: Refresh daemon control against representative JSON payloads and confirm the UI shows running, stopped, or unknown exactly as the structured status supports.

### Tests for User Story 4

- [X] T027 [P] [US4] Add structured daemon status builder and parser coverage in src/services/services_test.go
- [X] T028 [P] [US4] Add daemon unknown-state UI coverage in src/ui/ui_additional_test.go

### Implementation for User Story 4

- [X] T029 [US4] Switch daemon status command generation to JSON output in src/services/check_daemon_builder.go
- [X] T030 [US4] Implement structured daemon status parsing with unknown fallback in src/services/daemon_parser.go
- [X] T031 [US4] Render running, stopped, and unknown daemon states in src/ui/daemon_control.go

**Checkpoint**: User Story 4 is fully functional and independently testable.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final verification and documentation updates that span multiple stories.

- [X] T032 [P] Document registries, container export, build pull behavior, and daemon status changes in docs/user-guide.md and README.md
- [X] T033 Verify the implementation against the scenarios in specs/005-expand-container-workflows/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies; start immediately.
- **Foundational (Phase 2)**: Depends on Phase 1 and blocks all story work.
- **User Story 1 (Phase 3)**: Depends on Phase 2; recommended MVP slice.
- **User Story 2 (Phase 4)**: Depends on Phase 2; independent of US1 once shared routing/models are in place.
- **User Story 3 (Phase 5)**: Depends on Phase 2; independent of US1 and US2.
- **User Story 4 (Phase 6)**: Depends on Phase 2; independent of US1-US3.
- **Polish (Phase 7)**: Depends on completion of the stories you plan to ship.

### User Story Dependencies

- **US1**: No dependency on other user stories.
- **US2**: No dependency on other user stories.
- **US3**: No dependency on other user stories.
- **US4**: No dependency on other user stories.

### Within Each User Story

- Write the listed test updates before or alongside implementation and confirm they fail for the intended gap.
- Implement builders/parsers before the UI that consumes them.
- Complete the story’s screen wiring before moving to the next priority.

### Parallel Opportunities

- T003 can run in parallel with T001-T002 after screen naming is agreed.
- T005 and T006 can run in parallel in the foundational phase.
- For US1, T007-T010 can proceed in parallel across service and UI files.
- For US2, T014-T017 can proceed in parallel before converging on T018-T021.
- For US3, T022-T023 can proceed in parallel, then T024 and T025 can be split across service and UI files.
- For US4, T027-T028 can proceed in parallel, then T029-T030 can be split before T031.

---

## Parallel Example: User Story 1

```bash
# Service-side work for Registries
T007 Add registry list builder and parser coverage in src/services/services_test.go
T009 Implement the registry list command builder in src/services/registry_list_builder.go
T010 Implement the registry list JSON parser in src/services/registry_parser.go

# UI-side work for Registries
T008 Add Registries screen navigation and state coverage in src/ui/ui_additional_test.go
T011 Add the Registries entry point to image management in src/ui/image_list.go
T012 Implement registry loading, refresh, empty, and error states in src/ui/registries.go
```

## Parallel Example: User Story 2

```bash
# Service-side export workflow tasks
T014 Add export command-sequence and stopped-container guard coverage in src/services/services_test.go
T016 Implement the container export command builder in src/services/export_container_builder.go
T017 Implement the image save command builder in src/services/image_save_builder.go

# UI-side export workflow tasks
T015 Add container export submenu and screen-flow coverage in src/ui/ui_additional_test.go
T019 Add the stopped-only Export action to the container submenu in src/ui/container_submenu.go
T020 Implement destination-directory input, generated archive naming, and command preview in src/ui/container_export.go
```

## Parallel Example: User Story 3

```bash
# Service-side build pull work
T022 Add build pull-flag command coverage in src/services/services_test.go
T024 Extend build command generation with PullLatest support in src/services/build_image_builder.go

# UI-side build pull work
T023 Add build form checkbox state coverage in src/ui/ui_additional_test.go
T025 Add a default-enabled pull checkbox and preview state handling in src/ui/build.go
```

## Parallel Example: User Story 4

```bash
# Service-side daemon status work
T027 Add structured daemon status builder and parser coverage in src/services/services_test.go
T029 Switch daemon status command generation to JSON output in src/services/check_daemon_builder.go
T030 Implement structured daemon status parsing with unknown fallback in src/services/daemon_parser.go

# UI-side daemon status work
T028 Add daemon unknown-state UI coverage in src/ui/ui_additional_test.go
T031 Render running, stopped, and unknown daemon states in src/ui/daemon_control.go
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup.
2. Complete Phase 2: Foundational.
3. Complete Phase 3: User Story 1.
4. Validate Registries independently before taking on additional stories.

### Incremental Delivery

1. Finish Setup + Foundational once.
2. Deliver US1 as the MVP.
3. Add US2 for stopped-container export.
4. Add US3 for build pull visibility.
5. Add US4 for daemon status reliability.
6. Finish with Phase 7 documentation and verification.

### Parallel Team Strategy

1. One developer completes Phases 1-2.
2. After Phase 2, separate developers can own US1-US4 in parallel.
3. Reserve one final pass for Phase 7 documentation and quickstart validation.

---

## Notes

- [P] tasks touch different files and are safe to parallelize.
- Each story remains independently demonstrable after Phase 2 is complete.
- Keep command preview fidelity intact for every new workflow.
- Use `go test ./...` as the final automated verification pass before closing T033.