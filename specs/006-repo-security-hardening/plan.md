# Implementation Plan: Repository Security Hardening

**Branch**: `006-repo-security-hardening` | **Date**: 2026-04-01 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/006-repo-security-hardening/spec.md`

## Summary

Increase repository security posture by adding two automated GitHub-native controls: OSSF Scorecard as a required pull-request check and Dependabot updates for Go modules plus GitHub Actions on a monthly cadence. Keep implementation strictly repository-configuration based (`.github/workflows` and `.github/dependabot.yml`), with verification through pull request checks and generated dependency update PR behavior.

## Technical Context

**Language/Version**: YAML (GitHub Actions workflow and Dependabot configuration), Go 1.21 module context for dependency ecosystem detection  
**Primary Dependencies**: GitHub Actions runner (`ubuntu-latest`), OSSF Scorecard GitHub Action, Dependabot version updates engine  
**Storage**: Git repository configuration files only; no runtime data storage  
**Testing**: GitHub Actions check execution on pull requests/default branch, config lint via workflow parser in GitHub UI, and repository-level validation through observed check/PR generation behavior  
**Target Platform**: GitHub-hosted repository and workflow execution environment  
**Project Type**: Single-project CLI TUI repository with CI/security automation  
**Performance Goals**: Security check feedback available within standard PR check cycle; Dependabot update jobs execute within monthly schedule window  
**Constraints**: Preserve local-only app runtime model; do not introduce telemetry or cloud runtime dependencies beyond GitHub repository automation; keep security automation config versioned and reviewable  
**Scale/Scope**: Add one new security workflow file, one Dependabot config file, and related documentation/contracts for one repository

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Command-Safe TUI
✅ **COMPLIANT** - This feature does not alter TUI command execution or destructive-action paths. Existing command-safe behavior remains unchanged.

### Principle II: macOS 26.x + Apple Silicon Target
✅ **COMPLIANT** - Runtime target for the application is unchanged. Repository automation additions do not affect local macOS compatibility commitments.

### Principle III: Local-Only Operation
✅ **COMPLIANT** - Application runtime remains local-only; added GitHub automation applies to repository governance, not to app runtime telemetry or remote control surfaces.

### Principle IV: Clear Observability
✅ **COMPLIANT** - Scorecard check status becomes visible in pull request checks, and Dependabot update proposals become visible through standard pull request flows.

### Principle V: Tested Command Contracts
✅ **COMPLIANT** - No new app command builders/parsers are introduced. Governance requires explicit action mapping for feature behavior; this will be captured in a contract artifact for repository automation workflows.

### Workflow and Quality Gates
✅ **COMPLIANT** - Command/action mapping will be documented in `contracts/security-automation.md` as allowed by the constitution (spec-level table or contract document).

**GATE RESULT**: ✅ **PASS**

**Post-Phase 1 Re-check**: ✅ **PASS** - Phase 1 artifacts document exact workflow/config contracts and validation flow without constitutional violations.

## Project Structure

### Documentation (this feature)

```text
specs/006-repo-security-hardening/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── security-automation.md
└── tasks.md
```

### Source Code (repository root)

```text
.github/
├── dependabot.yml                 # NEW dependabot config for Go modules + GitHub Actions
└── workflows/
    └── scorecard.yml              # NEW OSSF Scorecard required check workflow

go.mod                             # Existing dependency manifest for Go ecosystem
```

**Structure Decision**: Keep the current single-project layout and implement changes as additive repository automation configuration under `.github/`, with no changes to runtime TUI code paths in `src/`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
