# Tasks: Repository Security Hardening

**Input**: Design documents from /specs/006-repo-security-hardening/
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/security-automation.md, quickstart.md

**Tests**: No code-level automated test suite changes are required by the specification; validation tasks focus on GitHub workflow/config behavior checks and merge-gate verification.

**Organization**: Tasks are grouped by user story so each story can be implemented, validated, and demonstrated independently.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Prepare repository paths and documentation anchors used by all stories.

- [X] T001 Create the initial Scorecard workflow scaffold file in .github/workflows/scorecard.yml
- [X] T002 Create a security automation operations note in docs/security-automation.md
- [X] T003 [P] Add a security automation section reference in README.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Define shared policy conventions required before story-specific implementation.

**⚠️ CRITICAL**: No user story work should start before this phase is complete.

- [X] T004 Define the canonical Scorecard required-check name and branch-protection mapping in docs/security-automation.md
- [X] T005 [P] Define repository validation checklist steps for security automation in specs/006-repo-security-hardening/quickstart.md
- [X] T006 [P] Align security automation contract details with planned implementation names in specs/006-repo-security-hardening/contracts/security-automation.md

**Checkpoint**: Shared security policy and validation conventions are finalized.

---

## Phase 3: User Story 1 - Establish Continuous Security Posture Visibility (Priority: P1) 🎯 MVP

**Goal**: Add OSSF Scorecard automation that runs on pull requests and default-branch pushes, with merge gating based on successful workflow completion.

**Independent Test**: Open a pull request and verify the Scorecard check appears and completes; confirm branch protection documentation states merge must be blocked when the check is non-success.

### Implementation for User Story 1

- [X] T007 [US1] Implement the Scorecard workflow triggers, permissions, and job in .github/workflows/scorecard.yml
- [X] T008 [US1] Set a stable Scorecard check/job display name for required-check binding in .github/workflows/scorecard.yml
- [X] T009 [US1] Document required-merge behavior for the Scorecard check in docs/security-automation.md
- [X] T010 [US1] Document PR and default-branch Scorecard verification steps in specs/006-repo-security-hardening/quickstart.md

**Checkpoint**: User Story 1 is fully functional and independently testable.

---

## Phase 4: User Story 2 - Automate Dependency Risk Reduction (Priority: P2)

**Goal**: Enable Dependabot for Go modules and GitHub Actions with monthly scheduling.

**Independent Test**: Confirm Dependabot configuration validates in GitHub and contains monthly rules for both gomod and github-actions ecosystems.

### Implementation for User Story 2

- [X] T011 [US2] Create Dependabot configuration schema and top-level metadata in .github/dependabot.yml
- [X] T012 [P] [US2] Add the monthly gomod update rule for root dependencies in .github/dependabot.yml
- [X] T013 [P] [US2] Add the monthly github-actions update rule in .github/dependabot.yml
- [X] T014 [US2] Document Dependabot ecosystem scope and cadence expectations in docs/security-automation.md
- [X] T015 [US2] Document Dependabot validation flow and expected PR behavior in specs/006-repo-security-hardening/quickstart.md

**Checkpoint**: User Story 2 is fully functional and independently testable.

---

## Phase 5: User Story 3 - Maintain Ongoing Security Hygiene (Priority: P3)

**Goal**: Ensure ongoing operational hygiene by documenting recurring review and triage workflow for Scorecard and Dependabot outputs.

**Independent Test**: Follow the operations guide to perform one recurring security review cycle and verify all required checks and update queues are covered.

### Implementation for User Story 3

- [X] T016 [US3] Add a monthly security review runbook for Scorecard and Dependabot triage in docs/security-automation.md
- [X] T017 [US3] Add merge-block troubleshooting guidance for failing/non-success Scorecard runs in docs/security-automation.md
- [X] T018 [US3] Add recurring maintenance verification steps to specs/006-repo-security-hardening/quickstart.md

**Checkpoint**: User Story 3 is fully functional and independently testable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final cross-story consistency and readiness checks.

- [X] T019 [P] Verify consistency of workflow/check names across .github/workflows/scorecard.yml and docs/security-automation.md
- [X] T020 [P] Verify consistency of Dependabot cadence and ecosystems across .github/dependabot.yml and docs/security-automation.md
- [ ] T021 Validate the full implementation against specs/006-repo-security-hardening/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies; start immediately.
- **Foundational (Phase 2)**: Depends on Phase 1 and blocks all user story work.
- **User Story 1 (Phase 3)**: Depends on Phase 2; recommended MVP.
- **User Story 2 (Phase 4)**: Depends on Phase 2; independent from US1 implementation details.
- **User Story 3 (Phase 5)**: Depends on Phase 2; can proceed once policy docs exist.
- **Polish (Phase 6)**: Depends on completion of all implemented stories.

### User Story Dependencies

- **US1**: No dependency on other user stories.
- **US2**: No dependency on other user stories.
- **US3**: No dependency on other user stories.

### Within Each User Story

- Implement configuration first.
- Document verification steps for that configuration.
- Validate story-specific acceptance before moving to the next story.

### Parallel Opportunities

- T003 can run in parallel with T001-T002.
- T005 and T006 can run in parallel in Phase 2.
- In US2, T012 and T013 can run in parallel once T011 creates the file skeleton.
- In Polish, T019 and T020 can run in parallel.

---

## Parallel Example: User Story 1

```bash
Task: T007 Implement Scorecard workflow in .github/workflows/scorecard.yml
Task: T009 Document required-merge behavior in docs/security-automation.md
```

## Parallel Example: User Story 2

```bash
Task: T012 Add monthly gomod update rule in .github/dependabot.yml
Task: T013 Add monthly github-actions update rule in .github/dependabot.yml
```

## Parallel Example: User Story 3

```bash
Task: T017 Add Scorecard merge-block troubleshooting guidance in docs/security-automation.md
Task: T018 Add recurring maintenance verification in specs/006-repo-security-hardening/quickstart.md
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup.
2. Complete Phase 2: Foundational.
3. Complete Phase 3: User Story 1.
4. Validate Scorecard visibility and merge-gate policy documentation.

### Incremental Delivery

1. Finish Setup + Foundational once.
2. Deliver US1 (Scorecard required check).
3. Deliver US2 (Dependabot ecosystems and cadence).
4. Deliver US3 (ongoing hygiene runbook).
5. Execute Phase 6 consistency and quickstart validation.

### Parallel Team Strategy

1. One contributor completes Phase 1-2.
2. After Phase 2:
   - Contributor A: US1 workflow implementation.
   - Contributor B: US2 Dependabot configuration.
   - Contributor C: US3 operational hygiene documentation.
3. Rejoin for Phase 6 cross-cutting validation.

---

## Notes

- [P] tasks touch different files and are safe to parallelize.
- Each user story remains independently implementable and testable after Phase 2.
- Branch protection setting changes are external to git-tracked files, so tasks document and validate the expected policy behavior explicitly.
