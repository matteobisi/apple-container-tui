# Binary Build Automation

This runbook describes how the repository builds and validates the `actui` binary in CI.

## Workflow Reference

- Workflow file: `.github/workflows/build-binary.yml`
- Workflow name: `Build Binary`
- Build command: `GOOS=darwin GOARCH=arm64 go build -o actui ./cmd/actui`

## Qualifying Triggers

- Push to `main` after merge.
- Manual dispatch for maintainer verification/recovery.

Qualifying updates include merged feature work and merged dependency updates.

## Artifact Contract

- Artifact name format: `actui-<os>-<arch>`
- Current artifact: `actui-darwin-arm64`
- Retention policy: explicit `retention-days` configured in workflow

## SBOM Generation

Every successful build produces a Software Bill of Materials (SBOM) alongside the binary. The SBOM captures the full Go module dependency graph resolved from `go.mod`/`go.sum` at build time.

| Property | Value |
|---|---|
| Generator | `anchore/sbom-action` (Syft) |
| Format | SPDX 2.3 JSON |
| Output file | `actui-darwin-arm64.spdx.json` |
| Workflow artifact name | `actui-darwin-arm64-sbom` |
| Retention policy | 14 days (same as binary artifact) |
| Release asset name | `actui-darwin-arm64.spdx.json` |

### SBOM Trigger and Flow

The SBOM is generated immediately after `go build` completes, before any upload steps:

```
 GOOS=darwin GOARCH=arm64 go build -o actui ./cmd/actui
  → anchore/sbom-action writes actui-darwin-arm64.spdx.json
    → upload-artifact uploads as 'actui-darwin-arm64-sbom'
      → Publish Release workflow downloads both artifacts
        → SBOM verified (jq format check) before release
          → actui-darwin-arm64.spdx.json attached to GitHub Release
```

### SBOM Verification

The publish workflow verifies the SBOM before attaching it to the release:

```sh
# Quick format check: confirms the file is valid JSON
jq empty release-assets/actui-darwin-arm64.spdx.json
```

A non-zero exit code (missing file or invalid JSON) will fail the publish job before the release is created.

### Scorecard Relevance

Attaching the SPDX 2.3 JSON file to the GitHub Release satisfies the OSSF Scorecard **SBOM** check (score: 0 → 10). Pinning all actions to immutable commit SHAs also improves the **Pinned-Dependencies** check score.

## Action Dependencies

The workflow uses GitHub Actions pinned to immutable commit SHAs (Node 24 compatible as of April 2026):
- `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
- `actions/setup-go@f111f3307d8850f501ac008e886eec1fd1932a34` (v5.3.0)
- `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
- `anchore/sbom-action@fd74a6fb98a204a1ad35bbfae0122c1a302ff88b` (v0.15.0, composite action)

Pinning to commit SHAs rather than mutable tags ensures immutable, reproducible builds and improves the OSSF Scorecard Pinned-Dependencies check.

## Diagnostics and Run Review

- Every run must show successful checkout, Go setup, build, and artifact upload steps.
- Failed runs must include actionable stderr/stdout in GitHub Actions logs.

Failure troubleshooting flow:

1. Confirm the failing step (`checkout`, `setup-go`, `build`, or `upload artifact`).
2. Read stderr/stdout from that step and classify cause (dependency, permissions, path, or runner).
3. Apply fix and rerun from workflow UI.
4. If artifact is missing or expired, verify `retention-days` and rerun.

## Release-Readiness Checklist

1. Build workflow succeeds for qualifying update.
2. Artifact uploads successfully with expected name.
3. Retention setting is present and matches policy.
4. Troubleshooting notes are updated for new failure modes.
5. Manual macOS 26.x verification is completed and recorded before release sign-off.

## Retention Troubleshooting

- If artifact disappears earlier than expected, check `retention-days` in workflow history.
- If retention is missing, update workflow and rerun validation.

## Validation Environment Profile

- Machine: Macbook M4
- OS: macOS 26.4
- Memory: 32GB
- Apple container version: 0.10.0

---

## Release Publication Automation

After a successful `Build Binary` run, the release workflow automatically publishes a GitHub Release with the produced binary. The full trigger chain is:

```
push to main
  → Build Binary workflow (.github/workflows/build-binary.yml)
    → job "Build actui binary" success
      → Publish Release workflow (.github/workflows/publish-release.yml) triggered via workflow_run
        → GitHub Release created with actui-darwin-arm64 and actui-darwin-arm64.spdx.json attached
          → GitHub provenance attestation generated for both release assets
```

Workflow file: `.github/workflows/publish-release.yml`  
Workflow name: `Publish Release`  
Trigger: `workflow_run` on `Build Binary`, type `completed`, gated on `conclusion == success`

