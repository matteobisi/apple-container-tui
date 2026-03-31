# Command Contracts: Expanded Container Workflows

**Feature**: 005-expand-container-workflows  
**Date**: 2026-03-31

## Overview

This feature extends the current actui command surface through existing Apple Container CLI operations. Every user-facing action below must map to the exact commands shown here, with preview visibility preserved before execution.

---

## 1. Registries Screen Contract

### User Action
Open the Registries screen from the image-management area of the TUI and refresh the visible registry list.

### Navigation Contract
- `ImageListScreen` adds a dedicated Registries entry point (recommended key: `g`) rather than introducing a new top-level navigation framework.
- `ActiveScreen` adds a `ScreenRegistries` value.
- Registries screen supports `r` to refresh and `esc` to return to the previous screen.

### Command Contract

| Operation | Command | Notes |
|-----------|---------|-------|
| Initial load | `container registry list --format json` | Primary read path |
| Refresh | `container registry list --format json` | Same command reused |

### Output Contract
- Expected output is a JSON array of objects.
- Known fields from local CLI sampling: `hostname`, `username`, `createdDate`, `modifiedDate`.
- Empty array is a valid success case and maps to an empty-state message.

### Failure Contract
- Command failure surfaces the existing formatted error path.
- Malformed JSON surfaces a readable parser error and leaves the screen navigable.

---

## 2. Stopped-Container Export Contract

### User Action
From a stopped container submenu, open an export screen, choose a destination directory, preview the exact workflow, and execute the export.

### Navigation Contract
- `ContainerSubmenuScreen` shows `Export container` only for stopped containers.
- Selecting export transitions to `ScreenContainerExport` with the chosen container attached.
- `ContainerExportScreen` uses a destination-directory input flow and returns to the submenu or list on completion/cancel.

### Command Sequence Contract

| Step | Command Shape | Required |
|------|---------------|----------|
| Export container state to temp image | `container export --image <generated-image-ref> <container-id>` | Yes |
| Save temp image to archive | `container image save --output <archive-path> <generated-image-ref>` | Yes |
| Remove temp image | `container image rm <generated-image-ref>` | Best effort |

### Preview Contract
- The UI must show the exact generated commands before execution, in order.
- Preview text must include the generated image reference and final archive path.
- Canceling the preview must leave the export request unexecuted.

### Success/Failure Contract
- If export or save fails, the workflow returns failure and surfaces stdout/stderr.
- If export and save succeed but cleanup fails, the workflow returns success with warning, preserving the archive path in the final result.
- The source container must remain unchanged on all failure paths.

---

## 3. Build `--pull` Contract

### User Action
Open the build form, review the default-enabled pull checkbox, optionally toggle it, and preview the final build command.

### Builder Contract
- `BuildImageBuilder` adds `PullLatest bool`.
- When `PullLatest == true`, command includes `--pull`.
- When `PullLatest == false`, command omits `--pull`.

### Command Shapes

| UI State | Command Shape |
|----------|---------------|
| Pull enabled | `container build --pull -t <tag> -f <file> <context-dir>` |
| Pull disabled | `container build -t <tag> -f <file> <context-dir>` |

### UI Contract
- Checkbox defaults to enabled whenever the build screen opens.
- Preview must reflect the exact state of the checkbox.
- Existing `enter=preview` behavior remains unchanged.

---

## 4. Daemon Status Contract

### User Action
Open or refresh daemon control and view a reliable `running`, `stopped`, or `unknown` state.

### Command Contract

| Operation | Command | Notes |
|-----------|---------|-------|
| Initial load | `container system status --format json` | Replaces text/table parsing |
| Refresh | `container system status --format json` | Same command reused |

### Output Contract
- Expected output is one JSON object.
- Known fields from local CLI sampling: `status`, `installRoot`, `appRoot`, `apiServerVersion`, `apiServerBuild`, `apiServerCommit`, `apiServerAppName`.
- `status` is the primary classification field.

### Parse Rules
- `status == "running"` -> `running`
- recognized non-running value -> `stopped`
- missing `status`, malformed JSON, or unrecognized values -> `unknown`

### UI Contract
- `DaemonControlScreen` must render a distinct unknown state rather than mapping everything non-running to stopped.
- Version metadata may be displayed or retained for future diagnostics, but state classification must not depend on free-form text search.

---

## 5. Automated Test Contract

The following command-facing behaviors require automated coverage:
- registry list builder emits `container registry list --format json`
- export workflow refuses running containers before command planning
- export workflow generates export/save/cleanup commands in the documented order
- build builder includes or omits `--pull` correctly
- daemon status builder emits `container system status --format json`
- daemon parser maps missing or malformed JSON to `unknown`