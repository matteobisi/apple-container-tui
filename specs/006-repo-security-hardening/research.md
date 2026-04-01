# Research: Repository Security Hardening

**Feature**: 006-repo-security-hardening  
**Date**: 2026-04-01  
**Status**: Complete

## Research Tasks

### 1. OSSF Scorecard Workflow Integration Strategy

**Decision**: Add a dedicated GitHub Actions workflow at `.github/workflows/scorecard.yml` that runs Scorecard for pull requests and default-branch pushes, and expose the workflow as the required merge gate.

**Rationale**:
- Matches FR-001/FR-003a and provides explicit check visibility in PR status.
- Keeps security policy as versioned code in-repo, aligned with FR-007.
- A dedicated workflow is easier to audit than embedding Scorecard in unrelated CI jobs.

**Alternatives considered**:
- Run Scorecard only manually: rejected because it does not satisfy continuous enforcement.
- Run Scorecard only on default branch: rejected because PR merge gating requires PR-context checks.
- Use numeric threshold policy immediately: rejected due to clarification selecting workflow-success gating only.

### 2. Merge Gate Policy for Scorecard

**Decision**: Define pass/fail for merge gating as successful Scorecard workflow completion, with no minimum numeric score threshold in this feature.

**Rationale**:
- Directly matches accepted clarification in spec (workflow success only).
- Enables enforcement with low setup friction while still surfacing findings.
- Avoids blocking legitimate PRs during initial rollout where baseline score is not yet tuned.

**Alternatives considered**:
- Require score >= 7.0 or higher: rejected by clarification.
- Keep Scorecard informational only: rejected because merge gating is required.

### 3. Dependabot Ecosystem and Cadence Scope

**Decision**: Configure Dependabot for `gomod` and `github-actions` ecosystems with monthly update cadence for both security and version updates.

**Rationale**:
- Matches FR-004 and accepted clarifications.
- Covers the repository's real dependency surfaces: Go modules and workflow actions.
- Monthly cadence controls PR volume while still providing recurring updates.

**Alternatives considered**:
- Include all detectable ecosystems: rejected to avoid unnecessary noise outside clarified scope.
- Daily or weekly cadence: rejected by clarification due to expected maintenance overhead.

### 4. Branch Protection and Required Check Binding Approach

**Decision**: Ensure the Scorecard workflow exposes a stable check name and document that repository branch protection must require that check before merge.

**Rationale**:
- Required-check enforcement is ultimately controlled in repository settings.
- Stable check naming prevents accidental policy drift when workflow internals change.
- Documentation in contracts/quickstart gives a clear implementation and validation path.

**Alternatives considered**:
- Rely only on workflow existence without branch protection binding: rejected because merge-gate requirement would not be guaranteed.
- Automate branch protection through external tooling: rejected as out of scope for this feature.

### 5. Validation Strategy for Security Automation

**Decision**: Validate through repository-native behavior checks: workflow appears and runs on PR/push, required-check policy blocks merges on failing status, and Dependabot opens monthly update PRs for both ecosystems.

**Rationale**:
- This feature is configuration-driven; behavior validation is more appropriate than runtime unit tests in Go code.
- Aligns directly with measurable outcomes SC-001 through SC-004.
- Keeps verification reproducible for maintainers.

**Alternatives considered**:
- Add Go tests for YAML files: rejected because it adds maintenance complexity for low confidence gain.
- Skip explicit validation steps: rejected because criteria are measurable and should be verified.
