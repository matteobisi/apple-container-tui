# Quickstart: Security Hardening Quick Wins

**Feature**: `012-security-hardening-quick-wins`
**Date**: 2026-04-04

This guide shows how to verify each of the two deliverables once implemented.

---

## Prerequisites

- Repository cloned and branch `012-security-hardening-quick-wins` checked out.
- GitHub CLI (`gh`) installed and authenticated to the repository.

---

## Deliverable 1: SECURITY.md

### Verify the file exists and is well-formed

```bash
# From repository root
ls -la SECURITY.md
```

Expected: file present at root (not in `docs/` or `.github/`).

```bash
# Check required sections are present
grep -E "Supported Versions|Report|Timeline|Disclosure" SECURITY.md
```

Expected: at least one match for each section heading.

### Scorecard verification

The OSSF Scorecard GitHub Action runs on schedule (check `.github/workflows/scorecard.yml` for the cron). After the PR is merged, wait for the next scheduled scan or trigger it manually:

```bash
gh workflow run scorecard.yml
```

After the run completes:
```bash
gh run view --log | grep "Security-Policy"
```

Expected output: `Security-Policy: 10 / 10`

---

## Deliverable 2: Release Provenance Attestation

### Verify workflow permissions and step

```bash
grep -A 10 "permissions:" .github/workflows/publish-release.yml
```

Expected: `id-token: write` and `attestations: write` are present at job level.

```bash
grep "attest-build-provenance" .github/workflows/publish-release.yml
```

Expected: action reference with SHA `a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` and comment `# v4.1.0, Node 24 compatible`.

### Trigger a release and verify attestation

After merging, trigger the build workflow (which triggers publish-release on success):

```bash
gh workflow run "Build Binary"
```

After the release is published, verify attestation exists:

```bash
# Download the release binary
gh release download <tag> -p "actui-darwin-arm64" -D /tmp/

# Verify provenance attestation
gh attestation verify /tmp/actui-darwin-arm64 --repo <owner>/apple-container-tui
```

Expected: attestation verification succeeds, showing the source repository, workflow, and commit SHA.

```bash
# Also verify SBOM attestation
gh release download <tag> -p "actui-darwin-arm64.spdx.json" -D /tmp/
gh attestation verify /tmp/actui-darwin-arm64.spdx.json --repo <owner>/apple-container-tui
```

Expected: both artifacts have valid provenance attestations.

### Scorecard verification

After the first attested release is published:

```bash
gh workflow run scorecard.yml
```

After completion:
```bash
gh run view --log | grep "Signed-Releases"
```

Expected: `Signed-Releases` score is non-zero (8–10/10).

---

## Pinned-Dependencies (bonus — no action required)

The Dockerfile was removed from the repository (it was from a different project). With no container images in the repo, the Pinned-Dependencies score reaches 10/10 automatically:

```bash
gh run view --log | grep "Pinned-Dependencies"
```

Expected: `Pinned-Dependencies: 10 / 10`

---

## Overall Scorecard Target

After both deliverables are merged and the first attested release exists:

```bash
gh run view --log | grep "Aggregate score"
```

Expected: aggregate score ≥ 7.5 (target range: 7.5–8.0).
