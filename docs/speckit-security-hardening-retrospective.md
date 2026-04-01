# Speckit to GitHub Hardening: End-to-End Security Upgrade for apple-container-tui

## Why We Did This

This repository started with no GitHub-native repository security automation and no branch protection on `main`.

The objective was to raise the baseline in a practical way:

- add OSSF Scorecard as a required merge check
- enable Dependabot for real dependency surfaces (`gomod`, `github-actions`)
- tune workflow permissions to least privilege
- enforce everything with branch protection on `main`

## Step 1: Speckit Feature `006-repo-security-hardening`

We used Speckit to drive the work from specification through implementation:

- spec: `specs/006-repo-security-hardening/spec.md`
- plan: `specs/006-repo-security-hardening/plan.md`
- research: `specs/006-repo-security-hardening/research.md`
- data model: `specs/006-repo-security-hardening/data-model.md`
- contract: `specs/006-repo-security-hardening/contracts/security-automation.md`
- quickstart: `specs/006-repo-security-hardening/quickstart.md`
- tasks: `specs/006-repo-security-hardening/tasks.md`

Main decisions captured in the spec clarifications:

- Scorecard must be required before merge
- Dependabot scope: `gomod` + `github-actions`
- Dependabot cadence: monthly
- Scorecard pass condition: workflow success (no numeric threshold)

Implementation landed in PR #1:

- https://github.com/matteobisi/apple-container-tui/pull/1

## Step 2: Initial Implementation Output

We introduced:

- `.github/workflows/scorecard.yml`
- `.github/dependabot.yml`
- `docs/security-automation.md`

The first pass worked functionally, but security tooling flagged an improvement area.

## Step 3: Follow-up Hardening `007-scorecard-hardening-followup`

After merge, GitHub security checks flagged broad token scope at workflow level.

Follow-up PR #10 addressed that:

- https://github.com/matteobisi/apple-container-tui/pull/10

Hardening changes:

1. Workflow-level token permissions set to empty (`permissions: {}`)
2. Required write permissions kept only at job scope
3. Actions pinned to immutable SHAs
4. Trigger noise removed by running `push` only on `main`

## Step 4: Branch Protection Enabled with GitHub CLI

Repository files alone are not enough; branch protection settings enforce merge policy.

Using `gh api`, we enabled protection on `main` with:

- required status check: `OSSF Scorecard`
- strict checks enabled
- PR review policy enabled
- required conversation resolution enabled
- force-push disabled
- branch deletion disabled

## Step 5: Dependabot Security Updates Enabled

Dependabot configuration existed in repo, but GitHub-level security updates needed explicit enablement.

We enabled:

- vulnerability alerts
- automated security fixes (Dependabot security updates)

Current live `security_and_analysis` baseline includes:

- `dependabot_security_updates`: enabled
- `secret_scanning`: enabled
- `secret_scanning_push_protection`: enabled

## Final Security Posture Snapshot

Strong points:

- `main` is protected
- `OSSF Scorecard` is enforced as required check
- Scorecard workflow is least-privilege and SHA-pinned
- Dependabot updates are configured for both relevant ecosystems
- Dependabot security updates are enabled
- secret scanning and push protection are enabled

Open policy choice (not a defect):

- required approvals currently set to `0`

For a solo-maintainer project, `0` is pragmatic. Move to `1` once there is at least one additional independent reviewer.

## Key Lessons

1. Start from specification, not from YAML edits.
2. Treat post-merge security findings as normal tuning, not failure.
3. Keep automation config and repository settings in sync.
4. Use `gh` for repeatable, auditable enforcement steps.

## Artifacts to Review

- Workflow: `.github/workflows/scorecard.yml`
- Dependabot: `.github/dependabot.yml`
- Ops guide: `docs/security-automation.md`
- Feature contract: `specs/006-repo-security-hardening/contracts/security-automation.md`
- Task closeout: `specs/006-repo-security-hardening/tasks.md`
