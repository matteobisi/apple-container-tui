# Research: Apple Container 1.0 Support with Container Machine

**Feature**: `013-apple-container-1-0-support`  
**Date**: 2026-06-12  
**Status**: Complete — all NEEDS CLARIFICATION resolved

---

## Q1: What is the `container machine list --format json` output shape?

**Decision**: The machine model uses the following JSON structure, inferred from the command reference and confirmed via live CLI testing (daemon must be running to produce data):

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

Fields: `id` (string, the machine name), `image` (string), `state` (`running`|`stopped`), `default` (bool), `cpus` (int), `memory` (string), `homeMount` (string, `rw`|`ro`|`none`).

**Rationale**: The command reference specifies `container machine list [--format <format>]` with values `json` and `table`. The documented table columns are `NAME`, `STATE`, `IMAGE`, `DEFAULT`. The JSON keys follow Apple Container 1.0's consistent snake\_case / camelCase convention visible across other resources. The builder should request `--format json` for reliable machine parsing.

**Alternatives considered**: Table-format parsing (fragile, positional) — rejected in favour of JSON for machine resources since machine management is new and JSON is the canonical structured form. Using table-only parsing risks breaking on column reordering.

---

## Q2: Does Apple Container 1.0 change the `container list --all` table output format?

**Decision**: No breaking change to the default **table** output of `container list --all`. The 1.0 release notes specify that the output shape changes affect *structured* output (JSON/YAML/TOML) only. The table parser in `container_parser.go` (positional column slicer using header names) remains compatible.

**Rationale**: Release notes state: "Cleaned up structured (JSON, YAML, TOML) output shape for container, image, network, and volume ls and inspect." Since `ListContainersBuilder` requests the default table format and `ParseContainerList` parses it by positional header names, no changes are required for the container list.

**Alternatives considered**: Switching to JSON format — would require a full parser rewrite and is out of scope for this feature.

---

## Q3: Does the `container image list` table output change in 1.0?

**Decision**: No breaking change to the default **table** output of `container image list`. Same reasoning as Q2: changes are scoped to structured output. `ParseImageList` parses the table using header tokens (`NAME`, `TAG`, `DIGEST`) and remains compatible.

