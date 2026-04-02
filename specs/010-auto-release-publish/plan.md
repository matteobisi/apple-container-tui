# Implementation Plan: Automated Binary Release Publishing

**Branch**: `010-auto-release-publish` | **Date**: 2026-04-02 | **Spec**: `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/010-auto-release-publish/spec.md`
**Input**: Feature specification from `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/010-auto-release-publish/spec.md`

## Summary

Add release automation that publishes a GitHub Release after successful binary build output is available, apply deterministic semantic version labeling (`0.1.0` start, patch increments for subsequent automated releases), and update maintainer documentation for build-to-release flow. Scope is automation and docs only; no application runtime/UI changes.

## Technical Context

**Language/Version**: GitHub Actions YAML workflows plus shell scripting on Ubuntu runners  
**Primary Dependencies**: `actions/checkout@v4`, `actions/setup-go@v5`, `actions/upload-artifact@v4`, GitHub CLI/API release actions (`actions/download-artifact`, `softprops/action-gh-release` or equivalent)  
**Storage**: GitHub Actions artifact storage and GitHub Releases assets; repository Markdown docs in `docs/`  
**Testing**: GitHub Actions run verification, dry-run/simulation via `workflow_dispatch`, and release evidence checks  
**Target Platform**: GitHub-hosted runners (`ubuntu-latest`) with release output consumed by macOS users  
**Project Type**: Existing Go CLI/TUI repository with CI/CD workflow extension only  
**Performance Goals**: Release publication completes within 10 minutes after successful qualifying build run  
**Constraints**: No app code changes, no duplicate release per source commit, explicit logs for failure diagnostics, preserve least-privilege permissions, and use GitHub Actions that are Node 24 compatible due to Node 20 deprecation timeline  
**Scale/Scope**: Add release-publish workflow integration and docs updates in `docs/binary-build-automation.md` and potentially `README.md`/`AGENTS.md`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Pre-Research Gate Assessment

- Principle I (Command-Safe TUI): PASS. No CLI/TUI command behavior is changed.
- Principle II (macOS 26.x + Apple Silicon): PASS. Release docs keep manual macOS verification as pre-release requirement.
- Principle III (Local-Only Operation): PASS. No telemetry or remote control surface is added to app runtime.
- Principle IV (Clear Observability): PASS. Release workflow requires actionable logs and troubleshooting guidance.
- Principle V (Tested Command Contracts): PASS. Workflow contract and run validation are explicitly planned.
- Workflow and Quality Gates: PASS. Command-mapping contract is provided via feature contract document and maintainer docs.

### Post-Design Gate Assessment

- PASS. Phase 0 and Phase 1 artifacts preserve constitution constraints and keep feature scope to automation and documentation.

## Project Structure

### Documentation (this feature)

```text
/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/010-auto-release-publish/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── release-publish-automation.md
└── tasks.md
```

### Source Code (repository root)

```text
.github/
└── workflows/
    ├── build-binary.yml
    └── [release publication workflow or integrated release job]

docs/
└── binary-build-automation.md

README.md
AGENTS.md

cmd/
└── actui/

src/
├── models/
├── services/
└── ui/

tests/
├── contract/
├── integration/
└── unit/
```

**Structure Decision**: Reuse existing single-project structure and implement release automation in GitHub workflows plus docs. No application source directories require modification for this feature.

## Phase 0 Research Output

Research decisions are in `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/010-auto-release-publish/research.md`.

## Phase 1 Design Output

- Data model: `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/010-auto-release-publish/data-model.md`
- Contract: `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/010-auto-release-publish/contracts/release-publish-automation.md`
- Quickstart: `/Users/matteo.bisi/GitRepo/Personal/apple-container-tui/specs/010-auto-release-publish/quickstart.md`

## Complexity Tracking

No constitution violations require justification.
