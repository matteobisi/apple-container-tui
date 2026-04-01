# Security Automation

This document defines repository security automation operations for OSSF Scorecard and Dependabot.

## Scope

- OSSF Scorecard workflow: .github/workflows/scorecard.yml
- Dependabot config: .github/dependabot.yml
- Ecosystems in scope: gomod, github-actions
- Dependabot cadence: monthly for security and version updates

## Canonical Required Check

- Workflow name: Repository Security
- Required job/check name: OSSF Scorecard
- Branch protection policy target: default branch
- Merge pass condition: successful workflow completion (no numeric score threshold)

## Branch Protection Mapping

1. Open repository Settings -> Branches.
2. Edit protection rule for the default branch.
3. Add required status check: OSSF Scorecard.
4. Enable merge blocking until required checks pass.

## Validation Steps

### Scorecard

1. Open a pull request and verify the OSSF Scorecard check appears.
2. Verify the same workflow runs on push to the default branch.
3. Confirm that a non-success Scorecard result blocks merge while branch protection is enabled.

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

### Scorecard job skipped on push

- Confirm push is to the repository default branch.
- Check workflow condition logic in the ossf-scorecard job.

### Dependabot not opening updates

- Confirm repository-level Dependabot is enabled.
- Confirm ecosystem directories in .github/dependabot.yml are valid.
- Wait for the monthly schedule window or trigger a manual check from GitHub UI.
