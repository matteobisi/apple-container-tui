# Tasks: Container 1.6 Support

**Input**: Design documents from `/specs/015-container-1-6-support/` (spec.md, plan.md, research.md, data-model.md, contracts/commands.md, quickstart.md)
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/commands.md, quickstart.md
**Branch**: `015-container-1-6-support`

**Tests**: Spec FR-010 explicitly requires automated checks for renamed cluster commands, removed start action, new create/load-image options, commit composition + guardrails, and extended machine creation. Test-update tasks are included for every story.

**Organization**: Tasks grouped by user story; each story is independently implementable and testable after Foundational completes.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story the task belongs to ([US1]–[US5])
- Exact file paths in every description

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Baseline before changes

- [ ] T001 [P] Record passing baseline: run `go build ./... && go test ./...` on a clean tree and note failures (if any) before editing
- [ ] T002 [P] Confirm backend: run `container --version` (expect 1.6.1) and `container cluster list` and save output shape for parser comparison

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Single-version target confirmation that all stories assume

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T003 Confirm 1.6.1-only target: verify `container k8s --help` fails with unknown command and `container cluster list` header matches the existing parser's expected columns (CLUSTER NODE ROLE STATE CPUS MEMORY ADDR PORTS); record result in research.md if it differs

**Checkpoint**: Foundation ready — 1.6.1-only target confirmed, user stories can begin

---

## Phase 3: User Story 1 - Keep managing local Kubernetes clusters after the rename (Priority: P1) 🎯 MVP

**Goal**: All Kubernetes actions work on Apple Container 1.6.1 via `container cluster`; the dead start action is gone.

**Independent Test**: With container 1.6.1 installed, press `k`, list clusters, create a cluster (defaults), delete it, load an image, write kubeconfig — every step shows a `container cluster ...` preview and reports its result; no start action appears in the submenu.

### Tests for User Story 1

- [ ] T004 [US1] Update builder unit tests in src/services/services_test.go: expect `cluster list`, `cluster create --name dev --cpus 4 --memory 8G --rm --node-image kindest/node:latest`, `cluster delete --name dev`, `cluster load-image --name dev demo:latest --platform linux/arm64`, `cluster write-config --name dev --kubeconfig /tmp/config`; delete the `k8s start` expectation

### Implementation for User Story 1

- [ ] T005 [P] [US1] Rename `k8s` to `cluster` in src/services/kubernetes_builders.go (list, create, delete, load-image, write-config) and delete `KubernetesStartBuilder`
- [ ] T006 [P] [US1] Update stale command comments in src/services/kubernetes_parser.go (`container k8s list` → `container cluster list`) and src/models/kubernetes_cluster.go
- [ ] T007 [US1] Remove the start option and its preview handler from src/ui/kubernetes_cluster_submenu.go (options become refresh / load-image / write-config / delete / create / back)
- [ ] T008 [P] [US1] Change the create-form name placeholder from `k8s-dev` to `cluster-dev` in src/ui/kubernetes_create.go (field list itself is extended in US2)
- [ ] T009 [US1] Update Kubernetes UI flow tests in src/ui/ui_flow_test.go: `cluster list/create/load-image/write-config` preview strings, no `Start cluster` option, error-path message still asserted via the generic unknown-command text

**Checkpoint**: US1 fully functional — Kubernetes workflow works end-to-end on 1.6.1 with defaults; `go test ./...` passes

---

## Phase 4: User Story 2 - Create and use multi-worker clusters (Priority: P1)

**Goal**: Cluster creation supports worker count and CNI choice; image loading supports per-node targets.

**Independent Test**: Create a cluster with workers=2 and CNI=NONE and see `--workers 2 --cni NONE` in the preview; load an image with node targets and see repeated `--node` flags; omitting the new fields behaves exactly like US1 defaults.

### Tests for User Story 2

- [ ] T010 [US2] Add builder unit tests in src/services/services_test.go: create with `--workers 2 --cni NONE` (and empty workers/CNI omitted); load-image with `--node n1 --node n2` (and empty nodes omitted)

### Implementation for User Story 2

- [ ] T011 [US2] Extend `KubernetesCreateBuilder` in src/services/kubernetes_builders.go with `Workers` (non-negative integer, `0`/empty omitted → no flag) and `CNI` (manifest path or `NONE`, empty omitted → no flag) fields plus validation
- [ ] T012 [US2] Extend `KubernetesLoadImageBuilder` in src/services/kubernetes_builders.go with `Nodes []string` (each non-empty trimmed token emitted as a repeated `--node <name>` flag; empty omitted)
- [ ] T013 [US2] Add Workers and CNI inputs to the create form in src/ui/kubernetes_create.go and pass them into `KubernetesCreateBuilder` on preview
- [ ] T014 [US2] Add a Nodes (comma-separated) input to src/ui/kubernetes_load_image.go and pass parsed node tokens into `KubernetesLoadImageBuilder` on preview
- [ ] T015 [US2] Extend Kubernetes UI flow tests in src/ui/ui_flow_test.go: create preview containing `--workers`/`--cni`, load-image preview containing `--node` flags, and default-omission cases

