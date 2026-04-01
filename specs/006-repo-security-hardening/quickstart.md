# Quickstart: Repository Security Hardening

**Feature**: 006-repo-security-hardening  
**Date**: 2026-04-01

## Goal

Implement and verify repository security automation by adding:
- OSSF Scorecard workflow as required merge gate
- Dependabot monthly updates for Go modules and GitHub Actions

## Prerequisites

- Write access to repository settings (to configure required status checks)
- GitHub Actions enabled for the repository
- Dependabot alerts/updates enabled at repository level
- Default branch protection available for update

## Implementation Steps

### 1. Add Scorecard GitHub Action

Create:
- `.github/workflows/scorecard.yml`

Expected behavior:
- Triggers on pull requests and default-branch pushes
- Produces a visible `OSSF Scorecard` check in PR checks

### 2. Add Dependabot Configuration

Create:
- `.github/dependabot.yml`

Expected behavior:
- Contains `gomod` and `github-actions` update rules
- Uses monthly schedule for both ecosystems

### 3. Bind Required Check in Branch Protection

In repository branch protection settings for the default branch:
- Add `OSSF Scorecard` as required status check
- Ensure merge is blocked when check is not successful

## Validation Steps

### Validate Scorecard Workflow

1. Open a test pull request.
2. Confirm the `Repository Security` workflow runs automatically.
3. Confirm `OSSF Scorecard` appears in the PR checks list.
4. Confirm merge cannot proceed if Scorecard check is not successful.

### Validate Dependabot Behavior

1. Confirm Dependabot configuration parses successfully in GitHub.
2. Confirm scheduled monthly checks are registered for `gomod` and `github-actions`.
3. Confirm update PRs are created when updates are available.

### Validate Concurrent Operation

1. Confirm Scorecard keeps running on PRs after Dependabot config is enabled.
2. Confirm Dependabot keeps creating update PRs after Scorecard workflow is enabled.
3. Confirm neither automation requires disabling the other.

## Completion Checklist

- [ ] `.github/workflows/scorecard.yml` added and valid
- [ ] Scorecard check visible on PRs and default-branch pushes
- [ ] Scorecard check enforced as required merge gate
- [ ] `.github/dependabot.yml` added and valid
- [ ] Dependabot rules present for `gomod` and `github-actions`
- [ ] Monthly cadence configured for both ecosystems
- [ ] Scorecard and Dependabot verified to operate concurrently
- [ ] Validation evidence captured in PR description or project notes
