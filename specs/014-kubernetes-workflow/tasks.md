# Tasks: Kubernetes Workflow

**Input**: Design documents from `/specs/014-kubernetes-workflow/`  
**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/kubernetes-commands.md](contracts/kubernetes-commands.md), [quickstart.md](quickstart.md)

**Tests**: Required by FR-013 and Constitution Principle V. Add focused service, UI-flow, and version-reporting tests before their related implementation tasks.

**Organization**: Tasks are grouped by user story so each increment can be implemented and tested independently.

## Phase 1: Setup

**Purpose**: Establish test fixtures used by the Kubernetes workflow.

- [x] T001 [P] Add representative `container k8s list` fixtures for empty, one-node, multi-node, wrapped-header, and malformed-row output in `src/services/services_test.go`
- [x] T002 [P] Add a reusable Kubernetes command-aware executor fixture for UI flow tests in `src/ui/ui_flow_test.go`

---

## Phase 2: Foundational

**Purpose**: Add shared Kubernetes data, command, parsing, and navigation contracts that block all workflow screens.

**CRITICAL**: Complete this phase before implementing user-story screens.

- [x] T003 Add `KubernetesCluster`, `KubernetesNode`, normalized cluster-state behavior, and command-input types in `src/models/kubernetes_cluster.go`
- [x] T004 Add test coverage for `container k8s list` table parsing, cluster grouping, state normalization, empty output, and malformed rows in `src/services/services_test.go`
- [x] T005 Implement tolerant terminal-table parsing and cluster aggregation in `src/services/kubernetes_parser.go`
- [x] T006 [P] Add the read-only `container k8s list` command builder in `src/services/kubernetes_list_builder.go`
- [x] T007 Add list-builder assertion to the Kubernetes command-contract tests in `src/services/services_test.go`
- [x] T008 Add Kubernetes active-screen constants, selected-cluster message data, and any dedicated form-return messages in `src/ui/messages.go`

**Checkpoint**: Shared models, parsing, command composition, and application routing are ready.

---

## Phase 3: User Story 1 - Open Kubernetes Operations (Priority: P1) MVP

**Goal**: Let a developer open the Kubernetes cluster list with `k`, view cluster/node state, enter the operation submenu, and return safely.

**Independent Test**: From the container list, press `k`, load a fixture-backed cluster list, enter a cluster submenu, traverse its options, and return to the container list without executing a state-changing command.

### Tests for User Story 1

- [x] T009 [US1] Add UI-flow tests for lowercase `k` navigation, cluster-list loading, selected-cluster submenu navigation, and `esc` back paths in `src/ui/ui_flow_test.go`

### Implementation for User Story 1

- [x] T010 [US1] Add lowercase `k` routing and the displayed Kubernetes key hint to the container list in `src/ui/container_list.go`
- [x] T011 [US1] Implement the keyboard-driven cluster/node table, direct list fetch, empty state, error state, refresh, create entry, and container-list back path in `src/ui/kubernetes_cluster_list.go`
- [x] T012 [US1] Implement the selected-cluster action submenu with list/refresh, start, load-image, write-config, delete, create, and back options in `src/ui/kubernetes_cluster_submenu.go`
- [x] T013 [US1] After T011-T012, register Kubernetes list/submenu state, selected-cluster propagation, initialization, updates, views, back-stack behavior, and loading-state reporting in `src/ui/app.go`

**Checkpoint**: The Kubernetes branch opens from `k`, displays parsed local cluster data, and supports navigation without executing operations.

---

## Phase 4: User Story 2 - Safely Manage a Kubernetes Cluster (Priority: P1)

**Goal**: Enable all supported Apple Container 1.4.1 Kubernetes commands using previews, type-to-confirm deletion, result rendering, and dry-run compatibility.

**Independent Test**: With a command-recording executor, create a cluster, preview/start it, load an image, write kubeconfig, and decline deletion; assert the exact command arguments and required approval behavior for every operation.

### Tests for User Story 2

- [x] T014 [P] [US2] Add command-builder validation and exact-argument tests for `create`, `start`, `delete`, `load-image`, and `write-config` in `src/services/services_test.go`
- [x] T015 [US2] Add UI-flow tests for create/load-image/write-config forms, preview acceptance/cancellation, start preview, delete type-to-confirm, result display, and dry-run-compatible command execution in `src/ui/ui_flow_test.go`

### Implementation for User Story 2

