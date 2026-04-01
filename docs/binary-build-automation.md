# Binary Build Automation

This runbook describes how the repository builds and validates the `actui` binary in CI.

## Workflow Reference

- Workflow file: `.github/workflows/build-binary.yml`
- Workflow name: `Build Binary`
- Build command: `go build -o actui ./cmd/actui`

## Qualifying Triggers

- Push to `main` after merge.
- Manual dispatch for maintainer verification/recovery.

Qualifying updates include merged feature work and merged dependency updates.

## Artifact Contract

- Artifact name format: `actui-<os>-<arch>`
- Current artifact: `actui-linux-amd64`
- Retention policy: explicit `retention-days` configured in workflow

## Action Dependencies

The workflow uses GitHub Actions pinned to specific versions (Node 24 compatible as of April 2026) for security and stability:
- `actions/checkout@v4.1.7` (Node 24 compatible)
- `actions/setup-go@v5.0.0` (Node 24 compatible)
- `actions/upload-artifact@v4.4.0` (Node 24 compatible)

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