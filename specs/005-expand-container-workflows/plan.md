# Implementation Plan: Expanded Container Workflows

**Branch**: `005-expand-container-workflows` | **Date**: 2026-03-31 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/005-expand-container-workflows/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Add four targeted workflow extensions without disturbing the current TUI architecture: a new Registries screen backed by `container registry list --format json`, a stopped-container export flow that previews and executes a narrow CLI command sequence, a visible build-form checkbox for `--pull` defaulted to enabled, and structured daemon status parsing through `container system status --format json`. The design preserves the existing screen-routing, builder/parser, preview-modal, and executor patterns, introducing only small additive models and services where the current seams are too narrow.

## Technical Context

**Language/Version**: Go 1.21+ (module target), validated in local Apple Container environment on macOS 26.x  
**Primary Dependencies**: Bubble Tea v1.2.4, Bubbles v0.20.0, Lipgloss v1.0.0, Cobra v1.8.1, Viper v1.19.0, Go standard library `encoding/json`  
**Storage**: Local filesystem only for exported OCI tar archives; no new persistent app-owned state required  
**Testing**: `go test ./...` for builders/parsers/UI flow tests plus manual verification against local Apple Container CLI  
**Target Platform**: macOS 26.x on Apple Silicon with Apple Container CLI installed  
**Project Type**: Single-project CLI TUI application  
**Performance Goals**: Screen transitions remain effectively instant; registry/status refresh complete in under 1 second on a healthy local runtime; command preview renders immediately before execution  
**Constraints**: Must preserve command preview visibility, stay local-only, use exact Apple Container CLI commands, avoid broad navigation/executor refactors, and keep export available only for stopped containers  
**Scale/Scope**: 2 new UI screens, 4-6 focused service/model additions or extensions, and targeted updates across existing UI/service tests

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Command-Safe TUI
✅ **COMPLIANT** - Every new user action maps directly to Apple Container CLI commands. Registries and daemon status are read-only queries. Build preview continues to show the exact `container build` invocation. Container export will preview the exact command sequence before execution rather than hiding orchestration behind an opaque action.

### Principle II: macOS 26.x + Apple Silicon Target
✅ **COMPLIANT** - All planned commands (`container registry list`, `container export`, `container image save`, `container system status --format json`) are local Apple Container CLI commands on the supported macOS target.

### Principle III: Local-Only Operation
✅ **COMPLIANT** - No telemetry, remote APIs, or background services are introduced. Registry data comes from local CLI output, export writes to a user-selected local directory, and daemon status stays local.

### Principle IV: Clear Observability
✅ **COMPLIANT** - Registry and daemon states will show explicit empty/error/unknown outcomes. Export will surface command-sequence progress and report save success separately from cleanup warnings. Existing stdout/stderr result rendering remains the primary feedback path.

### Principle V: Tested Command Contracts
✅ **COMPLIANT** - Plan includes automated tests for new command builders, JSON parsers, stopped-container export guardrails, build `--pull` composition, and screen navigation hooks. Manual verification remains required for live Apple Container execution.

### Workflow and Quality Gates
⚠️ **DOCUMENTED EXCEPTION** - The feature spec template in this repository does not currently include a command mapping table, while the constitution requires one at spec level. This plan mitigates the gap by creating `contracts/command-workflows.md` as the authoritative mapping artifact for this feature and keeping implementation aligned to that contract.

**GATE RESULT**: ✅ **PASS WITH DOCUMENTED MITIGATION** - No constitutional blocker remains after capturing command mappings in Phase 1 artifacts.

**Post-Phase 1 Re-check**: Still compliant. The design uses additive screens/builders/parsers, preserves command previews, and records the required command mapping in `contracts/command-workflows.md`.

## Project Structure

### Documentation (this feature)

```text
specs/005-expand-container-workflows/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── command-workflows.md
└── tasks.md
```

### Source Code (repository root)

```text
src/models/
├── command.go                  # Existing preview formatting primitive
├── daemon.go                   # Extend daemon state representation for unknown/structured metadata
├── image_reference.go          # Existing image naming utilities reused for generated export refs
└── registry_login.go           # NEW registry-list domain model

src/services/
├── build_image_builder.go      # Add PullLatest support -> --pull flag composition
├── check_daemon_builder.go     # Switch status command to --format json
├── daemon_parser.go            # Parse structured JSON with unknown fallback
├── image_delete_builder.go     # Reuse for best-effort export cleanup
├── registry_list_builder.go    # NEW builder for container registry list --format json
├── registry_parser.go          # NEW parser for registry list JSON
├── export_container_builder.go # NEW builder for container export --image
├── image_save_builder.go       # NEW builder for container image save --output
└── export_workflow.go          # NEW narrow sequential workflow orchestration for export

src/ui/
├── app.go                      # Register new screens and preserve routing behavior
├── messages.go                 # Add new ActiveScreen constants and navigation messages
├── image_list.go               # Add Registries entry point
├── container_submenu.go        # Add stopped-container export action
├── build.go                    # Add pull checkbox/default and preview visibility
├── daemon_control.go           # Display unknown/structured daemon status
├── registries.go               # NEW registries list screen
└── container_export.go         # NEW export destination / workflow screen

src/services/services_test.go   # Extend command builder/parser coverage
src/ui/ui_additional_test.go    # Extend navigation/build/daemon/export UI coverage
```

**Structure Decision**: Keep the existing single-project layout and extend it only at the same seams already used elsewhere in the repo: `ActiveScreen`-based routing in `src/ui`, command construction/parsing in `src/services`, and lightweight domain structs in `src/models`. The only new abstraction is a narrow export workflow service because the CLI export flow requires multiple exact commands.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Command mapping documented in contract artifact instead of spec body | Current repo spec template and prior feature specs do not contain a command mapping section, but this feature still needs an authoritative command contract | Blocking planning on template drift would halt delivery without improving the implementation; the contract document captures the same governance requirement with exact CLI commands |

