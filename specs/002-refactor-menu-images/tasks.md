# Tasks: Enhanced Menu Navigation and Image Management

**Feature**: 002-refactor-menu-images  
**Input**: Design documents from `/specs/002-refactor-menu-images/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

**Tests**: Contract tests are REQUIRED per constitution for command composition and destructive-action guardrails. Integration tests verify navigation flows. Unit tests cover business logic.

## Format: `- [ ] [ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story (US1, US2, US3)
- File paths use repository root as base

---

## Phase 1: Setup

**Purpose**: Minimal setup - project structure already exists

- [X] T001 Verify Go 1.21 and dependencies (Bubbletea v1.2.4, Bubbles v0.20.0) installed
- [X] T002 Review existing src/ui/ patterns (app.go, container_list.go, help.go, type_to_confirm.go)
- [X] T003 Review existing src/services/ command builder pattern

**Checkpoint**: Development environment ready

---

## Phase 2: Foundational (BLOCKING - Must Complete Before User Stories)

**Purpose**: Navigation infrastructure required by ALL user stories

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Define ViewType enum in src/ui/app.go (ContainerList, ContainerSubmenu, ContainerLogs, ContainerShell, ImageList, ImageSubmenu, ImageInspect, etc.)
- [X] T005 Add navigationStack []ViewType field to AppModel in src/ui/app.go
- [X] T006 Add selectedContainer *Container field to AppModel in src/ui/app.go
- [X] T007 Add selectedImage *Image field to AppModel in src/ui/app.go
- [X] T008 Implement pushView(ViewType) method in src/ui/app.go
- [X] T009 Implement popView() ViewType method in src/ui/app.go
- [X] T010 Update AppModel.Update() to handle navigation messages and route to active view in src/ui/app.go
- [X] T011 Update AppModel.View() to render active view based on currentView in src/ui/app.go
- [X] T012 Add BackToListMsg, BackToSubmenuMsg message types in src/ui/messages.go

**Testing**:
- [X] T013 Unit test navigation stack push/pop logic in tests/unit/navigation_test.go
- [X] T014 Verify Esc key pops from stack correctly in tests/unit/navigation_test.go

**Checkpoint**: Navigation infrastructure complete - can now implement user stories in parallel

---

## Phase 3: User Story 1 - Container Action Submenu (Priority: P1) 🎯 MVP

**Goal**: Replace Enter-to-toggle with Enter-to-submenu, provide context-sensitive container actions (start/stop/logs/shell)

**Independent Test**: Navigate container list → press Enter → verify submenu displays with context-appropriate options → select action → verify execution

### Tests for User Story 1 (Write FIRST, ensure they FAIL)

- [X] T015 [P] [US1] Contract test for container logs -f command in tests/contract/container_logs_test.go
- [X] T016 [P] [US1] Contract test for container exec -it with shell detection in tests/contract/container_exec_test.go
- [X] T017 [P] [US1] Unit test for shell detection probe sequence in tests/unit/shell_detector_test.go
- [X] T018 [US1] Integration test: Enter submenu → Select "Tail logs" → View logs → Esc → Back to submenu in tests/integration/navigation_flows_test.go
- [X] T019 [US1] Integration test: Enter submenu → "Enter container" → Shell detected → Interactive session → Exit → Back to submenu in tests/integration/navigation_flows_test.go
- [X] T020 [US1] Integration test: Shell detection with no available shells → Error message → Stay in submenu in tests/integration/navigation_flows_test.go

### Implementation for User Story 1

**Command Builders & Services**:
- [X] T021 [P] [US1] Create ContainerLogsBuilder in src/services/container_logs_builder.go (generates `container logs -f <name>`)
- [X] T022 [P] [US1] Create ContainerExecBuilder in src/services/container_exec_builder.go (generates `container exec -it <name> <shell>`)
- [X] T023 [P] [US1] Create ShellDetector in src/services/shell_detector.go (probe sequence: bash→sh→/bin/sh→/bin/bash→ash with caching)

**UI Components**:
- [X] T024 [US1] Create ContainerSubmenuModel in src/ui/container_submenu.go (arrow navigation, context-sensitive options based on container state)
- [X] T025 [US1] Implement ContainerSubmenuModel.Update() with option selection logic in src/ui/container_submenu.go
- [X] T026 [US1] Implement ContainerSubmenuModel.View() with menu rendering in src/ui/container_submenu.go
- [X] T027 [US1] Create ContainerLogsModel in src/ui/container_logs.go (viewport for log display, async streaming with tea.Cmd)
- [X] T028 [US1] Implement log streaming goroutine with buffered channel in src/ui/container_logs.go
- [X] T029 [US1] Handle log interruption (container stopped/removed) with message display in src/ui/container_logs.go
- [X] T030 [US1] Create ContainerShellModel in src/ui/container_shell.go (suspend TUI, exec shell, resume TUI wrapper)
- [X] T031 [US1] Implement shell exit detection and TUI restoration in src/ui/container_shell.go (detect shell process exit, restore TUI state, return to container submenu per FR-013)
- [X] T032 [US1] Implement shell detection integration in src/ui/container_shell.go (call ShellDetector, handle errors)

**Integration**:
- [X] T033 [US1] Modify ContainerListModel.Update() in src/ui/container_list.go - change Enter key from toggle to pushView(ContainerSubmenu)
- [X] T034 [US1] Add ContainerSubmenu routing in src/ui/app.go Update() method
- [X] T035 [US1] Add ContainerLogs routing in src/ui/app.go Update() method
- [X] T036 [US1] Add ContainerShell routing in src/ui/app.go Update() method
- [X] T037 [US1] Update main menu help text in src/ui/container_list.go - change "enter=toggle" to "enter=submenu"

**Error Handling**:
- [X] T038 [US1] Add error display for shell detection failure in src/ui/container_submenu.go
- [X] T039 [US1] Add error handling for log stream interruption in src/ui/container_logs.go

**Checkpoint**: Container submenu fully functional - can navigate, view logs, enter shell, handle errors

---

## Phase 4: User Story 2 - Image List View and Quick Actions (Priority: P2)

**Goal**: Add 'i' key for image list view, display images with NAME/TAG/DIGEST, support pull/build/prune operations

**Independent Test**: Press 'i' from main menu → verify image list displays with columns → press 'p'/'b'/'n' → verify operations execute and return to image list

### Tests for User Story 2 (Write FIRST, ensure they FAIL)

- [X] T040 [P] [US2] Contract test for container image list command in tests/contract/image_list_test.go
- [X] T041 [P] [US2] Contract test for container image prune command in tests/contract/image_prune_test.go
- [X] T042 [P] [US2] Unit test for image list parser (well-formed output) in tests/unit/image_parser_test.go
- [X] T043 [P] [US2] Unit test for image list parser (empty output) in tests/unit/image_parser_test.go
- [X] T044 [P] [US2] Unit test for image list parser (long names, truncation) in tests/unit/image_parser_test.go
- [X] T045 [US2] Integration test: Press 'i' → View images → Press 'p' → Pull → Return to image list → List refreshed in tests/integration/image_operations_test.go
- [X] T046 [US2] Integration test: Press 'i' → View images → Press 'n' → Type confirm → Prune → List refreshed in tests/integration/image_operations_test.go

### Implementation for User Story 2

**Data Models**:
- [X] T047 [P] [US2] Create Image entity in src/models/image.go (name, tag, digest fields with validation)

**Command Builders & Services**:
- [X] T048 [P] [US2] Create ImageListBuilder in src/services/image_list_builder.go (generates `container image list`)
- [X] T049 [P] [US2] Create ImagePruneBuilder in src/services/image_prune_builder.go (generates `container image prune`)
- [X] T050 [US2] Create image list parser in src/services/container_parser.go or new file (parse NAME/TAG/DIGEST columns, handle empty/malformed output)

**UI Components**:
- [X] T051 [US2] Create ImageListModel in src/ui/image_list.go (table display with 3 columns, arrow navigation)
- [X] T052 [US2] Implement ImageListModel.Update() with key handlers (p, b, n, r, Esc, Enter) in src/ui/image_list.go
- [X] T053 [US2] Implement ImageListModel.View() with table rendering and column width calculation in src/ui/image_list.go
- [X] T054 [US2] Implement async image list loading with tea.Cmd in src/ui/image_list.go
- [X] T055 [US2] Add empty state message ("No images found") in src/ui/image_list.go
- [X] T056 [US2] Implement image prune with type-to-confirm integration (confirmation word: "prune") in src/ui/image_list.go
- [X] T057 [US2] Add automatic list refresh after pull/build/prune operations in src/ui/image_list.go

**Integration**:
- [X] T058 [US2] Add 'i' key constant in src/ui/keys.go
- [X] T059 [US2] Add 'i' key handler in src/ui/container_list.go - pushView(ImageList)
- [X] T060 [US2] Add ImageList routing in src/ui/app.go Update() method
- [X] T061 [US2] Update main menu help text in src/ui/container_list.go - replace "p=pull, b=build" with "i=images"
- [X] T062 [US2] Modify ImagePullModel to return to ImageList instead of ContainerList in src/ui/image_pull.go
- [X] T063 [US2] Modify ImageBuildModel to return to ImageList instead of ContainerList in src/ui/build.go

**Error Handling**:
- [X] T064 [US2] Add error display for empty image list in src/ui/image_list.go
- [X] T065 [US2] Add error handling for "no unused images" prune result in src/ui/image_list.go
- [X] T066 [US2] Add error handling for daemon not running in src/ui/image_list.go

**Checkpoint**: Image list view fully functional - can browse images, pull/build/prune, list refreshes automatically

---

## Phase 5: User Story 3 - Image Details Submenu (Priority: P3)

**Goal**: Press Enter on image to open submenu with inspect and delete operations

**Independent Test**: Select image → press Enter → verify submenu displays → select "Inspect" → view JSON → Esc → select "Delete" → confirm → image removed

### Tests for User Story 3 (Write FIRST, ensure they FAIL)

- [X] T067 [P] [US3] Contract test for container image inspect | jq command in tests/contract/image_inspect_test.go
- [X] T068 [P] [US3] Contract test for container image rm command in tests/contract/image_delete_test.go
- [X] T069 [US3] Integration test: Select image → Enter → Select "Inspect" → View JSON → Scroll → Esc → Return to submenu in tests/integration/image_operations_test.go
- [X] T070 [US3] Integration test: Select image → Enter → Select "Delete" → Type confirm → Delete success → Return to list → List refreshed in tests/integration/image_operations_test.go
- [X] T071 [US3] Integration test: Select image in use → Delete → Error "in use by container" → Stay in submenu in tests/integration/image_operations_test.go

### Implementation for User Story 3

**Command Builders**:
- [X] T072 [P] [US3] Create ImageInspectBuilder in src/services/image_inspect_builder.go (generates `container image inspect <imageReference> | jq` where imageReference is NAME:TAG or NAME@DIGEST)
- [X] T073 [P] [US3] Create ImageDeleteBuilder in src/services/image_delete_builder.go (generates `container image rm <imageReference>` where imageReference is NAME:TAG or NAME@DIGEST)

**UI Components**:
- [X] T074 [US3] Create ImageSubmenuModel in src/ui/image_submenu.go (arrow navigation, 3 options: Inspect/Delete/Back)
- [X] T075 [US3] Implement ImageSubmenuModel.Update() with option selection logic in src/ui/image_submenu.go
- [X] T076 [US3] Implement ImageSubmenuModel.View() with menu rendering in src/ui/image_submenu.go
- [X] T077 [US3] Create ImageInspectModel in src/ui/image_inspect.go (viewport for JSON display, scrolling with arrow/pgup/pgdn/home/end keys)
- [X] T078 [US3] Implement async image inspect loading with tea.Cmd in src/ui/image_inspect.go
- [X] T079 [US3] Add scroll position indicator in src/ui/image_inspect.go
- [X] T080 [US3] Implement image delete with type-to-confirm integration (confirmation word: "delete") in src/ui/image_submenu.go

**Integration**:
- [X] T081 [US3] Add Enter key handler in ImageListModel to pushView(ImageSubmenu) in src/ui/image_list.go
- [X] T082 [US3] Add ImageSubmenu routing in src/ui/app.go Update() method
- [X] T083 [US3] Add ImageInspect routing in src/ui/app.go Update() method
- [X] T084 [US3] Implement navigation: ImageSubmenu → ImageInspect → Esc → ImageSubmenu in src/ui/image_inspect.go
- [X] T085 [US3] Implement navigation: ImageSubmenu → Delete success → ImageList (skip submenu) in src/ui/image_submenu.go

**Error Handling**:
- [X] T086 [US3] Add error display for "image in use by container(s)" in src/ui/image_submenu.go
- [X] T087 [US3] Add error handling for image not found in src/ui/image_inspect.go
- [X] T088 [US3] Add graceful degradation for jq not installed (display raw JSON) in src/ui/image_inspect.go

**Checkpoint**: Image submenu fully functional - can inspect and delete images with proper error handling

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Refinements affecting multiple user stories

- [X] T089 [P] Update README.md with new navigation pattern and image management features
- [X] T090 [P] Update docs/user-guide.md with submenu usage and image operations
- [X] T091 [P] Add help screen updates for new key bindings in src/ui/help.go
- [X] T092 Performance test: Image list with 100+ images displays in <2 seconds
- [X] T093 Performance test: Log streaming latency <100ms
- [X] T094 Code review: Ensure all new commands logged per Principle IV
- [X] T095 Code review: Verify all destructive actions (prune, delete) use type-to-confirm per Principle I
- [X] T096 Run full test suite (go test ./...)
- [X] T097 Validation: Run through all quickstart.md scenarios manually
- [X] T098 [P] Add logging for navigation state transitions (debug mode)

---

## Dependencies & Execution Order

### Phase Dependencies

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundational) ⚠️ BLOCKING
    ↓
    ├── Phase 3 (US1 - Container Submenu) 🎯 MVP
    ├── Phase 4 (US2 - Image List) (can start in parallel after foundational)
    └── Phase 5 (US3 - Image Submenu) (depends on Phase 4 completion)
    ↓