- [x] T016 [P] [US2] Implement validated `container k8s create` composition, including name, CPUs, memory, remove-on-stop, and node-image options, in `src/services/kubernetes_builders.go`
- [x] T017 [P] [US2] Implement validated `container k8s start --name <cluster>` composition in `src/services/kubernetes_builders.go`
- [x] T018 [P] [US2] Implement validated `container k8s delete --name <cluster>` composition in `src/services/kubernetes_builders.go`
- [x] T019 [P] [US2] Implement validated `container k8s load-image --name <cluster> <image> [--platform <platform>]` composition in `src/services/kubernetes_builders.go`
- [x] T020 [P] [US2] Implement validated `container k8s write-config --name <cluster> [--kubeconfig <path>]` composition in `src/services/kubernetes_builders.go`
- [x] T021 [US2] Implement the create form, optional resource/image inputs, command preview, result display, and return-to-list behavior in `src/ui/kubernetes_create.go`
- [x] T022 [US2] Implement the load-image form, selected-cluster target, image/platform validation, command preview, result display, and return-to-submenu behavior in `src/ui/kubernetes_load_image.go`
- [x] T023 [US2] Implement the write-config form, selected-cluster target, optional kubeconfig path, command preview, result display, and return-to-submenu behavior in `src/ui/kubernetes_write_config.go`
- [x] T024 [US2] Connect submenu start to preview, delete to `TypeToConfirmModal`, and the three form screens to their builders and command-result handling in `src/ui/kubernetes_cluster_submenu.go`
- [x] T025 [US2] After T021-T023, register the create, load-image, and write-config screen lifecycles and return routes in `src/ui/app.go`

**Checkpoint**: All six documented CLI operations are available from actui with the contracted command safety behavior.

---

## Phase 5: User Story 4 - Identify the Installed Release (Priority: P1)

**Goal**: Build published binaries with the exact GitHub release version and identify local builds as non-release builds.

**Independent Test**: Build once with the default value and once with `-ldflags "-X main.version=0.1.13"`; verify `actui --version` prints `dev` for the former and `0.1.13` for the latter. Review the CI artifact handoff to confirm publishing consumes the injected build tag unchanged.

### Tests for User Story 4

- [x] T026 [US4] Add a Go test covering default and linker-injected `actui --version` formatting in `cmd/actui/main_test.go`
- [x] T027 [US4] Add workflow-content assertions for release-tag computation before build, linker injection, version-artifact upload, and publish consumption in `tests/integration/release_version_workflow_test.go`

### Implementation for User Story 4

- [x] T028 [US4] Change the default CLI version to the documented non-release identifier and keep Cobra version output linker-overridable in `cmd/actui/main.go`
- [x] T029 [US4] Compute the next semantic release tag before compilation, inject its `v`-stripped value through Go linker flags, and upload the canonical tag as a release-version artifact in `.github/workflows/build-binary.yml`
- [x] T030 [US4] Download the release-version artifact and use its tag for idempotency and publication instead of recomputing a tag in `.github/workflows/publish-release.yml`

**Checkpoint**: Local source builds identify as development builds, and every CI release artifact is built and published from the same version value.

---

## Phase 6: User Story 3 - Understand Unsupported or Failed Kubernetes Actions (Priority: P2)

**Goal**: Surface unavailable CLI support, empty lists, partial output, and command failures without trapping the developer or misrepresenting state.

**Independent Test**: Drive the Kubernetes list and submenu with executors that return missing-command errors, empty output, partial table rows, and stderr failures; verify readable outcomes and working back navigation in each case.

### Tests for User Story 3

- [x] T031 [US3] Add UI-flow tests for unavailable Kubernetes commands, empty cluster lists, partial parsed data, failed actions with stderr, and back navigation after each outcome in `src/ui/ui_flow_test.go`

### Implementation for User Story 3

- [x] T032 [US3] Classify command-not-found and unsupported-command execution errors into an actionable Kubernetes-unavailable message in `src/ui/kubernetes_cluster_list.go`, using the executor error and stderr before invoking `src/services/kubernetes_parser.go`
- [x] T033 [US3] Render distinct unavailable, empty, partial-data, and command-failure states while preserving refresh and `esc` navigation in `src/ui/kubernetes_cluster_list.go`
- [x] T034 [US3] Render failed action names and formatted stdout/stderr results while retaining submenu navigation in `src/ui/kubernetes_cluster_submenu.go`

**Checkpoint**: Failure and compatibility paths remain clear, safe, and navigable.

---

## Phase 7: User Story 5 - Find Current Kubernetes Guidance (Priority: P2)

**Goal**: Document the Kubernetes workflow and release-version behavior for users and contributors.

**Independent Test**: Use only repository documentation to find the `k` entry point, identify all six supported commands and safeguards, and explain why a published binary version matches its release label.

### Tests for User Story 5

- [x] T035 [US5] Add documentation-link and required-keyword assertions for Kubernetes and release-version guidance in `tests/integration/documentation_test.go`

### Implementation for User Story 5

