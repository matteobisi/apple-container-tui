# Implementation Plan: Enhanced Menu Navigation and Image Management

**Branch**: `002-refactor-menu-images` | **Date**: 2026-02-16 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/002-refactor-menu-images/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This feature refactors the TUI menu structure to improve navigation and adds comprehensive image management capabilities. The primary changes include: (1) replacing Enter-to-toggle with Enter-to-submenu for containers, enabling contextual actions like start/stop, tail logs, and enter shell; (2) adding a new image list view accessible via 'i' key with pull/build/prune operations; (3) implementing image submenu for inspect and delete; (4) ensuring consistent arrow-key navigation across all submenus. This enhances discoverability and reduces cognitive load by organizing operations into logical navigation hierarchies rather than requiring users to memorize keyboard shortcuts.

## Technical Context

**Language/Version**: Go 1.21  
**Primary Dependencies**: Bubbletea v1.2.4 (TUI framework), Bubbles v0.20.0 (TUI components), Cobra v1.8.1 (CLI), Viper v1.19.0 (config)  
**Storage**: Local filesystem for logs and config (~/Library/Application Support/actui/), JSONL command logs  
**Testing**: Go standard testing (go test), contract tests in tests/contract/, integration tests in tests/integration/, unit tests in tests/unit/  
**Target Platform**: macOS 26.x on Apple Silicon (M-series)  
**Project Type**: Single project - terminal UI application  
**Performance Goals**: <2 seconds to display image lists (up to 100 images), <100ms log streaming latency, instant submenu navigation  
**Constraints**: No network dependencies beyond Apple Container CLI, keyboard-only navigation, must work in standard terminal emulators  
**Scale/Scope**: ~3,000 LOC currently, adding ~1,500 LOC for menu refactoring and image management features

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Command-Safe TUI
✅ **PASS** - All new operations (container logs -f, container exec -it, container image list, container image inspect, container image rm, container image prune) map directly to Apple Container CLI commands. Image prune uses type-to-confirm per destructive action requirement. Container submenu shows command preview before execution for start/stop actions (existing pattern). Shell detection and log streaming execute standard CLI commands with proper error handling.

### Principle II: macOS 26.x + Apple Silicon Target
✅ **PASS** - Feature is macOS-specific TUI enhancement. No cross-platform considerations introduced. All container commands target Apple Container CLI on macOS 26.x.

### Principle III: Local-Only Operation  
✅ **PASS** - No network dependencies added. All operations are local container/image management. Image list, inspect, prune, and container shell access are local-only operations using local daemon.

### Principle IV: Clear Observability
✅ **PASS** - All operations surface stdout/stderr (FR-027 pipes to jq, FR-009 streams logs, FR-033 shows shell detection errors, FR-034 shows log interruption messages). Error states explicitly defined (shell failure, log interruption, image in use, empty lists). Status visibility maintained through submenu context (running/stopped containers).

### Principle V: Tested Command Contracts
✅ **PASS** - Feature requires contract tests for new command builders (image list, image inspect, image prune, image delete, container logs streaming, container exec shell detection). Integration tests needed for submenu navigation flows. Unit tests for shell detection logic and navigation state management.

**Overall Assessment**: ✅ **COMPLIANT** - No constitution violations. All operations follow established command-preview-confirm pattern. Local-only, observable, testable.

---

## Post-Phase 1 Constitution Re-Check

*Re-evaluation after completing research, data model, and contracts*

### Principle I: Command-Safe TUI
✅ **PASS** - All command contracts documented (see contracts/). Every operation maps to Apple Container CLI command:
- container logs -f → container-submenu.md §1
- container exec -it → container-submenu.md §2 (with shell detection)
- container image list → image-list.md §1
- container image prune → image-list.md §4 (with type-to-confirm)
- container image inspect | jq → image-submenu.md §1
- container image rm → image-submenu.md §2 (with type-to-confirm)

Command preview maintained for start/stop (existing). Destructive operations (prune, delete) use type-to-confirm safety guard.

### Principle II: macOS 26.x + Apple Silicon Target
✅ **PASS** - No platform-specific changes in design. TUI enhancement works on existing platform. Shell detection sequence (bash/sh/ash) covers standard macOS containers.

### Principle III: Local-Only Operation
✅ **PASS** - All data sources local (container daemon, local image registry). No network calls except image pull/build (existing features). Navigation state stored in-memory only. Shell detection cache in-memory, not persisted.

### Principle IV: Clear Observability
✅ **PASS** - Observability enhanced:
- Log streaming shows real-time container output (FR-009, research.md §3)
- Image inspection displays full JSON metadata (FR-027, contracts/image-submenu.md §1)
- Shell detection errors clearly messaged (FR-033, contracts/container-submenu.md §2)
- All error scenarios documented in contracts with specific error messages
- Log interruption handled gracefully with user notification (FR-034, Clarification Q4)

### Principle V: Tested Command Contracts
✅ **PASS** - Test strategy defined in data-model.md and all contract files:
- Contract tests for 6 new command builders (logs, exec, image list/inspect/prune/delete)
- Unit tests for shell detection logic, image list parser, navigation state management
- Integration tests for full navigation flows (15 test scenarios across 3 contracts)
- Test fixtures for CLI output parsing (data-model.md Testing Strategy)

**Post-Design Assessment**: ✅ **COMPLIANT** - Design deepens compliance. All principles satisfied with concrete implementation details. No new concerns introduced during detailed design phase.

## Project Structure

### Documentation (this feature)

