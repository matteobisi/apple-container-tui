# Implementation Plan: Enhanced TUI Display with Table Layout

**Branch**: `003-tui-table-format` | **Date**: 2026-02-17 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/003-tui-table-format/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Enhance the TUI display with two main improvements: (1) Launch the application in full-screen mode using alternate screen buffer (like btop) to provide an isolated viewing experience that restores terminal state on exit, and (2) Replace the current unstructured container and image lists with properly formatted tables featuring column headers, aligned data columns, and visual separators. Technical approach involves integrating Bubbletea's alternate screen mode, implementing table rendering with dynamic column width calculation, and applying inverse video highlighting for row selection.

## Technical Context

**Language/Version**: Go 1.21  
**Primary Dependencies**: Bubbletea v1.2.4 (TUI framework), Lipgloss v1.0.0 (styling), Bubbles v0.20.0 (components)  
**Storage**: N/A (display-only feature)  
**Testing**: Go testing (existing test suite in tests/unit/, tests/integration/, tests/contract/)  
**Target Platform**: macOS 26.x on Apple Silicon (M-series)  
**Project Type**: Single project (TUI application)  
**Performance Goals**: <50ms table re-render on terminal resize, instant (<16ms) row selection highlight update  
**Constraints**: Minimum 80-character terminal width support, ANSI escape sequence compatibility, no external rendering dependencies  
**Scale/Scope**: 2 list views (containers, images), ~10-100 items per list typical, 12 UI screen types in codebase

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Principle I: Command-Safe TUI**  
✅ COMPLIANT - No changes to command execution logic; only affects display formatting. Existing command preview, dry-run, and confirmation flows remain unchanged.

**Principle II: macOS 26.x + Apple Silicon Target**  
✅ COMPLIANT - Feature targets same platform. Alternate screen buffer is standard ANSI functionality supported by macOS Terminal.app and iTerm2.

**Principle III: Local-Only Operation**  
✅ COMPLIANT - Pure display feature with no telemetry, cloud dependencies, or remote connections. All rendering happens locally in terminal.

**Principle IV: Clear Observability**  
✅ ENHANCED - Improves observability by making container and image data more scannable through structured table layout. Does not affect error display or logging.

**Principle V: Tested Command Contracts**  
✅ COMPLIANT - No changes to command composition or contracts. Display logic will be unit-testable through existing Go test infrastructure. Existing integration tests remain valid.

**Platform and Runtime Constraints**  
✅ COMPLIANT - No changes to Apple Container CLI interaction, daemon management, or privilege requirements. Configuration and logs unchanged.

**Workflow and Quality Gates**  
✅ COMPLIANT - No new destructive actions introduced. Manual verification on macOS 26.x required before release (standard process).

**Governance**  
✅ COMPLIANT - No constitution exceptions required. Feature aligns with all principles and constraints.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
src/
├── models/           # Data models (Container, Image, Command, Result, etc.)
├── services/         # Business logic and command builders
└── ui/              # TUI screens and components
    ├── app.go                    # Main application model (MODIFY: add alternate screen)
    ├── container_list.go         # Container list screen (MODIFY: table rendering)
    ├── image_list.go             # Image list screen (MODIFY: table rendering)
    ├── theme.go                  # Styling constants (MODIFY: add table styles)
    └── table.go                  # NEW: Reusable table component

tests/
├── contract/         # Command builder contract tests
├── integration/      # End-to-end integration tests
└── unit/            # Unit tests for services and UI components
    └── table_test.go            # NEW: Table component tests

cmd/actui/
└── main.go          # Entry point (MODIFY: enable alternate screen in tea.Program)
```

**Structure Decision**: Single project structure (Option 1). This feature adds one new file (`src/ui/table.go`) for the reusable table component and modifies 5 existing files (`app.go`, `container_list.go`, `image_list.go`, `theme.go`, `main.go`) to integrate alternate screen mode and table rendering. All modifications stay within the existing `src/ui/` directory following established patterns.

## Constitution Check (Post-Design Re-evaluation)

*Re-evaluated after Phase 1 design completion. All gates remain COMPLIANT.*

**Principle I: Command-Safe TUI**  
✅ COMPLIANT - Design confirmed: no command execution changes. Table rendering is purely presentational. Command preview, dry-run, and confirmation flows completely untouched.

**Principle II: macOS 26.x + Apple Silicon Target**  
✅ COMPLIANT - Design uses Bubbletea's `tea.WithAltScreen()` which leverages standard ANSI escape sequences supported by all macOS terminal emulators.

**Principle III: Local-Only Operation**  
✅ COMPLIANT - Design confirmed: zero network calls, no telemetry, no external services. All rendering is in-process string manipulation.

**Principle IV: Clear Observability**  
✅ ENHANCED - Table design with column headers, alignment, and visual separators significantly improves data scannability per success criteria (SC-003, SC-004). Error display paths unchanged.

**Principle V: Tested Command Contracts**  
✅ COMPLIANT - Design includes unit tests for table component (`table_test.go`) and integration tests for screen rendering. No command contract changes means existing contract tests remain valid without modification.

**Platform and Runtime Constraints**  
✅ COMPLIANT - Design confirmed: no changes to Apple Container CLI interaction. Alternate screen and table rendering are terminal-only concerns, well-isolated from backend logic.

**Workflow and Quality Gates**  
✅ COMPLIANT - No new destructive actions added. Design maintains all existing confirmation flows. Manual verification on macOS 26.x Terminal.app and iTerm2 will validate full-screen behavior per standard gates.

**Governance**  
✅ COMPLIANT - No constitution amendments required. Feature implementation fully aligns with all governing principles.

## Complexity Tracking

**No violations recorded.** All constitution checks passed in both initial and post-design evaluations.

---

## Phase 0 Deliverables ✅

- [research.md](research.md) - Technical decisions for alternate screen, table rendering, highlighting, resize handling, digest truncation, and header styling

## Phase 1 Deliverables ✅

- [data-model.md](data-model.md) - Display entities (TableColumn, TableRow, Table) with rendering algorithm and data mappings
- [contracts/table-interface.md](contracts/table-interface.md) - Table renderer interface specification and integration contracts
- [quickstart.md](quickstart.md) - Developer implementation guide with priority-ordered steps and troubleshooting

## Next Steps

**Ready for task breakdown and implementation.** Run:
```bash
/speckit.tasks
```

This will generate `tasks.md` with concrete implementation steps, dependency chains, and test cases based on the design artifacts created in this plan.
