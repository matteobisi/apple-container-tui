# Tasks: Automated Binary Build Workflow

**Input**: Design documents from `/specs/007-build-binary-action/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: No explicit TDD or test-first requirement was requested in the feature specification. This plan emphasizes implementation and verification tasks.

**Organization**: Tasks are grouped by user story to support independent implementation and validation.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Prepare target files and documentation surfaces.

- [x] T001 Review implementation scope and constraints in specs/007-build-binary-action/plan.md
- [x] T002 Create workflow scaffold in .github/workflows/build-binary.yml
- [x] T003 [P] Create runbook scaffold in docs/binary-build-automation.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Define shared behavior before user story work starts.

**CRITICAL**: Complete this phase before any user story implementation.

- [x] T004 Define qualifying merge and dispatch trigger rules in specs/007-build-binary-action/contracts/build-automation.md
- [x] T005 Implement trigger model and event filters in .github/workflows/build-binary.yml
- [x] T006 [P] Define least-privilege token permissions in .github/workflows/build-binary.yml
- [x] T007 [P] Define artifact retention policy and audit expectations in specs/007-build-binary-action/contracts/build-automation.md
- [x] T008 Define validation evidence format in specs/007-build-binary-action/quickstart.md

**Checkpoint**: Foundation complete; user stories can begin.

---

## Phase 3: User Story 1 - Automatically produce release-ready binaries after meaningful repository updates (Priority: P1)

**Goal**: Build and publish `actui` binaries for qualifying merged feature/dependency updates.

**Independent Test**: Merge a qualifying update to `main` (or run manual dispatch), confirm a successful workflow run, verify artifact availability, and verify actionable diagnostics on failure.

### Implementation for User Story 1

- [x] T009 [US1] Implement checkout and Go setup steps in .github/workflows/build-binary.yml
- [x] T010 [US1] Implement build command execution for ./cmd/actui in .github/workflows/build-binary.yml
- [x] T011 [US1] Implement deterministic artifact naming and upload in .github/workflows/build-binary.yml
- [x] T012 [US1] Implement explicit artifact retention window in .github/workflows/build-binary.yml
- [x] T013 [US1] Add run summary and failure diagnostics guidance in docs/binary-build-automation.md
- [x] T014 [US1] Validate local parity build command in docs/user-guide.md
- [x] T015 [US1] Record workflow syntax and retention verification steps in specs/007-build-binary-action/quickstart.md

**Checkpoint**: User Story 1 is functional and independently demonstrable.

---

## Phase 4: User Story 2 - Preserve repeatable build knowledge in repository documentation (Priority: P2)

**Goal**: Provide maintainers a complete operational runbook for workflow usage and troubleshooting.

**Independent Test**: A maintainer unfamiliar with setup can follow docs only and complete verification without external guidance.

### Implementation for User Story 2

- [x] T016 [US2] Document trigger behavior and qualifying updates in docs/binary-build-automation.md
- [x] T017 [US2] Document artifact retrieval and retention policy in docs/binary-build-automation.md
- [x] T018 [US2] Document retention troubleshooting for expired artifacts in docs/binary-build-automation.md
- [x] T019 [US2] Document release-readiness checklist in docs/binary-build-automation.md
- [x] T020 [US2] Document build failure troubleshooting flow in docs/binary-build-automation.md
- [x] T021 [US2] Add runbook cross-reference in docs/user-guide.md

**Checkpoint**: User Story 2 documentation is independently usable.

---

## Phase 5: User Story 3 - Capture reference validation environment for auditability (Priority: P3)

**Goal**: Capture reference validation environment details for reproducibility and incident analysis.

**Independent Test**: A maintainer can find and use explicit environment profile details during verification.

### Implementation for User Story 3

- [x] T022 [US3] Add validation environment profile (Macbook M4, macOS 26.4, 32GB RAM, Apple container 0.10.0) in docs/binary-build-automation.md
- [x] T023 [US3] Add environment-profile acceptance details in specs/007-build-binary-action/contracts/build-automation.md
- [x] T024 [US3] Align verification steps with environment profile in specs/007-build-binary-action/quickstart.md

**Checkpoint**: User Story 3 is independently verifiable.

---

## Phase 6: Polish and Cross-Cutting Concerns

**Purpose**: Final consistency checks and release-readiness gates.

- [x] T025 [P] Add README pointer to build automation runbook in README.md
- [x] T026 Validate consistency across specs/007-build-binary-action/spec.md
- [ ] T027 Validate quickstart end-to-end in specs/007-build-binary-action/quickstart.md
- [x] T028 Run and record mandatory manual verification on macOS 26.x in specs/007-build-binary-action/quickstart.md
- [x] T029 Capture final implementation evidence in specs/007-build-binary-action/tasks.md

---

## Dependencies and Execution Order

### Phase Dependencies

- Setup (Phase 1): No dependencies.
- Foundational (Phase 2): Depends on Setup and blocks all story phases.
- User Story phases (Phase 3 to 5): Depend on Foundational completion.
- Polish (Phase 6): Depends on completion of selected story phases.

### User Story Dependencies

- US1 (P1): Starts after Foundational; independent of US2 and US3.
- US2 (P2): Starts after Foundational; depends on US1 workflow behavior being available for documentation.
- US3 (P3): Starts after Foundational; depends on US2 documentation surfaces.

### Task-Level Dependency Highlights

- T005 depends on T004.
- T006 and T007 depend on T005.
- T009 to T012 depend on T005 and T006.
- T013 to T015 depend on T010 to T012.
- T016 to T021 depend on T013.
- T022 to T024 depend on T019 and T021.
- T025 to T029 depend on completion of T009 to T024.

---

## Parallel Opportunities

- Setup: T003 can run in parallel with T001 to T002.
- Foundational: T006 and T007 can run in parallel after T005.
- Polish: T025 can run in parallel with T026 to T027 once story tasks complete.

## Parallel Example: User Story 1

```bash
Task T009: Implement checkout and Go setup in .github/workflows/build-binary.yml
Task T014: Validate local parity build command in docs/user-guide.md
```

## Parallel Example: User Story 2

```bash
Task T017: Document artifact retrieval and retention policy in docs/binary-build-automation.md
Task T021: Add runbook cross-reference in docs/user-guide.md
```

## Parallel Example: User Story 3

```bash
Task T023: Add environment-profile acceptance details in specs/007-build-binary-action/contracts/build-automation.md
Task T024: Align verification steps in specs/007-build-binary-action/quickstart.md
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Setup (Phase 1).
2. Complete Foundational (Phase 2).
3. Complete User Story 1 (Phase 3).
4. Validate workflow run, artifact availability, diagnostics, and retention behavior.

