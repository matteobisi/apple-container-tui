# Tasks: Enhanced TUI Display with Table Layout

**Feature**: 003-tui-table-format  
**Input**: Design documents from `/specs/003-tui-table-format/`  
**Tech Stack**: Go 1.21, Bubbletea v1.2.4, Lipgloss v1.0.0

**Organization**: Tasks grouped by user story to enable independent implementation and testing. Tests not included per feature specification (no TDD requirement).

## Phase 1: Setup

**Purpose**: Verify environment and dependencies

- [X] T001 Verify Go 1.21+ installed and Bubbletea v1.2.4 + Lipgloss v1.0.0 in go.mod

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core table component that MUST be complete before user story implementation

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T002 Create src/ui/table.go with TableColumn, TableRow, and Table type definitions
- [X] T003 Implement NewTable() constructor function in src/ui/table.go
- [X] T004 Implement SetRows() method in src/ui/table.go
- [X] T005 Implement calculateColumnWidths() method with priority-based allocation in src/ui/table.go
- [X] T006 Implement padOrTruncate() helper function for cell formatting in src/ui/table.go
- [X] T007 Implement Render() method with header, separator, and row rendering in src/ui/table.go
- [X] T008 Implement TruncateDigest() helper function in src/ui/table.go

**Checkpoint**: Table component complete and ready for integration - user story implementation can now begin

---

## Phase 3: User Story 1 - Clean Screen Experience (Priority: P1) 🎯 MVP

**Goal**: Launch application in full-screen mode with alternate screen buffer (like btop), restoring terminal state on exit

**Independent Test**: Launch actui, verify screen clears completely showing only TUI, quit with 'q', confirm previous terminal content restored without TUI in scrollback

### Implementation for User Story 1

- [X] T009 [US1] Modify cmd/actui/main.go to add tea.WithAltScreen() option to tea.NewProgram() call

**Verification Steps**:
1. Run `go run cmd/actui/main.go`
2. Confirm screen clears completely on launch
3. Press 'q' to quit
4. Confirm previous terminal state restored
5. Scroll up in terminal - confirm no TUI output in scrollback

**Checkpoint**: User Story 1 complete - application now launches in clean full-screen mode

---

## Phase 4: User Story 2 - Tabular Container View (Priority: P2)

**Goal**: Display container list as formatted table with Name, State, Base Image columns, headers, alignment, and visual separators

**Independent Test**: Launch actui, verify containers displayed in table with bold headers, separator line below headers, aligned columns (Name/State/Base Image), row highlighting on navigation, visual separator between table and menu

### Implementation for User Story 2

- [X] T010 [US2] Modify src/ui/container_list.go View() method to create Table with 3 columns (Name/State/Base Image)
- [X] T011 [US2] Add loop in src/ui/container_list.go View() to populate TableRow structs from m.containers
- [X] T012 [US2] Call table.Render() with m.width and m.cursor in src/ui/container_list.go View()
- [X] T013 [US2] Add horizontal separator line after table using strings.Repeat("─", m.width) in src/ui/container_list.go View()
- [X] T014 [US2] Update WindowSizeMsg handling in src/ui/container_list.go Update() to store width (if not already present)

**Verification Steps**:
1. Run `go run cmd/actui/main.go`
2. Container list should display as table with:
   - Bold column headers: Name | State | Base Image
   - Horizontal separator line below headers
   - Container data in aligned columns
   - Selected row with inverse video highlighting
   - Horizontal separator between table and keyboard menu
3. Use arrow keys to navigate - selection should move and highlight should follow
4. Resize terminal - table should re-render with adjusted column widths
5. Test with no containers - should show headers with "No items found"

**Checkpoint**: User Story 2 complete - container list now displays as professional table format

---

## Phase 5: User Story 3 - Tabular Image View (Priority: P3)

**Goal**: Display image list as formatted table with Name, Tag, Digest columns (digest truncated to 12 chars), headers, alignment, and visual separators

**Independent Test**: Navigate to images view (press 'i'), verify images displayed in table with bold headers, separator line below headers, aligned columns (Name/Tag/Digest), digest showing 12 characters, row highlighting on navigation, visual separator between table and menu

### Implementation for User Story 3

- [X] T015 [US3] Modify src/ui/image_list.go View() method to create Table with 3 columns (Name/Tag/Digest)
- [X] T016 [US3] Add loop in src/ui/image_list.go View() to populate TableRow structs from m.images with TruncateDigest() on digest field
- [X] T017 [US3] Call table.Render() with m.width and m.cursor in src/ui/image_list.go View()
- [X] T018 [US3] Add horizontal separator line after table using strings.Repeat("─", m.width) in src/ui/image_list.go View()
- [X] T019 [US3] Update WindowSizeMsg handling in src/ui/image_list.go Update() to store width (if not already present)

**Verification Steps**:
1. Run `go run cmd/actui/main.go`
2. Press 'i' to navigate to images view
3. Image list should display as table with:
   - Bold column headers: Name | Tag | Digest
   - Horizontal separator line below headers
   - Image data in aligned columns
   - Digest values truncated to 12 characters (no "sha256:" prefix)
   - Selected row with inverse video highlighting
   - Horizontal separator between table and keyboard menu
4. Use arrow keys to navigate - selection should move and highlight should follow
5. Resize terminal - table should re-render with adjusted column widths
6. Test with long repository names - should truncate with ellipsis
7. Test with no images - should show headers with "No items found"

**Checkpoint**: User Story 3 complete - image list now displays as professional table format matching container list

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and documentation

