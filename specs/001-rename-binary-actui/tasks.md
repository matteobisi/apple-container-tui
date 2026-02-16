---
description: "Implementation tasks for binary rename from apple-tui to actui"
---

# Tasks: Rename Binary from apple-tui to actui

**Input**: Design documents from `/specs/001-rename-binary-actui/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/rename-inventory.md

**Tests**: Test tasks are NOT included as tests were not explicitly requested in the specification. The existing test suite will be updated to reflect the new binary name and verify functionality.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `src/`, `tests/` at repository root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

**Status**: ✅ No setup required - project structure already exists

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**Status**: ✅ No foundational changes required - this is a rename operation with no architectural dependencies

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Build Produces Correctly Named Binary (Priority: P1) 🎯 MVP

**Goal**: Rename the binary from "apple-tui" to "actui" so that the build process produces the correctly named executable

**Independent Test**: Run `go build -o actui ./cmd/actui` and verify the output binary is named "actui"

### Implementation for User Story 1

- [X] T001 [US1] Rename directory from cmd/apple-tui/ to cmd/actui/ using git mv
- [X] T002 [US1] Update .gitignore lines 23-24 to replace "apple-tui" with "actui"
- [X] T003 [US1] Verify build succeeds: run `go build -o actui ./cmd/actui`
- [X] T004 [US1] Verify binary name: confirm output file is named "actui" not "apple-tui"

**Checkpoint**: At this point, User Story 1 should be fully functional - build produces "actui" binary

---

## Phase 4: User Story 2 - Documentation Reflects Correct Binary Name (Priority: P2)

**Goal**: Update all user-facing documentation to reference "actui" instead of "apple-tui"

**Independent Test**: Search documentation files for "apple-tui" and verify all references changed to "actui"

### Implementation for User Story 2

- [X] T005 [P] [US2] Update README.md build commands: change "go build -o apple-tui ./cmd/apple-tui" to "go build -o actui ./cmd/actui"
- [X] T006 [P] [US2] Update README.md execution examples: change "./apple-tui" to "./actui" (lines 81, 87)
- [X] T007 [P] [US2] Update README.md config paths: change "~/.config/apple-tui/" to "~/.config/actui/" (lines 101, 106)
- [X] T008 [P] [US2] Update README.md config paths: change "~/Library/Application Support/apple-tui/" to "~/Library/Application Support/actui/" (lines 102, 106, 120)
- [X] T009 [P] [US2] Update docs/user-guide.md build commands: change "go build -o apple-tui ./cmd/apple-tui" to "go build -o actui ./cmd/actui" (line 28)
- [X] T010 [P] [US2] Update docs/user-guide.md execution examples: change "./apple-tui" references to "./actui" (lines 34, 40, 54)
- [X] T011 [P] [US2] Update docs/user-guide.md config paths: change "~/.config/apple-tui/" to "~/.config/actui/" (line 175)
- [X] T012 [P] [US2] Update docs/user-guide.md config paths: change "~/Library/Application Support/apple-tui/" to "~/Library/Application Support/actui/" (lines 176, 180, 194)
- [X] T013 [US2] Verify documentation consistency: run `grep -r "apple-tui" README.md docs/` and confirm zero matches

**Checkpoint**: At this point, User Stories 1 AND 2 should both work - build produces "actui" and docs reference "actui"

---

## Phase 5: User Story 3 - Clean Removal of Old Binary Name References (Priority: P3)

**Goal**: Remove all code references to "apple-tui" to ensure consistency across the entire codebase

**Independent Test**: Run comprehensive search `grep -r "apple-tui" src/ .gitignore` and verify zero matches in active files

### Implementation for User Story 3

- [X] T014 [P] [US3] Update src/services/config_manager.go line 25: change "apple-tui" to "actui" in .config path
- [X] T015 [P] [US3] Update src/services/config_manager.go line 26: change "apple-tui" to "actui" in Application Support path
- [X] T016 [P] [US3] Update src/services/config_manager.go line 28: change "apple-tui" to "actui" in writePath
- [X] T017 [P] [US3] Update src/services/log_writer.go line 38: change "apple-tui" to "actui" in log path
- [X] T018 [P] [US3] Update src/services/services_test.go line 264: change "apple-tui" to "actui" in test config path
- [X] T019 [P] [US3] Update src/services/services_test.go line 305: change "apple-tui" to "actui" in test log path check
- [X] T020 [P] [US3] Update src/services/services_test.go line 317: change "apple-tui" to "actui" in test log path
- [X] T021 [US3] Verify code consistency: run `grep -r "apple-tui" src/ .gitignore` and confirm zero matches
- [X] T022 [US3] Verify historical specs preserved: confirm specs/001-apple-container-tui/ still references old name (should not be changed)

**Checkpoint**: All user stories should now be independently functional - no "apple-tui" references in active code

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and quality checks

- [X] T023 [P] Run full test suite: execute `go test ./...` and verify all tests pass
- [X] T024 [P] Verify config creation: run `./actui` and confirm config created at `~/Library/Application Support/actui/config`
- [X] T025 [P] Verify log creation: run `./actui` and confirm logs created at `~/Library/Application Support/actui/command.log`
- [X] T026 Clean up old binary: remove any `apple-tui` binary from repository root if present
- [X] T027 Run comprehensive verification: execute verification checklist from quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: ✅ Complete (no work required)
- **Foundational (Phase 2)**: ✅ Complete (no work required)
- **User Stories (Phase 3-5)**: Can proceed in priority order or in parallel
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start immediately - No dependencies
- **User Story 2 (P2)**: Can start immediately - Independent of US1 (different files)
- **User Story 3 (P3)**: Can start immediately - Independent of US1/US2 (different files)

**All user stories are independent and can be worked in parallel!**

### Within Each User Story

#### User Story 1 (Sequential)
1. T001: Directory rename (must be first)
2. T002: .gitignore update (can be parallel with T001)
3. T003: Build verification (depends on T001)
4. T004: Binary verification (depends on T003)

#### User Story 2 (Parallel)
- T005-T012: All documentation updates can run in parallel
- T013: Documentation verification (depends on T005-T012)

#### User Story 3 (Parallel)
- T014-T020: All code updates can run in parallel
- T021-T022: Verification (depends on T014-T020)

### Parallel Opportunities

- **Cross-story parallelism**: All three user stories can be worked simultaneously by different developers
- **Within US2**: All 8 documentation file updates (T005-T012) can run in parallel
- **Within US3**: All 7 code file updates (T014-T020) can run in parallel
- **Polish Phase**: All verification tasks (T023-T025) can run in parallel

---

## Parallel Example: All User Stories

```bash
# Team with 3 developers can work all stories simultaneously:

