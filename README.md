# Apple Container TUI

## Overview

Apple Container TUI is a keyboard-first terminal UI for managing Apple Container
operations on macOS. It lets you list containers, start/stop/delete them,
pull images, build from Containerfiles, and manage the daemon with safe command
previews and confirmations.

## Features

- Container list with start/stop and refresh
- Safe delete with type-to-confirm
- Image pull and build workflows
- Daemon start/stop controls
- Command preview before execution
- Dry-run mode for safe practice
- JSONL command logs with rotation

## Installation

### Prerequisites

- macOS 26.x on Apple Silicon
- Apple Container CLI installed and in PATH
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

TBD.
