# Contract: Binary Build Automation Workflow

## Purpose

Define the expected contract for repository automation that builds and publishes the `actui` binary after qualifying merged changes, and establish the documentation obligations for maintainers.

## Workflow Contract

- Workflow Name: `Build Binary`
- Required Trigger Events:
  - Push to `main` after pull request merge.
  - Manual dispatch by maintainers for verification/recovery.
- Qualifying merged updates:
  - Feature changes merged via pull request.
  - Dependency updates (automated or manual) merged via pull request.
- Required Behavior:
  - Checkout repository source at triggering commit.
  - Setup Go toolchain compatible with project requirements.
  - Build binary from `./cmd/actui`.
  - Publish artifact(s) from successful runs with stable naming.
  - Configure explicit artifact retention window for audit and troubleshooting workflows.
  - Expose run logs and status for maintainer diagnosis.
- Security Contract:
  - Use least-privilege token permissions.
  - Do not request broad write scopes unless explicitly justified in future plans.

## Inputs

- Repository source at trigger commit.
- Workflow trigger metadata (event type, actor, commit SHA).

## Outputs

- Success or failure status for each run.
- Binary artifact package for successful runs.
- Log output sufficient for troubleshooting failures.

## Acceptance Contract

- A merged feature change to `main` results in one successful build run and artifact availability.
- A merged dependency update to `main` results in one successful build run and artifact availability.
- Failed runs present actionable error context in run logs.
- Retention configuration is visible in workflow definition and validated via run evidence.
- Release sign-off includes documented manual verification on macOS 26.x.

## Validation Evidence Format

Each release-candidate validation record must include:

- Workflow run URL
- Trigger type (`push_main` or `manual_dispatch`)
- Commit SHA
- Artifact name and `retention-days` value
- Manual verification status on macOS 26.x (`pass` or `fail`) and notes

## Documentation Contract

A maintainer-facing document under `docs/` must include:

- Trigger behavior and when automation runs.
- Expected artifacts and where to retrieve them.
- Validation checklist for release readiness.
- Troubleshooting steps for common build failures.
- Reference validation machine profile:
  - Macbook M4
  - macOS 26.4
  - 32GB RAM
  - Apple container 0.10.0

## Non-Goals

- Defining release tagging/publishing policy.
- Replacing local developer build workflows.
- Introducing remote telemetry or external orchestration dependencies.