**Checkpoint**: US1 AND US2 both work — single-node defaults unchanged, multi-worker options previewed and executed correctly

---

## Phase 5: User Story 3 - Commit a container filesystem to an image (Priority: P1)

**Goal**: Container submenu offers a guarded commit action that builds `container commit <id> <reference>`.

**Independent Test**: Select a container, choose commit, enter a reference, approve the `container commit <id> <reference>` preview, and see success or a readable error; dry-run creates nothing.

### Tests for User Story 3

- [ ] T016 [P] [US3] Add commit builder unit tests in src/services/services_test.go: `container commit <id> <reference>` exact args, plus empty-id and empty-reference validation errors (API per T017 contract)

### Implementation for User Story 3

- [ ] T017 [US3] Create `CommitContainerBuilder{ContainerID, Reference}` in src/services/commit_container_builder.go building `container commit <id> <reference>` with required-token validation on both fields (new file, follows one-builder-per-file convention)
- [ ] T018 [US3] Add the commit action to src/ui/container_submenu.go: new `Commit container` option for stopped and running containers, image-reference text input following existing form textinput patterns, command preview with explicit approve/cancel, result reporting and dry-run via the shared executor
- [ ] T019 [US3] Add commit UI flow coverage in src/ui/ui_flow_test.go (or src/ui/ui_additional_test.go if the submenu flow lives there): commit preview shows `container commit <id> <reference>`; empty reference shows a validation message before preview

**Checkpoint**: All P1 stories (US1–US3) independently functional — rename, multi-worker, and commit each testable alone

---

## Phase 6: User Story 4 - Create machines with current 1.6 options (Priority: P2)

**Goal**: Machine creation accepts CPUs, memory, and home-mount mode, matching `machine set` and the 1.6.1 CLI.

**Independent Test**: Create a machine with CPUs/memory/home-mount and see all three flags in the preview; invalid CPUs (`0`) or bad home-mount show a validation message before preview; edit-resources flow unchanged.

### Tests for User Story 4

- [ ] T020 [US4] Add machine-create builder unit tests in src/services/services_test.go: `container machine create <img> --name <n> --cpus <c> --memory <m> --home-mount <mode>` with empty optionals omitted, plus invalid-cpus and invalid-home-mount validation errors

### Implementation for User Story 4

- [ ] T021 [US4] Extend `MachineCreateBuilder` in src/services/machine_create_builder.go with `CPUs` (positive integer when provided), `Memory` (non-empty token when provided), `HomeMount` (`ro|rw|none` when provided) fields; omit empty optionals from args
- [ ] T022 [US4] Add CPUs, Memory, and HomeMount inputs to the create form in src/ui/machine_create.go and pass them into `MachineCreateBuilder` on preview
- [ ] T023 [US4] Extend machine UI flow coverage in src/ui/ui_flow_test.go (or src/ui/ui_additional_test.go): create preview carries the new flags; invalid input shows a validation message before preview

**Checkpoint**: US4 works alongside US1–US3 — machine create/edit consistent with the installed CLI

---

## Phase 7: User Story 5 - Find current 1.6 guidance in README and docs (Priority: P2)

**Goal**: README and internal guides describe the 1.6.1 backend, renamed cluster commands, removed start, commit, multi-worker options, and machine options with no contradictions.

**Independent Test**: Using only the updated docs, describe the supported Container version, the `k` entry point and its operations, and the commit action; `grep -rn "k8s" README.md docs/ src/` returns no command references.

### Implementation for User Story 5

- [ ] T024 [P] [US5] Update README.md: Apple Container 1.6.1 as supported backend (prerequisites + verify command), Kubernetes operations as the `container cluster` family with no start action, commit action in features, `--workers`/`--cni`/`--node` and machine CPUs/memory/home-mount options
- [ ] T025 [P] [US5] Update docs/user-guide.md: replace `container k8s ...` command references with `container cluster ...`, remove the start action, document multi-worker create / node-targeted load-image / commit / extended machine create, bump version references from 1.0/1.5.0-era text to 1.6.1
- [ ] T026 [P] [US5] Update docs/ai-menu-map.md: Kubernetes branch owns `container cluster` commands via src/services/kubernetes_builders.go with no start path, commit action owned by src/ui/container_submenu.go plus the new builder file, machine-create ownership includes the new fields
- [ ] T027 [US5] Sweep stale references across README.md, docs/, and src/ (`k8s` command strings, `k8s-dev` placeholder remnants, start-cluster text, pre-1.6 version claims) and fix any misses from T005/T024–T026; verify `src/ui/help.go` backend text matches 1.6.1

