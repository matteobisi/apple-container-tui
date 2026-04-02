# Tasks: Automated Binary Release Publishing

**Input**: Design documents from `/specs/010-auto-release-publish/`
**Prerequisites**: plan.md ✅ | spec.md ✅ | research.md ✅ | data-model.md ✅ | contracts/ ✅ | quickstart.md ✅

**Tests**: No test tasks — tests were not requested in the feature specification.

**Tech stack**: GitHub Actions YAML workflows + shell scripting on ubuntu-latest runners
**Scope boundary**: No application source code changes; output is `.github/workflows/publish-release.yml` + docs updates only

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: User story this task belongs to (US1, US2, US3)
- Exact file paths are included in descriptions

---

## Phase 1: Setup (Repository Readiness)

**Purpose**: Confirm interface contract between existing build workflow and new release workflow before writing any code

- [x] T001 Audit `.github/workflows/build-binary.yml` to confirm exact workflow name (`Build Binary`), artifact upload action used (`actions/upload-artifact@v4`), and artifact output name (`actui-linux-amd64`) — values needed for `workflow_run` trigger and `download-artifact` step
- [x] T002 Verify that repository Settings → Actions → General → Workflow permissions allow `contents: write` scoping, or confirm it can be overridden per-job in the workflow YAML

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Confirm all cross-cutting decisions that the release workflow and docs must satisfy before user stories begin

**⚠️ CRITICAL**: Complete before writing workflow or docs

- [x] T003 Confirm which `actions/download-artifact` version tag (v3 or v4) supports download from a `workflow_run` context and is Node 24 compatible — document result as inline comment in workflow file
- [x] T004 Confirm the GitHub release action to use (`softprops/action-gh-release` or GitHub CLI `gh release create`) and that the chosen option is Node 24 compatible — document the choice in the workflow file header comment

**Checkpoint**: Interface constants known, actions chosen → user story implementation can begin

---

## Phase 3: User Story 1 — Publish Release After Qualified Build (Priority: P1) 🎯 MVP

**Goal**: A successful `Build Binary` run automatically triggers a new GitHub Release with `actui-linux-amd64` attached.

**Independent Test**: Push to main → build completes successfully → release workflow runs → one new release appears in GitHub Releases with binary attached and no pre-existing tag collision.

### Implementation for User Story 1

- [x] T005 [US1] Create `.github/workflows/publish-release.yml` with workflow name `Publish Release`, `workflow_run` trigger on `Build Binary`, filter `types: [completed]`, initial permissions block (`contents: write`, all other scopes omitted or read), and a single job `publish`
- [x] T006 [US1] Add guard step in `.github/workflows/publish-release.yml` to exit early (`fail` or `cancel`) when `github.event.workflow_run.conclusion != 'success'`, with an explicit log message explaining the skip
- [x] T007 [US1] Add artifact download step in `.github/workflows/publish-release.yml` using the Node 24 compatible action confirmed in T003, scoped to artifact `actui-linux-amd64` from the triggering `workflow_run`; add error-exit step if artifact is not found
- [x] T008 [US1] Add idempotency guard step in `.github/workflows/publish-release.yml` that queries existing tags (`gh release list` or GitHub API), detects if the computed tag already exists, logs the outcome, and skips publication if duplicate detected
- [x] T009 [US1] Add release creation step in `.github/workflows/publish-release.yml` using the action confirmed in T004 to create a release with computed `tag_name`, title derived from tag, auto-generated notes, and upload `actui-linux-amd64` binary as release asset
- [x] T010 [US1] Add stage log outputs throughout `.github/workflows/publish-release.yml` for: artifact resolution result, version selection result, idempotency check result, and publication outcome — each log must identify the step and contain enough context for diagnostics

**Checkpoint**: At this point, US1 is functional. A qualifying build triggers a release with a working (but possibly hardcoded-start) tag. US2 will wire in dynamic version increment.

---

## Phase 4: User Story 2 — Apply Predictable Version Labels (Priority: P2)

**Goal**: Each automated release receives a deterministic tag (`v0.1.0` first run; `v0.1.1`, `v0.1.2`, … on subsequent runs) derived from existing tags, with conflict handling logged.

**Independent Test**: Trigger three successive qualifying builds and verify tags follow `v0.1.0 → v0.1.1 → v0.1.2`; rerun the third build and verify no new duplicate tag is created.

### Implementation for User Story 2

- [x] T011 [US2] Add version computation step in `.github/workflows/publish-release.yml` (before the idempotency guard T008) that: lists all existing `v*` tags sorted by semver, extracts the latest patch number, increments patch by 1, falls back to `v0.1.0` when no prior automated tag exists, and outputs the computed tag via `$GITHUB_OUTPUT` (`tag_name`)
- [x] T012 [US2] Wire `${{ steps.compute-version.outputs.tag_name }}` into the idempotency guard step (T008) and release creation step (T009) in `.github/workflows/publish-release.yml`, replacing any hardcoded tag reference
- [x] T013 [US2] Add conflict-strategy log in `.github/workflows/publish-release.yml` for when computed tag equals an existing tag: log computed tag, existing tag list, and the skip/find-next decision, then exit cleanly without error

**Checkpoint**: At this point, US1 and US2 are both functional. Automatic patch increment is live and idempotency is enforced.

---

## Phase 5: User Story 3 — Document and Communicate Release Process (Priority: P3)

**Goal**: `docs/binary-build-automation.md` describes the complete build-to-release lifecycle; `README.md`/`AGENTS.md` are updated if maintainer entry points changed.

