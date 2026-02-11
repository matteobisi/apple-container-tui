<!--
Sync Impact Report
- Version: UNSET -> 0.1.0
- Modified principles:
	- Template placeholder 1 -> I. Command-Safe TUI
	- Template placeholder 2 -> II. macOS 26.x + Apple Silicon Target
	- Template placeholder 3 -> III. Local-Only Operation
	- Template placeholder 4 -> IV. Clear Observability
	- Template placeholder 5 -> V. Tested Command Contracts
- Added sections: None
- Removed sections: None
- Templates requiring updates:
	- .specify/templates/plan-template.md ✅ unchanged
	- .specify/templates/spec-template.md ✅ unchanged
	- .specify/templates/tasks-template.md ✅ updated
	- .specify/templates/commands/*.md ⚠ not found
- Deferred TODOs: None (resolved 2026-02-11)
-->
# Container TUI Constitution

## Core Principles

### I. Command-Safe TUI
All user actions MUST map to Apple Container CLI commands with a visible command
preview. Destructive actions (delete, stop daemon, delete stopped containers)
MUST require explicit confirmation and MUST support a dry-run mode that shows
the exact command without executing it.

### II. macOS 26.x + Apple Silicon Target
The project MUST run on macOS 26.x on Apple Silicon (M-series). Cross-platform
support is out of scope unless the constitution is amended.

### III. Local-Only Operation
The TUI MUST operate locally with no telemetry, no cloud dependencies, and no
remote control surface. All state must be stored locally in user space.

### IV. Clear Observability
Every action MUST surface success, failure, stdout, and stderr in a readable
format. Log output MUST be retained locally for troubleshooting, with paths
documented in the UI or help text.

### V. Tested Command Contracts
Command composition, argument validation, and destructive-action guardrails MUST
have automated tests. Integration tests MUST cover at least start, stop, pull,
and delete flows against a local Apple Container environment when feasible.

## Platform and Runtime Constraints

- The CLI backend MUST be Apple Container, following the official command
	reference for syntax.
- The TUI MUST be able to start and stop the Apple Container daemon.
- Configuration files and logs MUST live in user-writable locations and MUST
	NOT require elevated privileges for normal operation.
- No background services beyond Apple Container itself.

## Workflow and Quality Gates

- Feature specs MUST include a command mapping table for each user action.
- Any new destructive action MUST include a dry-run path and confirmation copy.
- Manual verification MUST be performed on macOS 26.x before release.
- Releases MUST follow semantic versioning with clear changelogs.

## Governance

- This constitution supersedes other project practices.
- Amendments require updating this file, recording rationale, and noting any
	migration or behavior changes.
- All plans and reviews MUST include a constitution compliance check.
- Exceptions require written justification in the implementation plan's
	Complexity Tracking section.

**Version**: 0.1.0 | **Ratified**: 2026-02-11 | **Last Amended**: 2026-02-11
