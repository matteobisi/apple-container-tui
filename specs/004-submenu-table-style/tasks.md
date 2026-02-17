# Tasks: Consistent Table Styling for Submenus

**Feature**: 004-submenu-table-style  
**Input**: Design documents from `/specs/004-submenu-table-style/`  
**Prerequisites**: ✅ plan.md, ✅ spec.md, ✅ research.md, ✅ data-model.md, ✅ contracts/, ✅ quickstart.md

**Tests**: Manual testing only (display-only feature, no command contracts). Per constitution, automated tests not required for pure UI presentation changes.

**Organization**: Tasks grouped by user story for independent implementation and testing.

---

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1=Container Submenu, US2=Image Submenu, US3=Help Screen)
- All file paths are absolute from repository root

---

## Phase 1: Setup

**Purpose**: Verify environment and dependencies

- [ ] T001 Verify Go 1.21+ installed and Bubbletea v1.2.4 + Lipgloss v1.0.0 in go.mod
- [ ] T002 Review Feature 003 styling patterns in src/ui/table.go and src/ui/container_list.go for reference

**Checkpoint**: Environment ready, styling patterns understood

---

## Phase 2: User Story 1 - Container Submenu Visual Consistency (Priority: P1) 🎯 MVP

**Goal**: Apply bold headers, horizontal separators, and inverse video selection to container submenu screen

**Independent Test**: Navigate to any container submenu (press Enter on a container), verify bold headers, separators span full width, and arrow key navigation uses inverse video highlighting without cursor prefix

### Implementation for User Story 1

- [ ] T003 [US1] Add width field (int) to ContainerSubmenuScreen struct in src/ui/container_submenu.go
- [ ] T004 [US1] Add WindowSizeMsg handler to Update() method in src/ui/container_submenu.go to store terminal width
- [ ] T005 [US1] Refactor View() method in src/ui/container_submenu.go: Add bold "Container Details" header before container info section
- [ ] T006 [US1] Refactor View() method in src/ui/container_submenu.go: Add horizontal separator (strings.Repeat("─", m.width)) after container info with blank lines before/after
- [ ] T007 [US1] Refactor View() method in src/ui/container_submenu.go: Add bold "Available Actions" header before action list
- [ ] T008 [US1] Refactor View() method in src/ui/container_submenu.go: Replace cursor prefix pattern with inverse video selection (lipgloss.NewStyle().Reverse(true)) for action items
- [ ] T009 [US1] Refactor View() method in src/ui/container_submenu.go: Add horizontal separator after action list with blank lines before/after
- [ ] T010 [US1] Manual test: Build application (go build -o actui cmd/actui/main.go) and verify container submenu displays correctly

**Checkpoint**: Container submenu fully styled with bold headers, separators, and inverse video selection - User Story 1 complete

---

## Phase 3: User Story 2 - Image Submenu Visual Consistency (Priority: P2)

**Goal**: Apply bold headers, horizontal separators, and inverse video selection to image submenu screen

**Independent Test**: Navigate to any image submenu (press 'i' then Enter on an image), verify bold headers, separators span full width, and arrow key navigation uses inverse video highlighting

### Implementation for User Story 2

- [ ] T011 [P] [US2] Add width field (int) to ImageSubmenuScreen struct in src/ui/image_submenu.go
- [ ] T012 [P] [US2] Add WindowSizeMsg handler to Update() method in src/ui/image_submenu.go to store terminal width
- [ ] T013 [US2] Refactor View() method in src/ui/image_submenu.go: Add bold "Image Details" header before image info section
- [ ] T014 [US2] Refactor View() method in src/ui/image_submenu.go: Add horizontal separator (strings.Repeat("─", m.width)) after image info with blank lines before/after
- [ ] T015 [US2] Refactor View() method in src/ui/image_submenu.go: Add bold "Available Actions" header before action list
- [ ] T016 [US2] Refactor View() method in src/ui/image_submenu.go: Replace cursor prefix pattern with inverse video selection for action items
- [ ] T017 [US2] Refactor View() method in src/ui/image_submenu.go: Add horizontal separator after action list with blank lines before/after
- [ ] T018 [US2] Manual test: Navigate to image submenu and verify visual consistency with container submenu styling

**Checkpoint**: Image submenu fully styled - User Story 2 complete

---

## Phase 4: User Story 3 - Help Screen Visual Consistency (Priority: P3)

**Goal**: Reorganize help screen into 4 sections with bold headers and horizontal separators between each section

**Independent Test**: Press '?' from any screen, verify 4 section headers (Navigation, Container Actions, Image Actions, General) are bold and separated by horizontal lines

### Implementation for User Story 3

- [ ] T019 [US3] Add width field (int) to HelpScreen struct in src/ui/help.go if not present
- [ ] T020 [US3] Add WindowSizeMsg handler to Update() method in src/ui/help.go if not present
- [ ] T021 [US3] Refactor View() method in src/ui/help.go: Reorganize content into "Navigation" section with bold header
- [ ] T022 [US3] Refactor View() method in src/ui/help.go: Add horizontal separator after Navigation section
- [ ] T023 [US3] Refactor View() method in src/ui/help.go: Reorganize content into "Container Actions" section with bold header
- [ ] T024 [US3] Refactor View() method in src/ui/help.go: Add horizontal separator after Container Actions section
- [ ] T025 [US3] Refactor View() method in src/ui/help.go: Reorganize content into "Image Actions" section with bold header
- [ ] T026 [US3] Refactor View() method in src/ui/help.go: Add horizontal separator after Image Actions section
- [ ] T027 [US3] Refactor View() method in src/ui/help.go: Reorganize content into "General" section with bold header
- [ ] T028 [US3] Refactor View() method in src/ui/help.go: Add horizontal separator after General section
- [ ] T029 [US3] Manual test: Open help screen and verify all 4 sections display with bold headers and separators

