# Implementation Plan: Security Hardening Quick Wins

**Branch**: `012-security-hardening-quick-wins` | **Date**: 2026-04-04 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/012-security-hardening-quick-wins/spec.md`

## Summary

Improve the OSSF Scorecard overall score from ~5.9 to ~7.5–8.0 by implementing two low-effort security hardening items: a standard `SECURITY.md` security policy file (Security-Policy 0→10), and release build provenance attestation via `actions/attest-build-provenance` in the publish workflow (Signed-Releases 0→8–10). A previously planned Dockerfile pinning deliverable is no longer needed — the Dockerfile was removed from the repository (it belonged to a different project), which resolves the Pinned-Dependencies containerImage finding for free (7→10).

## Technical Context

**Language/Version**: Go 1.24 (TUI application); YAML (GitHub Actions workflows); Markdown (SECURITY.md)
**Primary Dependencies**: `actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` (v4.1.0, Node 24 compatible)
**Storage**: N/A — no runtime storage changes
**Testing**: Manual verification via Scorecard scan; `gh attestation verify` for provenance
**Target Platform**: GitHub Actions (ubuntu-latest); repository metadata (GitHub)
**Project Type**: CI/CD pipeline amendment + repository policy document
**Performance Goals**: N/A
**Constraints**: All workflow actions must be SHA-pinned per AGENTS.md policy; Node 24 compatible only
**Scale/Scope**: 2 files changed/created: `SECURITY.md` (new), `.github/workflows/publish-release.yml` (amended)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Applies? | Status | Notes |
|---|---|---|---|
| I. Command-Safe TUI | No | PASS | No TUI code changes; deliverables are CI/CD and docs only |
| II. macOS 26.x + Apple Silicon | No | PASS | Workflow runs on ubuntu-latest (CI); no production code changes |
| III. Local-Only Operation | No | PASS | No telemetry or remote control surface added; provenance is CI-only |
| IV. Clear Observability | Marginal | PASS | Workflow attestation step outputs are visible in GH Actions logs |
| V. Tested Command Contracts | Yes | PASS | Workflow contract document provided (`contracts/workflow-contract.md`) |
| Workflow: Command mapping artifact | Marginal | PASS | No new user-facing commands; workflow contract covers CI interface |
| Workflow: Destructive action guardrails | No | PASS | No destructive actions introduced |
| Workflow: Manual verification on macOS | No | PASS | CI/CD-only changes; nothing to verify on macOS |

**Post–Phase 1 re-check**: PASS. No constitution violations. The workflow contract in `contracts/workflow-contract.md` documents the publish-release.yml amendment. No complexity tracking needed.

## Project Structure

### Documentation (this feature)

```text
specs/012-security-hardening-quick-wins/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── workflow-contract.md  # Phase 1 output
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
SECURITY.md                                  # NEW — security policy (Deliverable 1)
.github/workflows/publish-release.yml        # AMENDED — attestation step (Deliverable 2)
```

**Structure Decision**: No new source directories. Both deliverables are single-file changes at existing locations (repository root and `.github/workflows/`).

## Deliverables

### Deliverable 1: SECURITY.md (Security-Policy 0 → 10)

- **File**: `SECURITY.md` at repository root
- **Purpose**: Standard GitHub security policy enabling responsible vulnerability disclosure
- **Scorecard impact**: Security-Policy check 0 → 10 (MEDIUM risk)
- **Content**: Supported versions table, reporting method (GitHub private security advisory), response timeline, disclosure policy
- **Placement rationale**: Root is most visible for human contributors and has highest Scorecard detection confidence (see [research.md](research.md) Decision 1)

### Deliverable 2: Release Provenance Attestation (Signed-Releases 0 → 8–10)

- **File**: `.github/workflows/publish-release.yml` (amendment)
- **Purpose**: Attach signed SLSA build provenance to release artifacts for supply-chain verification
- **Scorecard impact**: Signed-Releases check 0 → 8–10 (HIGH risk); effective after first attested release
- **Changes**:
  - Add `id-token: write` and `attestations: write` to job-level permissions
  - Append `actions/attest-build-provenance` step after `gh release create`
  - Attest both `actui-darwin-arm64` and `actui-darwin-arm64.spdx.json`
  - Step gated by same idempotency condition as publish step
- **Action pinning**: `actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32 # v4.1.0, Node 24 compatible`

### Removed: Dockerfile Base Image SHA Pinning

The Dockerfile was deleted from the repository — it belonged to a different project (Microsoft markitdown). The Pinned-Dependencies containerImage finding resolves for free (7 → 10) with no deliverable needed.

## Expected Scorecard Impact

| Check | Before | After | Change driver |
|---|---|---|---|
| Security-Policy | 0/10 | 10/10 | `SECURITY.md` at root |
| Pinned-Dependencies | 7/10 | 10/10 | Dockerfile removed (no containerImage to pin) |
| Signed-Releases | 0/10 | 8–10/10 | Provenance attestation on next release |
| **Overall (estimated)** | **5.9** | **~7.5–8.0** | Meets SC-001 target |

## Complexity Tracking

No constitution violations to justify. Feature is minimal: one new file and one workflow step addition.
