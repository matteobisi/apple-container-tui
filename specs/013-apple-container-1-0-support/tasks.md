# Tasks: Apple Container 1.0 Support with Container Machine

**Input**: Design documents from `specs/013-apple-container-1-0-support/`  
**Prerequisites**: plan.md ✅ | spec.md ✅ | research.md ✅ | data-model.md ✅ | contracts/machine-commands.md ✅

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel with other [P] tasks in the same phase (different files, no unmet dependencies)
- **[Story]**: Which user story this task belongs to (US1–US6)
- Setup and Foundational phases have no story label

---

## Phase 1: Setup

**Purpose**: Verify the development baseline is green before any changes.

- [X] T001 Run `go test ./...` from repo root to confirm all existing tests pass before any changes

**Checkpoint**: All existing tests green — safe to proceed.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: New types and screen-ID constants that ALL machine user stories (US1, US2, US3) depend on. No machine screen can be implemented until this phase is complete.

- [X] T002 Create `ContainerMachine`, `MachineState` constants, and `MachineEditInput` types in `src/models/machine.go`
- [X] T003 Add `ScreenMachineList`, `ScreenMachineSubmenu`, `ScreenMachineInspect`, `ScreenMachineLogs`, `ScreenMachineEditResources` constants to `src/ui/messages.go`
- [X] T004 [P] Add machine screen fields (`machineList`, `machineSub`, `machineInspect`, `machineLogs`, `machineEditRes`) and empty routing switch-cases to `src/ui/app.go` (depends on T003)
- [X] T005 [P] Add `M` key binding that navigates to `ScreenMachineList` in `src/ui/container_list.go` (depends on T003)

**Checkpoint**: Foundation ready — machine user stories can now be implemented.

---

## Phase 3: User Story 5 — Version String Fix (Priority: P1) 🎯

**Goal**: `actui --version` reports `actui version 0.1.12` matching the published GitHub release.

**Independent Test**: Build from source and run `./actui --version` — output must read `actui version 0.1.12`.

- [X] T006 [US5] Update `version` constant from `"0.4"` to `"0.1.12"` in `cmd/actui/main.go`

**Checkpoint**: User Story 5 complete and independently verifiable with `go run ./cmd/actui --version`.

---

## Phase 4: User Story 6 — Apple Container 1.0 Compatibility Fix (Priority: P2)

**Goal**: Existing Registries screen renders real data against Apple Container 1.0's new JSON output shape.

**Independent Test**: With Apple Container 1.0 running, open actui → press `i` → press `g`; all configured registries appear with non-empty hostname column.

- [X] T007 [US6] Add JSON struct tags to `RegistryLogin` in `src/models/registry_login.go`: `name`→`Hostname`, `username`→`Username`, change `CreatedDate`/`ModifiedDate` to `time.Time` with `json:"creationDate"` and `json:"modificationDate"` tags
- [X] T008 [US6] Add regression test covering Apple Container 1.0 registry JSON format in `src/services/services_test.go` (depends on T007)

**Checkpoint**: User Story 6 complete — run `go test ./src/services/...` to verify.

---

## Phase 5: User Story 1 — Browse and Manage Container Machines (Priority: P1) 🎯

**Goal**: Full machine list and management submenu accessible via `M` from the container list.

**Independent Test**: Press `M` from the container list — machine table appears. Press `enter` on a machine — submenu shows contextual actions. Stop/start/delete/set-default actions preview commands and execute correctly.

### Service Layer (all parallel — different files)

- [X] T009 [P] [US1] Create `MachineListBuilder` in `src/services/machine_list_builder.go` building `container machine list --format json`
- [X] T010 [P] [US1] Create `ParseMachineList` function in `src/services/machine_parser.go` unmarshalling JSON array into `[]models.ContainerMachine`
- [X] T011 [P] [US1] Create `MachineStopBuilder` in `src/services/machine_stop_builder.go` building `container machine stop <id>`
- [X] T012 [P] [US1] Create `MachineStartBuilder` in `src/services/machine_start_builder.go` building `container machine run -n <id>`
- [X] T013 [P] [US1] Create `MachineDeleteBuilder` in `src/services/machine_delete_builder.go` building `container machine delete <id>`
- [X] T014 [P] [US1] Create `MachineSetDefaultBuilder` in `src/services/machine_set_default_builder.go` building `container machine set-default <id>`
- [X] T015 [P] [US1] Add unit tests for `ParseMachineList` and all US1 builders in `src/services/services_test.go`

### UI Layer

