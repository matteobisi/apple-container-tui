# Implementation Plan: Container 1.6 Support

**Branch**: `015-container-1-6-support` | **Date**: 2026-11-02 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/015-container-1-6-support/spec.md` (update actui from Apple Container 1.5.0 to 1.6.1 + refresh README and internal docs)

## Summary

Apple Container 1.6.0/1.6.1 breaks actui's Kubernetes workflow (`container k8s` renamed to `container cluster`, `start` removed) and adds new capabilities actui should expose. The plan: rename all five Kubernetes builders from `k8s` to `cluster` and delete the start builder; remove the start option from the cluster submenu and its tests; extend the create form with Workers + CNI and the load-image form with node targets; add a `container commit` builder plus a commit action in the container submenu; extend machine creation with CPUs/memory/home-mount; update all builder/UI-flow tests; then update README.md, `docs/user-guide.md`, `docs/ai-menu-map.md`, and stale code comments so no `k8s`/start references remain. Parser logic is unchanged (verified identical table shape on live 1.6.1 output). See [research.md](./research.md), [data-model.md](./data-model.md), [contracts/commands.md](./contracts/commands.md), [quickstart.md](./quickstart.md).

## Technical Context

**Language/Version**: Go 1.24.2
**Primary Dependencies**: bubbletea v1.3.10, bubbles v1.0.0, lipgloss v1.1.0 (TUI); cobra v1.10.2, viper v1.21.0 (CLI/config)
**Storage**: N/A (local TOML config in `~/Library/Application Support/actui/`, JSONL command logs with rotation)
**Testing**: `go test ./...` (existing `src/services/services_test.go`, `src/models/models_test.go`, `src/ui/ui_flow_test.go`, `src/ui/ui_additional_test.go`); `go vet ./...`; `gofmt -l`
**Target Platform**: macOS 26.x on Apple Silicon, Apple Container CLI 1.6.1 backend
**Project Type**: CLI terminal UI (single binary `actui`)
**Performance Goals**: Interactive TUI responsiveness (no new perf-sensitive paths; all new work is command-string composition plus form inputs)
**Constraints**: Local-only, no telemetry; every state-changing action MUST show a command preview and require approval; destructive removals keep type-to-confirm; all operations honor `--dry-run`; single-version target (1.6.1, no `k8s` fallback)
**Scale/Scope**: ~100 Go files across `src/models`, `src/services`, `src/ui`; 5 spec user stories; docs updates to README + user guide + AI menu map

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **I. Command-Safe TUI** — Every new/changed action (cluster create/delete/load-image/write-config, container commit, machine create) gets a visible command preview with explicit approval; cluster/machine delete keeps type-to-confirm; dry-run shows the command without executing. No direct-execution paths.
- [x] **II. macOS 26.x + Apple Silicon** — No platform change; backend version bump only.
- [x] **III. Local-Only** — No network calls added; all state stays in user-space config/logs. No telemetry.
- [x] **IV. Clear Observability** — Success/failure/stdout/stderr surface through existing result rendering; cluster list parser keeps its tolerant unknown-state fallback; partial-failure CLI errors are shown readably, never swallowed.
- [x] **V. Tested Command Contracts** — Builder tests for renamed `cluster` commands, removed start, new workers/CNI/node/commit/machine-create options; UI-flow tests for submenu options, create/load/commit previews, and validation messages.
- [x] **Workflow gates** — Explicit command mapping in [contracts/commands.md](./contracts/commands.md); destructive actions keep dry-run + confirm copy; manual verification on macOS 26.x per [quickstart.md](./quickstart.md); semver release via existing automation.

Post-design re-check: no new violations introduced. The start-builder removal *reduces* dead-code surface. No complexity exceptions needed.

## Project Structure

### Documentation (this feature)

```text
specs/015-container-1-6-support/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── commands.md      # Phase 1 output (exact CLI mapping)
├── checklists/
│   └── requirements.md  # Spec quality gate (all pass)
└── tasks.md             # Phase 2 output (/speckit.tasks - NOT created here)
```

### Source Code (repository root)

```text
src/
├── models/
│   └── kubernetes_cluster.go   # unchanged (table shape identical)
├── services/
│   ├── kubernetes_builders.go  # k8s -> cluster; +workers/+cni; +nodes; delete start
│   ├── kubernetes_parser.go    # comment update only
│   ├── commit_container_builder.go  # NEW: container commit
│   ├── machine_create_builder.go    # +cpus/+memory/+home-mount
│   └── services_test.go        # updated expectations
├── ui/
│   ├── container_submenu.go    # +commit action (preview flow)
│   ├── kubernetes_cluster_submenu.go  # remove start option
│   ├── kubernetes_create.go    # +workers/+cni fields
│   ├── kubernetes_load_image.go     # +nodes field
│   ├── machine_create.go       # +cpus/+memory/+home-mount fields
│   ├── ui_flow_test.go / ui_additional_test.go  # updated flows
│   ├── messages.go / app.go    # only if commit needs a screen (direct preview: no)
│   └── help.go                 # version/backend text if applicable
cmd/actui/                      # unchanged entry points
docs/
├── ai-menu-map.md              # cluster ownership, no start, commit action
└── user-guide.md               # 1.6.1 workflows
README.md                       # 1.6.1 backend, cluster rename, commit, workers/CNI, machine opts
```

**Structure Decision**: Single-project layout (existing `src/{models,services,ui}` + `cmd/actui` + `docs/`); no new top-level directories except the NEW single builder file following the existing one-builder-per-file convention.

## Complexity Tracking

No constitution violations; this section intentionally left empty.