```text
specs/002-refactor-menu-images/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   ├── container-submenu.md
│   ├── image-list.md
│   └── image-submenu.md
├── checklists/
│   └── requirements.md  # Already created during clarification
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Single project structure (existing)
cmd/
└── actui/
    └── main.go          # CLI entry point

src/
├── models/              # Data structures
│   ├── container.go     # Existing
│   ├── image.go         # NEW: Image entity for image list
│   └── ...
├── services/            # Business logic
│   ├── container_parser.go                    # Existing
│   ├── image_list_builder.go                  # NEW: Build image list command
│   ├── image_inspect_builder.go               # NEW: Build image inspect command
│   ├── image_delete_builder.go                # NEW: Build image delete command
│   ├── image_prune_builder.go                 # NEW: Build image prune command
│   ├── container_logs_builder.go              # NEW: Build container logs -f command
│   ├── container_exec_builder.go              # NEW: Build container exec command with shell detection
│   ├── shell_detector.go                      # NEW: Shell detection logic
│   └── ...
└── ui/                  # TUI components
    ├── app.go                                  # Main application model (UPDATE: add navigation state)
    ├── container_list.go                       # Existing (UPDATE: change Enter behavior)
    ├── container_submenu.go                    # NEW: Container action submenu
    ├── container_logs.go                       # NEW: Live log streaming view
    ├── container_shell.go                      # NEW: Interactive shell session wrapper
    ├── image_list.go                           # NEW: Image list view
    ├── image_submenu.go                        # NEW: Image action submenu
    ├── image_inspect.go                        # NEW: Image inspection view
    ├── keys.go                                 # Existing (UPDATE: add image list key)
    └── ...

tests/
├── contract/            # Command generation tests
│   ├── image_list_test.go                     # NEW
│   ├── image_inspect_test.go                  # NEW
│   ├── image_prune_test.go                    # NEW
│   ├── container_logs_test.go                 # NEW
│   └── container_exec_test.go                 # NEW
├── integration/         # End-to-end flow tests
│   ├── image_operations_test.go               # Existing (UPDATE: add inspect/delete/prune)
│   └── navigation_flows_test.go               # NEW: Submenu navigation tests
└── unit/
    ├── shell_detector_test.go                 # NEW
    └── ...
```

**Structure Decision**: Using existing single-project Go structure. New files added to src/models/ (Image entity), src/services/ (6 new command builders + shell detector), and src/ui/ (6 new views + updates to existing). Tests follow established pattern: contract tests for command builders, integration tests for workflows, unit tests for business logic. Total additions: ~12 new files, ~3 file updates.

## Complexity Tracking

No constitution violations. All requirements align with established principles.

---

## Planning Phase Completion Summary

**Date Completed**: 2026-02-16  
**Command**: `/speckit.plan`  
**Status**: ✅ **COMPLETE**

### Artifacts Generated

| Phase | Artifact | Location | Lines | Status |
|-------|----------|----------|-------|--------|
| **Setup** | Implementation Plan | plan.md | 159 | ✅ Complete |
| **Phase 0** | Technical Research | research.md | 318 | ✅ Complete |
| **Phase 1** | Data Model | data-model.md | 347 | ✅ Complete |
| **Phase 1** | Container Submenu Contract | contracts/container-submenu.md | 275 | ✅ Complete |
| **Phase 1** | Image List Contract | contracts/image-list.md | 288 | ✅ Complete |
| **Phase 1** | Image Submenu Contract | contracts/image-submenu.md | 265 | ✅ Complete |
| **Phase 1** | Developer Quickstart | quickstart.md | 429 | ✅ Complete |
| **Phase 1** | Agent Context Update | .github/agents/copilot-instructions.md | Updated | ✅ Complete |

**Total Specification Size**: 2,081 lines of detailed planning documentation

### Decisions Recorded

**Research Decisions** (6):
1. Navigation state management using Bubbletea nested models with stack
2. Shell detection via sequential probing with caching
3. Log streaming via background goroutine with tea.Cmd pattern
4. Image list parsing using column-based text parsing
5. JSON display using bubbles viewport component
6. Type-to-confirm reuse from existing component

**Data Model Entities** (6):
1. Image - Container image metadata
2. Navigation State - View stack and selected items
3. Shell Detection Result - Shell availability cache
4. Container Submenu Option - Context-sensitive menu items
5. Image Submenu Option - Image action menu items
6. Log Stream State - Active streaming session state

**Command Contracts** (6 new commands):
1. `container logs -f` - Live log streaming
2. `container exec -it` - Interactive shell with auto-detection
3. `container image list` - Local image enumeration
4. `container image inspect | jq` - JSON metadata display
5. `container image prune` - Unused image cleanup
6. `container image rm` - Single image deletion

### Constitution Compliance

**Pre-Planning Check**: ✅ ALL PRINCIPLES PASS  
**Post-Design Check**: ✅ ALL PRINCIPLES PASS (compliance deepened)

No violations. No exceptions required. Feature aligns fully with Container TUI Constitution v0.2.0.

### Next Steps

**Ready for Phase 2**: Task decomposition via `/speckit.tasks` command

The planning phase is complete. All technical unknowns resolved, all entities modeled, all command contracts defined, developer quickstart guide created. Implementation can now begin with clear specifications and test-driven approach.

**Estimated Implementation Scope**:
- New files: 12
- Modified files: 3
- New LOC: ~1,500
- Test cases: 15+ integration scenarios, 10+ contract tests, 8+ unit tests
- Total tasks: 98 (includes 1 task added post-analysis for shell exit handling)

---

**Plan Complete** | **Branch**: 002-refactor-menu-images | **Next**: `/speckit.tasks`
