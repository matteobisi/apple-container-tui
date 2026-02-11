# Apple Container TUI User Guide

## Overview

Apple Container TUI is a keyboard-first terminal UI for listing containers,
starting or stopping them, pulling images, building from Containerfiles,
and managing the Apple Container daemon. Every action shows a command preview
before execution, with optional dry-run mode for safe practice.

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
go build -o apple-tui ./cmd/apple-tui
```

## Run

```bash
./apple-tui
```

Dry-run mode (preview only, no execution):

```bash
./apple-tui --dry-run
```

## Navigation

Main keys:

- `up/down` or `j/k` to move
- `enter` to toggle start/stop on containers
- `q` to quit
- `?` for help

## Workflow: List and Start/Stop Containers

1. Launch `./apple-tui`
2. Use arrow keys to select a container
3. Press `s` to start, `t` to stop, or `enter` to toggle
4. Confirm the command preview

ASCII screenshot:

```
+------------------------------------------------------+
| Containers                                           |
|                                                      |
| > web-api [stopped] nginx:latest                     |
|   worker  [running] alpine:latest                    |
|                                                      |
| Keys: up/down, enter=toggle, s=start, t=stop, ...    |
|                                                      |
| Command Preview                                      |
|                                                      |
| container start web-api                              |
|                                                      |
| Confirm (y/n)                                        |
+------------------------------------------------------+
```

## Workflow: Pull an Image

1. Press `p` from the main screen
2. Enter the image reference
3. Confirm the preview to start pulling

ASCII screenshot:

```
+------------------------------------------------------+
| Pull Image                                           |
|                                                      |
| Image reference: nginx:latest                        |
|                                                      |
| Progress: [===========                  ]            |
|                                                      |
| Command Preview                                      |
|                                                      |
| container image pull nginx:latest                    |
|                                                      |
| Confirm (y/n)                                        |
+------------------------------------------------------+
```

## Workflow: Build from Containerfile/Dockerfile

1. Press `b` from the main screen
2. Select a build file in the file picker
3. Enter a tag
4. Confirm the preview to build

ASCII screenshot:

```
+------------------------------------------------------+
| Build Image                                          |
|                                                      |
| File: ./Containerfile                                |
| Tag: my-app:latest                                   |
|                                                      |
| Progress: [========                     ]            |
|                                                      |
| Command Preview                                      |
|                                                      |
| container build -t my-app:latest -f ./Containerfile . |
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
2. Press `s` to start or `t` to stop
3. Confirm in the yes/no prompt

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

- `~/.config/apple-tui/config`
- `~/Library/Application Support/apple-tui/config`

Writes go to:

- `~/Library/Application Support/apple-tui/config`

Example TOML:

```toml
default_build_file = "Containerfile"
confirm_destructive_actions = true
theme_mode = "auto"
refresh_on_focus = false
log_retention_days = 7
```

Logs are stored at:

- `~/Library/Application Support/apple-tui/command.log`

## Troubleshooting

- "apple container CLI not found" -> install from https://github.com/apple/container
- "daemon is not running" -> open Daemon screen and start it
- Build errors -> ensure a Containerfile or Dockerfile exists in the chosen folder
- Permission errors -> confirm your user can run the CLI commands without sudo

## Help Screen

Press `?` from any screen to review shortcuts and paths.