If the build fails or is cancelled, the `Publish Release` workflow job is skipped automatically (no release is created).

### Release Provenance Attestation

After a new release is created, the same workflow generates keyless provenance attestations for both release assets using GitHub's built-in OIDC identity and attestations store.

| Property | Value |
|---|---|
| Action | `actions/attest-build-provenance` |
| Pinned version | `a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` (v4.1.0) |
| Attested assets | `actui-darwin-arm64`, `actui-darwin-arm64.spdx.json` |
| Required permissions | `id-token: write`, `attestations: write` |
| Verification command | `gh attestation verify <artifact> --repo matteobisi/apple-container-tui` |

The attestation step is gated by the same idempotency check as release publication, so reruns that skip duplicate release creation do not generate duplicate attestation attempts.

### Version Labeling Policy

Automated releases follow semantic versioning with a fixed patch-increment strategy:

| Release | Tag |
|---------|-----|
| First automated release | `v0.1.0` |
| Second | `v0.1.1` |
| Third | `v0.1.2` |
| … | … |

Rules:
- Tag prefix: `v`
- Starting version: `v0.1.0` (when no prior automated release exists)
- Increment strategy: patch component only (`PATCH + 1`), major and minor stay fixed
- Version is computed at release time by querying existing `gh release list` output and finding the highest `v*.*.*` tag

To change the version increment strategy (e.g., bump minor), update the version-computation step in `.github/workflows/publish-release.yml` and update this doc before the change takes effect.

### Duplicate and Rerun Behavior

Release publication is idempotent. If the `Publish Release` job runs more than once for the same computed tag (e.g., due to a workflow rerun), it:

1. Computes the next version tag as normal
2. Checks whether a release with that tag already exists (`gh release view`)
3. If the tag exists: logs a `::notice::` annotation and exits cleanly without publishing a duplicate
4. If the tag does not exist: publishes as normal

Expected log indicator when a duplicate is detected:

```
Notice: Release 'v0.1.2' already exists. Skipping publication (idempotency enforced).
```

This is not an error. The workflow run will show as completed successfully.

### Troubleshooting Release Publication

**Artifact not found**

```
Error: Artifact 'actui-darwin-arm64' not found in release-assets/
```

- The `Build Binary` run did not upload the artifact, or the artifact has expired.
- Verify the triggering build run completed `Upload artifact` step successfully.
- Check `retention-days` in `build-binary.yml`; if the artifact expired, re-trigger a qualifying build.

**Permission error (contents: write missing)**

```
gh: HTTP 403: Resource not accessible by integration
```

- The `GITHUB_TOKEN` lacks `contents: write` permission.
- Verify the `Publish Release` workflow has `permissions: contents: write` set at the job level.
- Check repository Settings → Actions → General → Workflow permissions for the default token scope.

**Attestation permission error**

```
Error: Resource not accessible by integration
```

- Verify the `publish` job grants both `id-token: write` and `attestations: write`.
- Confirm the attestation step runs after the release is created and is not skipped by idempotency.
- Re-run the workflow only after confirming a new release will actually be published.

**Duplicate tag skip (not an error)**

The `::notice::` log line is informational. If you need a new release despite the tag existing (e.g., re-attach an asset), delete the existing release and tag in GitHub UI, then rerun the `Publish Release` workflow manually from Actions.

**Locating run logs**

1. Go to the repository → Actions tab.
2. Select the `Publish Release` workflow.
3. Open the run that corresponds to the failing or skipped job.
4. Expand step logs: `Compute next version tag`, `Check for duplicate release`, `Publish release` for detailed output.

### Operator Validation Checklist for Release Automation

Use this checklist after enabling or modifying the release automation:

1. `Publish Release` workflow is enabled in repository Actions settings.
2. Artifact name in `publish-release.yml` matches `actui-darwin-arm64` (matches `build-binary.yml`).
3. Trigger a qualifying push to `main` and confirm `Build Binary` succeeds.
4. Confirm `Publish Release` workflow starts after `Build Binary` completes.
5. Verify the new release appears in GitHub Releases with the expected version tag (e.g., `v0.1.0`).
6. Verify `actui-darwin-arm64` is attached as a release asset.
7. Verify `actui-darwin-arm64.spdx.json` is attached as a release asset.
8. Verify both assets with `gh attestation verify` against `matteobisi/apple-container-tui`.
9. Rerun the same `Publish Release` run and confirm the duplicate-tag notice appears and no second release or attestation is created.
10. Review workflow logs to confirm all stage log lines are present (artifact, version, idempotency, publication, attestation).
