# Apple Container TUI

## Overview

Apple Container TUI is a keyboard-first terminal UI for managing [Apple Container](https://github.com/apple/container)
operations on macOS. It lets you list containers, start/stop/delete them,
pull images, build from Containerfiles, and manage the daemon with safe command
previews and confirmations.

**This repository was created as a proof-of-concept to test and document the
capabilities of spec-kit, a structured software development workflow. The full
process and insights will be documented on [www.msbiro.net](https://www.msbiro.net).**

### Development Process

This project demonstrates a hybrid AI-assisted development approach:

- **Specification Phase**: All spec-kit artifacts (plan, spec, tasks, data model,
  contracts) were created using Claude Sonnet 4.5
- **Implementation Phase**: Code implementation was performed using GPT-5.2-codex
- **Validation Phase**: Final testing, validation, and project oversight by the
  author

This workflow showcases how different AI models can be leveraged for their
strengths across different phases of software development.

## Features

- Container list with container action submenus (start/stop/logs/shell)
- Safe delete with type-to-confirm
- Image management screen (`i`) with list/pull/build/prune
- Image submenu with inspect/delete
- Daemon start/stop controls
- Command preview before execution
- Dry-run mode for safe practice
- JSONL command logs with rotation

## Interface

The TUI provides a clean, keyboard-driven interface with tabular layouts for easy scanning:

**Container List View**:
```
Containers

Name                                             State      Base Image                                            
───────────────────────────────────────────────────────────────────────────────────────────────────────────────
buildkit                                         stopped    ghcr.io/apple/container-builder-shim/builder:0.7.0
cba13176-5dae-497f-a74b-381671056c3b             stopped    markitdown:latest
17007fa5-710a-4ac2-98e0-7923cb26153f             stopped    docker.io/library/ubuntu:latest
808552b9-d78e-4448-a691-927c3848b4b5             stopped    docker.io/library/ubuntu:latest
───────────────────────────────────────────────────────────────────────────────────────────────────────────────

Keys: up/down, enter=submenu, s=start, t=stop, d=delete(!), i=images, r=refresh, m=manage, ?=help, q=quit
```

**Image List View**:
```
Images

Name                                             Tag        Digest
───────────────────────────────────────────────────────────────────────────────────────────────────────────────
docker.io/library/ubuntu                         latest     c2a6e3c0        
markitdown                                       latest     8f3d91a2
───────────────────────────────────────────────────────────────────────────────────────────────────────────────

Keys: up/down=navigate, enter=submenu, p=pull, b=build, n=image-prune, r=refresh, esc=back
```

The application launches in full-screen mode (alternate screen buffer), clearing the terminal
for an immersive experience and restoring your previous terminal content on exit.

All operations show command previews before execution, and destructive actions
require explicit confirmation for safety.

## Installation

### Prerequisites

- macOS 26.x on Apple Silicon
- Apple Container CLI installed and in PATH (tested with version 0.9.0)
- Go 1.21+ if building from source

Verify the CLI is available:

```bash
container system version
```

### Build from source

```bash
go mod download
go build -o actui ./cmd/actui
```

## Quick Start

Launch the TUI:

```bash
./actui
```

Dry-run mode (preview only, no execution):

```bash
./actui --dry-run
```

Helpful keys:

- `?` for help
- `i` to open image management
- `p`/`b`/`n` for pull/build/prune inside image view
- `m` to manage the daemon

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

## Development

Format and lint:

```bash
gofmt -w ./...
golangci-lint run ./...
```

Run tests:

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details. You are free to use, modify, and distribute this software with
proper attribution to the original author.
