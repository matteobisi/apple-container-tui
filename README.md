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

- Container list with start/stop and refresh
- Safe delete with type-to-confirm
- Image pull and build workflows
- Daemon start/stop controls
- Command preview before execution
- Dry-run mode for safe practice
- JSONL command logs with rotation

## Interface

The TUI provides a clean, keyboard-driven interface for container management:

```
Containers

> buildkit [stopped] ghcr.io/apple/container-builder-shim/builder:0.7.0
  cba13176-5dae-497f-a74b-381671056c3b [stopped] markitdown:latest
  17007fa5-710a-4ac2-98e0-7923cb26153f [stopped] docker.io/library/ubuntu:latest
  808552b9-d78e-4448-a691-927c3848b4b5 [stopped] docker.io/library/ubuntu:latest

Keys: up/down, enter=toggle, s=start, t=stop, d=delete(!), r=refresh, p=pull, b=build, m=manage, ?=help, q=quit
```

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
go build -o apple-tui ./cmd/apple-tui
```

## Quick Start

Launch the TUI:

```bash
./apple-tui
```

Dry-run mode (preview only, no execution):

```bash
./apple-tui --dry-run
```

Helpful keys:

- `?` for help
- `p` to pull an image
- `b` to build an image
- `m` to manage the daemon

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
