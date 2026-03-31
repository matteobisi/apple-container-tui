# Research: Expanded Container Workflows

**Feature**: 005-expand-container-workflows  
**Date**: 2026-03-31  
**Status**: Complete

## Research Tasks

### 1. Registries Screen Data Source and Display Contract

**Decision**: Use `container registry list --format json` as the sole data source for the Registries screen and display `hostname` as the primary identifier with `username` as secondary context.

**Rationale**:
- The local CLI exposes `container registry list` specifically for registry logins and supports stable JSON output.
- This matches the clarified scope: show registries explicitly configured or exposed by the current runtime, not inferred from image history.
- `hostname` and `username` are enough to distinguish entries without inventing app-owned registry metadata.
- Sample local output confirms the JSON payload shape: `[ { "hostname": "registry.sighup.io", "username": "...", "createdDate": ..., "modifiedDate": ... } ]`.

**Alternatives considered**:
- Infer registries from image names or pull history: rejected because it violates the clarified source-of-truth decision and can show stale or unrelated registries.
- Persist registry records inside actui config: rejected because it creates app-owned state for data the CLI already manages.
- Surface raw created/modified timestamps prominently: rejected for v1 because the timestamp semantics are not yet normalized in the UI and are not required for identifying entries.


### 2. Container Export Must Stay Command-Safe Despite Requiring Multiple CLI Steps

**Decision**: Implement container export as a narrow workflow that previews and executes this exact sequence for stopped containers: `container export --image <generated-image> <container-id>`, then `container image save --output <generated-tar> <generated-image>`, then best-effort `container image rm <generated-image>` cleanup.

**Rationale**:
- The Apple Container CLI does not export containers directly to a file path. It exports a container state to an image, and separately supports saving an image to an OCI-compatible tar archive.
- This sequence satisfies the clarified UX requirement: user chooses a destination directory, the app generates the archive filename, and the feature produces a local export artifact.
- Using a generated intermediate image name keeps the flow deterministic and avoids asking the user to invent an implementation detail.
- Best-effort cleanup avoids leaving temporary export images behind in normal cases while still preserving a successful tar export if cleanup fails.

**Alternatives considered**:
- Export directly to a destination path in one command: rejected because the CLI does not support it.
- Reuse the existing file picker by generalizing it into a directory picker: rejected for this feature because it widens the UI surface unnecessarily and changes prior architecture more than a dedicated export screen would.
- Skip cleanup of the generated intermediate image: rejected because it leaves behind unexpected image clutter after a file export workflow.


### 3. Build `--pull` Should Be an Additive Builder and Screen Change

**Decision**: Extend `BuildScreen` and `BuildImageBuilder` with a boolean `PullLatest` option defaulted to `true`, and include `--pull` in the generated command only when that option is enabled.

**Rationale**:
- The local CLI explicitly supports `container build --pull`.
- The current build flow already uses one screen state object and one builder; adding a boolean is the smallest change that preserves the existing preview and execution path.
- Defaulting to enabled matches the accepted clarification and makes the safer freshness behavior visible before submission.
- Omitting `--pull` when disabled keeps the command preview honest and easy to reason about.

**Alternatives considered**:
- Persist the last checkbox value in config: rejected because the spec explicitly wants a deterministic default and does not require new persisted preferences.
- Always include `--pull` regardless of UI state: rejected because it would make the checkbox misleading and violate command preview accuracy.


### 4. Daemon Status Should Use Structured JSON Instead of Text Heuristics

**Decision**: Switch daemon status collection to `container system status --format json` and classify state from structured fields, primarily `status`, with an explicit `unknown` fallback when required fields are missing or unrecognized.

**Rationale**:
- The CLI supports a stable JSON output contract and local sampling confirms fields such as `status`, `installRoot`, `appRoot`, and `apiServerVersion`.
- Current parsing uses string matching (`running` and `not running`), which is fragile and cannot represent `unknown` accurately.
- Structured parsing enables clear status coloring in the UI while keeping version metadata available for display or debugging.
- Falling back to `unknown` is safer than guessing from partial text and aligns with the clarified requirement.

**Alternatives considered**:
- Improve string matching on the table output: rejected because it remains brittle and loses typed fields the CLI already provides.
- Support text and JSON paths equally: rejected for initial implementation because it doubles parser complexity without a demonstrated need.


### 5. Preserve the Existing Architecture by Adding One Focused Workflow Service

**Decision**: Keep the current `ActiveScreen` routing, command builder/parser service pattern, and `CommandExecutor` interface intact. Add one focused export workflow service to plan and run the multi-command export sequence instead of introducing a general batch executor abstraction.

**Rationale**:
- The repo already cleanly separates UI state from command composition and output parsing.
- Registries, build pull, and daemon status all fit naturally into existing patterns with small builders/parsers/screen changes.
- Export is the only feature that crosses the single-command boundary; a narrow workflow service handles that case without forcing every executor and preview flow to understand command batches globally.
- This matches the user's goal of adding features without undoing the existing architecture.

**Alternatives considered**:
- Add a generic multi-command executor interface everywhere: rejected because only one feature needs sequencing and the broader refactor adds risk for little gain.
- Handle export by shelling out through a compound string command: rejected because it weakens command preview fidelity and bypasses existing typed command builders.


## Summary of Design Decisions

| Area | Decision | Why |
|------|----------|-----|
| Registries | `container registry list --format json`; show hostname + username | Matches runtime-owned source of truth with stable JSON |
| Container export | Export -> save -> cleanup command sequence with generated names | Satisfies destination-directory UX while staying command-safe |
| Build form | Add `PullLatest` bool defaulted to true | Small additive change to current screen/builder pattern |
| Daemon status | Use `container system status --format json` | Replaces brittle text parsing with typed status classification |
| Architecture | Add one export workflow service, keep other seams unchanged | Preserves current design and limits refactor scope |

**State transitions**:

## Derived Naming Rules

### Generated Export Image Reference

**Pattern**: `actui-export/<container-slug>:<timestamp>`

**Rules**:

### Generated Archive Path

**Pattern**: `<destination-directory>/<container-slug>-<timestamp>.oci.tar`

**Rules**:

## View-State Notes


## Test-Relevant Invariants
- Running containers MUST NOT expose an export option in the submenu.
- Disabling `PullLatest` MUST remove `--pull` from the previewed build command.
- Missing `status` in daemon JSON MUST produce `unknown`, never `running` or `stopped`.
- Registry parsing MUST tolerate empty arrays and missing optional fields without crashing.