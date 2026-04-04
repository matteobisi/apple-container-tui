# Research: Security Hardening Quick Wins

**Feature**: `012-security-hardening-quick-wins`
**Date**: 2026-04-04
**Phase**: 0 — all NEEDS CLARIFICATION items resolved

---

## Decision 1: SECURITY.md Placement

**Decision**: Place `SECURITY.md` at the repository root (`/SECURITY.md`).

**Rationale**: OSSF Scorecard's Security-Policy check accepts `SECURITY.md`, `docs/SECURITY.md`, or `.github/SECURITY.md`. Repository root is the most visible location for human contributors (immediately browsable on GitHub's default view) and has the highest Scorecard detection confidence. GitHub itself displays a "Security Policy" link in the repository sidebar only when the file is at root or in `.github/`.

**Alternatives considered**:
- `.github/SECURITY.md` — equally valid for Scorecard; rejected because root placement is more discoverable without navigating folders.
- `docs/SECURITY.md` — valid for Scorecard; rejected because it buries the policy alongside feature documentation.

---

## Decision 2: Dockerfile Removed — No Pinning Needed

**Decision**: Skip the Dockerfile base image SHA256 pinning deliverable entirely.

**Rationale**: The Dockerfile was deleted from the repository. It was from a different project (Microsoft markitdown) and did not belong in this Go TUI repository. With the Dockerfile removed, there are no containerImage references in the repository at all, so the OSSF Scorecard Pinned-Dependencies containerImage sub-check resolves to 10/10 automatically — no action required.

**Previous state**: The Dockerfile referenced `python:3.13-slim-bullseye` by tag only, which was the sole source of the Pinned-Dependencies 7/10 score. The plan originally included a SHA256 pinning deliverable for this file.

**Current state**: Dockerfile deleted (commit on main). Pinned-Dependencies reaches 10/10 for free.

---

## Decision 3: Provenance Attestation Action

**Decision**: Use `actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` (v4.1.0).

**Rationale**:
- `actions/attest-build-provenance` v4.1.0 is the latest stable release as of 2026-04-04 (published 2026-02-26).
- The action is a composite wrapper over `actions/attest@59d89421af93a897026c735860bf21b6eb4f7b26` (v4.1.0), which declares `using: node24` — confirming Node 24 compatibility per AGENTS.md policy.
- Tag `v4.1.0` is a lightweight tag pointing directly to commit `a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` (verified via GitHub API: `type: commit`).
- The action requires only `id-token: write` and `attestations: write` permissions — no external secrets or third-party services.

**SHA resolution**:
```
Endpoint: GET /repos/actions/attest-build-provenance/git/ref/tags/v4.1.0
Response: { "object": { "sha": "a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32", "type": "commit" } }
```

**Node 24 compatibility**:
```
Composite: actions/attest-build-provenance (no runtime declared, delegates to child)
Child:     actions/attest@59d89421af93a897026c735860bf21b6eb4f7b26
Child runs: using: node24   ✓
```

**Pinned reference for AGENTS.md**:
```
actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32 # v4.1.0, Node 24 compatible
```

**Alternatives considered**:
- `actions/attest` directly — the underlying action; `attest-build-provenance` is the opinionated wrapper that generates SLSA provenance predicates, which is the correct choice for Scorecard Signed-Releases credit.
- `sigstore/cosign-installer` + manual signing — heavier approach requiring key management; rejected in favour of the GitHub-native keyless signing path.

---

## Decision 4: Attestation Subject Strategy

**Decision**: Use multi-line `subject-path` to attest both release artifacts (`actui-darwin-arm64` and `actui-darwin-arm64.spdx.json`) in one step invocation, satisfying FR-008.

**Rationale**: `actions/attest-build-provenance` v4.1.0 accepts a multi-line `subject-path` and can attest up to 1024 subjects per invocation. Attesting both in one step keeps the workflow minimal and ensures both attestations share the same provenance record (same run, same commit).

**Alternatives considered**:
- Two separate attestation steps (one for binary, one for SBOM) — valid but verbose; the single-step approach is simpler and explicitly supported.

---

## Decision 5: Permissions Model

**Decision**: Add `id-token: write` and `attestations: write` at the `publish` job level, alongside existing `actions: read` and `contents: write`.

**Rationale**:
- `id-token: write` enables the OIDC token exchange that backs keyless signing via GitHub's Sigstore infrastructure.
- `attestations: write` enables storing the attestation bundle in the GitHub repository attestations store.
- Both permissions are job-scoped, consistent with the existing least-privilege pattern in publish-release.yml.
- No step-level permission overrides are needed.

**Alternatives considered**:
- Workflow-level permissions — rejected; the existing workflow uses `permissions: {}` at top level with job-level overrides, which is the correct least-privilege pattern and must not be changed.

---

## Decision 6: OSSF Scorecard Score Impact (Updated)

| Check | Before | After | Notes |
|---|---|---|---|
| Security-Policy | 0/10 | 10/10 | `SECURITY.md` at root detected immediately |
| Pinned-Dependencies | 7/10 | 10/10 | Dockerfile removed from repo — no containerImage to pin |
| Signed-Releases | 0/10 | 8–10/10 | Provenance attestation; score after first attested release |
| Overall (estimated) | 5.9 | ~7.5–8.0 | Meets SC-001 target |

**Note on Pinned-Dependencies**: The original plan included a Dockerfile base image SHA pinning deliverable to close this gap. The Dockerfile has since been deleted (it was from the wrong project — Microsoft markitdown, not this Go TUI). With no Dockerfile in the repository, there are zero unpinned containerImage dependencies, and the score reaches 10/10 automatically.

**Note on Signed-Releases**: Scorecard awards Signed-Releases credit only after at least one release with attestation exists. The score is non-zero after the first publish following the merged workflow change.
