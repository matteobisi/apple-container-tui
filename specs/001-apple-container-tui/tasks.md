# Tasks: Apple Container TUI

**Input**: Design documents from `/specs/001-apple-container-tui/`
**Prerequisites**: plan.md (✓), spec.md (✓), research.md (✓), data-model.md (✓), contracts/ (✓)

**Tests**: Tests are REQUIRED for command composition, argument validation, and destructive-action guardrails per the constitution. Additional tests are optional unless explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `- [ ] [ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `src/`, `tests/` at repository root (per plan.md)
- Paths shown below use single project structure

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create Go module and initialize project with go mod init
- [X] T002 [P] Add Bubbletea v1.2.4, Lipgloss v1.0.0, Bubbles v0.20.0 dependencies to go.mod
- [X] T003 [P] Add Cobra and Viper dependencies for CLI and config management to go.mod
- [X] T004 Create project directory structure: src/models/, src/services/, src/ui/, cmd/apple-tui/
- [X] T005 [P] Create test directory structure: tests/contract/, tests/integration/, tests/unit/
- [X] T006 [P] Create config/ directory and default.toml template file
- [X] T007 [P] Create docs/ directory and user-guide.md stub
- [X] T008 Create main entry point in cmd/apple-tui/main.go with version flag support
- [X] T009 [P] Setup .gitignore for Go projects (binaries, vendor/, IDE files)
- [X] T010 [P] Create README.md skeleton with project name and placeholder sections

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T011 Create Command struct in src/models/command.go with Executable and Args fields
- [X] T012 Create CommandBuilder interface in src/services/command_builder.go with Validate() and Build() methods
- [X] T013 Create CommandExecutor interface in src/services/executor.go with Execute(cmd Command) method
- [X] T014 Implement DryRunExecutor in src/services/dry_run_executor.go that echoes commands without executing
- [X] T015 Implement RealExecutor in src/services/real_executor.go using os/exec.Command for safe execution
- [X] T016 Create Result struct in src/models/result.go with ExitCode, Stdout, Stderr, Duration, Status fields
- [X] T017 [P] Create UserConfig model in src/models/config.go with TOML tags for all config fields
- [X] T018 [P] Create ConfigManager service in src/services/config_manager.go for dual-path config loading (~/.config and ~/Library/Application Support)
- [X] T019 [P] Create DaemonStatus model in src/models/daemon.go with Running, Version, LastChecked fields
- [X] T020 Create startup CLI check function in src/services/cli_checker.go that verifies 'container' binary exists in PATH (using command -v container or equivalent) and runs 'container system version', exiting with clear error message if either check fails
- [X] T021 [P] Setup Bubbletea app initialization in src/ui/app.go with Update(), View(), and Init() methods
- [X] T022 [P] Create keyboard handler registry in src/ui/keys.go for common key bindings
- [X] T023 Contract test for DryRunExecutor in tests/contract/dry_run_test.go verifying no actual execution
- [X] T024 Contract test for RealExecutor command building in tests/contract/executor_test.go verifying os/exec.Command usage

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Manage Containers From a Menu (Priority: P1) 🎯 MVP

**Goal**: List containers and start or stop a selected container from a TUI menu

**Independent Test**: Can list containers, select one, and start or stop it with visible status feedback

### Tests for User Story 1 ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T025 [P] [US1] Contract test for ListContainersBuilder in tests/contract/list_containers_test.go verifying 'container list --all' command
- [X] T026 [P] [US1] Contract test for StartContainerBuilder in tests/contract/start_container_test.go verifying containerID validation
- [X] T027 [P] [US1] Contract test for StopContainerBuilder in tests/contract/stop_container_test.go verifying containerID validation
- [X] T028 [P] [US1] Unit test for Container model validation in tests/unit/container_test.go
- [X] T029 [P] [US1] Integration test for list/start/stop workflow in tests/integration/container_lifecycle_test.go (requires Apple Container CLI)

### Implementation for User Story 1

- [X] T030 [P] [US1] Create Container model in src/models/container.go with ID, Name, Image, Status, Created, Ports fields
- [X] T031 [P] [US1] Create PortMapping model in src/models/port_mapping.go with HostPort, ContainerPort, Protocol fields
- [X] T032 [US1] Create ListContainersBuilder in src/services/list_containers_builder.go implementing CommandBuilder interface
- [X] T033 [US1] Create StartContainerBuilder in src/services/start_container_builder.go with containerID validation
- [X] T034 [US1] Create StopContainerBuilder in src/services/stop_container_builder.go with containerID validation
- [X] T035 [US1] Create container list parser in src/services/container_parser.go to parse 'container list --all' table output
- [X] T036 [US1] Create ContainerListScreen in src/ui/container_list.go with Bubbletea Model, Update, and View methods
- [X] T037 [US1] Add keyboard navigation (arrow keys, Enter) to ContainerListScreen for selecting containers
- [X] T038 [US1] Create CommandPreviewModal in src/ui/command_preview.go showing command string and confirmation prompt
- [X] T039 [US1] Integrate ContainerListScreen into main app initialization in cmd/apple-tui/main.go
- [X] T040 [US1] Implement start action handler in src/ui/container_list.go that builds command, shows preview, executes on confirm
- [X] T041 [US1] Implement stop action handler in src/ui/container_list.go that builds command, shows preview, executes on confirm
- [X] T042 [US1] Create result display component in src/ui/result_display.go showing success/error with stdout/stderr
- [X] T043 [US1] Add manual refresh keybinding ('r' key) to ContainerListScreen that re-fetches container list

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Pull and Build Images (Priority: P2)

**Goal**: Pull images and build from a selected container file without leaving the TUI

**Independent Test**: Can pull a known image and build from a chosen file, with progress and final status shown

### Tests for User Story 2 ⚠️

- [X] T044 [P] [US2] Contract test for PullImageBuilder in tests/contract/pull_image_test.go verifying reference validation
- [X] T045 [P] [US2] Contract test for BuildImageBuilder in tests/contract/build_image_test.go verifying file path and tag validation
- [X] T046 [P] [US2] Unit test for ImageReference model validation in tests/unit/image_reference_test.go
- [X] T047 [P] [US2] Unit test for BuildSource model validation and file existence check in tests/unit/build_source_test.go
- [X] T048 [P] [US2] Integration test for pull and build workflow in tests/integration/image_operations_test.go

### Implementation for User Story 2

- [X] T049 [P] [US2] Create ImageReference model in src/models/image_reference.go with Registry, Repository, Tag, Digest fields
- [X] T050 [P] [US2] Create BuildSource model in src/models/build_source.go with FilePath, FileType, WorkingDirectory, Exists fields
- [X] T051 [US2] Create PullImageBuilder in src/services/pull_image_builder.go with reference format validation
- [X] T052 [US2] Create BuildImageBuilder in src/services/build_image_builder.go with file existence and tag validation
- [X] T053 [US2] Create file type auto-detect function in src/services/build_file_detector.go (Containerfile then Dockerfile)
- [X] T054 [US2] Create ImagePullScreen in src/ui/image_pull.go with text input for image reference
- [X] T055 [US2] Create FilePickerScreen in src/ui/file_picker.go using Bubbles filepicker component
- [X] T056 [US2] Create BuildScreen in src/ui/build.go with tag input and progress display
- [X] T057 [US2] Add navigation to ImagePullScreen from main menu ('p' key)
- [X] T058 [US2] Add navigation to FilePickerScreen from main menu ('b' key)
- [X] T059 [US2] Implement pull action handler in src/ui/image_pull.go with command preview and progress display
- [X] T060 [US2] Implement build action handler in src/ui/build.go with streaming output to Bubbles viewport
- [X] T061 [US2] Create progress bar component in src/ui/progress.go for pull/build operations using Bubbles progress
- [X] T062 [US2] Integrate ImagePullScreen and BuildScreen into main app router in src/ui/app.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Safe Destructive and Daemon Actions (Priority: P3)

**Goal**: Delete stopped containers and start or stop the daemon with safety checks

**Independent Test**: Can delete a stopped container and start or stop the daemon with confirmation prompts and clear results

### Tests for User Story 3 ⚠️

- [X] T063 [P] [US3] Contract test for DeleteContainerBuilder in tests/contract/delete_container_test.go verifying no --force flag
- [X] T064 [P] [US3] Contract test for StartDaemonBuilder in tests/contract/start_daemon_test.go
- [X] T065 [P] [US3] Contract test for StopDaemonBuilder in tests/contract/stop_daemon_test.go
- [X] T066 [P] [US3] Unit test for type-to-confirm validation in tests/unit/confirmation_test.go
- [X] T067 [P] [US3] Unit test for destructive action metadata in tests/unit/destructive_actions_test.go
- [X] T068 [P] [US3] Integration test for delete workflow in tests/integration/delete_container_test.go

### Implementation for User Story 3

- [X] T069 [US3] Create DeleteContainerBuilder in src/services/delete_container_builder.go with stopped-only validation
- [X] T070 [US3] Create StartDaemonBuilder in src/services/start_daemon_builder.go (no arguments)
- [X] T071 [US3] Create StopDaemonBuilder in src/services/stop_daemon_builder.go (no arguments)
- [X] T072 [US3] Create CheckDaemonStatusBuilder in src/services/check_daemon_builder.go for 'container system status'
- [X] T073 [US3] Create daemon status parser in src/services/daemon_parser.go to parse running/stopped state
- [X] T074 [US3] Create TypeToConfirmModal in src/ui/type_to_confirm.go with text input requiring exact match
- [X] T075 [US3] Create YesNoConfirmModal in src/ui/yes_no_confirm.go for daemon start/stop confirmations
- [X] T076 [US3] Create DaemonControlScreen in src/ui/daemon_control.go with start/stop options
- [X] T077 [US3] Add delete action handler to ContainerListScreen in src/ui/container_list.go with type-to-confirm modal
- [X] T078 [US3] Add daemon control navigation from main menu ('m' key for management)
- [X] T079 [US3] Implement daemon start handler in src/ui/daemon_control.go with yes/no confirmation
- [X] T080 [US3] Implement daemon stop handler in src/ui/daemon_control.go with warning and yes/no confirmation
- [X] T081 [US3] Add visual indicators (⚠️ red) for destructive actions in all UI components
- [X] T082 [US3] Create daemon status check on app startup in cmd/apple-tui/main.go showing guidance if not running

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T083 [P] Create HelpScreen in src/ui/help.go showing keyboard shortcuts, config paths, and log paths
- [X] T084 [P] Add help screen navigation ('?' key) from any screen
- [X] T085 [P] Implement Lipgloss theme support in src/ui/theme.go for light/dark/auto modes
- [X] T086 [P] Create command log writer in src/services/log_writer.go for JSON lines format to ~/Library/Application Support/apple-tui/command.log
- [X] T087 [P] Implement log rotation in src/services/log_writer.go based on UserConfig.LogRetentionDays
- [X] T088 [P] Add window resize handler in src/ui/app.go using Bubbletea tea.WindowSizeMsg
- [X] T089 [P] Create error formatter in src/services/error_formatter.go with common pattern matching
- [X] T090 [P] Add loading spinner component in src/ui/spinner.go using Bubbles spinner for long operations
- [X] T091 [P] Implement status bar in src/ui/status_bar.go showing current screen and live command preview
- [X] T092 [P] Add version display in HelpScreen and --version flag output
- [X] T093 Code cleanup: Run gofmt and golangci-lint across all source files
- [X] T094 Code cleanup: Add godoc comments to all exported types and functions
- [X] T095 Documentation: Complete user-guide.md in docs/ with screenshots (ASCII art) of all workflows
- [X] T096 Documentation: Complete README.md with installation instructions, quick start, and usage examples (builds on T010 skeleton)
- [X] T097 [P] Performance: Add container list caching to avoid redundant CLI calls
- [X] T098 [P] Security: Validate all user inputs to prevent command injection (verify os/exec.Command usage)
- [X] T099 Run validation checklist from quickstart.md: all FRs, constitution principles, user stories
- [X] T100 Run full test suite: go test ./... and verify >70% coverage

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Independent of US1 (can build/pull without container mgmt)
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Uses Container model from US1 but can be implemented independently

### Within Each User Story

- Tests (required by constitution) MUST be written and FAIL before implementation
- Models before services (builders depend on model validation)
- Services before UI (UI components use service builders)
- Core UI before action handlers
- Story complete before moving to next priority

### Parallel Opportunities

- **Setup phase**: T002, T003, T005, T006, T007, T009, T010 can all run in parallel
- **Foundational phase**: T017+T018 (config), T019 (daemon), T023+T024 (tests) can run in parallel
- **Within each user story**: All tests marked [P] can run in parallel, all models marked [P] can run in parallel
- **Across user stories**: Once Foundational completes, US1, US2, US3 can ALL be worked on in parallel by different developers

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task T025: "Contract test for ListContainersBuilder in tests/contract/list_containers_test.go"
Task T026: "Contract test for StartContainerBuilder in tests/contract/start_container_test.go"
Task T027: "Contract test for StopContainerBuilder in tests/contract/stop_container_test.go"
Task T028: "Unit test for Container model validation in tests/unit/container_test.go"
# All can be written in parallel

# Launch all models for User Story 1 together:
Task T030: "Create Container model in src/models/container.go"
Task T031: "Create PortMapping model in src/models/port_mapping.go"
# Both can be implemented in parallel
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T010)
2. Complete Phase 2: Foundational (T011-T024) - CRITICAL - blocks all stories
3. Complete Phase 3: User Story 1 (T025-T043)
4. **STOP and VALIDATE**: Test User Story 1 independently using quickstart.md workflows
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add Polish (Phase 6) → Final release
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (T025-T043)
   - Developer B: User Story 2 (T044-T062)
   - Developer C: User Story 3 (T063-T082)
3. Stories complete and integrate independently
4. Team reconvenes for Polish phase

---

## Task Count Summary

- **Phase 1 (Setup)**: 10 tasks
- **Phase 2 (Foundational)**: 14 tasks (BLOCKING)
- **Phase 3 (User Story 1 - P1)**: 19 tasks (5 tests + 14 implementation)
- **Phase 4 (User Story 2 - P2)**: 19 tasks (5 tests + 14 implementation)
- **Phase 5 (User Story 3 - P3)**: 20 tasks (6 tests + 14 implementation)
- **Phase 6 (Polish)**: 18 tasks

**Total**: 100 tasks

**Parallel opportunities**: 43 tasks marked [P] can run in parallel within their phase

**Independent test criteria per story**:
- US1: Can list, select, start, and stop containers with command preview
- US2: Can pull an image and build from Containerfile/Dockerfile with progress
- US3: Can delete stopped container and manage daemon with type-to-confirm safety

**Suggested MVP scope**: Phases 1-3 only (24 foundational + 19 US1 = 43 tasks)

---

## Notes

- [P] tasks = different files, no dependencies within phase
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- All tests must fail before implementation (TDD per constitution where applicable)
- Commit after each task or logical group of parallel tasks
- Stop at any checkpoint to validate story independently
- Constitution-required tests (command contracts, destructive safeguards) are mandatory
- Integration tests require Apple Container CLI installed and running
