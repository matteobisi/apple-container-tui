# Quickstart: Binary Build Automation and Documentation

## 1. Preconditions

- Work from branch `007-build-binary-action`.
- Ensure repository CI can run GitHub Actions workflows.
- Confirm local build baseline works:

```bash
go build -o actui ./cmd/actui
```

## 2. Implement Workflow

1. Add a new workflow file under `.github/workflows/` for binary build automation.
2. Configure triggers:
   - `push` on `main`
   - `workflow_dispatch`
3. Add build steps:
   - checkout
   - setup Go
   - build `./cmd/actui`
   - upload artifact
4. Set explicit artifact retention configuration for uploaded artifacts.
5. Keep workflow permissions minimal and explicit.

## 3. Implement Documentation

1. Add or update a document under `docs/` describing:
   - when workflow runs
   - expected artifacts
   - validation process
   - troubleshooting
2. Include required validation machine profile:
   - Macbook M4, macOS 26.4, 32GB RAM, Apple container 0.10.0

## 4. Validate

1. Trigger workflow from a test merge path (or manual dispatch where needed).
2. Verify run status is successful.
3. Verify artifact is downloadable and correctly named.
4. Verify retention settings are present and match documented policy.
5. Follow docs walkthrough to ensure instructions are complete for another maintainer.
6. Perform and record mandatory manual verification on macOS 26.x before release sign-off.

## 5. Validation Evidence Format

Record the following for release-candidate verification:

1. Workflow run URL
2. Trigger type (`push_main` or `manual_dispatch`)
3. Commit SHA
4. Artifact name and `retention-days`
5. Manual verification result on macOS 26.x and notes

## 6. Latest Manual Verification Record

- Date: 2026-04-01
- Trigger type: local manual verification
- Host OS version: macOS 26.4
- Build command: `go build -o /tmp/actui-verify ./cmd/actui`
- Build result: pass
- Apple container version check: `container 0.10.0`
- Notes: Workflow run URL and commit SHA are pending CI execution from a merged qualifying change.

## 7. Prepare for Tasks Phase

- Confirm plan, research, data model, contract, and quickstart artifacts are complete.
- Proceed to `/speckit.tasks` to generate an ordered implementation backlog.