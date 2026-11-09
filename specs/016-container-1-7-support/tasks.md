# Tasks: Container 1.7 Support

**Input**: Design documents from `/specs/016-container-1-7-support/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/
**Branch**: `016-container-1-7-support`

**Tests**: Spec FR-008 explicitly requires automated checks for the renamed command composition, removed actions, new create options, and updated parsing — test tasks are included per story. Constitution V likewise mandates builder/parser/UI-flow test coverage.

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story (US1, US2, US3, US4)
- Exact file paths in every description

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm baseline builds before touching cluster code

- [X] T001 Verify baseline builds and tests pass: run `go build ./...`, `go vet ./...`, `go test ./...` from repo root and record status in task notes

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Service-layer contract changes every story depends on — MUST complete before UI phases

**⚠️ CRITICAL**: No UI work (Phases 4+) can begin until this phase is complete

- [X] T002 [US1] Delete `KubernetesLoadImageBuilder` struct and its `Build` method from `src/services/kubernetes_builders.go`
- [X] T003 [US1] Add `SharedImageStore bool` field to `KubernetesCreateBuilder` in `src/services/kubernetes_builders.go`, appending bare `--shared-image-store` flag only when true
- [X] T004 [US1] Extend `KubernetesCreateBuilder` workers validation in `src/services/kubernetes_builders.go` to accept bare `<n>` or `count=<n>,cpus=<n>,memory=<size>` specs and emit repeatable `--workers` flags (comma/space-separated groups each become one `--workers` flag)
- [X] T005 [US2] Extend `ParseKubernetesList` in `src/services/kubernetes_parser.go` to detect the 1.7.1 aggregate shape (`NAME` first header: `NAME NODES K8S VERSION STATE`) alongside the existing per-node shape (`CLUSTER` first header), skipping unparseable rows without error
- [X] T006 Remove `KubernetesLoadImageInput` struct from `src/models/kubernetes_cluster.go`

**Checkpoint**: `go build ./...` passes; all cluster command composition lives behind the 1.7.1 contract

---

## Phase 3: User Story 1 - Manage Kubernetes clusters without the removed load-image command (Priority: P1) 🎯 MVP

**Goal**: No load-image action offered anywhere; creation exposes the `--shared-image-store` replacement; every remaining operation previews a supported 1.7.1 command

**Independent Test**: With Apple Container 1.7.1 installed, press `k`, list clusters, create a cluster with the shared image store option, write kubeconfig, delete the cluster — each showing a command preview using only supported 1.7.1 subcommands, with no load-image action offered (spec US1)

### Tests for User Story 1

- [X] T007 [P] [US1] Rewrite `TestKubernetesBuilders` in `src/services/services_test.go`: remove load-image assertions, assert `Build()` never emits `load-image`, assert `--shared-image-store` present only when enabled, assert worker-spec forms (`2`, `count=2,cpus=2,memory=4G`) compose correctly
- [X] T008 [P] [US1] Update UI flow tests in `src/ui/ui_flow_test.go`: remove the node-targeted load-image preview test (~line 327), assert cluster submenu options contain no load-image action
- [X] T009 [P] [US1] Update `tests/integration/documentation_test.go`: replace required string `container cluster load-image` with `container cluster list --nodes` and `--shared-image-store`, replace `1.6.1` with `1.7.1`, add `container cluster load-image` to the forbidden list

### Implementation for User Story 1

- [X] T010 [US1] Remove the `Load image` option and its `load-image` select branch from `src/ui/kubernetes_cluster_submenu.go`
- [X] T011 [US1] Delete screen file `src/ui/kubernetes_load_image.go`
- [X] T012 [US1] Remove `ScreenKubernetesLoadImage` constant from `src/ui/messages.go`
- [X] T013 [US1] Remove `kubernetesLoadImage` field, init wiring, screen-change `SetCluster` propagation, `Update`/`View`/`statusBar`/`isLoading`/`initForActive` branches for the load-image screen in `src/ui/app.go`
- [X] T014 [US1] Add shared-image-store toggle input and wire it into the create preview in `src/ui/kubernetes_create.go`

**Checkpoint**: `grep -rn "load-image\|LoadImage\|ScreenKubernetesLoadImage" src/ --include='*.go'` returns nothing; US1 independently testable per its Independent Test

---

## Phase 4: User Story 2 - Read cluster topology in the one-row-per-cluster list (Priority: P1)

**Goal**: Cluster list renders the 1.7.1 aggregate shape with per-node detail available, state correct, empty states readable

**Independent Test**: With one multi-node cluster present, list clusters shows each once with correct state, node detail matches `container cluster list --nodes` output (spec US2)

### Tests for User Story 2

- [X] T015 [P] [US2] Extend `TestParseKubernetesList` in `src/services/services_test.go` with the 1.7.1 aggregate table fixture (`NAME NODES K8S VERSION STATE` header) asserting one cluster per row with correct name/state, plus the `--nodes` per-node fixture asserting node rows parse as before

### Implementation for User Story 2

- [X] T016 [US2] Verify `src/ui/kubernetes_cluster_list.go` renders parsed aggregates without change; fetch per-node detail via `container cluster list --nodes` where node topology is required (submenu detail), keeping the empty-state message and error formatting intact

**Checkpoint**: US1 + US2 both work; `k` list matches live `container cluster list` / `--nodes` output

---

## Phase 5: User Story 3 - Size worker nodes at creation time (Priority: P2)

**Goal**: Cluster creation accepts worker count with per-group CPU/memory sizing, carried into the previewed command with validation before execution

**Independent Test**: Create a cluster with worker count + CPUs + memory, verify preview carries the spec, invalid specs show validation messages before execution (spec US3)

### Tests for User Story 3

- [X] T017 [P] [US3] Add worker-spec validation tests in `src/services/services_test.go`: invalid specs (negative count, unparseable count, empty token handling) return errors from `KubernetesCreateBuilder.Build()`; empty workers omitted from args

### Implementation for User Story 3

- [X] T018 [US3] Update the workers field in `src/ui/kubernetes_create.go` (label, placeholder e.g. `count=2,cpus=2,memory=4G`, pass-through to builder) so the preview shows `--workers <spec>` and builder validation errors surface in `errorMsg` before preview

**Checkpoint**: US3 complete; create preview for `count=2,cpus=2,memory=4G` + shared store matches contract `container cluster create ... --workers count=2,cpus=2,memory=4G --shared-image-store`

---

## Phase 6: User Story 4 - Find current 1.7.1 guidance in README and docs (Priority: P2)

**Goal**: README + internal guides + constitution describe the 1.7.1 backend, macOS 27 prerequisite, and current cluster operation set with no conflicting instructions

**Independent Test**: Using only updated README + guides, describe supported Container/macOS versions, the `k` operation set, and the image-loading replacement (spec US4)

- [X] T019 [P] [US4] Update `README.md`: `1.6.1` → `1.7.1` backend (prereq + machine line), `macOS 26.x` → `macOS 27.x`, Kubernetes bullet (remove load-image/per-node targets, add `--shared-image-store` + worker sizing + one-row list), feature list load-image references
- [X] T020 [P] [US4] Update `docs/user-guide.md`: prerequisites (1.7.1 + macOS 27.x), cluster action list and command mapping (no load-image; add `--shared-image-store`, `--nodes`, worker sizing), registries/backend version strings
- [X] T021 [P] [US4] Update `docs/ai-menu-map.md`: Kubernetes branch section (remove `KubernetesLoadImage` row + load-image ownership, update builder flags `--workers` sizing/`--shared-image-store`, list parsing both shapes), version strings `1.6.1` → `1.7.1`
- [X] T022 [P] [US4] Update `docs/binary-build-automation.md` macOS 26.x verification references → 27.x and `.github/agents/copilot-instructions.md` macOS 26.x validation note → 27.x
- [X] T023 [US4] Amend `.specify/memory/constitution.md` principle II (`macOS 26.x` → `macOS 27.x`), overview `macOS 26.x` reference, and manual-verification gate `macOS 26.x` → `macOS 27.x`, with amendment rationale per governance rules (depends on T019–T022 for consistent wording)

**Checkpoint**: `grep -rni "1\.6\.1\|macOS 26\|load-image" README.md docs/ .specify/memory/constitution.md` returns no stale user-facing references (historical `specs/015-*` excluded)

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final validation across all stories

- [X] T024 Run full quickstart validation per `specs/016-container-1-7-support/quickstart.md` (live 1.7.1 scenarios 1–6) and record results
- [X] T025 Run automated gate from repo root: `gofmt -l ./...` clean, `go vet ./...` clean, `go test ./...` all pass; fix any failures

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all UI story work (Phases 3–5 UI tasks)
- **User Stories**: US1 (Phase 3) → US2 (Phase 4) → US3 (Phase 5) in priority order; US2 implementation (T016) depends on foundational parser (T005); US3 form work (T018) depends on foundational builder fields (T003–T004)
- **Docs (Phase 6)**: Depends on code contract settling (Phases 2–5) for accurate command strings; T023 depends on T019–T022
- **Polish (Phase 7)**: Depends on all stories complete

### Within Each User Story

- Builder/parser changes before UI wiring before docs-test updates
- Test updates written alongside implementation (constitution V); failing-then-passing verified via `go test` per phase

### Parallel Opportunities

- T007, T008, T009 (different test files) can run in parallel
- T010–T014 touch different files but share the load-image removal intent — T011/T012/T013 (deletion + constant + wiring) can run in parallel after T010
- T019–T022 (four doc files) can run in parallel
- T015 (parser test) can run in parallel with US1 UI tasks once T005 lands

---

## Parallel Example: User Story 1 tests

```bash
# Launch all US1 test updates together (different files, no dependencies):
Task: "Rewrite TestKubernetesBuilders in src/services/services_test.go per T007"
Task: "Update UI flow tests in src/ui/ui_flow_test.go per T008"
Task: "Update tests/integration/documentation_test.go strings per T009"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001)
2. Complete Phase 2: Foundational (T002–T006)
3. Complete Phase 3: US1 tests + removal + create toggle (T007–T014)
4. **STOP and VALIDATE**: `go build ./... && go test ./...`; manual `k` workflow shows no load-image and previews only 1.7.1 commands
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → command contract correct
2. + US1 → dead workflow gone, replacement exposed (MVP!)
3. + US2 → list topology correct on 1.7.1 shapes
4. + US3 → worker sizing usable
5. + US4 docs + constitution → release-ready guidance
6. Polish gate → merge

---

## Notes

- Single-version target: 1.7.1 only, no `load-image` fallback (spec assumption)
- Installed-CLI help is authoritative over release-note wording (`--shared-image-store` on `create`, no `cluster run` on 1.7.1 — see research.md)
- Historical `specs/015-*` files intentionally untouched — they describe the prior release
- Commit after each task or logical group; stop at any checkpoint to validate