**Checkpoint**: Documentation review passes — README plus affected guides agree on backend version, entry points, operations, and safeguards

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Final validation across all stories

- [ ] T028 Run the automated gate: `gofmt -l ./...` clean, `go vet ./...` clean, `go test ./...` all pass; fix any regressions
- [ ] T029 Run specs/015-container-1-6-support/quickstart.md scenarios 1–6 against Apple Container 1.6.1 (rename, multi-worker create, node-targeted load, commit incl. dry-run, machine options, docs grep) and record outcomes

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — T001, T002 in parallel
- **Foundational (Phase 2)**: Depends on Setup — T003 BLOCKS all user stories
- **User Stories (Phase 3–7)**: All depend on Foundational completion
  - US1 → US2 → US3 → US4 sequential in priority order recommended (US2 extends US1's builders/forms; US3/US4 touch different files and can overlap once US1 lands)
  - US5 docs can be drafted in parallel with code phases (all [P]), final sweep T027 after code is complete
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1 (P1)**: After Foundational — no dependencies on other stories; MVP (restores broken Kubernetes workflow)
- **US2 (P1)**: After US1 (extends `KubernetesCreateBuilder`/create form and load-image form from US1)
- **US3 (P1)**: After Foundational — new files only (`commit_container_builder.go` + submenu action); no dependency on US1/US2 code
- **US4 (P2)**: After Foundational — touches only machine files; independent of US1–US3
- **US5 (P2)**: After code phases for final sweep (T027); doc drafts (T024–T026) can start any time after Foundational

### Within Each User Story

- Builder/tests contract first (contracts/commands.md defines exact args), then UI wiring, then UI-flow tests
- New-field validation before preview in every form (workers, nodes, reference, cpus, home-mount)
- Story checkpoint (independent test) before moving to next priority

### Parallel Opportunities

- T001 + T002 (setup, independent checks)
- T005 + T006 + T008 (US1, three different files)
- T016 (new test file section) can be written from the T017 contract in parallel with T017
- T020's expectations derive from contracts/commands.md and can be drafted alongside T021
- T024 + T025 + T026 (three different docs, no shared files)
- US3 and US4 code phases touch disjoint file sets and can be worked in parallel by two developers after US1

---

## Parallel Example: User Story 1

```bash
# Launch independent US1 file changes together:
Task: "Rename k8s to cluster in src/services/kubernetes_builders.go, delete KubernetesStartBuilder (T005)"
Task: "Update stale command comments in src/services/kubernetes_parser.go and src/models/kubernetes_cluster.go (T006)"
Task: "Change create-form placeholder k8s-dev to cluster-dev in src/ui/kubernetes_create.go (T008)"
```

## Parallel Example: User Story 5 (docs)

```bash
# Launch all doc updates together (different files, contracts define content):
Task: "Update README.md for 1.6.1 backend, cluster family, commit, workers/CNI, machine opts (T024)"
Task: "Update docs/user-guide.md cluster commands and workflows (T025)"
Task: "Update docs/ai-menu-map.md ownership for cluster/commit/machine flows (T026)"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001–T002)
2. Complete Phase 2: Foundational (T003 — 1.6.1-only target confirmed)
3. Complete Phase 3: User Story 1 (T004–T009)
4. **STOP and VALIDATE**: `k` workflow end-to-end on 1.6.1 + `go test ./...` green
5. Mergeable increment: Kubernetes workflow unbroken (the P1 fix that motivated the feature)

### Incremental Delivery

1. Setup + Foundational → target confirmed
2. Add US1 → test independently → MVP (Kubernetes works on 1.6.1)
3. Add US2 → test independently → multi-worker clusters usable
4. Add US3 → test independently → commit usable
5. Add US4 → test independently → machine create current
6. Add US5 → docs review passes → feature complete

### Parallel Team Strategy

With two developers after Foundational:

1. Both complete Setup + Foundational together
2. Dev A: US1 → US2 (Kubernetes builders/forms are sequential)
3. Dev B: US3 + US4 (disjoint files: submenu/commit vs machine) then US5 docs
4. Joint: Phase 8 polish + quickstart validation on 1.6.1

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps each task to its user story for traceability
- T004/T010/T016/T020 test tasks are REQUIRED by spec FR-010 (constitution V: tested command contracts)
- Commit after each task or logical group; stop at any checkpoint to validate that story independently
- `container cluster list` parser needs NO logic change (verified identical table shape) — T005/T006 cover the rename + comments only
- No new screens or `messages.go`/`app.go` routing changes planned (commit is a submenu preview action); if implementation proves otherwise, update docs/ai-menu-map.md per AGENTS.md
- AGENTS.md working rules apply: non-main branch (done: `015-container-1-6-support`), commit changes, open PR to `main`, update docs/ai-menu-map.md + README sections for affected functionality