Phase 6 (Polish)
```

### Critical Path

**Minimum Viable Product (MVP)**: Complete Phase 1 → Phase 2 → Phase 3 (US1)

**User Story Dependencies**:
- **US1 (T015-T039)**: Independent after Phase 2 - delivers container submenu functionality
- **US2 (T040-T066)**: Independent after Phase 2 - delivers image list functionality
- **US3 (T067-T088)**: **Depends on US2** (needs ImageListModel for integration) - delivers image submenu

**Recommended Implementation Order**:
1. Complete Phase 1 (Setup) - T001-T003
2. Complete Phase 2 (Foundational) - T004-T014 (⚠️ BLOCKING)
3. Complete Phase 3 (US1) - T015-T039 (🎯 MVP - delivers immediate value)
4. Complete Phase 4 (US2) - T040-T066 (independent, can be done in parallel with US1 if team capacity allows)
5. Complete Phase 5 (US3) - T067-T088 (requires US2 complete)
6. Complete Phase 6 (Polish) - T089-T098

### Parallel Opportunities

**Within Phase 2 (Foundational)**:
- T004-T007 (model field additions) can be done together
- T008-T009 (methods) can be done after fields complete
- T013-T014 (tests) can be written in parallel after implementation

**Within Phase 3 (US1)**:
- T015-T017 (tests) can all be written in parallel
- T021-T023 (command builders) can be developed in parallel
- T024, T027, T030 (UI models) can be developed in parallel after command builders

**Within Phase 4 (US2)**:
- T039-T043 (tests) can all be written in parallel
- T046 (Image model) independent
- T047-T048 (command builders) can be developed in parallel
- Many UI tasks can proceed in parallel with different team members

**Within Phase 5 (US3)**:
- T066-T067 (tests) can be written in parallel
- T071-T072 (command builders) can be developed in parallel
- T073 and T076 (UI models) can be developed in parallel

**Within Phase 6 (Polish)**:
- T088-T090 (documentation) can all be done in parallel
- T091-T092 (performance tests) can run in parallel
- T093-T094 (code review) can be done in parallel

### Task Dependencies Within User Stories

**US1 Flow**:
```
T015-T020 (tests written first) → 
T021-T023 (command builders) → 
T024-T032 (UI components, can partially overlap) → 
T033-T037 (integration) → 
T038-T039 (error handling)
```

**US2 Flow**:
```
T040-T046 (tests written first) → 
T047 (Image model) → 
T048-T050 (builders & parsers) → 
T051-T057 (UI components) → 
T058-T063 (integration) → 
T064-T066 (error handling)
```

**US3 Flow** (requires US2 complete):
```
T067-T071 (tests written first) → 
T072-T073 (command builders) → 
T074-T080 (UI components) → 
T081-T085 (integration with ImageList from US2) → 
T086-T088 (error handling)
```

---

## Implementation Strategy

### MVP Delivery (Fastest Value)

**Goal**: Get container submenu working ASAP

**Minimum Path**: T001-T003 → T004-T014 → T015-T039

**Result**: Users can navigate container list → press Enter → access submenu with start/stop/logs/shell

**Estimated Effort**: ~60% of total feature (container submenu is most complex due to shell detection and log streaming)

### Incremental Delivery

**After MVP**, choose one of:

1. **Option A - Image Management Next**: Complete T040-T066 (US2) + T067-T088 (US3)
   - Delivers complete image management workflow
   - Both image stories together = coherent feature set

2. **Option B - Simple then Complex**: Complete T040-T066 (US2) first, defer T067-T088 (US3)
   - Image list + pull/build/prune delivers value
   - Image inspect/delete can be separate release

**Final Polish**: T089-T098 after all desired user stories complete

### Testing Strategy

**TDD Approach** (per constitution):
1. Write contract tests FIRST (verify command generation)
2. Write integration tests (verify navigation flows)
3. Implement until tests pass
4. Write unit tests for complex logic (shell detection, parsing)

**Test Execution Order**:
- Contract tests: Fast, run frequently during development
- Unit tests: Fast, run frequently
- Integration tests: Slower (full TUI flows), run before committing

---

## Verification Checklist

**Before marking Phase 2 complete**:
- [ ] Navigation stack push/pop works correctly
- [ ] Esc key returns to previous view
- [ ] Can navigate between existing views (container list, help, daemon control) using new infrastructure

**Before marking US1 complete**:
- [ ] Enter key on container opens submenu (not toggle)
- [ ] Submenu shows context-appropriate options (Start for stopped, Stop for running)
- [ ] Can view live container logs and Esc to return
- [ ] Shell detection works for standard containers (bash, sh variants)
- [ ] Shell detection failure shows error and stays in submenu
- [ ] All contract tests pass for container logs and exec commands

**Before marking US2 complete**:
- [ ] 'i' key navigates to image list from main menu
- [ ] Image list displays NAME, TAG, DIGEST columns correctly
- [ ] Empty image list shows appropriate message
- [ ] Pull/build operations return to image list (not container list)
- [ ] Image prune uses type-to-confirm and refreshes list
- [ ] List displays in <2 seconds for 100 images

**Before marking US3 complete**:
- [ ] Enter key on image opens image submenu
- [ ] Image inspect displays formatted JSON with scrolling
- [ ] Image delete uses type-to-confirm
- [ ] "Image in use" error handled correctly
- [ ] Successful delete returns to image list with refresh

**Before marking feature complete**:
- [ ] All 98 tasks checked off
- [ ] Full test suite passes (go test ./...)
- [ ] Manual validation of all quickstart.md scenarios completed
- [ ] Documentation updated (README, user guide, help screen)
- [ ] Constitution compliance verified (command logging, type-to-confirm on destructive actions)

---

**Total Tasks**: 98  
**Estimated Distribution**: Setup (3), Foundational (11), US1 (25), US2 (27), US3 (22), Polish (10)  
**MVP Path**: 39 tasks (Setup + Foundational + US1)  
**Full Feature**: 98 tasks