- [X] T020 [P] Run full manual validation per quickstart.md testing checklist
- [X] T021 [P] Test resize behavior at 80-character minimum width
- [X] T022 [P] Test edge cases: empty lists, very long names, constrained terminal width
- [X] T023 [P] Update README.md with screenshots or description of new table views (if applicable)
- [X] T024 Verify constitution compliance: all command execution paths unchanged, display-only changes

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - verify environment
- **Foundational (Phase 2)**: Table component creation - BLOCKS all user stories (US1, US2, US3)
- **User Story 1 (Phase 3)**: Can start immediately after Setup (independent of table component)
- **User Story 2 (Phase 4)**: Depends on Foundational phase completion
- **User Story 3 (Phase 5)**: Depends on Foundational phase completion
- **Polish (Phase 6)**: Depends on desired user stories being complete

### User Story Dependencies

```
Setup (T001)
    ↓
Foundational (T002-T008) ← Creates table component
    ↓
    ├─→ US1 (T009) ← Independent, can start after Setup
    ├─→ US2 (T010-T014) ← Requires table component
    └─→ US3 (T015-T019) ← Requires table component
        ↓
Polish (T020-T024)
```

- **User Story 1**: Independent - only modifies main.go, no table component dependency
- **User Story 2**: Depends on table component (T002-T008) but independent of US1/US3
- **User Story 3**: Depends on table component (T002-T008) but independent of US1/US2

### Parallel Opportunities

**After Foundational Phase Completes**:
- US1 can run in parallel with US2 (different files: main.go vs container_list.go)
- US2 and US3 can run in parallel (different files: container_list.go vs image_list.go)
- All Polish tasks (T020-T024) marked [P] can run in parallel

**Optimal Execution with Single Developer**:
1. Complete T001 (Setup)
2. Complete T002-T008 (Foundational) - ~30 minutes
3. Complete T009 (US1) - ~5 minutes → TEST independently
4. Complete T010-T014 (US2) - ~20 minutes → TEST independently
5. Complete T015-T019 (US3) - ~20 minutes → TEST independently
6. Complete T020-T024 (Polish) in parallel - ~15 minutes

**Total Estimated Time**: ~1.5 hours

---

## Parallel Example: After Foundational Phase

```bash
# All user stories can proceed in parallel after table component is ready:

# Developer A or Time Slot 1:
Task T009: "Modify cmd/actui/main.go to add tea.WithAltScreen()"

# Developer B or Time Slot 2 (can run parallel with A):
Task T010: "Modify src/ui/container_list.go View() method..."
Task T011: "Add loop in src/ui/container_list.go View()..."
Task T012: "Call table.Render() in src/ui/container_list.go..."
Task T013: "Add horizontal separator in src/ui/container_list.go..."
Task T014: "Update WindowSizeMsg handling in src/ui/container_list.go..."

# Developer C or Time Slot 3 (can run parallel with B):
Task T015: "Modify src/ui/image_list.go View() method..."
Task T016: "Add loop in src/ui/image_list.go View()..."
Task T017: "Call table.Render() in src/ui/image_list.go..."
Task T018: "Add horizontal separator in src/ui/image_list.go..."
Task T019: "Update WindowSizeMsg handling in src/ui/image_list.go..."
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

**Fastest Path to Value**:
1. T001 (Setup verification)
2. T009 (Enable alternate screen)
3. Validate: Launch app, verify clean screen, quit, verify restore
4. **DEMO READY**: Professional full-screen TUI experience

**Note**: US1 doesn't require table component, so can deliver value immediately without T002-T008

### Incremental Delivery (Recommended)

**Add Table Features Progressively**:
1. Complete T001 (Setup)
2. Complete T002-T008 (Foundational table component)
3. **Milestone A**: Complete T009 (US1) → Clean screen mode working
4. **Milestone B**: Complete T010-T014 (US2) → Container table ready → Test independently
5. **Milestone C**: Complete T015-T019 (US3) → Image table ready → Test independently
6. Complete T020-T024 (Polish)

**Each milestone delivers independently testable value**

### Parallel Team Strategy

**With 2+ Developers**:
1. Team: Complete T001 together (5 min)
2. Developer A: Complete T002-T008 (Foundational, 30 min)
3. Once T002-T008 done:
   - Developer B: T009 (US1, 5 min) + T010-T014 (US2, 20 min)
   - Developer C: T015-T019 (US3, 20 min)
4. Team: T020-T024 (Polish, parallel, 15 min)

**Total Team Time**: ~50 minutes (vs 1.5 hours sequential)

---

## File Summary

### New Files (2)
- `src/ui/table.go` - Reusable table component with rendering logic

### Modified Files (5)
- `cmd/actui/main.go` - Add tea.WithAltScreen() option
- `src/ui/container_list.go` - Integrate table rendering in View()
- `src/ui/image_list.go` - Integrate table rendering in View()

**No test files required** - Feature specification does not request test implementation

---

## Notes

- **[P] marker**: Tasks can run in parallel (different files, no dependencies within phase)
- **[Story] label**: Maps task to user story for traceability (US1/US2/US3)
- **Independent testing**: Each user story can be validated independently
- **No breaking changes**: All existing command execution, navigation, and modal flows remain unchanged
- **Visual only**: This is purely a display enhancement - no business logic changes
- **Estimated effort**: Per quickstart.md, 1.5-2 hours total implementation time

---

## Quick Reference

**Priority Sequence**: US1 (P1: Clean screen) → US2 (P2: Container table) → US3 (P3: Image table)

**Critical Path**: T001 → T002-T008 → {T009, T010-T014, T015-T019} → T020-T024

**Fastest to Demo**: T001 → T009 (10 minutes for clean screen mode)

**Full Feature**: T001 → T002-T008 → T009-T019 → T020-T024 (~1.5 hours)

**Manual Verification**: See quickstart.md "Testing Checklist" section for detailed validation steps