- [X] T016 [US1] Create `MachineListScreen` in `src/ui/machine_list.go`: table with Name/State/Image/Default columns, empty-state message, `enter` to open submenu, `esc` returns to container list (depends on T009, T010)
- [X] T017 [US1] Create `MachineSubmenuScreen` skeleton in `src/ui/machine_submenu.go`: state-aware option set; Stop (running only), Start (stopped only), Delete with type-to-confirm, Set as Default with command preview (depends on T011, T012, T013, T014)
- [X] T018 [US1] Wire `machineList` and `machineSub` into `Update` and `View` routing in `src/ui/app.go` (depends on T016, T017)
- [X] T019 [US1] Add machine list and submenu navigation tests in `src/ui/ui_flow_test.go` (depends on T018)

**Checkpoint**: User Story 1 complete — navigate to Machines, list machines, open submenu, stop/start/delete/set-default all work.

---

## Phase 6: User Story 2 — View Machine Details and Logs (Priority: P2)

**Goal**: Inspect and Logs actions available from the machine submenu.

**Independent Test**: From machine submenu select Inspect — detail view shows configuration. Select Logs — log output displays. `esc` from each returns to submenu.

### Service Layer (parallel)

- [X] T020 [P] [US2] Create `MachineInspectBuilder` in `src/services/machine_inspect_builder.go` building `container machine inspect <id>`
- [X] T021 [P] [US2] Create `MachineLogsBuilder` in `src/services/machine_logs_builder.go` building `container machine logs <id>`

### UI Layer

- [X] T022 [US2] Create `MachineInspectScreen` in `src/ui/machine_inspect.go`: raw JSON display with scroll, `esc`/`q` returns to submenu (depends on T020)
- [X] T023 [US2] Create `MachineLogsScreen` in `src/ui/machine_logs.go`: reuse ContainerLogs rendering pattern, `esc`/`q` returns to submenu (depends on T021)
- [X] T024 [US2] Add Inspect and Logs actions to `MachineSubmenuScreen` in `src/ui/machine_submenu.go` (depends on T022, T023)
- [X] T025 [US2] Wire `machineInspect` and `machineLogs` into `Update` and `View` routing in `src/ui/app.go` (depends on T022, T023)

**Checkpoint**: User Story 2 complete — inspect and logs reachable from submenu and navigate back correctly.

---

## Phase 7: User Story 3 — Edit Machine Resources (Priority: P2)

**Goal**: Edit Resources action in submenu opens a pre-filled form; confirms with `container machine set` command preview; shows advisory that changes require restart.

**Independent Test**: Select Edit Resources from submenu — form pre-filled with current cpus/memory/home-mount. Change values, confirm — command preview shows correct `container machine set -n <id> cpus=… memory=… home-mount=…`. Advisory note visible.

### Service Layer

- [X] T026 [P] [US3] Create `MachineSetBuilder` in `src/services/machine_set_builder.go` building `container machine set -n <id> cpus=<N> memory=<M> home-mount=<mode>`

### UI Layer

- [X] T027 [US3] Create `MachineEditResourcesScreen` in `src/ui/machine_edit_resources.go`: three-field form (CPUs, Memory, HomeMount) pre-filled from selected machine, command preview before submission, advisory text "Changes take effect after next stop and restart" (depends on T026)
- [X] T028 [US3] Add Edit Resources action to `MachineSubmenuScreen` in `src/ui/machine_submenu.go` (depends on T027)
- [X] T029 [US3] Wire `machineEditRes` into `Update` and `View` routing in `src/ui/app.go` (depends on T027)

**Checkpoint**: User Story 3 complete — resource editing form works end-to-end with command preview.

---

## Phase 8: User Story 4 — Create Container Machine (Priority: P3) ⚠️ DEFERRED

**Goal**: Create a new container machine from an image reference directly from actui.

**Note**: This phase is explicitly deferred per research.md (Q9). Implement only after US1–US3 are validated. All tasks in this phase are independent of each other after T030.

**Independent Test**: Press create hotkey from machine list — form accepts image reference and optional name. Confirm — `container machine create` command preview shown and executed. New machine appears in list on return.

- [X] T030 [US4] Add `ScreenMachineCreate` constant to `src/ui/messages.go` and empty routing case to `src/ui/app.go`
- [X] T031 [P] [US4] Create `MachineCreateBuilder` in `src/services/machine_create_builder.go` building `container machine create <image> [--name <name>]`
- [X] T032 [US4] Create `MachineCreateScreen` in `src/ui/machine_create.go`: image reference input (required) + optional name field, command preview before execution, returns to machine list on success (depends on T031)
- [X] T033 [US4] Wire create hotkey (e.g. `c`) in `src/ui/machine_list.go` and `machineCreate` routing in `src/ui/app.go` (depends on T032)

**Checkpoint**: User Story 4 complete — create flow functional end-to-end.

---

