# Workflow Contract: SBOM Generation for Binary Builds

**Feature**: `011-sbom-binary-build`  
**Date**: 2026-04-04  
**Contract type**: CI/CD workflow interface contract

This document defines the interface contracts between the two GitHub Actions workflows involved in this feature. It specifies what each workflow produces (outputs) and consumes (inputs), their triggers, and the invariants that must hold at runtime.

---

## Contract 1: build-binary.yml (Producer)

### Trigger Contract

| Trigger | Condition |
|---|---|
| `push` | Branch: `main` only |
| `workflow_dispatch` | Manual; no conditions |

### Permission Contract

| Permission | Level | Reason |
|---|---|---|
| `contents` | `read` | Read repository source for checkout and SBOM generation |

No additional permissions are required for SBOM generation without attestation.

### Step Execution Order (invariant)

```
1. Checkout repository
2. Setup Go
3. Build actui binary               → produces: ./actui
4. Generate SBOM                    → produces: ./actui-linux-amd64.spdx.json
5. Upload binary artifact           → produces: workflow artifact actui-linux-amd64
6. Upload SBOM artifact             → produces: workflow artifact actui-linux-amd64-sbom
```

**Invariant**: Step 4 (SBOM generation) MUST execute only after Step 3 (binary build) succeeds. If Step 3 fails, Steps 4–6 MUST NOT run. The workflow run MUST be marked as failed if Step 4 fails, even if Step 3 succeeded.

### Output Contract

| Output | Name | Type | Contents | Retention |
|---|---|---|---|---|
| Binary artifact | `actui-linux-amd64` | Workflow artifact | Single file: `actui` executable (renamed to `actui-linux-amd64` by publish) | 14 days |
| SBOM artifact | `actui-linux-amd64-sbom` | Workflow artifact | Single file: `actui-linux-amd64.spdx.json` (SPDX 2.3 JSON) | 14 days |

**Pairing invariant**: Both artifacts MUST be produced in the same workflow run. A run that produces only the binary and not the SBOM MUST be considered failed (FR-008).

### Action Pinning Contract

All actions in this workflow MUST be pinned to an immutable commit SHA with a version comment. No bare version tags are permitted.

| Action | Required from | Notes |
|---|---|---|
| `actions/checkout` | `11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2) | From AGENTS.md |
| `actions/setup-go` | SHA resolved at implementation via `git ls-remote` | v5 series |
| `anchore/sbom-action` | SHA resolved at implementation via `git ls-remote` | Latest stable; composite action, no Node.js requirement |
| `actions/upload-artifact` (binary) | `ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2) | From AGENTS.md |
| `actions/upload-artifact` (SBOM) | `ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2) | Same SHA as binary upload |

---

## Contract 2: publish-release.yml (Consumer)

### Trigger Contract (unchanged)

| Trigger | Condition |
|---|---|
| `workflow_run` | Workflow: `Build Binary`; type: `completed`; conclusion: `success` |

**Gate invariant**: If the triggering `Build Binary` run did not produce BOTH the binary artifact and the SBOM artifact, publish MUST fail with a clear error message.

### Permission Contract (unchanged — already satisfies all requirements)

| Permission | Level | Reason |
|---|---|---|
| `actions` | `read` | Download artifacts from triggering workflow run |
| `contents` | `write` | Create releases and upload release assets |

### Input Contract

| Input | Source artifact | Expected filename after download | Expected location |
|---|---|---|---|
| Binary | `actui-linux-amd64` (from triggering run) | `actui` | `release-assets/actui` → renamed to `release-assets/actui-linux-amd64` |
| SBOM | `actui-linux-amd64-sbom` (from triggering run) | `actui-linux-amd64.spdx.json` | `release-assets/actui-linux-amd64.spdx.json` |

### Step Execution Order (additions highlighted)

```
1. Checkout repository
2. Download binary artifact             ← actui-linux-amd64
3. **Download SBOM artifact**           ← actui-linux-amd64-sbom   [NEW]
4. Verify and rename binary             → release-assets/actui-linux-amd64
5. **Verify SBOM**                      → release-assets/actui-linux-amd64.spdx.json  [NEW]
6. Compute next version tag
7. Check for duplicate release
8. Publish release                      → attaches both assets     [MODIFIED]
```

### Release Asset Contract

The `gh release create` command MUST attach both assets:

```sh
gh release create "$TAG" \
  release-assets/actui-linux-amd64 \
  release-assets/actui-linux-amd64.spdx.json \
  --title "$TITLE" \
  --generate-notes \
  --target "$BUILD_SHA"
```

| Release asset name | File | Scorecard relevance |
|---|---|---|
| `actui-linux-amd64` | Binary executable | Existing; unchanged |
| `actui-linux-amd64.spdx.json` | SBOM (SPDX 2.3 JSON) | Required for Scorecard SBOM check |

### Invariants

1. SBOM verification MUST run before the release is published. A missing or malformed SBOM MUST block publication.
2. The SBOM release asset name MUST end in `.spdx.json` so the Scorecard SBOM check can detect it.
3. The SBOM and binary MUST come from the same triggering workflow run (enforced by using the same `run-id` for both `download-artifact` steps).

---

## Scorecard Compliance Summary

| Scorecard check | Current state | After this feature |
|---|---|---|
| `SBOM` | 0 (no SBOM in releases) | 10/10 (SBOM attached to every release) |
| `Pinned-Dependencies` | Partial (build-binary.yml uses bare tags) | Improved (all actions SHA-pinned in build-binary.yml) |
| `Token-Permissions` | Maintained | Maintained (least-privilege permissions unchanged) |
