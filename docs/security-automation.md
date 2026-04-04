# Security Automation

This document defines repository security automation operations for OSSF Scorecard and Dependabot.

## Scope

- OSSF Scorecard workflow: .github/workflows/scorecard.yml
- Dependabot config: .github/dependabot.yml
- Ecosystems in scope: gomod, github-actions
- Dependabot cadence: monthly for security and version updates
- Workflow actions are pinned to immutable commit SHAs where practical
- SBOM (SPDX 2.3 JSON) is generated on every build and attached to GitHub Releases — satisfies the Scorecard `SBOM` check (10/10)
- Root `SECURITY.md` publishes the vulnerability disclosure policy — satisfies the Scorecard `Security-Policy` check (10/10)
- Release provenance attestations are generated for published assets via GitHub Actions — improves the Scorecard `Signed-Releases` check

## Canonical Required Check

- Workflow name: Repository Security
- Required job/check name: OSSF Scorecard
- Branch protection policy target: default branch
- Merge pass condition: successful workflow completion (no numeric score threshold)
- Workflow-level token permissions: none by default; permissions are granted only at job scope

## Branch Protection Mapping

1. Open repository Settings -> Branches.
2. Edit protection rule for the default branch.
3. Add required status check: OSSF Scorecard.
4. Enable merge blocking until required checks pass.

## Validation Steps

### Scorecard

1. Open a pull request and verify the OSSF Scorecard check appears.
2. Verify the same workflow runs only on push to the default branch.
3. Confirm that a non-success Scorecard result blocks merge while branch protection is enabled.
4. Confirm the repository root contains `SECURITY.md`.
5. After the next published release, confirm the Scorecard `Signed-Releases` check is non-zero.

### Security Policy

1. Verify `SECURITY.md` exists at the repository root.
2. Confirm it documents a private reporting channel, supported versions, and a response timeline.
3. Confirm GitHub shows the repository Security Policy link after merge.

### Release Provenance

1. Confirm `.github/workflows/publish-release.yml` grants `id-token: write` and `attestations: write` at job scope.
2. Confirm the workflow uses `actions/attest-build-provenance` pinned by SHA.
3. After the next release, verify the released binary with `gh attestation verify <artifact> --repo matteobisi/apple-container-tui`.
4. Repeat verification for the SBOM asset.

### Dependabot

1. Verify .github/dependabot.yml parses in GitHub.
2. Confirm rules exist for gomod and github-actions.
3. Confirm both rules use monthly interval.
4. Confirm Dependabot PRs are generated when updates are available.

### Concurrent Operation (FR-008)

1. Confirm Scorecard continues to run for PRs after Dependabot is enabled.
2. Confirm Dependabot continues to open PRs after Scorecard workflow is enabled.
3. Confirm neither automation requires disabling the other.

## Monthly Review Runbook

1. Review the latest default-branch Scorecard run and note any failures.
2. Triage open Dependabot PRs by ecosystem.
3. Merge or close each Dependabot PR with rationale.
4. Re-run/verify required checks on merged update PRs.
5. Record exceptions or postponements in project notes.

## Troubleshooting

### Scorecard check missing from required checks list

- Ensure at least one workflow run has completed for the default branch.
- Confirm the required check name matches exactly: OSSF Scorecard.
- Confirm workflow file remains at .github/workflows/scorecard.yml.

### StepSecurity reports broad token permissions

- Keep workflow-level `permissions` empty and grant only job-level permissions required by Scorecard.
- Avoid setting `security-events: write` at workflow scope.
- Re-run the workflow after permission changes and verify the SARIF upload step still succeeds.

### Scorecard or security review flags unpinned actions

- Pin third-party and GitHub-maintained actions by full commit SHA instead of floating tags.
- Keep a trailing comment with the human-readable tag for maintenance, for example `# v4.2.2`.
- When updating an action, refresh both the SHA and the tag comment together.

### Security policy file not detected

- Confirm the file is named exactly `SECURITY.md`.
- Confirm it is located at the repository root, `.github/SECURITY.md`, or `docs/SECURITY.md`.
- After merge, wait for the next Scorecard run or trigger one manually.

### Attestation verification fails

- Confirm the `Publish Release` workflow granted both `id-token: write` and `attestations: write`.
- Confirm the release run created a new release rather than skipping due to idempotency.
- Verify the asset name passed to `gh attestation verify` matches the published release asset exactly.

### No Scorecard push run on a feature branch

- This is expected: push-triggered Scorecard runs are limited to `main`.
- Pull requests still run the `OSSF Scorecard` check normally.

### Dependabot not opening updates

- Confirm repository-level Dependabot is enabled.
- Confirm ecosystem directories in .github/dependabot.yml are valid.
- Wait for the monthly schedule window or trigger a manual check from GitHub UI.
