# Implementation Plan: Rename Binary from apple-tui to actui

**Branch**: `001-rename-binary-actui` | **Date**: 2026-02-16 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-rename-binary-actui/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Rename the Go binary output from "apple-tui" to "actui" to better reflect the project's purpose as a general-purpose container TUI. This is a clean break with no backward compatibility - all references in source code, build configuration, documentation, tests, and directory structure must be updated. The Go module name remains unchanged per specification.

## Technical Context

**Language/Version**: Go 1.21
**Primary Dependencies**: 
  - github.com/charmbracelet/bubbletea (TUI framework)
  - github.com/charmbracelet/bubbles (TUI components)
  - github.com/spf13/cobra (CLI framework)
  - github.com/spf13/viper (configuration)
**Storage**: Local filesystem (config.toml, JSONL logs)
**Testing**: Go standard testing (testing package)
**Target Platform**: macOS 26.x (Apple Silicon)
**Project Type**: Single binary CLI/TUI application
**Performance Goals**: N/A (rename operation)
**Constraints**: Must not change Go module name (container-tui), only binary name  
**Scale/Scope**: ~50 source files, comprehensive rename across codebase

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Initial Check (Pre-Phase 0)

| Principle | Impact | Compliance |
|-----------|--------|------------|
| **I. Command-Safe TUI** | No impact - binary name doesn't affect command safety or dry-run behavior | ✅ PASS |
| **II. macOS 26.x + Apple Silicon** | No impact - binary name doesn't affect platform compatibility | ✅ PASS |
| **III. Local-Only Operation** | No impact - binary name doesn't affect operational scope | ✅ PASS |
| **IV. Clear Observability** | No impact - logging and error handling unchanged | ✅ PASS |
| **V. Tested Command Contracts** | Tests must be updated to reference new binary name in documentation | ✅ PASS (with updates) |

**Platform/Runtime Constraints**: No impact - Apple Container CLI backend unchanged  
**Workflow/Quality Gates**: Standard testing applies - all existing tests must pass with updated references

**Overall Status**: ✅ **ALL GATES PASS** - This is a low-risk cosmetic change with no functional impact

### Post-Design Re-check (After Phase 1)

**Re-evaluation completed**: 2026-02-16

| Principle | Post-Design Impact | Compliance |
|-----------|-------------------|------------|
| **I. Command-Safe TUI** | Confirmed: No changes to command safety mechanisms | ✅ PASS |
| **II. macOS 26.x + Apple Silicon** | Confirmed: No platform changes | ✅ PASS |
| **III. Local-Only Operation** | Confirmed: Still local-only, no new dependencies | ✅ PASS |
| **IV. Clear Observability** | Confirmed: Log paths updated but format unchanged | ✅ PASS |
| **V. Tested Command Contracts** | Confirmed: Test paths updated, contracts unchanged | ✅ PASS |

**Design Artifacts Review**:
- research.md: No new technologies or patterns introduced
- data-model.md: No data model changes
- contracts/: File path changes only, no architectural changes
- quickstart.md: Build and test procedures unchanged

**Final Status**: ✅ **ALL GATES PASS** - Constitution compliance maintained throughout design phase

## Project Structure

### Documentation (this feature)

```text
specs/001-rename-binary-actui/
├── spec.md              # Feature specification (completed)
├── plan.md              # This file (/speckit.plan output)
├── research.md          # Phase 0 output (minimal - naming conventions)
├── data-model.md        # Phase 1 output (N/A - no data entities for rename)
├── quickstart.md        # Phase 1 output (build/test instructions)
├── contracts/           # Phase 1 output (file change inventory)
│   └── rename-inventory.md  # List of all files requiring changes
└── tasks.md             # Phase 2 output (/speckit.tasks - not created by /speckit.plan)
```

### Source Code (repository root)

```text
container-tui/              # Root directory
├── cmd/
│   └── apple-tui/          # TO BE RENAMED → actui/
│       └── main.go         # Binary entry point - references need updating
├── src/
│   ├── models/             # Domain models - minimal/no changes
│   ├── services/           # Business logic - minimal/no changes  
│   └── ui/                 # TUI screens - minimal/no changes
├── tests/
│   ├── contract/           # Contract tests - may reference binary name
│   ├── integration/        # Integration tests - may reference binary name
│   └── unit/               # Unit tests - minimal changes
├── docs/
│   └── user-guide.md       # Documentation - needs actui references
├── config/
│   └── default.toml        # Configuration - check for binary name references
├── go.mod                  # Module definition - UNCHANGED per spec
├── Dockerfile              # Build config - likely needs binary name update
└── README.md               # Primary documentation - needs comprehensive updates
```

**Structure Decision**: Single project structure (Go CLI/TUI application). The primary change is renaming `cmd/apple-tui/` to `cmd/actui/` and updating all references throughout the codebase. No architectural changes required.

## Complexity Tracking

**No violations detected** - All constitution checks pass. This is a straightforward rename operation with no architectural complexity.

---

## Planning Summary

**Status**: ✅ **COMPLETE** - Ready for implementation via `/speckit.tasks`

### Artifacts Generated

| Artifact | Status | Location |
|----------|--------|----------|
| **Implementation Plan** | ✅ Complete | [plan.md](./plan.md) |
| **Research** | ✅ Complete | [research.md](./research.md) |
| **Data Model** | ✅ Complete (N/A) | [data-model.md](./data-model.md) |
| **Contracts** | ✅ Complete | [contracts/rename-inventory.md](./contracts/rename-inventory.md) |
| **Quickstart** | ✅ Complete | [quickstart.md](./quickstart.md) |
| **Agent Context** | ✅ Updated | `.github/agents/copilot-instructions.md` |

### Key Decisions

1. **Directory Rename**: Use `git mv` for `cmd/apple-tui/` → `cmd/actui/`
2. **Path Updates**: Config and log paths change from `apple-tui` to `actui`
3. **No Migration**: Clean break, users manually migrate config if needed
4. **Module Name**: Unchanged (`container-tui`)
5. **Scope**: 8 files require changes, ~30+ string replacements

### Constitution Compliance

- ✅ All gates passed (initial and post-design)
- ✅ No architectural complexity introduced
- ✅ No new dependencies or technologies
- ✅ Existing test contracts maintained

### Risk Assessment

**Overall Risk**: 🟢 **LOW**

- No functional changes
- Mechanical search-and-replace operation
- Well-defined scope and verification steps
- All existing tests will catch regressions

### Next Steps

**Ready for**: `/speckit.tasks` - Task breakdown and implementation

**Implementation estimate**: 30-60 minutes
- Directory rename: 2 minutes
- Code updates: 10 minutes  
- Documentation updates: 15 minutes
- Testing and verification: 20 minutes
- Commit and cleanup: 5 minutes

**Blocking issues**: None identified
