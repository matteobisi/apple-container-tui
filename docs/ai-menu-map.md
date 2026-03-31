# AI Menu Map

## Purpose

This document is for AI agents and contributors, not end users.

Use it to answer three questions before editing the TUI:

1. Where does this feature belong in the current menu structure?
2. Which files own that workflow?
3. Which routing and service files must change together?

The canonical sources of truth remain [src/ui/messages.go](src/ui/messages.go) for screen ids and [src/ui/app.go](src/ui/app.go) for navigation.

## Screen Graph

```text
ContainerList
|- ContainerSubmenu
|  |- ContainerLogs
|  |- ContainerShell
|  `- ContainerExport
|- ImageList
|  |- ImageSubmenu
|  |  `- ImageInspect
|  |- ImagePull
|  |- Registries
|  `- FilePicker -> Build
|- DaemonControl
`- Help
```

## Routing Backbone

- Screen ids are declared in [src/ui/messages.go](src/ui/messages.go)
- Screen instances are held and updated in [src/ui/app.go](src/ui/app.go)
- `screenChangeMsg` is the primary navigation message
- `BackToListMsg` returns from nested flows to their parent list view
- `BackToSubmenuMsg` returns from nested container or image detail flows to the active submenu
- New screens are not complete until they are added to both the message definitions and the `AppModel` switch logic

## Screen Ownership

| Screen | Purpose | Entry Trigger | Main UI File | Main Service Files |
|--------|---------|---------------|--------------|--------------------|
| ContainerList | Root container browser and quick actions | App start | [src/ui/container_list.go](src/ui/container_list.go) | [src/services/list_containers_builder.go](src/services/list_containers_builder.go), [src/services/container_parser.go](src/services/container_parser.go) |
| ContainerSubmenu | Contextual actions for the selected container | `enter` on container row | [src/ui/container_submenu.go](src/ui/container_submenu.go) | start, stop, logs, shell, export related services under [src/services](src/services) |
| ContainerLogs | Tail logs for selected container | `logs` action from container submenu | [src/ui/container_logs.go](src/ui/container_logs.go) | [src/services/container_logs_builder.go](src/services/container_logs_builder.go) |
| ContainerShell | Open shell into a running container | `shell` action from container submenu | [src/ui/container_shell.go](src/ui/container_shell.go) | [src/services/container_exec_builder.go](src/services/container_exec_builder.go) |
| ContainerExport | Export stopped container to OCI archive | `export` action from container submenu | [src/ui/container_export.go](src/ui/container_export.go) | [src/services/export_workflow.go](src/services/export_workflow.go), [src/services/export_container_builder.go](src/services/export_container_builder.go), [src/services/image_save_builder.go](src/services/image_save_builder.go), [src/services/image_delete_builder.go](src/services/image_delete_builder.go) |
| ImageList | Root image browser and image-level actions | `i` from container list | [src/ui/image_list.go](src/ui/image_list.go) | [src/services/image_list_builder.go](src/services/image_list_builder.go) |
| ImageSubmenu | Contextual actions for selected image | `enter` on image row | [src/ui/image_submenu.go](src/ui/image_submenu.go) | image inspect and delete services under [src/services](src/services) |
| ImageInspect | Render detailed image metadata | `inspect` action from image submenu | [src/ui/image_inspect.go](src/ui/image_inspect.go) | [src/services/image_inspect_builder.go](src/services/image_inspect_builder.go) |
| ImagePull | Pull image workflow with preview | `p` from image list | [src/ui/image_pull.go](src/ui/image_pull.go) | [src/services/pull_image_builder.go](src/services/pull_image_builder.go) |
| Registries | View runtime-managed registry logins | `g` from image list | [src/ui/registries.go](src/ui/registries.go) | [src/services/registry_list_builder.go](src/services/registry_list_builder.go), [src/services/registry_parser.go](src/services/registry_parser.go) |
| FilePicker | Select build source file | `b` from image list | [src/ui/file_picker.go](src/ui/file_picker.go) | build source detection services under [src/services](src/services) |
| Build | Build image from selected file | file chosen in file picker | [src/ui/build.go](src/ui/build.go) | [src/services/build_image_builder.go](src/services/build_image_builder.go), [src/services/build_file_detector.go](src/services/build_file_detector.go) |
| DaemonControl | Daemon status and start/stop actions | `m` from container list | [src/ui/daemon_control.go](src/ui/daemon_control.go) | [src/services/check_daemon_builder.go](src/services/check_daemon_builder.go), [src/services/daemon_parser.go](src/services/daemon_parser.go), [src/services/start_daemon_builder.go](src/services/start_daemon_builder.go), [src/services/stop_daemon_builder.go](src/services/stop_daemon_builder.go) |
| Help | Global help and version display | `?` from most screens | [src/ui/help.go](src/ui/help.go) | none |