**Independent Test**: Read `docs/binary-build-automation.md` without opening any workflow YAML — all trigger conditions, version-labeling rules, duplicate handling, and troubleshooting steps must be fully understood from docs alone.

### Implementation for User Story 3

- [x] T014 [P] [US3] Update `docs/binary-build-automation.md` with a "Release Publication Automation" section describing the full trigger chain: push to main → `Build Binary` workflow → `build-binary.yml` job success → `Publish Release` workflow trigger via `workflow_run`, including workflow name references
- [x] T015 [P] [US3] Update `docs/binary-build-automation.md` with a "Version Labeling Policy" subsection documenting: starting version `v0.1.0`, patch increment strategy, tag prefix `v`, and worked examples (`v0.1.0 → v0.1.1 → v0.1.2`)
- [x] T016 [P] [US3] Update `docs/binary-build-automation.md` with a "Duplicate and Rerun Behavior" subsection explaining idempotency guarantee: reruns of the same source commit do not produce a new release, with expected workflow log indicators
- [x] T017 [US3] Update `docs/binary-build-automation.md` with a "Troubleshooting Release Publication" subsection covering: artifact not found, permission errors (`contents: write` missing), duplicate tag skip (not an error), and how to locate run logs
- [x] T018 [US3] Update `docs/binary-build-automation.md` with an "Operator Validation Checklist" for the release automation: workflow enabled, artifact name matches contract, first release tag verified, duplicate rerun test passed, logs reviewed
- [x] T019 [P] [US3] Review `README.md` and add a brief reference to `docs/binary-build-automation.md` for the release automation process if README currently has no pointer to the release workflow
- [x] T020 [P] [US3] Review `AGENTS.md` and update if the release automation workflow changes the normal operator entry points for the repository (no changes needed — AGENTS.md is scoped to TUI screen navigation only)

**Checkpoint**: All three user stories are independently complete and testable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Validation, final checks, and commit readiness

- [ ] T021 [P] Validate `.github/workflows/publish-release.yml` YAML syntax with `actionlint` or GitHub Actions workflow linter (or via `gh workflow list` after push to branch)
- [ ] T022 Run quickstart.md validation steps: trigger qualifying build path, confirm release workflow fires, verify release tag and asset, test duplicate protection, verify docs match observed behavior
- [ ] T023 [P] Commit all planning artifacts (`specs/010-auto-release-publish/`) and implementation files (`.github/workflows/publish-release.yml`, `docs/binary-build-automation.md`, and any updated `README.md`/`AGENTS.md`) to branch `010-auto-release-publish`
- [ ] T024 Open pull request from `010-auto-release-publish` to `main` per AGENTS.md workflow rules (non-main branch, PR required, no direct push to main)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 completion — blocks all user story work
- **User Story 1 (Phase 3)**: Depends on Phase 2 — can start once actions are confirmed
- **User Story 2 (Phase 4)**: Depends on Phase 3 completion — version step is wired into US1 workflow
- **User Story 3 (Phase 5)**: Depends on Phases 3 + 4 — docs describe the final implemented workflow behavior
- **Polish (Phase 6)**: Depends on all user stories complete

### User Story Dependencies

- **US1 (P1)**: Can start after Foundational (Phase 2) — no dependency on US2 or US3
- **US2 (P2)**: Depends on US1 workflow skeleton existing — adds version computation into the same file
- **US3 (P3)**: Depends on US1 + US2 being complete — documents the finalized implementation

### Parallel Opportunities Within Phases

- T001 and T002 (Phase 1) can be done in parallel
- T003 and T004 (Phase 2) can be done in parallel
- T005–T010 (Phase 3) are sequential within the same file (each builds on previous)
- T011–T013 (Phase 4) are sequential within the same file
- T014–T015–T016 (Phase 5) can be done in parallel (distinct subsections)
- T017, T019, T020 (Phase 5) can be done in parallel (separate docs files)
- T021, T023 (Phase 6) can be done in parallel after T022

---

## Parallel Example: User Story 3 Docs (T014–T020)

```bash
# T014, T015, T016, T019, T020 can be written in parallel (different sections/files)
#
# Stream 1: docs/binary-build-automation.md sections
#   T014 → Release Publication Automation section
#   T015 → Version Labeling Policy subsection
#   T016 → Duplicate and Rerun Behavior subsection
#   T017 → Troubleshooting subsection
#   T018 → Operator Validation Checklist
#
# Stream 2: top-level entry point files (fully independent)
#   T019 → README.md review and update
#   T020 → AGENTS.md review and update
```

---

## Implementation Strategy

**MVP scope** (suggested minimum for initial PR): US1 (Phase 3) — a working `publish-release.yml` that triggers, downloads the artifact, checks for duplicates, and publishes with a deterministic starting tag. This alone satisfies FR-001, FR-002, FR-003, FR-004, FR-006, and SC-001.

**Incremental delivery order**:
1. Phase 1 + 2: Setup and research (T001–T004) — read-only, low risk
2. Phase 3 (T005–T010): Workflow skeleton + artifact publish — core deliverable
3. Phase 4 (T011–T013): Version auto-increment — adds US2 value on top of US1
4. Phase 5 (T014–T020): Docs — US3 value, can partially overlap with Phases 3/4
5. Phase 6 (T021–T024): Validate and submit PR

**Total tasks**: 24
**Tasks per user story**: US1 = 6 tasks (T005–T010) | US2 = 3 tasks (T011–T013) | US3 = 7 tasks (T014–T020)
**Setup/Foundational tasks**: 4 (T001–T004)
**Polish tasks**: 4 (T021–T024)
