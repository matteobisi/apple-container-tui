# Tasks: Security Hardening Quick Wins

**Input**: Design documents from `/specs/012-security-hardening-quick-wins/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/workflow-contract.md, quickstart.md

**Tests**: Not requested — no test tasks generated.

**Organization**: Tasks are grouped by user story. US1 and US2 are fully independent and can be implemented in parallel.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[US1]**: User Story 1 — Add SECURITY.md (Security-Policy 0→10)
- **[US2]**: User Story 2 — Attest build provenance (Signed-Releases 0→8–10)
- Exact file paths included in descriptions

---

## Phase 1: Setup

**Purpose**: Verify preconditions for this feature

- [x] T001 Verify Dockerfile no longer exists in repository root (confirms Pinned-Dependencies 7→10 resolved with no action needed)

---

## Phase 2: User Story 1 — Add Security Policy File (Priority: P1) 🎯 MVP

**Goal**: Create a `SECURITY.md` at the repository root so security researchers can find vulnerability reporting instructions, and the OSSF Scorecard Security-Policy check awards 10/10.

**Independent Test**: Verify `SECURITY.md` exists at repository root, contains vulnerability reporting instructions and a response timeline, and that `grep -E "Supported Versions|Report|Timeline|Disclosure" SECURITY.md` matches all required sections.

### Implementation for User Story 1

- [x] T002 [P] [US1] Create SECURITY.md at repository root with all required sections per data-model.md: supported-versions table, reporting method (GitHub private vulnerability reporting / Security Advisories), response timeline (concrete window e.g. 5 business days), and disclosure policy — file path: SECURITY.md

**Checkpoint**: SECURITY.md is present at root. Run `grep -E "Supported Versions|Report|Timeline" SECURITY.md` to confirm required sections. Scorecard Security-Policy check should report 10/10 on next scan.

---

## Phase 3: User Story 2 — Attest Build Provenance for Release Artifacts (Priority: P2)

**Goal**: Append a provenance attestation step to `.github/workflows/publish-release.yml` so that every release binary and SBOM gets a signed SLSA provenance attestation, and the Scorecard Signed-Releases check awards partial-to-full credit.

**Independent Test**: Run `grep "attest-build-provenance" .github/workflows/publish-release.yml` to confirm the action is present with SHA `a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32`. Verify `id-token: write` and `attestations: write` appear in the publish job permissions. After a release, run `gh attestation verify <artifact> --repo <owner>/apple-container-tui` to confirm attestation.

### Implementation for User Story 2

- [x] T003 [P] [US2] Add `id-token: write` and `attestations: write` permissions to the `publish` job permissions block (alongside existing `actions: read` and `contents: write`) in .github/workflows/publish-release.yml
- [x] T004 [US2] Append "Attest build provenance" step after the "Publish release" step in .github/workflows/publish-release.yml — use `actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32 # v4.1.0, Node 24 compatible`, with `subject-path` listing both `release-assets/actui-darwin-arm64` and `release-assets/actui-darwin-arm64.spdx.json`, gated by `if: steps.idempotency-check.outputs.skip != 'true'` (same condition as publish step) — per contracts/workflow-contract.md

**Checkpoint**: Workflow file has correct permissions and the attestation step in the right position. `grep -c "attest-build-provenance" .github/workflows/publish-release.yml` returns 1. All action references are SHA-pinned per AGENTS.md.

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Documentation updates and final validation

- [x] T005 [P] Update docs/security-automation.md with provenance attestation documentation: describe the new attestation step, `gh attestation verify` usage, and Scorecard impact
- [x] T006 [P] Update README.md security-related sections to mention SECURITY.md and release provenance verification (if applicable sections exist)
- [x] T007 Run quickstart.md validation steps: verify SECURITY.md sections present, verify workflow permissions and attestation step, confirm all action SHAs match pinning contract in contracts/workflow-contract.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — verification only
- **User Story 1 (Phase 2)**: No dependencies on other phases — can start immediately after Setup
- **User Story 2 (Phase 3)**: No dependencies on US1 — can start immediately after Setup
- **Polish (Phase 4)**: Depends on US1 and US2 completion

### User Story Dependencies

- **User Story 1 (P1)**: Fully independent — single new file (`SECURITY.md`), no interaction with US2
- **User Story 2 (P2)**: Fully independent — amends existing workflow file, no interaction with US1

### Within User Story 2

- T003 (permissions) and T004 (attestation step) modify the same file (`.github/workflows/publish-release.yml`)
- T003 SHOULD be applied before T004 to avoid merge conflicts, but both changes target different YAML sections
- T004 depends on T003 being complete (the attestation step requires the `id-token` and `attestations` permissions to function)

### Parallel Opportunities

- T002 (US1) and T003 (US2) can run in parallel — different files, no dependencies
- T005 and T006 can run in parallel — different documentation files
- US1 and US2 are fully independent stories and can be implemented concurrently

---

## Parallel Example: Full Feature

```
# After Setup (T001):

# US1 and US2 can start simultaneously:
  ├── T002 [US1] Create SECURITY.md
  └── T003 [US2] Add permissions to publish-release.yml
       └── T004 [US2] Append attestation step to publish-release.yml

# After both stories complete:
  ├── T005 Update docs/security-automation.md
  ├── T006 Update README.md
  └── T007 Run quickstart.md validation
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (verify Dockerfile removal)
2. Complete Phase 2: User Story 1 — create SECURITY.md
3. **STOP and VALIDATE**: Scorecard Security-Policy check should reach 10/10
4. Merge independently — delivers SC-002, SC-006 immediately

### Incremental Delivery

1. Setup → verification only
2. Add User Story 1 (SECURITY.md) → merge and verify independently → Scorecard Security-Policy 10/10
3. Add User Story 2 (attestation) → merge and verify after next release → Scorecard Signed-Releases non-zero
4. Each story adds Scorecard score improvement without depending on the other (SC-005)
5. Both together (plus Dockerfile removal) target overall score 5.9 → ~7.5–8.0 (SC-001)

### Key Constraints

- All workflow actions MUST be SHA-pinned per AGENTS.md policy — no bare version tags
- Only Node 24 compatible action versions are permitted
- The attestation step MUST fail the workflow if it errors — a release must never be published without its attestation (contracts/workflow-contract.md invariant)

---

## Notes

- Total deliverables: 2 files (1 new, 1 amended)
- `SECURITY.md` — new file at repository root
- `.github/workflows/publish-release.yml` — amended (permissions + 1 new step)
- The Dockerfile removal (Pinned-Dependencies 7→10) is already done — T001 is a verification-only step
- Signed-Releases score only updates after the first release published with attestation
- No runtime code changes; no Go code modified; no TUI behavior affected
