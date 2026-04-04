# Data Model: Security Hardening Quick Wins

**Feature**: `012-security-hardening-quick-wins`
**Date**: 2026-04-04

This feature introduces no new runtime entities (no database tables, no in-memory models, no config schema changes). It delivers two static file changes:

1. A new documentation file (`SECURITY.md`)
2. A CI/CD workflow modification (build provenance attestation step)

The entities below model the content structure and invariants of each deliverable.

---

## Entity 1: SecurityPolicy

**File**: `SECURITY.md` (repository root)

| Field | Type | Constraint | Description |
|---|---|---|---|
| `supported_versions` | Markdown table | Required | Lists project versions with security-fix support status |
| `reporting_method` | Markdown section | Required (at least one) | How to privately disclose a vulnerability (GitHub Security Advisory) |
| `response_timeline` | Markdown section | Required | Maintainer's commitment: acknowledgment window and fix/disclosure timeline |
| `disclosure_policy` | Markdown section | Optional | Coordinated disclosure process and publication timeline |

**Validation rules**:
- File MUST be named exactly `SECURITY.md` (case-sensitive).
- File MUST reside at the repository root (i.e., `/SECURITY.md` relative to repo root).
- The `reporting_method` section MUST reference GitHub private vulnerability reporting as the primary channel (consistent with FR-002).
- The `response_timeline` section MUST state at least one concrete time window (e.g., "within 5 business days") to satisfy FR-003 and SC-006.

**State transitions**: None — this is a static document.

---

## Entity 2: ProvenanceAttestationStep

**File**: `.github/workflows/publish-release.yml`

| Field | Type | Constraint | Description |
|---|---|---|---|
| `action` | string | `actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` | Pinned to commit SHA per AGENTS.md |
| `action_version_comment` | string | `# v4.1.0, Node 24 compatible` | Inline version comment required by AGENTS.md |
| `subject_path` | string (multi-line) | Both release assets listed | Attests `actui-darwin-arm64` and `actui-darwin-arm64.spdx.json` |
| `required_permissions` | map | `id-token: write`, `attestations: write` | Job-level; added to existing permission set |
| `position_in_job` | enum | `after publish-release` | MUST run after `gh release create`; step conditional matches `publish` step |

**Job `publish` permissions after change**:
```yaml
permissions:
  actions:      read    # existing: download artifacts from triggering run
  contents:     write   # existing: create releases and upload assets
  id-token:     write   # new: OIDC token for keyless signing
  attestations: write   # new: store attestation in GitHub attestations store
```

**Step execution invariant**: The attestation step MUST be guarded by the same `if: steps.idempotency-check.outputs.skip != 'true'` condition as the `publish release` step. This ensures attestation is only attempted when a new release was actually published.

**Subject paths (both required per FR-008)**:
```
release-assets/actui-darwin-arm64
release-assets/actui-darwin-arm64.spdx.json
```

**Validation rules**:
- The action MUST be referenced by full commit SHA, not by version tag.
- The `subject-path` input MUST reference the files as they exist on disk in `release-assets/` (already present after the artifact download and rename steps).
- No `github-token` input override is required; the action defaults to `${{ github.token }}`.
