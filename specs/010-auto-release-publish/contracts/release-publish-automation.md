# Contract: Automated Release Publication Workflow

## Purpose

Define the contract for publishing GitHub Releases automatically after successful binary build runs, including version-labeling behavior and documentation obligations.

## Workflow Contract

- Build Workflow Dependency: `Build Binary` must complete successfully before release publication is attempted.
- Release Trigger: workflow chaining from build completion event (`workflow_run`) with success gating.
- Required Behavior:
  - Resolve the triggering build run and artifact metadata.
  - Determine release tag from version-labeling policy.
  - Prevent duplicate publication for already released source/tag.
  - Create one release record and upload one binary asset per qualified run.
  - Emit actionable logs for each stage (artifact lookup, version resolution, publish step).

## Version Labeling Contract

- Initial automated release tag: `v0.1.0`
- Increment strategy: patch increment for each subsequent automated release (`v0.1.1`, `v0.1.2`, ...)
- Duplicate handling:
  - If computed tag already exists, apply documented conflict strategy and log outcome.
  - Publication MUST remain idempotent for reruns of the same source commit.

## Inputs

- Successful `Build Binary` run metadata
- Artifact name and retrieval reference
- Existing release/tag state from repository
- Version-labeling policy configuration

## Outputs

- Published GitHub Release with deterministic tag and title
- Uploaded binary release asset (`actui-linux-amd64`)
- Workflow logs containing run linkage, computed version, and publication result

## Security Contract

- Use least-privilege permissions.
- `contents: write` is allowed only where required to create releases and upload assets.
- No additional external secrets required unless explicitly documented.
- Selected GitHub Actions must be Node 24 compatible; Node 20-only action versions are not acceptable for this workflow.

## Acceptance Contract

- A successful qualifying build produces exactly one release publication attempt.
- A failed build produces no release publication.
- Duplicate reruns do not create duplicate release records for the same source commit.
- Published release contains the expected binary asset.
- Logs and docs provide enough context for maintainers to diagnose failures.

## Documentation Contract

The runbook at `docs/binary-build-automation.md` must include:

- End-to-end trigger flow from build to release publication
- Version-labeling policy and examples
- Duplicate handling behavior and rerun expectations
- Failure troubleshooting path for release publication stage
- Operator checklist for validation and sign-off

When maintainer entry points change materially, `README.md` and/or `AGENTS.md` must include pointers to the updated runbook.

## Non-Goals

- Modifying application runtime code, TUI screens, or command handlers
- Multi-platform packaging expansion beyond current binary artifact
- Manual release note authoring workflows outside automated baseline generation
