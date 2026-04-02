# Quickstart: Build-to-Release Automation (Docs + Workflow Scope)

## 1. Preconditions

- Work from branch `010-auto-release-publish`.
- Existing build workflow `Build Binary` is already enabled and producing artifact `actui-linux-amd64`.
- Repository permissions allow release creation from GitHub Actions (`contents: write` where needed).
- Selected GitHub Actions are Node 24 compatible (avoid Node 20-only action versions).
- No application code changes are required for this feature.

## 2. Implement Release Workflow

1. Add or update release workflow under `.github/workflows/`.
2. Configure trigger chain from successful completion of `Build Binary`.
3. Add guard condition to run only on successful build conclusion.
4. Download or resolve the binary artifact from triggering build run.
5. Compute next version tag according to policy:
   - Start at `v0.1.0`
   - Increment patch for each new automated release.
6. Ensure idempotency:
   - Prevent duplicate publication for reruns of same source commit/tag.
7. Publish release and attach `actui-linux-amd64` asset.
8. Emit stage logs for artifact resolution, version selection, and publication result.

## 3. Update Documentation

1. Update `docs/binary-build-automation.md` with:
   - Build-to-release trigger chain
   - Version-labeling policy and examples
   - Duplicate handling/rerun behavior
   - Release publication troubleshooting
2. Update `README.md` and/or `AGENTS.md` references only if maintainer entry points changed.

## 4. Validate

1. Trigger a qualifying build path and confirm build success.
2. Confirm release workflow runs after build completion.
3. Verify a release is created with expected version tag.
4. Verify `actui-linux-amd64` is attached to the published release.
5. Rerun with same source reference and verify duplicate protection behavior.
6. Confirm troubleshooting guidance matches observed logs and failure modes.

## 5. Validation Evidence Format

Capture the following for feature sign-off:

1. Build run URL and conclusion
2. Release workflow run URL and conclusion
3. Published tag and release URL
4. Asset name attached (`actui-linux-amd64`)
5. Duplicate-handling verification result
6. Docs update references (`docs/binary-build-automation.md`, and `README.md`/`AGENTS.md` if changed)

## 6. Operator Notes

- This feature intentionally excludes changes to app runtime behavior.
- If version strategy changes (e.g., minor bump cadence), update policy docs and contract before implementation changes.

## 7. Prepare for Tasks Phase

- Confirm `plan.md`, `research.md`, `data-model.md`, `contracts/`, and `quickstart.md` are complete.
- Proceed to `/speckit.tasks` to generate implementation backlog.