## Primary Entry Points

### Container Root

- File: [src/ui/container_list.go](src/ui/container_list.go)
- Normal actions:
  - `enter` opens [src/ui/container_submenu.go](src/ui/container_submenu.go)
  - `s` previews start
  - `t` previews stop
  - `d` opens delete confirmation
  - `i` opens [src/ui/image_list.go](src/ui/image_list.go)
  - `m` opens [src/ui/daemon_control.go](src/ui/daemon_control.go)
  - `?` opens [src/ui/help.go](src/ui/help.go)

### Container Action Branch

- File: [src/ui/container_submenu.go](src/ui/container_submenu.go)
- Option set depends on container state:
  - Running container: stop, logs, shell
  - Stopped container: start, logs, export
- If a new container-specific workflow is added, this is the default insertion point

### Image Branch

- File: [src/ui/image_list.go](src/ui/image_list.go)
- Primary actions:
  - `enter` opens [src/ui/image_submenu.go](src/ui/image_submenu.go)
  - `p` opens [src/ui/image_pull.go](src/ui/image_pull.go)
  - `g` opens [src/ui/registries.go](src/ui/registries.go)
  - `b` opens [src/ui/file_picker.go](src/ui/file_picker.go)
  - `n` opens prune confirmation
- If a new image-level workflow is added, this is the default insertion point

### Build Branch

- File selection starts in [src/ui/file_picker.go](src/ui/file_picker.go)
- Selected file routes into [src/ui/build.go](src/ui/build.go)
- When a new build option is needed, it usually belongs in [src/ui/build.go](src/ui/build.go) plus [src/services/build_image_builder.go](src/services/build_image_builder.go)

## Change Recipes

### Add a New Screen

1. Add a new screen constant in [src/ui/messages.go](src/ui/messages.go)
2. Add the screen state field to `AppModel` in [src/ui/app.go](src/ui/app.go)
3. Initialize it in `NewAppModel` in [src/ui/app.go](src/ui/app.go)
4. Register its init/update/view branches in [src/ui/app.go](src/ui/app.go)
5. Add the screen file under [src/ui](src/ui)
6. Add or extend tests in [src/ui/ui_flow_test.go](src/ui/ui_flow_test.go) or [src/ui/ui_additional_test.go](src/ui/ui_additional_test.go)
7. Update this file

### Add a New Container Action

1. Start in [src/ui/container_submenu.go](src/ui/container_submenu.go)
2. Decide whether the action is a direct preview/command or a dedicated screen
3. If dedicated screen, follow the new screen recipe
4. Add command builders/parsers in [src/services](src/services)
5. If the action depends on running vs stopped state, update `buildOptions()` in [src/ui/container_submenu.go](src/ui/container_submenu.go)

### Add a New Image Workflow

1. Start in [src/ui/image_list.go](src/ui/image_list.go)
2. Decide whether it is a direct action, submenu action, or multi-step screen
3. If it requires selecting a file, route through [src/ui/file_picker.go](src/ui/file_picker.go)
4. Add builders/parsers in [src/services](src/services)
5. Update tests and this map

### Add a New Build Option

1. Update form state and key handling in [src/ui/build.go](src/ui/build.go)
2. Update preview generation in [src/ui/build.go](src/ui/build.go)
3. Update command building in [src/services/build_image_builder.go](src/services/build_image_builder.go)
4. Add tests for both the builder and the UI toggle/preview path

### Add a New Daemon Capability

1. Start in [src/ui/daemon_control.go](src/ui/daemon_control.go)
2. Add or update command builders in [src/services](src/services)
3. If status shape changes, update [src/services/daemon_parser.go](src/services/daemon_parser.go)
4. Keep unknown-state fallback behavior intact unless intentionally changing it

## Guardrails For Future Agents

- Prefer matching existing flow shapes over inventing a new navigation pattern
- Keep previews and confirmations consistent with the rest of the TUI
- Additive changes are safer than broad route reshuffles in this codebase
- If a feature can live inside an existing screen, prefer that over adding a new top-level branch
- When behavior changes, update both tests and this document so the map stays trustworthy

## Maintenance Rule

Update this file whenever one of these changes happens:

- a new screen is added or removed
- an entry key or source menu changes for an existing workflow
- ownership of a workflow moves to a different file
- a service dependency changes enough that a future editor would start in the wrong place