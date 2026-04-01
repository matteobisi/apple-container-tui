# Data Model: Repository Security Hardening

**Feature**: 006-repo-security-hardening  
**Date**: 2026-04-01

## Overview

This feature introduces repository-level security automation artifacts rather than application runtime entities. The model below defines the configuration and status objects maintainers interact with when validating Scorecard and Dependabot behavior.

## Entities

### SecurityCheckPolicy

**Purpose**: Represents merge-gating policy for security checks on pull requests.

**Fields**:
- `checkName` (string, required): Stable Scorecard check name shown in GitHub checks.
- `requiredForMerge` (bool, required): Must be `true` for protected branches in this feature.
- `passCondition` (enum, required): `workflow_success`.
- `numericScoreThreshold` (number, optional): Not set for this feature.
- `appliesTo` ([]string, required): Target branches/patterns under branch protection.

**Validation rules**:
- `checkName` MUST match the check produced by `.github/workflows/scorecard.yml`.
- `requiredForMerge` MUST be true for the default branch protection rule.
- `passCondition` MUST be `workflow_success` per accepted clarification.
- `numericScoreThreshold` MUST be null/absent in this feature scope.

**Relationships**:
- Enforced via repository branch protection settings.
- Depends on `ScorecardWorkflowDefinition` producing the expected check run.

---

### ScorecardWorkflowDefinition

**Purpose**: Represents the GitHub Actions workflow configuration that runs OSSF Scorecard.

**Fields**:
- `path` (string, required): `.github/workflows/scorecard.yml`.
- `triggers` ([]enum, required): `pull_request`, `push` to default branch.
- `permissions` (map, required): Least-privilege permissions required by Scorecard.
- `jobName` (string, required): Human-readable and stable for branch protection binding.
- `resultVisibility` (enum, required): `github_checks`.

**Validation rules**:
- Workflow MUST trigger for PRs and default-branch pushes.
- Job/check naming MUST remain stable across routine edits.
- Permissions MUST be explicit and no broader than needed.

**Relationships**:
- Produces results consumed by `SecurityCheckPolicy`.
- Provides observability for `SecurityPostureResult`.

---

### DependabotUpdatePolicy

**Purpose**: Represents Dependabot configuration scope and schedule for this repository.

**Fields**:
- `path` (string, required): `.github/dependabot.yml`.
- `version` (int, required): Dependabot config schema version.
- `updates` ([]DependabotEcosystemRule, required): One rule per ecosystem.

**Validation rules**:
- MUST include exactly these ecosystems in scope: `gomod`, `github-actions`.
- Each ecosystem rule MUST use monthly schedule.
- Rules SHOULD include package-ecosystem directory and open PR limit aligned with project policy.

**Relationships**:
- Generates `DependencyUpdateProposal` pull requests over time.

---

### DependabotEcosystemRule

**Purpose**: Defines update behavior for one dependency ecosystem.

**Fields**:
- `packageEcosystem` (enum, required): `gomod` or `github-actions`.
- `directory` (string, required): `/` for repo-root files.
- `scheduleInterval` (enum, required): `monthly`.
- `openPullRequestsLimit` (int, optional): PR throttling value.

**Validation rules**:
- `packageEcosystem` MUST be one of the two clarified values.
- `scheduleInterval` MUST be `monthly`.
- `directory` MUST reference the actual file location for that ecosystem.

**Relationships**:
- Child of `DependabotUpdatePolicy`.
- Drives creation of `DependencyUpdateProposal` instances.

---

### DependencyUpdateProposal

**Purpose**: Represents an automated Dependabot pull request proposing dependency updates.

**Fields**:
- `ecosystem` (enum, required): `gomod` or `github-actions`.
- `targetFiles` ([]string, required): Files changed by the update.
- `createdAt` (datetime, required): PR creation timestamp.
- `status` (enum, required): `open`, `merged`, `closed`.

**Validation rules**:
- Proposal MUST map to one configured ecosystem rule.
- Proposal metadata MUST include enough context for maintainers to triage.

**Relationships**:
- Produced by `DependabotUpdatePolicy`.
- Reviewed by maintainers as part of routine repo maintenance.

---

### SecurityPostureResult

**Purpose**: Represents one visible outcome of a Scorecard workflow execution.

**Fields**:
- `workflowRunId` (string, required): GitHub run identifier.
- `commitSha` (string, required): Commit under evaluation.
- `conclusion` (enum, required): `success`, `failure`, `cancelled`, `timed_out`, `neutral`.
- `score` (number, optional): Scorecard score when emitted.
- `findingsSummary` (string, optional): High-level finding context for maintainers.

**Validation rules**:
- `conclusion` MUST be visible in checks UI for each relevant PR/push.
- For merge-gating in this feature, `conclusion` MUST be `success`.

**Relationships**:
- Produced by `ScorecardWorkflowDefinition`.
- Evaluated against `SecurityCheckPolicy` during merge decision.

## State Transitions

### SecurityPostureResult
- `queued` -> `in_progress` -> `success|failure|cancelled|timed_out|neutral`
- Merge eligibility path: only `success` satisfies `SecurityCheckPolicy.passCondition`.

### DependencyUpdateProposal
- `open` -> `merged` when accepted.
- `open` -> `closed` when rejected/deferred.
- New proposals appear in subsequent monthly cycles if updates remain available.

## Test-Relevant Invariants

- Scorecard workflow check must exist and be visible on every PR and default-branch push.
- Required-check merge policy must reject merges when Scorecard check is non-success.
- Dependabot must have exactly two ecosystem rules in this feature: `gomod` and `github-actions`.
- Both ecosystem rules must use `monthly` cadence.