### Incremental Delivery

1. Ship US1 automation.
2. Add US2 documentation completeness.
3. Add US3 environment auditability.
4. Complete Phase 6 validation and release-readiness checks.

### Suggested MVP Scope

- MVP: Phase 1 + Phase 2 + Phase 3 (through T015).
- Post-MVP: Documentation depth and auditability in Phases 4 to 6.

---

## Notes

- All tasks use strict checklist format with task ID and file path.
- Story labels are used only in user story phases.
- [P] marks tasks with safe parallel execution potential.
- Commit per task or tight logical batch for traceability.

## Implementation Evidence

- Workflow scaffold implemented in `.github/workflows/build-binary.yml` with pinned actions, push/manual triggers, build step, artifact upload, and explicit `retention-days: 14`.
- Runbook created in `docs/binary-build-automation.md` with trigger scope, retention policy, diagnostics, troubleshooting flow, and release checklist.
- User guide cross-reference added in `docs/user-guide.md`.
- README runbook pointer added in `README.md`.
- Contract and quickstart aligned in `specs/007-build-binary-action/contracts/build-automation.md` and `specs/007-build-binary-action/quickstart.md`.
- Local manual verification evidence captured on 2026-04-01: macOS `26.4`, build command `go build -o /tmp/actui-verify ./cmd/actui` passed, Apple container version `0.10.0`.