# Data Model: Apple Container 1.0 Support with Container Machine

**Feature**: `013-apple-container-1-0-support`  
**Date**: 2026-06-12

---

## New Entity: ContainerMachine

### Definition

Represents a persistent Linux VM managed by `container machine`. Corresponds to one entry in `container machine list --format json`.

### Fields

| Field | Type | Source JSON key | Description |
|-------|------|-----------------|-------------|
| `ID` | `string` | `id` | Machine name/identifier |
| `Image` | `string` | `image` | OCI image reference used to create the machine |
| `State` | `MachineState` | `state` | Current runtime state |
| `IsDefault` | `bool` | `default` | Whether this is the default machine |
| `CPUs` | `int` | `cpus` | Number of virtual CPUs |
| `Memory` | `string` | `memory` | Memory allocation (e.g. `"8G"`) |
| `HomeMount` | `string` | `homeMount` | Home directory mount mode: `rw`, `ro`, or `none` |

### MachineState

```
type MachineState string

const (
    MachineStateRunning MachineState = "running"
    MachineStateStopped MachineState = "stopped"
    MachineStateUnknown MachineState = "unknown"
)
```

**State transitions**:
- `stopped` → start (`container machine run`) → `running`
- `running` → stop (`container machine stop`) → `stopped`
- Any state + delete (`container machine delete`) → removed

### Validation Rules

- `ID` must be non-empty (required for all targeted commands).
- `State` is parsed from the `state` JSON field; unrecognised values map to `MachineStateUnknown`.
- `HomeMount` is one of `rw`, `ro`, `none`; defaults to `rw` if absent.
- `CPUs` must be ≥ 1; `Memory` is a string (e.g. `"2G"`) — not validated further at model level.

---

## Modified Entity: RegistryLogin

### Current State (broken in 1.0)

```go
type RegistryLogin struct {
    Hostname     string
    Username     string
    CreatedDate  int64
    ModifiedDate int64
}
```

No JSON struct tags → Go uses exported field names as JSON keys. The 1.0 `container registry list --format json` output uses lowercase camelCase keys with ISO-8601 dates.

### Required Fix

Add JSON struct tags and change date fields to `time.Time`:

| Field | JSON key (1.0 output) | New struct tag |
|-------|-----------------------|----------------|
| `Hostname` | `name` | `json:"name"` |
| `Username` | `username` | `json:"username"` |
| `CreatedDate` | `creationDate` (RFC3339 string) | Change type to `time.Time`, tag `json:"creationDate"` |
| `ModifiedDate` | `modificationDate` (RFC3339 string) | Change type to `time.Time`, tag `json:"modificationDate"` |

**Note**: The `id` field in the JSON equals `name` for registry entries. `name` is the authoritative hostname.

### Impact

- `registry_parser.go`: `ParseRegistryList` continues to unmarshal `[]RegistryLogin` — struct tag fix is sufficient; no parser logic change needed.
- `models_test.go` / `services_test.go`: update any test fixtures using old numeric date fields.

---

## New Actions: MachineAction

Defines what operations can be performed on a machine from the submenu, based on state:

| Action | Applicable State | CLI Command |
|--------|-----------------|-------------|
| Inspect | any | `container machine inspect <id>` |
| Logs | any | `container machine logs <id>` |
| Start / Open Shell | stopped | `container machine run -n <id>` |
| Stop | running | `container machine stop <id>` |
| Edit Resources | stopped (effective after restart) | `container machine set -n <id> cpus=… memory=… home-mount=…` |
| Set as Default | any | `container machine set-default <id>` |
| Delete | any (with confirmation) | `container machine delete <id>` |

---

## MachineEditInput

Represents the input values collected from the resource-editing form:

| Field | Type | Constraints |
|-------|------|-------------|
| `Name` | `string` | Set from selected machine; read-only in form |
| `CPUs` | `string` | Numeric string; parsed as int; ≥ 1 |
| `Memory` | `string` | e.g. `"2G"`, `"512M"`; passed through to CLI verbatim |
| `HomeMount` | `string` | One of `rw`, `ro`, `none`; cycle-select or text input |