## Final Phase: Polish & Cross-Cutting Concerns

- [X] T034 [P] Update `docs/ai-menu-map.md`: add machine screens to screen graph and screen ownership table; add M hotkey entry to Container Root entry points section
- [X] T035 Run `go test ./...` from repo root; fix any compilation errors or test failures
- [X] T036 Run manual validation using `specs/013-apple-container-1-0-support/quickstart.md` on macOS with Apple Container 1.0

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately
- **Phase 2 (Foundational)**: Depends on Phase 1 — **BLOCKS US1, US2, US3, US4**
- **Phase 3 (US5)**: Independent — can run in parallel with Phase 2
- **Phase 4 (US6)**: Independent — can run after Phase 1; no dependency on Foundational
- **Phase 5 (US1)**: Depends on Phase 2 (Foundational complete)
- **Phase 6 (US2)**: Depends on Phase 2 + Phase 5 (machine submenu shell in place)
- **Phase 7 (US3)**: Depends on Phase 2 + Phase 5 (machine submenu shell in place)
- **Phase 8 (US4)**: Depends on Phase 2 + Phase 5 (machine list screen in place) — deferred
- **Final Phase**: Depends on all desired stories complete

### User Story Dependencies

- **US5 (P1)**: Independent of everything — single-line fix
- **US6 (P2)**: Independent of machine feature — fix existing registry model
- **US1 (P1)**: Depends on Foundational (Phase 2) only
- **US2 (P2)**: Depends on US1 (machine submenu must exist to add Inspect/Logs actions)
- **US3 (P2)**: Depends on US1 (machine submenu must exist to add Edit Resources action)
- **US4 (P3)**: Depends on US1 (machine list must exist for create hotkey) — deferred

### Within Each User Story

- Service layer tasks always before their dependent UI screens
- UI screens before wiring into app.go
- App.go wiring before navigation tests
- Models before service builders (T002 before T009–T014)

### Parallel Opportunities

- T004 and T005 (Foundational): both depend on T003; can run in parallel with each other
- T009–T015 (US1 services): all different files, no interdependencies — full parallel
- T020 and T021 (US2 services): different files — parallel
- T003 (Foundational) and T006 (US5): completely independent — parallel

---

## Parallel Execution Examples

### User Story 1 Service Layer (after Foundational complete)

```
Simultaneously start all of:
  T009  machine_list_builder.go
  T010  machine_parser.go
  T011  machine_stop_builder.go
  T012  machine_start_builder.go
  T013  machine_delete_builder.go
  T014  machine_set_default_builder.go
  T015  services_test.go (machine tests)

Then once all complete:
  T016  machine_list.go        \
  T017  machine_submenu.go      } parallel (different files)
  Then: T018  app.go wiring
  Then: T019  ui_flow_test.go
```

### Foundational + Quick Fixes (parallel start)

```
Simultaneously:
  T002  src/models/machine.go
  T006  cmd/actui/main.go (US5 version fix)
  T007  src/models/registry_login.go (US6 compat fix)
```

---

## Implementation Strategy

### MVP Scope (US5 + US1 + US6 = 19 tasks)

1. Complete **Phase 1** (T001) — baseline green
2. Complete **Phase 2** (T002–T005) — foundations
3. Complete **Phase 3** (T006) — version fix (can overlap with Phase 2)
4. Complete **Phase 4** (T007–T008) — compat fix (can overlap with Phase 2)
5. Complete **Phase 5** (T009–T019) — machine list + submenu
6. **STOP and VALIDATE**: Run quickstart validations 1–9 and `go test ./...`

### Full Delivery

Add Phase 6 (US2) → Phase 7 (US3) → validate → add Phase 8 (US4, deferred) → Polish.

### Single Developer Recommended Order

T001 → T002 → T003 → T004 → T005 → T006 → T007 → T008 → T009 → T010 → T011 → T012 → T013 → T014 → T015 → T016 → T017 → T018 → T019 → T020 → T021 → T022 → T023 → T024 → T025 → T026 → T027 → T028 → T029 → T030 → T031 → T032 → T033 → T034 → T035 → T036

---

## Notes

- All machine builders must accept the `services.CommandExecutor` interface — no direct `exec.Command` calls
- Destructive actions (T013 delete, T033 create — indirectly) must use type-to-confirm pattern from `src/ui/type_to_confirm.go`
- Stop/Start/SetDefault must use yes/no confirm pattern from `src/ui/yes_no_confirm.go`
- Command previews via `src/ui/command_preview.go`
- Machine submenu (T017) state-awareness: show Stop only when `MachineStateRunning`; show Start only when `MachineStateStopped`; show all other actions regardless of state
- US4 (T030–T033) is gated — do not begin until US1 is fully validated
