# Apple Container TUI User Guide

## Overview

Apple Container TUI is a keyboard-first terminal UI for listing containers,
starting or stopping them, exporting stopped containers, pulling images,
building from Containerfiles, browsing runtime registries, and managing the
Apple Container daemon. Every action shows a command preview before execution,
with optional dry-run mode for safe practice.

## Prerequisites

- macOS 26.x on Apple Silicon
- Apple Container CLI installed and in PATH
- Go 1.21+ if building from source

Verify the CLI is available:

```bash
container system version
```

## Installation

Build from source:

```bash
go mod download
go build -o actui ./cmd/actui
```

## Run

```bash
./actui
```

Dry-run mode (preview only, no execution):

```bash
./actui --dry-run
```

## Navigation

Main keys:

- `up/down` or `j/k` to move
- `enter` to open submenu on selected item
- `q` to quit
- `?` for help

## Workflow: List and Start/Stop Containers

1. Launch `./actui`
2. Use arrow keys to select a container
3. Press `enter` to open the container submenu
4. Choose `Start container`, `Stop container`, `Tail container log`, `Enter container`, or `Export container` when the container is stopped
5. Confirm command previews where applicable

ASCII screenshot:

```
+------------------------------------------------------+
| Containers                                           |
|                                                      |
| > web-api [stopped] nginx:latest                     |
|   worker  [running] alpine:latest                    |
|                                                      |
| Keys: up/down, enter=submenu, s=start, t=stop, ...   |
|                                                      |
| Command Preview                                      |
|                                                      |
| container start web-api                              |
|                                                      |
| Confirm (y/n)                                        |
+------------------------------------------------------+
```

## Workflow: Image Management (List, Pull, Build, Prune)

1. Press `i` from the main screen to open image list
2. Press `p` to pull, `g` to view runtime registries, `b` to build, or `n` to prune unused images
3. For prune, type `prune` to confirm
4. Image list refreshes automatically after operations

ASCII screenshot:

```
+------------------------------------------------------+
| Images                                               |
|                                                      |
| > ubuntu                     latest   sha256:abc...  |
|   alpine                     latest   sha256:def...  |
|                                                      |
| Keys: up/down, enter=submenu, p=pull, g=registries, b=build, n=prune |
|                                                      |
| Command Preview                                      |
|                                                      |
| container image list                                 |
|                                                      |
| Confirm (y/n)                                        |
+------------------------------------------------------+
```

## Workflow: Image Inspect/Delete

1. Open image list with `i`
2. Select image and press `enter`
3. Choose `Inspect image` to view metadata or `Delete image` to remove it
4. For delete, type `delete` to confirm

## Workflow: Browse Registries

1. Press `i` from the main screen to open the image list
2. Press `g` to open the Registries screen
3. Review runtime-managed registry entries by hostname and username
4. Press `r` to refresh or `esc` to return to the image list

## Workflow: Export a Stopped Container

1. Select a stopped container and press `enter`
2. Choose `Export container`
3. Enter a destination directory
4. Review the export and save preview
5. Confirm to create an OCI tar archive in the selected directory
6. Choose whether to delete the temporary exported image after the archive is written

## Workflow: Build from Containerfile/Dockerfile

1. Press `i` from the main screen, then `b` in image list
2. Select a build file in the file picker
3. Enter a tag
4. Leave `Pull latest base images` enabled or toggle it with `p`
5. Confirm the preview to build

ASCII screenshot:

```
+------------------------------------------------------+
| Build Image                                          |
|                                                      |
| File: ./Containerfile                                |
| Tag: my-app:latest                                   |
| [x] Pull latest base images                          |
|                                                      |
| Progress: [========                     ]            |
|                                                      |
| Command Preview                                      |
|                                                      |
| container build --pull -t my-app:latest -f ./Containerfile . |
|                                                      |
| Confirm (y/n)                                        |
+------------------------------------------------------+
```

## Workflow: Delete a Stopped Container

1. Select a stopped container
2. Press `d`
3. Type the exact name/ID to confirm

ASCII screenshot:

```
+------------------------------------------------------+
| Delete Container                                     |
|                                                      |
| Command: container delete web-api                    |
| Confirm by typing: web-api                           |
|                                                      |
| Type to confirm: web-api                             |
|                                                      |
| Press enter to confirm, esc to cancel                |
+------------------------------------------------------+
```

## Workflow: Start/Stop the Daemon

1. Press `m` for daemon management
2. Review the structured daemon state as `running`, `stopped`, or `unknown`
3. Press `s` to start or `t` to stop
4. Confirm in the yes/no prompt

ASCII screenshot:

```
+------------------------------------------------------+
| Daemon Control                                       |
|                                                      |
| Status: stopped                                      |
|                                                      |
| Actions:                                             |
|   s - start daemon                                   |
|   t - stop daemon (!)                                |
|                                                      |
| Confirm (y/n)                                        |
+------------------------------------------------------+
```

## Configuration

Config is read from:

- `~/.config/actui/config`
- `~/Library/Application Support/actui/config`

Writes go to:

- `~/Library/Application Support/actui/config`

Example TOML:

```toml
default_build_file = "Containerfile"
confirm_destructive_actions = true
theme_mode = "auto"
refresh_on_focus = false
log_retention_days = 7
```

Logs are stored at:

- `~/Library/Application Support/actui/command.log`

## Troubleshooting

- "apple container CLI not found" -> install from https://github.com/apple/container
- "daemon is not running" -> open Daemon screen and start it
- "daemon status unknown" -> refresh the daemon screen; if it persists, inspect `container system status --format json` directly
- Build errors -> ensure a Containerfile or Dockerfile exists in the chosen folder
- Export errors -> confirm the destination directory exists and is writable
- Permission errors -> confirm your user can run the CLI commands without sudo

## Help Screen

Press `?` from any screen to review shortcuts and paths.
