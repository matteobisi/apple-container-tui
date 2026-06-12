# Command Contracts: Machine Management

**Feature**: `013-apple-container-1-0-support`  
**Date**: 2026-06-12

This document maps each actui machine management action to its Apple Container CLI command. All commands use `container machine` (no alias). Commands that modify state are previewed before execution; destructive commands require confirmation.

---

## Machine List

| Property | Value |
|----------|-------|
| **actui trigger** | `M` from ContainerList screen |
| **CLI command** | `container machine list --format json` |
| **Output format** | JSON array of machine objects |
| **Builder** | `MachineListBuilder` |
| **Parser** | `ParseMachineList` in `machine_parser.go` |
| **Error handling** | Show empty-state message when result is empty; surface CLI errors in result display |

### Example output

```json
[
  {
    "id":        "dev",
    "image":     "alpine:latest",
    "state":     "running",
    "default":   true,
    "cpus":      4,
    "memory":    "8G",
    "homeMount": "rw"
  }
]
```

---

## Machine Inspect

| Property | Value |
|----------|-------|
| **actui trigger** | "Inspect" from MachineSubmenu |
| **CLI command** | `container machine inspect <id>` |
| **Output format** | JSON (raw display, no parsing beyond pretty-print) |
| **Builder** | `MachineInspectBuilder` |
| **Preview shown** | No (read-only) |
| **Error handling** | Surface CLI errors in result display |

### Arguments

| Argument | Source |
|----------|--------|
| `<id>` | `ContainerMachine.ID` of selected machine |

---

## Machine Logs

| Property | Value |
|----------|-------|
| **actui trigger** | "Logs" from MachineSubmenu |
| **CLI command** | `container machine logs <id>` |
| **Output format** | Plain text (streamed) |
| **Builder** | `MachineLogsBuilder` |
| **Preview shown** | No (read-only) |
| **UI pattern** | Reuse ContainerLogs screen pattern |

### Arguments

| Argument | Source |
|----------|--------|
| `<id>` | `ContainerMachine.ID` of selected machine |

---

## Machine Stop

| Property | Value |
|----------|-------|
| **actui trigger** | "Stop" from MachineSubmenu (running machines only) |
| **CLI command** | `container machine stop <id>` |
| **Output format** | None (exit code) |
| **Builder** | `MachineStopBuilder` |
| **Preview shown** | Yes — command preview with yes/no confirmation |
| **Dry-run safe** | Yes |
| **Error handling** | Surface stderr; refresh machine list on success |

### Arguments

| Argument | Source |
|----------|--------|
| `<id>` | `ContainerMachine.ID` of selected machine |

---

## Machine Start (Run)

| Property | Value |
|----------|-------|
| **actui trigger** | "Start" from MachineSubmenu (stopped machines only) |
| **CLI command** | `container machine run -n <id>` |
| **Output format** | None (boots machine and returns) |
| **Builder** | `MachineStartBuilder` |
| **Preview shown** | Yes — command preview with yes/no confirmation |
| **Dry-run safe** | Yes |
| **Note** | `container machine run` boots the machine if stopped; no separate `machine start` command exists |

### Arguments

| Argument | Source |
|----------|--------|
| `-n <id>` | `ContainerMachine.ID` of selected machine |

---

## Machine Edit Resources

| Property | Value |
|----------|-------|
| **actui trigger** | "Edit Resources" from MachineSubmenu |
| **CLI command** | `container machine set -n <id> cpus=<N> memory=<M> home-mount=<mode>` |
| **Output format** | None (exit code) |
| **Builder** | `MachineSetBuilder` |
| **Preview shown** | Yes — command preview before submission |
| **Dry-run safe** | Yes |
| **Note** | Changes take effect after next stop+start; UI shows advisory message |

### Arguments

| Argument | Source |
|----------|--------|
| `-n <id>` | `ContainerMachine.ID` of selected machine |
| `cpus=<N>` | User input, pre-filled from `ContainerMachine.CPUs` |
| `memory=<M>` | User input, pre-filled from `ContainerMachine.Memory` |
| `home-mount=<mode>` | User input, pre-filled from `ContainerMachine.HomeMount`; one of `rw`, `ro`, `none` |

---

## Machine Set Default

| Property | Value |
|----------|-------|
| **actui trigger** | "Set as Default" from MachineSubmenu |
| **CLI command** | `container machine set-default <id>` |
| **Output format** | None (exit code) |
| **Builder** | `MachineSetDefaultBuilder` |
| **Preview shown** | Yes — command preview with yes/no confirmation |
| **Dry-run safe** | Yes |

### Arguments

| Argument | Source |
|----------|--------|
| `<id>` | `ContainerMachine.ID` of selected machine |

---

## Machine Delete

| Property | Value |
|----------|-------|
| **actui trigger** | "Delete" from MachineSubmenu |
| **CLI command** | `container machine delete <id>` |
| **Output format** | None (exit code) |
| **Builder** | `MachineDeleteBuilder` |
| **Preview shown** | Yes — type-to-confirm destructive action (machine name must be typed) |
| **Dry-run safe** | Yes |
| **Error handling** | Surface stderr; return to machine list on success |

### Arguments

| Argument | Source |
|----------|--------|
| `<id>` | `ContainerMachine.ID` of selected machine |

---

## Registry List (Compatibility Fix)

| Property | Value |
|----------|-------|
| **actui trigger** | `g` from ImageList screen (existing) |
| **CLI command** | `container registry list --format json` (unchanged) |
| **Change** | `RegistryLogin` struct needs JSON tags to match 1.0 output shape |
| **Builder** | `RegistryListBuilder` — no change needed |
| **Parser** | `ParseRegistryList` — no logic change, struct tag fix propagates |

### 1.0 JSON field mapping

| JSON field | `RegistryLogin` field | Action |
|------------|----------------------|--------|
| `name` | `Hostname` | Add `json:"name"` tag |
| `username` | `Username` | Add `json:"username"` tag |
| `creationDate` | `CreatedDate` | Change type to `time.Time`, add `json:"creationDate"` tag |
| `modificationDate` | `ModifiedDate` | Change type to `time.Time`, add `json:"modificationDate"` tag |
