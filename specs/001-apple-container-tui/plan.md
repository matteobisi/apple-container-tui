# Implementation Plan: Apple Container TUI

**Branch**: `001-apple-container-tui` | **Date**: 2026-02-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-apple-container-tui/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a terminal user interface (TUI) for managing Apple Container operations on macOS 26.x with Apple Silicon. The TUI provides keyboard-driven menus for listing, starting, stopping, deleting containers, pulling images, building from container files, and managing the daemon - all with command previews and safety confirmations for destructive actions.

## Technical Context

**Language/Version**: Go 1.21+ (chosen for optimal balance of productivity, performance, binary distribution, and TUI library maturity)  
**Primary Dependencies**: Bubbletea v1.2.4 (TUI framework), Lipgloss v1.0.0 (styling), Bubbles v0.20.0 (UI components), Cobra (CLI), Viper (config management)  
**Storage**: Local files (user config in ~/.config/apple-tui/ and ~/Library/Application Support/apple-tui/)  
**Testing**: Go's built-in `go test` with strategy pattern for dry-run/execute separation, contract tests for command building, integration tests for start/stop/pull/delete flows (required per constitution principle V)  
**Target Platform**: macOS 26.x on Apple Silicon (M-series CPUs)
**Project Type**: single (standalone CLI/TUI application)  
**Performance Goals**: Interactive TUI response <100ms for UI updates (Go delivers 15-50ms), command execution time depends on Apple Container CLI  
**Constraints**: Local-only operation, no network dependencies, keyboard-only navigation, no elevated privileges for normal operation, must handle terminal resize gracefully  
**Scale/Scope**: Single-user local tool, ~10 concurrent container operations, support for dozens of containers in list view

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. Command-Safe TUI
- ✅ **PASS**: FR-008 requires command preview before execution
- ✅ **PASS**: FR-008 requires dry-run mode support
- ✅ **PASS**: FR-004 requires explicit confirmation for destructive actions (delete)
- ✅ **PASS**: User Story 3 covers daemon start/stop with confirmation
- ✅ **PASS** (Phase 1): contracts/cli-commands.md specifies command preview pattern and dry-run executor
- ✅ **PASS** (Phase 1): contracts/cli-commands.md specifies type-to-confirm for delete operations

### II. macOS 26.x + Apple Silicon Target
- ✅ **PASS**: FR-012 explicitly requires macOS 26.x on Apple Silicon
- ✅ **PASS**: No cross-platform requirements in scope
- ✅ **PASS** (Phase 1): Go 1.21+ selected with native ARM64 support, GOOS=darwin GOARCH=arm64

### III. Local-Only Operation
- ✅ **PASS**: No remote control surface in requirements
- ✅ **PASS**: FR-011 specifies local config storage only
- ✅ **PASS**: No telemetry or cloud dependencies mentioned
- ✅ **PASS** (Phase 1): data-model.md confirms local file storage for config and logs
- ✅ **PASS** (Phase 1): quickstart.md documents log location: ~/Library/Application Support/apple-tui/command.log

### IV. Clear Observability
- ✅ **PASS**: FR-009 requires stdout/stderr display for all actions
- ✅ **PASS**: FR-009 requires success/failure visibility
- ✅ **PASS** (Phase 1): quickstart.md specifies log path and format (JSON lines)
- ✅ **PASS** (Phase 1): data-model.md defines CommandRun entity capturing stdout/stderr/exit code
- ✅ **PASS** (Phase 1): quickstart.md Help screen documents config and log paths

### V. Tested Command Contracts
- ✅ **PASS** (Phase 1): contracts/cli-commands.md specifies unit test patterns for command builders
- ✅ **PASS** (Phase 1): contracts/cli-commands.md specifies contract tests for destructive actions
- ✅ **PASS** (Phase 1): contracts/cli-commands.md specifies integration test patterns (optional)
- ⚠️ **DEFER TO TASKS**: Actual test implementation will be in tasks.md

**Status**: ✅ All gates pass. All Phase 1 verification items resolved. Test implementation deferred to task phase.

## Project Structure

### Documentation (this feature)

```text
specs/001-apple-container-tui/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── cli-commands.md  # Apple Container CLI command mapping
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
src/
├── models/          # Container, ImageReference, BuildSource, DaemonStatus, CommandRun, UserConfig
├── services/        # Apple Container CLI wrapper, config manager, command builder
├── ui/              # TUI screens, menus, keyboard handlers, layout managers
└── main            # Entry point, CLI check, app initialization

tests/
├── contract/        # Command composition and CLI argument validation tests
├── integration/     # End-to-end container operation tests (if Apple Container available)
└── unit/            # Model and service unit tests

docs/
└── user-guide.md    # Usage documentation

config/
└── default.toml     # Default configuration template
```

**Structure Decision**: Single project structure chosen because this is a standalone TUI application with no separate backend/frontend or mobile components. All functionality runs in a single process with the TUI framework handling UI and the services layer wrapping Apple Container CLI commands.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No constitution violations detected. All principles are satisfied by the current design.