- [x] T036 [US5] Document Kubernetes availability, `k` navigation, supported operations, safety confirmations, dry-run behavior, and development/release version output in `README.md`
- [x] T037 [US5] Add Kubernetes screen graph, ownership table, `k` entry point, command-service dependencies, and change recipe guidance in `docs/ai-menu-map.md`
- [x] T038 [US5] Document the Kubernetes workflow, key bindings, validation expectations, empty/error behavior, and safeguards in `docs/user-guide.md`
- [x] T039 [US5] Update version-labeling and build-to-publish artifact handoff documentation in `docs/binary-build-automation.md`

**Checkpoint**: User-facing and contributor documentation agree with the delivered workflow and CI version policy.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Verify the complete feature against its command contract, quality gates, and macOS runtime behavior.

- [x] T040 Run formatting on Kubernetes Go source and tests with `gofmt` for `src/models/kubernetes_cluster.go`, `src/services/kubernetes_*.go`, `src/ui/kubernetes_*.go`, `src/services/services_test.go`, `src/ui/ui_flow_test.go`, and `cmd/actui/main_test.go`
- [x] T041 Run focused service, UI, CLI-version, workflow, and documentation tests with `go test ./src/services ./src/ui ./cmd/actui ./tests/integration`
- [x] T042 Run the complete automated suite with `go test ./...`
- [ ] T043 Perform the macOS 26.x Apple Silicon manual quickstart, including `container k8s --help`, cluster list, every preview, dry-run, deletion cancellation, and release-version comparison, following `specs/014-kubernetes-workflow/quickstart.md`
- [x] T044 Re-check the constitution table and command mapping against the implementation, then record any remaining manual-validation evidence in `specs/014-kubernetes-workflow/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts immediately.
- **Foundational (Phase 2)**: Depends on T001-T002 and blocks Kubernetes story phases.
- **User Story 1 (Phase 3)**: Depends on T003-T008.
- **User Story 2 (Phase 4)**: Depends on the navigation and selection behavior from User Story 1.
- **User Story 4 (Phase 5)**: Depends only on setup; it can run in parallel with Kubernetes implementation after T001-T002.
- **User Story 3 (Phase 6)**: Depends on User Stories 1 and 2 so it can enhance actual list/action surfaces.
- **User Story 5 (Phase 7)**: Depends on the final navigation, command, and version behavior from User Stories 1, 2, and 4.
- **Polish (Phase 8)**: Depends on all desired story phases.

### User Story Dependencies

```text
Setup -> Foundational -> US1 -> US2 -> US3 -> US5
Setup ------------------------------> US4 ---^  ^
```

- **US1** provides the Kubernetes route, cluster selection, and operation entry points.
- **US2** adds command execution safely to the US1 controls.
- **US4** is independent of the Kubernetes UI and can be implemented in parallel after setup.
- **US3** hardens the delivered Kubernetes UI and operations.
- **US5** documents the final behavior of US1, US2, and US4.

### Parallel Opportunities

- T002 and T003 can proceed in parallel.
- T007 can be implemented while T005-T006 establish parser behavior.
- T016 can start independently of UI tests in T017.
- T018-T022 are independent builder files and can proceed in parallel after T016 defines their expected contracts.
- T026-T030 (versioning) can be assigned separately from the Kubernetes workflow after setup.
- T036-T039 edit distinct documentation files and can proceed in parallel once feature behavior stabilizes.

## Parallel Example: User Story 2

```text
Task: "Implement validated container k8s create composition in src/services/kubernetes_create_builder.go"
Task: "Implement validated container k8s start composition in src/services/kubernetes_start_builder.go"
Task: "Implement validated container k8s delete composition in src/services/kubernetes_delete_builder.go"
Task: "Implement validated container k8s load-image composition in src/services/kubernetes_load_image_builder.go"
Task: "Implement validated container k8s write-config composition in src/services/kubernetes_write_config_builder.go"
```

## Implementation Strategy

### MVP First

1. Complete Setup and Foundational work.
2. Implement User Story 1 to deliver the `k` entry point, list, submenu, and back navigation.
3. Run T011 and manually confirm the navigation path.
4. Demonstrate the non-mutating Kubernetes discovery workflow before adding command execution.

### Incremental Delivery

1. Add US1 for discoverable Kubernetes navigation and status display.
2. Add US2 for the supported, guarded cluster commands.
3. Deliver US4 independently to fix binary-version identity as soon as CI handoff is ready.
4. Add US3 failure handling and US5 documentation after the behaviors are stable.
5. Complete cross-cutting validation before release.

## Notes

- Every task uses the required checkbox, sequential task ID, optional `[P]` marker only for independently editable work, and `[US#]` label for user-story work.
- The plan deliberately excludes undocumented `stop` and `inspect` operations because the installed Apple Container 1.4.1 CLI does not provide them.
- Do not add a Kubernetes screen without updating both `src/ui/messages.go` and `src/ui/app.go`, then refresh `docs/ai-menu-map.md`.