**Checkpoint**: Help screen reorganized with 4 sections, bold headers, and separators - User Story 3 complete

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Optional enhancement and validation across all screens

- [ ] T030 [P] Optional: Apply same styling to daemon control screen (src/ui/daemon_control.go) if it has submenu-like structure
- [ ] T031 Build final binary (go build -o actui cmd/actui/main.go) and verify no compilation errors
- [ ] T032 Execute validation checklist from specs/004-submenu-table-style/quickstart.md for all updated screens
- [ ] T033 [P] Test terminal resize behavior: Verify separators adjust width dynamically in all updated screens
- [ ] T034 [P] Test narrow terminal (80 chars): Verify layout remains readable in all updated screens
- [ ] T035 Complete application walk-through: Navigate through main list → container submenu → back → images → image submenu → back → help screen, verify visual consistency
- [ ] T036 [P] Update README.md with screenshots or examples of new submenu styling if applicable
- [ ] T037 Create VALIDATION.md in specs/004-submenu-table-style/ documenting manual test results

**Checkpoint**: All screens styled consistently, validated, and documented

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - start immediately
- **User Story 1 (Phase 2)**: Depends on Setup - Start after T002
- **User Story 2 (Phase 3)**: Depends on Setup only - Can start in parallel with US1 after T002
- **User Story 3 (Phase 4)**: Depends on Setup only - Can start in parallel with US1/US2 after T002
- **Polish (Phase 5)**: Depends on desired user stories being complete

### User Story Dependencies

- ✅ **User Story 1 (P1) - Container Submenu**: INDEPENDENT - No dependencies on other stories, can implement and test standalone
- ✅ **User Story 2 (P2) - Image Submenu**: INDEPENDENT - No dependencies on US1, can implement and test standalone, marked [P] for parallel execution
- ✅ **User Story 3 (P3) - Help Screen**: INDEPENDENT - No dependencies on US1/US2, can implement and test standalone

**All user stories can be implemented in parallel by different developers or sequentially by priority (P1 → P2 → P3)**

### Within Each User Story

**Container Submenu (US1)**:
- T003 (width field) → T004 (WindowSizeMsg handler) → Parallelizable from here
- T005-T009 (View() refactoring) can be done in sequence or as single combined refactor
- T010 (manual test) must be last

**Image Submenu (US2)**:
- Same pattern as US1: T011-T012 sequential, T013-T017 can be combined, T018 last
- **Entire US2 marked [P]** because it works on different file than US1

**Help Screen (US3)**:
- T019-T020 first (infrastructure), then T021-T028 (section reorganization), T029 last (test)

### Parallel Opportunities

**Maximum Parallelism** (3 developers):
- Developer 1: US1 (T003-T010) - Container submenu
- Developer 2: US2 (T011-T018) - Image submenu (marked [P])
- Developer 3: US3 (T019-T029) - Help screen

**Sequential Implementation** (1 developer following priorities):
- Week 1: T001-T002 (setup) → T003-T010 (US1) → T010 validation
- Week 2: T011-T018 (US2) → T018 validation
- Week 3: T019-T029 (US3) → T029 validation
- Week 4: T030-T037 (polish)

**Recommended MVP Scope**: User Story 1 only (T001-T010) provides immediate value by updating the most frequently used screen (container submenu)

---

## Implementation Strategy

### MVP First (Recommended)
Focus on **User Story 1 (Container Submenu)** as the MVP. This provides:
- Immediate visual consistency for primary workflow
- Pattern validation before extending to other screens
- Early feedback opportunity

**MVP Deliverable**: T001-T010 complete → Container submenu fully styled → Manual validation passed

### Incremental Delivery
After MVP validation:
1. **Immediate Next**: User Story 2 (Image Submenu) - applies validated pattern to second-most-used screen
2. **Final Polish**: User Story 3 (Help Screen) - completes visual consistency across all screens
3. **Optional Enhancement**: Daemon control screen (T030)

### Risk Mitigation
- Each user story is independently testable (acceptance scenarios in spec.md)
- Manual testing after each story ensures no regression
- Lipgloss patterns already validated in Feature 003 (low technical risk)
- No command logic changes (zero risk to core functionality)

---

## Task Summary

- **Total Tasks**: 37
- **Setup**: 2 tasks (T001-T002)
- **User Story 1 (Container Submenu)**: 8 tasks (T003-T010)
- **User Story 2 (Image Submenu)**: 8 tasks (T011-T018) - All marked [P] for parallel
- **User Story 3 (Help Screen)**: 11 tasks (T019-T029)
- **Polish**: 8 tasks (T030-T037) - 4 marked [P] for parallel

**Parallel Task Count**: 13 tasks marked [P] (35% can run in parallel with proper staffing)

**MVP Scope**: 10 tasks (T001-T010) delivers functional container submenu with new styling

**Estimated Effort**:
- User Story 1: ~2-3 hours (most critical path)
- User Story 2: ~1.5-2 hours (pattern reuse)
- User Story 3: ~2-3 hours (more complex reorganization)
- Polish: ~1-2 hours (validation and documentation)
- **Total**: ~7-10 hours for complete feature

---

**Next Action**: Start with T001 (environment verification), then proceed to T003 (Container Submenu implementation) for MVP delivery.