Developer A (User Story 1):
  T001: Rename cmd/apple-tui/ → cmd/actui/
  T002: Update .gitignore
  T003: Verify build
  T004: Verify binary name

Developer B (User Story 2 - all in parallel):
  T005: Update README.md build command
  T006: Update README.md execution examples
  T007: Update README.md .config paths
  T008: Update README.md Application Support paths
  T009: Update user-guide.md build command
  T010: Update user-guide.md execution examples
  T011: Update user-guide.md .config paths
  T012: Update user-guide.md Application Support paths
  T013: Verify documentation

Developer C (User Story 3 - all in parallel):
  T014: Update config_manager.go .config path
  T015: Update config_manager.go Application Support path
  T016: Update config_manager.go writePath
  T017: Update log_writer.go log path
  T018: Update services_test.go config path
  T019: Update services_test.go log path check
  T020: Update services_test.go log path
  T021: Verify code consistency
  T022: Verify historical specs
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete T001-T004: User Story 1
2. **STOP and VALIDATE**: Run build, verify "actui" binary created
3. Deploy/demo if ready (binary renamed, but docs/code still reference old name)

### Incremental Delivery (Recommended)

1. Complete User Story 1 (T001-T004) → Binary renamed ✓
2. Complete User Story 2 (T005-T013) → Documentation updated ✓
3. Complete User Story 3 (T014-T022) → Code cleaned up ✓
4. Complete Polish (T023-T027) → Full validation ✓

Each story adds value incrementally:
- After US1: Build works with new name
- After US2: Users see correct documentation
- After US3: Codebase fully consistent

### Parallel Team Strategy

With multiple developers or parallelizable work:

1. Launch all three user stories simultaneously (independent files)
2. Within each story, parallelize file updates
3. Converge on Polish phase for final validation

**Estimated Time**:
- Sequential: 30-60 minutes
- Parallel (3 developers): 15-20 minutes

---

## Notes

- **[P] tasks**: Different files, can run in parallel
- **[Story] label**: Maps task to user story (US1, US2, US3)
- **Historical specs**: Do NOT change specs/001-apple-container-tui/ - preserved as documentation
- **Module name**: Do NOT change go.mod module name - explicitly out of scope
- **Config migration**: Users must manually migrate config files - no automatic migration

---

## Success Criteria Validation

After completing all tasks, verify:

- ✅ **SC-001**: Build process produces binary named "actui" (T003-T004)
- ✅ **SC-002**: Zero occurrences of "apple-tui" in active code/config/docs (T013, T021)
- ✅ **SC-003**: README and docs reference "actui" consistently (T005-T013)
- ✅ **SC-004**: Build completes without errors (T003, T023)
- ✅ **SC-005**: Users following README use "actui" command (T005-T013)
- ✅ **SC-006**: All test suites pass with updated references (T023)
