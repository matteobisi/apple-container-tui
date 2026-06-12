# Quickstart Validation Guide: Apple Container 1.0 Support with Container Machine

**Feature**: `013-apple-container-1-0-support`  
**Date**: 2026-06-12

This guide describes how to validate the feature end-to-end after implementation. It is not a test suite — it is a runbook for manual verification on a macOS machine with Apple Container 1.0 installed.

---

## Prerequisites

- macOS 26.x, Apple Silicon
- Apple Container 1.0 installed (`container system version` reports `1.0.0`)
- `actui` binary built from this branch (`go build ./cmd/actui`)
- Container system started (`container system start`)
- At least one container machine created for testing: `container machine create alpine:latest --name dev`

---

## Validation 1: Version String

**Goal**: Confirm `actui --version` reports the correct version.

```bash
./actui --version
```

**Expected output**:
```
actui version 0.1.12
```

---

## Validation 2: Registry List Compatibility

**Goal**: Confirm the Registries screen renders correctly against Apple Container 1.0 JSON output.

1. Start `./actui`
2. Press `i` to open the Image List screen
3. Press `g` to open the Registries screen
4. Verify registry entries are listed with a non-empty hostname column

**Expected**: Each configured registry appears as a row with its hostname shown (e.g. `registry.sighup.io`). No "unknown" or empty rows.

**Failure indicator**: All rows empty or screen shows "no registries" when `container registry list --format json` returns entries.

---

## Validation 3: Container List and Image List Compatibility

**Goal**: Confirm existing screens parse Apple Container 1.0 output correctly.

1. Start the container system if not already running: `container system start`
2. Ensure at least one container exists: `container run -d --name test-box alpine:latest sleep 300`
3. Start `./actui`
4. Verify the container list shows `test-box` with its image and status columns populated

Then:
5. Press `i` to open Image List
6. Verify at least one image appears with Name and Tag columns populated

**Expected**: Both screens display data without empty columns or parse errors.

---

## Validation 4: Machine List Screen

**Goal**: Confirm the new Machines screen opens and lists machines.

1. Start `./actui`
2. From the container list, press `M` (shift+m)
3. Verify the Machines screen opens with a table showing: Name, State, Image, Default columns
4. Verify the `dev` machine appears with state `running` and `DEFAULT` marked

**Expected**: Machine list table is populated. Empty-state message appears when no machines exist.

---

## Validation 5: Machine Submenu — Read-Only Actions

**Goal**: Confirm Inspect and Logs actions work.

1. From the Machines screen, press `enter` on the `dev` machine
2. Verify the Machine Submenu appears with actions: Inspect, Logs, Stop, Edit Resources, Set as Default, Delete
3. Select **Inspect** — verify machine configuration is shown (CPUs, memory, home-mount, image, state)
4. Press `esc` to return to submenu
5. Select **Logs** — verify log output is displayed (may be empty for a fresh machine)
6. Press `esc` to return

**Expected**: Both detail views open and navigate back correctly.

---

## Validation 6: Machine Submenu — Stop

**Goal**: Confirm Stop shows a command preview before executing.

1. From the Machines screen, select the running `dev` machine
2. In the submenu, select **Stop**
3. Verify a command preview appears: `container machine stop dev`
4. Confirm — verify machine stops (state changes to `stopped` in the list)

**Expected**: Command preview shown; machine stops; list refreshes with updated state.

---

## Validation 7: Machine Submenu — Start

**Goal**: Confirm Start (run) works on a stopped machine.

1. With `dev` in stopped state, press `enter` on it in the machine list
2. Verify submenu shows **Start** (not Stop)
3. Select **Start** — verify command preview: `container machine run -n dev`
4. Confirm — verify machine boots (state returns to `running`)

**Expected**: Command preview shown; machine starts; list refreshes.

---

## Validation 8: Edit Resources

**Goal**: Confirm resource editing previews the correct `container machine set` command.

1. Stop the `dev` machine if running
2. From the submenu, select **Edit Resources**
3. Verify a form opens pre-filled with current cpus, memory, and home-mount values
4. Change CPUs to `2` and Memory to `4G`
5. Confirm — verify the command preview shows: `container machine set -n dev cpus=2 memory=4G home-mount=rw`
6. Verify an advisory note: "Changes take effect after next stop and restart"

**Expected**: Form pre-filled, command preview correct, advisory shown.

---

## Validation 9: Machine Delete (Destructive Action)

**Goal**: Confirm delete requires type-to-confirm.

1. Create a disposable machine: `container machine create alpine:latest --name to-delete`
2. Open actui, navigate to Machines, select `to-delete`
3. In the submenu, select **Delete**
4. Verify a type-to-confirm prompt appears requesting the machine name
5. Type `to-delete` and confirm
6. Verify machine is removed from the list

**Expected**: Type-to-confirm gate; machine deleted; list updated.

---

## Validation 10: Navigation and Back Flows

**Goal**: Confirm `esc`/`q` navigation returns to expected parent screens.

| From screen | Action | Expected destination |
|-------------|--------|---------------------|
| MachineList | press `esc` | ContainerList |
| MachineSubmenu | press `esc` | MachineList |
| MachineInspect | press `esc` | MachineSubmenu |
| MachineLogs | press `esc` | MachineSubmenu |
| MachineEditResources | press `esc` | MachineSubmenu |

---

## Dry-Run Smoke Test

```bash
./actui --dry-run
```

Navigate through all new machine screens. Verify no CLI commands are actually executed (all previews show `[dry-run]` indicator or equivalent). Existing dry-run behaviour must be preserved for all machine operations.
