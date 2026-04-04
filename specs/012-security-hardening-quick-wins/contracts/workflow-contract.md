# Workflow Contract: Security Hardening Quick Wins

**Feature**: `012-security-hardening-quick-wins`
**Date**: 2026-04-04
**Contract type**: CI/CD workflow interface contract (amendment to publish-release.yml)

This document defines the interface change to `publish-release.yml` introduced by this feature. It is an amendment to the existing contract established in `011-sbom-binary-build`.

---

## Contract: publish-release.yml (amended)

### Change Summary

A provenance attestation step is appended to the `publish` job after the `gh release create` step. No existing steps are modified. No new workflow files are added.

---

### Permission Contract (updated)

| Permission | Level | Before | After | Reason for change |
|---|---|---|---|---|
| `actions` | job | `read` | `read` | Unchanged; still required for cross-run artifact download |
| `contents` | job | `write` | `write` | Unchanged; still required for release creation |
| `id-token` | job | absent | `write` | New; required for OIDC token exchange (keyless signing) |
| `attestations` | job | absent | `write` | New; required to write attestation bundle to GitHub store |

The top-level `permissions: {}` setting is unchanged. All permissions remain scoped to the `publish` job.

---

### Step Execution Order (invariant, amended)

```
1.  Checkout repository
2.  Download artifact actui-darwin-arm64
3.  Download SBOM artifact actui-darwin-arm64-sbom
4.  Verify and rename artifact
5.  Verify SBOM artifact
6.  Compute next version tag
7.  Check for duplicate release (idempotency)
8.  Publish release               → ghost step if skip=true
9.  Attest build provenance       → ghost step if skip=true   [NEW]
```

**Invariant**: Step 9 MUST only run when `steps.idempotency-check.outputs.skip != 'true'`. A failed attestation in Step 9 MUST cause the workflow run to fail, so that a release is never published without a corresponding attestation.

---

### New Step Contract: Attest Build Provenance

| Property | Value |
|---|---|
| Step name | `Attest build provenance` |
| `if` condition | `steps.idempotency-check.outputs.skip != 'true'` |
| Action | `actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32 # v4.1.0, Node 24 compatible` |
| `subject-path` | `release-assets/actui-darwin-arm64\nrelease-assets/actui-darwin-arm64.spdx.json` |
| Node runtime | `node24` (via `actions/attest@59d89421af93a897026c735860bf21b6eb4f7b26`) |
| GITHUB_TOKEN | Default (`${{ github.token }}`); no explicit input required |

**Input invariant**: The files at `subject-path` MUST exist on disk when this step runs. They are guaranteed to be present after Steps 4 and 5 succeed (binary renamed to `actui-darwin-arm64`, SBOM confirmed at `actui-darwin-arm64.spdx.json`, both under `release-assets/`).

**Output contract**: On success, the step attaches a signed SLSA provenance attestation to both listed artifacts in the GitHub repository attestations store. The attestation is publicly queryable via `gh attestation verify <artifact> --repo <owner/repo>`.

---

### Action Pinning Contract (cumulative for publish-release.yml)

All actions in this workflow MUST be pinned to an immutable commit SHA with a version comment. No bare version tags are permitted.

| Action | SHA | Version | Node |
|---|---|---|---|
| `actions/checkout` | `11bd71901bbe5b1630ceea73d27597364c9af683` | v4.2.2 | 24 |
| `actions/download-artifact` | `cc203385981b70ca67e1cc392babf9cc229d5806` | v4.1.9 | 24 |
| `actions/attest-build-provenance` | `a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` | v4.1.0 | 24 (composite → `actions/attest` node24) |

---

### Idempotency Guarantee (unchanged)

The idempotency check in Step 7 still gates both Step 8 (publish release) and Step 9 (attest). If a release already exists, neither the release creation nor the attestation is re-attempted, preventing duplicate attestation records.
