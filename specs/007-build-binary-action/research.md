# Research: Automated Binary Build Workflow

## Decision 1: Trigger automation on merges to main and manual dispatch

- Decision: Run the binary build workflow on pushes to `main` (which represent merged pull requests) and allow optional manual dispatch for recovery/verification runs.
- Rationale: This captures merged feature and dependency updates while avoiding duplicate runs from pull request events, and preserves an explicit rerun path for maintainers.
- Alternatives considered:
  - Trigger on `pull_request` only: rejected because artifact publication should reflect integrated code, not pre-merge state.
  - Trigger on both `push` and `pull_request`: rejected to avoid duplicate build runs and noisy status reporting.

## Decision 2: Build with repository Go version and produce deterministic artifact naming

- Decision: Use `actions/setup-go` with the repository Go version and output a consistently named binary artifact (`actui-<os>-<arch>`).
- Rationale: Aligns CI builds with project toolchain expectations and makes artifact retrieval predictable for maintainers.
- Alternatives considered:
  - Pin a fixed Go patch version independent from repository setup: rejected to reduce drift from project-defined toolchain requirements.
  - Use ad-hoc artifact names per run: rejected because predictable naming aids operational clarity.

## Decision 3: Keep workflow permissions minimal and explicit

- Decision: Use least-privilege workflow permissions and avoid broad write scopes for this build automation workflow.
- Rationale: Maintains repository security posture and mirrors hardening standards already adopted for existing workflows.
- Alternatives considered:
  - Default token permissions: rejected due to unnecessary scope and reduced auditability.

## Decision 4: Document operational flow in docs with explicit validation machine profile

- Decision: Add a dedicated docs page that explains triggers, outputs, validation checklist, troubleshooting, and reference machine profile (Macbook M4, macOS 26.4, 32GB RAM, Apple container 0.10.0).
- Rationale: Satisfies maintainability requirement and provides concrete reproducibility context for future maintainers.
- Alternatives considered:
  - Inline docs only in README: rejected because operational runbook content is too detailed for README scope.

## Decision 5: Validate with workflow run checks plus local build parity command

- Decision: Define validation around successful GitHub workflow execution and local parity command `go build -o actui ./cmd/actui`.
- Rationale: Couples CI confidence with the documented local verification path already used by maintainers.
- Alternatives considered:
  - CI-only validation: rejected because local parity steps improve troubleshooting and onboarding.

## Decision 6: Configure explicit artifact retention for auditability

- Decision: Set an explicit artifact retention window in the workflow and validate retention behavior in run evidence.
- Rationale: FR-005 requires maintainers to audit recent build behavior and investigate regressions.
- Alternatives considered:
  - Rely on GitHub default retention: rejected because default values can change and are less explicit for auditing.

## Decision 7: Enforce manual macOS 26.x verification before release sign-off

- Decision: Add a required manual verification step on macOS 26.x and capture evidence in quickstart/release notes.
- Rationale: The constitution requires manual verification on macOS 26.x before release.
- Alternatives considered:
  - CI-only validation: rejected because it would violate constitution quality gates.