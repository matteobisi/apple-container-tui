# Implementation Plan: Apple Container 1.0 Support with Container Machine

**Branch**: `013-apple-container-1-0-support` | **Date**: 2026-06-12 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `specs/013-apple-container-1-0-support/spec.md`

## Summary

Add a new Container Machine management workflow to actui — a dedicated Machines screen accessible via `M` from the container list, with a contextual submenu covering inspect, logs, stop/start, resource editing, set-default, and delete. Simultaneously fix a breaking 1.0 compatibility issue in the registry parser (`RegistryLogin` model has no JSON struct tags matching the new 1.0 output shape) and correct the `actui --version` string from `0.4` to `0.1.12`. Container machine create (P3) is deferred.

## Technical Context

**Language/Version**: Go 1.24.2  
**Primary Dependencies**: charmbracelet/bubbletea v1.3.10, charmbracelet/bubbles v1.0.0, charmbracelet/lipgloss v1.1.0, spf13/cobra v1.10.2  
**Storage**: N/A (stateless TUI; reads from Apple Container CLI output)  
**Testing**: `go test ./...` (unit + UI flow tests in `src/ui/`, service builder/parser tests in `src/services/`)  
**Target Platform**: macOS 26.x, Apple Silicon  
**Project Type**: CLI/TUI desktop application  
**Performance Goals**: Sub-100ms screen render; machine list fetch bounded by CLI execution time  
**Constraints**: Keyboard-only navigation; no mouse; no background goroutines beyond bubbletea command pipeline; all CLI calls through `services.CommandExecutor` interface for dry-run safety  
**Scale/Scope**: ~5 new UI screens; ~8 new service builders/parser; 2 model fixes

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Status |
|-----------|-------|--------|
| I. Command-Safe TUI | All machine actions mapped to CLI commands in [contracts/machine-commands.md](contracts/machine-commands.md); destructive delete uses type-to-confirm; stop/start/set use command preview with yes/no | ✅ PASS |
| II. macOS 26.x + Apple Silicon Target | No cross-platform code introduced; target unchanged | ✅ PASS |
| III. Local-Only Operation | No network calls; no telemetry; machine data from local CLI only | ✅ PASS |
| IV. Clear Observability | All machine actions surface success/failure via result display; machine inspect shows full JSON; logs via existing log viewer pattern | ✅ PASS |
| V. Tested Command Contracts | New builders and parsers require unit tests; machine list parser must have a table/JSON fixture test; registry parser fix must include a regression test for 1.0 JSON format | ✅ PASS (required by tasks) |

**Post-design re-check**: All contracts defined in [contracts/machine-commands.md](contracts/machine-commands.md). Data model in [data-model.md](data-model.md). No violations.

## Project Structure

### Documentation (this feature)

```text
specs/013-apple-container-1-0-support/
├── plan.md              ← this file
├── research.md          ← Phase 0 output
├── data-model.md        ← Phase 1 output
├── quickstart.md        ← Phase 1 output
├── contracts/
│   └── machine-commands.md   ← Phase 1 output
└── tasks.md             ← Phase 2 output (/speckit.tasks - NOT created here)
```

### Source Code (repository root)

```text
cmd/actui/
└── main.go                          # version const: "0.4" → "0.1.12"

src/models/
├── machine.go                       # NEW: ContainerMachine, MachineState, MachineEditInput
└── registry_login.go                # FIX: add JSON struct tags; change date fields to time.Time

src/services/
├── machine_list_builder.go          # NEW: container machine list --format json
├── machine_inspect_builder.go       # NEW: container machine inspect <id>
├── machine_logs_builder.go          # NEW: container machine logs <id>
├── machine_stop_builder.go          # NEW: container machine stop <id>
├── machine_start_builder.go         # NEW: container machine run -n <id>
├── machine_set_builder.go           # NEW: container machine set -n <id> ...
├── machine_set_default_builder.go   # NEW: container machine set-default <id>
├── machine_delete_builder.go        # NEW: container machine delete <id>
├── machine_parser.go                # NEW: ParseMachineList (JSON → []ContainerMachine)
└── services_test.go                 # UPDATE: registry 1.0 fixture test; machine parser tests

src/ui/
├── messages.go                      # ADD: ScreenMachineList, ScreenMachineSubmenu,
│                                    #      ScreenMachineInspect, ScreenMachineLogs,
│                                    #      ScreenMachineEditResources
├── machine_list.go                  # NEW: machine list screen
├── machine_submenu.go               # NEW: machine submenu screen
├── machine_inspect.go               # NEW: machine inspect detail view
├── machine_logs.go                  # NEW: machine logs view
├── machine_edit_resources.go        # NEW: resource editing form
├── container_list.go                # UPDATE: add M hotkey → ScreenMachineList
├── app.go                           # UPDATE: add machine screen fields + switch-case routing
└── ui_flow_test.go                  # UPDATE: add machine navigation tests

docs/
└── ai-menu-map.md                   # UPDATE: screen graph, screen ownership table

tests/
├── unit/                            # existing
└── integration/                     # existing
```

**Structure Decision**: Single-project layout matching existing repository convention. All new files follow the established `*_builder.go` / `*_parser.go` pattern in `src/services/` and the `*_screen` naming convention in `src/ui/`.