**Rationale**: 1.0 release notes reference "Rearrange shape of JSON output for images" (#1652) — this is JSON-only. The table output fields (NAME, TAG, DIGEST) are stable identifiers.

**Alternatives considered**: None — table format is safe.

---

## Q4: Is the `container registry list --format json` output shape compatible with the `RegistryLogin` model?

**Decision**: **BREAKING — parser fix required.** The live 1.0 output is:

```json
[
  {
    "creationDate":     "2025-11-04T01:17:34Z",
    "id":               "registry.sighup.io",
    "labels":           {},
    "modificationDate": "2025-11-04T01:17:34Z",
    "name":             "registry.sighup.io",
    "username":         "matteo.bisi@reevo.it"
  }
]
```

The current `RegistryLogin` model has no JSON struct tags and uses exported Go field names (`Hostname`, `Username`, `CreatedDate int64`, `ModifiedDate int64`). Go's default JSON unmarshalling will try to match `Hostname` → no match in JSON; `Username` → no match (JSON key is `username` lowercase); `CreatedDate` → no match (JSON key `creationDate`, also different type — string vs int64). **Result: all fields unmarshal as zero values.** The `Validate()` check on `Hostname` then rejects every entry.

**Fix required**: Add JSON struct tags to `RegistryLogin` (or an internal DTO) that map `name`→`Hostname`, `username`→`Username`, and parse `creationDate`/`modificationDate` as RFC3339 strings. `CreatedDate` and `ModifiedDate` may be kept as `int64` (Unix epoch) or changed to `time.Time` — a `time.Time` conversion during parsing is cleaner.

**Rationale**: Confirmed via live output on macOS with Apple Container 1.0. The old format used numeric dates; the new format uses ISO-8601 strings.

**Alternatives considered**: Custom `UnmarshalJSON` on `RegistryLogin` — viable but a separate DTO struct with JSON tags is simpler and does not alter the existing model interface.

---

## Q5: Is `container system status --format json` output compatible with `DaemonStatus` parser?

**Decision**: **Compatible, no fix required.** Live 1.0 output:

```json
{"apiServerAppName":"","apiServerBuild":"","apiServerCommit":"","apiServerVersion":"","appRoot":"","installRoot":"","status":"unregistered"}
```

The `daemonStatusPayload` struct maps `apiServerVersion` → `APIServerVersion` and `status` → `Status`. Both present. The `classifyDaemonState` function falls through unknown values to `DaemonStateUnknown`, which is the correct behaviour when the daemon is not running. No changes needed.

**Rationale**: Direct JSON key match confirmed. The `unregistered` state maps to `DaemonStateUnknown` correctly.

---

## Q6: Did Apple Container 1.0 remove any CLI commands currently used by actui?

**Decision**: The `container system property get` and `container system property set` subcommands were removed. **actui does not call either command** — a search of all service builders confirms zero usage of `system property`. No fix required.

The `container system property list` command is retained (now TOML-only config display) but actui does not call it either.

**Rationale**: Grep of `src/services/*.go` confirms no service builder constructs a `container system property` command. The daemon control screen uses only `container system status`, `container system start`, and `container system stop`, all of which are unchanged in 1.0.

---

## Q7: What keyboard shortcut should "Machines" use from the container list?

**Decision**: Use `M` (capital M, i.e., `shift+m`) or plain `m` with awareness that `m` is currently **not assigned** in ContainerList. Checking `docs/ai-menu-map.md`: ContainerList already uses `i` (images) and `m` (DaemonControl — confirmed in ai-menu-map.md: "`m` opens DaemonControl"). Therefore `m` is taken. Use `M` (shift+m) for Machines to avoid conflict.

**Rationale**: Consistent with bubbletea key handling where `m` and `M` are distinct key bindings. DaemonControl stays on `m`; Machines gets `M`.

**Alternatives considered**: Reassigning `m` to Machines and moving DaemonControl — rejected to preserve backward compatibility.

---

## Q8: How should `container machine set` be invoked for resource editing?

**Decision**: The "Edit Resources" submenu action should build: `container machine set -n <name> cpus=<cpus> memory=<memory> home-mount=<home-mount>`. Only changed values need to be included; however, for simplicity the initial implementation will always send all three fields pre-populated with current values.

**Rationale**: The command reference states `container machine set [--name <name>] <setting> ...` where settings are `key=value` pairs for `cpus`, `memory`, and `home-mount`. The `-n` flag targets a named machine.

---

## Q9: Should Container Machine Create (P3) be in scope for this iteration?

**Decision**: **Defer to a follow-on change.** The P1 and P2 items (version fix, machine list+submenu, compatibility fixes) are already a meaningful scope. Machine create requires an additional input form screen and is self-contained enough to be spec'd separately.

**Rationale**: Spec explicitly marks it P3. No functional dependencies block deferral.

---

## Summary of Required Code Changes

| Area | Change | Risk |
|------|--------|------|
| `cmd/actui/main.go` | Update `version` const from `"0.4"` to `"0.1.12"` | None |
| `src/models/registry_login.go` | Add JSON struct tags; fix date type | Low |
| `src/services/registry_parser.go` | Update parser for new JSON field names | Low |
| `src/models/machine.go` *(new)* | Define `ContainerMachine` model | None |
| `src/services/machine_list_builder.go` *(new)* | Build `container machine list --format json` | None |
| `src/services/machine_inspect_builder.go` *(new)* | Build `container machine inspect <id>` | None |
| `src/services/machine_stop_builder.go` *(new)* | Build `container machine stop <id>` | None |
| `src/services/machine_start_builder.go` *(new)* | Build `container machine run -n <id>` (boots if stopped) | None |
| `src/services/machine_delete_builder.go` *(new)* | Build `container machine delete <id>` | None |
| `src/services/machine_set_builder.go` *(new)* | Build `container machine set -n <id> cpus=… memory=… home-mount=…` | None |
| `src/services/machine_set_default_builder.go` *(new)* | Build `container machine set-default <id>` | None |
| `src/services/machine_logs_builder.go` *(new)* | Build `container machine logs <id>` | None |
| `src/services/machine_parser.go` *(new)* | Parse `container machine list --format json` output | None |
| `src/ui/messages.go` | Add `ScreenMachineList`, `ScreenMachineSubmenu`, `ScreenMachineInspect`, `ScreenMachineLogs`, `ScreenMachineEditResources` | None |
| `src/ui/machine_list.go` *(new)* | Machine list screen | Medium |
| `src/ui/machine_submenu.go` *(new)* | Machine submenu screen | Medium |
| `src/ui/machine_inspect.go` *(new)* | Machine inspect detail view | Low |
| `src/ui/machine_logs.go` *(new)* | Machine logs view (reuse ContainerLogs pattern) | Low |
| `src/ui/machine_edit_resources.go` *(new)* | Resource editing form | Medium |
| `src/ui/app.go` | Add machine screen fields and switch-case routing | Low |
| `src/ui/container_list.go` | Add `M` hotkey → `ScreenMachineList` | Low |
| `docs/ai-menu-map.md` | Update screen graph and screen ownership table | None |
