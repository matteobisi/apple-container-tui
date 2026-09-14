# Implementation Plan: Kubernetes Workflow

**Branch**: `014-kubernetes-workflow` | **Date**: 2026-09-14 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `/specs/014-kubernetes-workflow/spec.md`

## Summary

Add a local Kubernetes cluster-management branch to actui. Lowercase `k` opens a cluster-and-node list sourced from `container k8s list`; `enter` opens a cluster submenu for create, start, delete, load-image, and write-config. Forms collect only parameters the installed Container 1.4.1 CLI requires, while state-changing commands use existing preview and confirmation patterns.

Replace the fixed binary version with a build-time-injected value. CI calculates the next release tag before compilation, passes it to Go through linker flags, and persists the same tag for publishing so every released binary reports its matching GitHub version.

## Technical Context

**Language/Version**: Go 1.24.2  
**Primary Dependencies**: Bubble Tea v1.3.10, Bubbles v1.0.0, Lipgloss v1.1.0, Cobra v1.10.2  
**Storage**: N/A; stateless TUI consuming local Apple Container CLI output  
**Testing**: `go test ./...`; builder/parser tests in `src/services`, UI-flow tests in `src/ui`  
**Target Platform**: macOS 26.x, Apple Silicon, Apple Container 1.4.1; Kubernetes available from Container 1.2.0  
**Project Type**: CLI/TUI desktop application  
**Performance Goals**: Under 100 ms screen rendering; cluster list fetch bounded by CLI time  
**Constraints**: Keyboard-only and local-only; CLI calls through `services.CommandExecutor`; no background work beyond Bubble Tea commands; Kubernetes list has no JSON/format option and is terminal-table output  
**Scale/Scope**: One top-level workflow, five new UI screens, six builders, one parser/model, CI version injection, and documentation changes

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Status |
|-----------|-------|--------|
| I. Command-Safe TUI | Each action is defined in [contracts/kubernetes-commands.md](contracts/kubernetes-commands.md). Create, start, load-image, and write-config use previews; delete uses type-to-confirm; all respect dry-run. | PASS |
| II. macOS 26.x + Apple Silicon Target | The workflow invokes only the local Apple Container Kubernetes CLI on the existing target platform. | PASS |
| III. Local-Only Operation | Local Container-managed clusters only; no remote, cloud, or telemetry capabilities. | PASS |
| IV. Clear Observability | List, preview, result, unavailable, and failure states show readable status/output with working back navigation. | PASS |
| V. Tested Command Contracts | Builders, tolerant table parsing, destructive confirmation, flow navigation, and release/local version reporting have targeted automated tests. | PASS |
| Scope boundaries | Constitution v0.4.0 permits only local Container Kubernetes cluster lifecycle support; remote clusters and workload orchestration remain excluded. | PASS |

**Post-design re-check**: The 2026-09-15 amendment is limited to local cluster management. The command contract, data model, and quickstart meet the remaining constitutional requirements with no exception.

## Project Structure

### Documentation (this feature)

```text
specs/014-kubernetes-workflow/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── kubernetes-commands.md
└── tasks.md              # Created by /speckit.tasks
```

### Source Code (repository root)

```text
cmd/actui/
└── main.go                              # Linker-injectable version default

src/models/
└── kubernetes_cluster.go                 # NEW: cluster and node data

src/services/
├── kubernetes_list_builder.go            # NEW: container k8s list
├── kubernetes_create_builder.go          # NEW: container k8s create
├── kubernetes_start_builder.go           # NEW: container k8s start
├── kubernetes_delete_builder.go          # NEW: container k8s delete
├── kubernetes_load_image_builder.go      # NEW: container k8s load-image
├── kubernetes_write_config_builder.go    # NEW: container k8s write-config
├── kubernetes_parser.go                  # NEW: parse terminal-table output
└── services_test.go                      # UPDATE: builders and parser coverage

src/ui/
├── container_list.go                     # UPDATE: k entry point
├── messages.go                           # UPDATE: Kubernetes screen ids/messages
├── app.go                                # UPDATE: screen registration/routing
├── kubernetes_cluster_list.go            # NEW: cluster/node list
├── kubernetes_cluster_submenu.go         # NEW: cluster actions
├── kubernetes_create.go                  # NEW: create form
├── kubernetes_load_image.go              # NEW: load-image form
├── kubernetes_write_config.go            # NEW: write-config form
└── ui_flow_test.go                       # UPDATE: navigation, form, confirmation flows

.github/workflows/
├── build-binary.yml                      # UPDATE: compute/inject/persist release tag
└── publish-release.yml                   # UPDATE: consume persisted tag

docs/
├── ai-menu-map.md                        # UPDATE: Kubernetes navigation ownership
├── binary-build-automation.md            # UPDATE: version injection flow
└── user-guide.md                         # UPDATE: Kubernetes walkthrough

README.md                                 # UPDATE: Kubernetes support and version behavior
.specify/memory/constitution.md           # UPDATE: permit bounded local Kubernetes support
```

**Structure Decision**: Extend the existing single Go CLI using its established list/submenu/form, command-builder, preview, and confirmation patterns. No new runtime dependency is needed.

## Complexity Tracking

No constitutional exceptions. Constitution v0.4.0 includes the bounded local Kubernetes workflow.
