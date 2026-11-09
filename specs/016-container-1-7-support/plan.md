# Implementation Plan: Container 1.7 Support

**Branch**: `016-container-1-7-support` | **Date**: 2026-11-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/016-container-1-7-support/spec.md` (update actui from Apple Container 1.6.1 to 1.7.1 + refresh README and internal docs, macOS 26 -> 27 prerequisite)

## Summary

Apple Container 1.7.1 makes three breaking changes to the `container cluster` family actui wraps (`load-image` removed in favor of `--shared-image-store` on create; `list` one-row-per-cluster with `--nodes` for per-node detail; lifecycle commands refuse cluster nodes) and adds worker resource sizing (`--workers count=<n>,cpus=<n>,memory=<size>`). The plan: delete the load-image builder/screen/submenu option/app wiring/tests; add a shared-image-store toggle plus extended workers validation to the create builder and form; extend the list parser to accept both the aggregate and per-node table shapes; update all builder/UI-flow/docs tests; then sweep README.md, `docs/user-guide.md`, `docs/ai-menu-map.md`, `docs/binary-build-automation.md`, copilot instructions, and the constitution (26.x -> 27.x) so no stale version or load-image references remain. All other flows (containers, images, machines, registries, daemon) verified unchanged against 1.7.1 help output. See [research.md](./research.md), [data-model.md](./data-model.md), [contracts/commands.md](./contracts/commands.md), [quickstart.md](./quickstart.md).

## Technical Context

**Language/Version**: Go 1.24.2
**Primary Dependencies**: bubbletea v1.3.10, bubbles v1.0.0, lipgloss v1.1.0 (TUI); cobra v1.10.2, viper v1.21.0 (CLI/config)
**Storage**: N/A (local TOML config in `~/Library/Application Support/actui/`, JSONL command logs with rotation)
**Testing**: `go test ./...` (existing `src/services/services_test.go`, `src/models/models_test.go`, `src/ui/ui_flow_test.go`, `src/ui/ui_additional_test.go`, `tests/integration/documentation_test.go`); `go vet ./...`; `gofmt -l`
**Target Platform**: macOS 27.x on Apple Silicon, Apple Container CLI 1.7.1 backend
**Project Type**: CLI terminal UI (single binary `actui`)
**Performance Goals**: Interactive TUI responsiveness (no new perf-sensitive paths; all new work is command-string composition plus form inputs and parser shape detection)
**Constraints**: Local-only, no telemetry; every state-changing action MUST show a command preview and require approval; destructive removals keep type-to-confirm; all operations honor `--dry-run`; single-version target (1.7.1, no `load-image` fallback)
**Scale/Scope**: ~100 Go files across `src/models`, `src/services`, `src/ui`; 4 spec user stories; docs updates to README + user guide + AI menu map + build automation docs + copilot instructions + constitution

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **I. Command-Safe TUI** — Every changed action (cluster create/list/delete/write-config) keeps its visible command preview with explicit approval; cluster delete keeps type-to-confirm; dry-run shows the command without executing. Load-image deletion removes a dead path; no direct-execution paths added.
- [x] **II. macOS 27.x + Apple Silicon** — Platform target amended 26.x -> 27.x per FR-011 (spec assumption + verified `sw_vers` 27.2 + upstream macOS 27 runners); constitution file updated with rationale and migration notes as part of this change.
- [x] **III. Local-Only** — No network calls added; all state stays in user-space config/logs. No telemetry.
- [x] **IV. Clear Observability** — Success/failure/stdout/stderr surface through existing result rendering; cluster list parser keeps its tolerant unknown-state fallback across both table shapes; partial-failure CLI errors shown readably, never swallowed.
- [x] **V. Tested Command Contracts** — Builder tests for removed load-image (assert absence), shared-image-store + worker-sizing composition, dual-shape list parsing; UI-flow tests for submenu options, create preview, and validation messages; docs test strings updated.
- [x] **Workflow gates** — Explicit command mapping in [contracts/commands.md](./contracts/commands.md); destructive actions keep dry-run + confirm copy; manual verification on macOS 27.x per [quickstart.md](./quickstart.md); semver release via existing automation.

Post-design re-check: no new violations introduced. The load-image removal *reduces* dead-code surface (also drops one screen constant + app wiring). No complexity exceptions needed.

## Project Structure

### Documentation (this feature)

```text
specs/016-container-1-7-support/
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
│   └── kubernetes_cluster.go   # remove KubernetesLoadImageInput (load-image deleted)
├── services/
│   ├── kubernetes_builders.go  # delete LoadImageBuilder; +SharedImageStore/+worker-spec on Create
│   ├── kubernetes_parser.go    # accept aggregate + per-node list shapes
│   └── services_test.go        # updated expectations
├── ui/
│   ├── kubernetes_cluster_submenu.go  # remove load-image option
│   ├── kubernetes_create.go    # +shared-image-store toggle + worker-spec field
│   ├── kubernetes_load_image.go     # DELETED (screen + inputs + preview)
│   ├── kubernetes_cluster_list.go   # unchanged (renders parsed aggregates)
│   ├── ui_flow_test.go / ui_additional_test.go  # updated flows
│   ├── messages.go         # remove ScreenKubernetesLoadImage
│   └── app.go              # remove load-image wiring
cmd/actui/                      # unchanged entry points
tests/integration/
└── documentation_test.go  # required/forbidden strings -> 1.7.1
docs/
├── ai-menu-map.md              # cluster ownership, no load-image, create opts
├── user-guide.md               # 1.7.1 workflows
└── binary-build-automation.md  # 26.x -> 27.x verification refs
.github/agents/
└── copilot-instructions.md     # 26.x validation note -> 27.x
.specify/memory/
└── constitution.md             # II. 26.x -> 27.x + overview touch-ups
README.md                       # 1.7.1 backend, macOS 27 prereq, cluster changes
```

**Structure Decision**: Single-project layout (existing `src/{models,services,ui}` + `cmd/actui` + `docs/`); one screen file deleted following the reverse of the new-screen recipe (constant + state + wiring + tests + map updated together).

## Complexity Tracking

No constitution violations; this section intentionally left empty.
