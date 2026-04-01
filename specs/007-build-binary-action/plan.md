# Implementation Plan: Automated Binary Build Workflow

**Branch**: `007-build-binary-action` | **Date**: 2026-04-01 | **Spec**: `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/007-build-binary-action/spec.md`
**Input**: Feature specification from `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/007-build-binary-action/spec.md`

## Summary

Deliver one GitHub Actions workflow that builds `actui` on qualifying merged updates, publishes stable binary artifacts with explicit retention, and adds maintainer documentation that includes troubleshooting and mandatory manual verification on macOS 26.x using the reference machine profile.

## Technical Context

**Language/Version**: Go 1.21+ for build command; GitHub Actions YAML for automation
**Primary Dependencies**: `actions/checkout`, `actions/setup-go`, `actions/upload-artifact`, Go toolchain from `go.mod`
**Storage**: GitHub Actions artifact storage and repository Markdown docs in `docs/`
**Testing**: Workflow run validation in GitHub Actions; local parity command `go build -o actui ./cmd/actui`
**Target Platform**: GitHub-hosted CI runners; manual verification on macOS 26.x Apple Silicon
**Project Type**: Existing Go CLI/TUI project with CI and docs additions
**Performance Goals**: Workflow starts within 5 minutes of qualifying merge and produces artifacts in a single run
**Constraints**: Least-privilege workflow permissions, explicit artifact retention, no regression to required security checks
**Scale/Scope**: Add one workflow file and one docs runbook update plus spec artifacts for this feature

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Pre-Research Gate Assessment

- Principle I (Command-Safe TUI): PASS. No runtime destructive flow changes are introduced.
- Principle II (macOS 26.x + Apple Silicon): PASS. Plan includes required macOS 26.x manual verification.
- Principle III (Local-Only Operation): PASS. No telemetry or remote control surface introduced.
- Principle IV (Clear Observability): PASS. Workflow output and diagnostics are explicit deliverables.
- Principle V (Tested Command Contracts): PASS. Validation tasks cover build command behavior and workflow outcomes.
- Workflow and Quality Gates: PASS with required deliverables: command contract artifact in `contracts/`, and explicit manual macOS verification before release.

### Post-Design Gate Assessment

- PASS. Research and design artifacts include retention strategy and mandatory manual verification coverage.

## Project Structure

### Documentation (this feature)

```text
/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/007-build-binary-action/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── build-automation.md
└── tasks.md
```

### Source Code (repository root)

```text
.github/
└── workflows/
    └── build-binary.yml

docs/
└── binary-build-automation.md

cmd/
└── actui/
    └── main.go

src/
├── models/
├── services/
└── ui/

tests/
├── contract/
├── integration/
└── unit/
```

**Structure Decision**: Use the existing single-project layout. Add one workflow under `.github/workflows/` and maintainer runbook docs under `docs/`.

## Phase 0 Research Output

Research decisions are in `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/007-build-binary-action/research.md`.

## Phase 1 Design Output

- Data model: `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/007-build-binary-action/data-model.md`
- Contract: `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/007-build-binary-action/contracts/build-automation.md`
- Quickstart: `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/007-build-binary-action/quickstart.md`

## Complexity Tracking

No constitution violations require justification.
