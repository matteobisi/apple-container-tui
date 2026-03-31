# Quickstart: Expanded Container Workflows

**Feature**: 005-expand-container-workflows  
**Date**: 2026-03-31  
**For**: Developers implementing and validating the feature without broad architectural changes

## Prerequisites

- macOS 26.x on Apple Silicon
- Apple Container CLI installed and available as `container`
- At least one logged-in registry for Registries-screen validation
- At least one stopped container for export validation
- A build file (`Containerfile` or `Dockerfile`) available for build-form validation
- Go 1.21+ and repo dependencies installed

## Implementation Order

### 1. Add Registry Listing End-to-End

Files likely touched:
- `src/models/registry_login.go` (new)
- `src/services/registry_list_builder.go` (new)
- `src/services/registry_parser.go` (new)
- `src/ui/messages.go`
- `src/ui/app.go`
- `src/ui/image_list.go`
- `src/ui/registries.go` (new)
- `src/ui/help.go`

Validation target:
- New Registries screen loads data from `container registry list --format json`
- Empty, error, and loaded states are all visible and navigable

### 2. Add Stopped-Container Export Flow

Files likely touched:
- `src/services/export_container_builder.go` (new)
- `src/services/image_save_builder.go` (new)
- `src/services/export_workflow.go` (new)
- `src/ui/container_submenu.go`
- `src/ui/container_export.go` (new)
- `src/ui/messages.go`
- `src/ui/app.go`

Validation target:
- Running containers do not show Export
- Stopped containers open the export screen
- Preview shows all generated commands before execution
- Successful export produces an OCI tar in the chosen directory

### 3. Extend Build Form with Pull Checkbox

Files likely touched:
- `src/services/build_image_builder.go`
- `src/ui/build.go`
- `src/ui/help.go`

Validation target:
- Checkbox is visible and defaults to enabled
- Preview contains `--pull` when enabled and omits it when disabled

### 4. Move Daemon Status to Structured Parsing

Files likely touched:
- `src/services/check_daemon_builder.go`
- `src/services/daemon_parser.go`
- `src/models/daemon.go`
- `src/ui/daemon_control.go`

Validation target:
- `running`, `stopped`, and `unknown` all render correctly
- Missing/invalid JSON no longer produces misleading `running` or `stopped`

## Suggested Verification Flow

### Automated

Run the full test suite after each vertical slice and once again at the end:

```bash
go test ./...
```

Focus additions:
- builder tests in `src/services/services_test.go`
- parser tests for registry list and daemon status JSON
- UI routing/state tests in `src/ui/ui_additional_test.go`

### Manual

Launch the app:

```bash
go run cmd/actui/main.go
```

#### Registries screen
1. Open the Images view.
2. Trigger the Registries screen from the new image-domain shortcut.
3. Confirm loaded rows show hostname and username.
4. Confirm `r` refreshes and `esc` returns cleanly.

#### Container export
1. Open a stopped container submenu.
2. Confirm Export appears.
3. Open export flow, choose a destination directory, and inspect the preview.
4. Execute the workflow and confirm a tar archive appears in the destination directory.
5. Repeat with a running container and confirm Export is not offered.

#### Build pull checkbox
1. Open the build flow from the image screen.
2. Confirm the pull checkbox starts enabled.
3. Preview once with it enabled and once disabled.
4. Verify only the enabled preview contains `--pull`.

#### Daemon status
1. Open daemon control.
2. Confirm normal status loads through the JSON path.
3. Exercise refresh.
4. Validate unknown-state rendering by injecting malformed or incomplete output in tests if a live manual path is inconvenient.

## Acceptance Checklist

- [ ] Registries screen uses runtime-managed registry data only
- [ ] Export is available only for stopped containers
- [ ] Export preview shows exact command sequence and generated output path
- [ ] Build checkbox defaults to enabled and controls `--pull`
- [ ] Daemon status uses structured JSON and supports `unknown`
- [ ] `go test ./...` passes
- [ ] Manual verification completed on macOS 26.x with Apple Container CLI