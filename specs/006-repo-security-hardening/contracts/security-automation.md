# Security Automation Contract: Repository Hardening

**Feature**: 006-repo-security-hardening  
**Date**: 2026-04-01

## Overview

This contract defines repository-level automation required to satisfy security hardening behavior in the spec: OSSF Scorecard as a required merge-gating check and Dependabot monthly updates for Go modules and GitHub Actions.

## 1. Scorecard Workflow Contract

### Configuration Contract
- File MUST exist at `.github/workflows/scorecard.yml`.
- Workflow MUST trigger on:
  - pull requests
  - pushes to the default branch
- Workflow MUST expose a stable job/check identity suitable for branch protection required-check binding.
- Canonical required-check name MUST be `OSSF Scorecard`.
- Workflow permissions MUST be explicitly declared and least-privilege.
- Workflow-level token permissions SHOULD be empty or more restrictive than job-level permissions; write scopes MUST be granted only to the job that needs them.
- External GitHub Actions references SHOULD be pinned to immutable commit SHAs rather than floating version tags.

### Execution Contract
- Every pull request MUST show a Scorecard check run.
- Every default-branch push MUST show a Scorecard check run.
- Workflow success is the only pass condition required by this feature (no numeric score threshold).

### Merge Gate Contract
- Branch protection for the default branch MUST require the Scorecard check to pass before merge.
- Branch protection for the default branch MUST require the `OSSF Scorecard` check to pass before merge.
- Any non-success conclusion (`failure`, `cancelled`, `timed_out`, `neutral`) MUST block merge until a successful run is present.

### Observability Contract
- Maintainers MUST be able to inspect Scorecard check status from standard GitHub checks UI.
- Workflow logs MUST be available through the run details for troubleshooting.

## 2. Dependabot Configuration Contract

### Configuration Contract
- File MUST exist at `.github/dependabot.yml`.
- Config schema version MUST be declared as v2.
- Exactly these ecosystems are in scope for this feature:
  - `gomod`
  - `github-actions`
- Both ecosystem rules MUST use monthly scheduling.

### Proposal Contract
- Dependabot MUST open update pull requests when updates are available for each configured ecosystem.
- Update PRs MUST identify the ecosystem and changed dependency context so maintainers can triage.
- Update cadence target is one run per ecosystem per 31-day cycle.

## 3. Non-Goals and Boundaries

- No additional ecosystems are required in this feature.
- No minimum Scorecard numeric threshold is enforced in this feature.
- No automation of branch protection settings via external tooling is required in this feature.

## 4. Coexistence Contract

- Scorecard workflow and Dependabot update automation MUST operate concurrently.
- Enabling Scorecard MUST NOT disable or impede Dependabot update checks.
- Enabling Dependabot MUST NOT disable or impede Scorecard check execution.

## 5. Validation Contract

### Required Verification Steps
1. Confirm `.github/workflows/scorecard.yml` exists and is valid in GitHub Actions.
2. Confirm Scorecard workflow runs on PR and default-branch push events.
3. Confirm branch protection includes Scorecard as required check and blocks merge on non-success.
4. Confirm `.github/dependabot.yml` exists and validates.
5. Confirm Dependabot rule entries exist for both `gomod` and `github-actions` with monthly schedule.
6. Confirm at least one Dependabot PR can be produced in each ecosystem when updates are available.
7. Confirm Scorecard and Dependabot both run successfully in the same implementation window without disabling either automation.
