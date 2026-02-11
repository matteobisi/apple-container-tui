<!--
Sync Impact Report
- Version: 0.1.0 → 0.2.0
- Modified principles: None
- Added sections:
	- Project Overview (describes implemented TUI functionality and scope)
- Removed sections: None
- Templates requiring updates:
	- .specify/templates/plan-template.md ✅ unchanged (constitution checks already present)
	- .specify/templates/spec-template.md ✅ unchanged (user scenarios align with principles)
	- .specify/templates/tasks-template.md ✅ unchanged (test requirements align with Principle V)
	- .specify/templates/commands/*.md ✅ N/A (no command files present in repository)
- Deferred TODOs: None
- Amendment rationale: Added Project Overview section to document implemented scope
  and capabilities now that spec 001-apple-container-tui is complete. This provides
  context for governance decisions and ensures future amendments reference a clear
  baseline of what the project does.
-->
# Container TUI Constitution

## Project Overview

**Container TUI** is a keyboard-first terminal user interface for managing Apple
Container operations on macOS 26.x (Apple Silicon). It provides a command-safe
wrapper around the Apple Container CLI, enabling developers to perform container
lifecycle operations, image management, and daemon control without memorizing
command syntax.

### Implemented Capabilities

The TUI implements the following workflows as defined in spec `001-apple-container-tui`:

- **Container Lifecycle**: List containers with status; start, stop, and delete
  containers with command previews and confirmations for destructive actions.
- **Image Operations**: Pull images by reference; build images from Containerfile
  or Dockerfile with auto-detection and progress feedback.
- **Daemon Management**: Start and stop the Apple Container daemon with safety
  confirmations and status visibility.
- **Safety Features**: Dry-run mode for all operations; command preview before
  execution; type-to-confirm for destructive actions; JSONL command logging with
  automatic rotation.

### Scope Boundaries

- **In Scope**: Local Apple Container CLI wrapper with keyboard navigation,
  command safety, and observability.
- **Out of Scope**: Remote container management, cloud integrations, telemetry,
  multi-platform support (Windows/Linux), container orchestration (Kubernetes),
  container networking configuration, volume management beyond CLI defaults.

This overview serves as the baseline for all governance and amendment decisions.

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

**Version**: 0.2.0 | **Ratified**: 2026-02-11 | **Last Amended**: 2026-02-11
