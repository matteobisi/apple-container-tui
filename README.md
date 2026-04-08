# Apple Container TUI

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Build](https://github.com/matteobisi/apple-container-tui/actions/workflows/build-binary.yml/badge.svg)](https://github.com/matteobisi/apple-container-tui/actions/workflows/build-binary.yml)
[![Latest Release](https://img.shields.io/github/v/release/matteobisi/apple-container-tui)](https://github.com/matteobisi/apple-container-tui/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/matteobisi/apple-container-tui)](go.mod)
[![Contributors](https://img.shields.io/github/contributors/matteobisi/apple-container-tui)](https://github.com/matteobisi/apple-container-tui/graphs/contributors)
[![Project Status: Proof of Concept](https://img.shields.io/badge/status-proof--of--concept-orange)](https://www.msbiro.net)

> Keyboard-first terminal UI for managing [Apple Container](https://github.com/apple/container) on macOS.

## Overview

Apple Container TUI (`actui`) lets you list, start, stop, delete, and export containers; pull and build images; browse registries; and control the daemon — all from the terminal with safe command previews and confirmation prompts before any destructive action.

This repository is a proof-of-concept showcasing [spec-kit](https://www.msbiro.net), a structured AI-assisted development workflow. Full process notes are published on [www.msbiro.net](https://www.msbiro.net).

## Features

- Container list with action submenus (start / stop / logs / shell / export)
- Safe delete with type-to-confirm
- Image management (`i`) — list, pull, build, prune, inspect, delete
- Dedicated registries view (`g`)
- Build form with `--pull` toggle (enabled by default)
- Daemon start/stop with structured status (`running` / `stopped` / `unknown`)
- Command preview before every execution
- Dry-run mode for safe practice
- JSONL command logs with rotation

## Quick Start

### Prerequisites

- macOS 26.x on Apple Silicon
- [Apple Container CLI](https://github.com/apple/container) installed and in `PATH`
- Go 1.21+ if building from source

Verify the CLI is available:

```bash
container system version
```

### Install

```bash
go mod download
go build -o actui ./cmd/actui
```

Or download the pre-built `actui-darwin-arm64` binary from the [latest release](https://github.com/matteobisi/apple-container-tui/releases/latest).

### Run

```bash
./actui            # normal mode
./actui --dry-run  # preview only, no commands executed
```

## Key Bindings

| Context | Key | Action |
|---|---|---|
| Anywhere | `?` | Help screen |
| Anywhere | `q` | Quit |
| Container list | `enter` | Open container submenu |
| Container list | `s` / `t` | Start / Stop selected container |
| Container list | `d` | Delete (type-to-confirm) |
| Container list | `i` | Open image management |
| Container list | `m` | Daemon management |
| Container list | `r` | Refresh |
| Image list | `p` | Pull image |
| Image list | `b` | Build from Containerfile |
| Image list | `g` | Browse registries |
| Image list | `n` | Prune unused images |
| Image list | `esc` | Back to container list |

For full workflow walkthroughs and ASCII screenshots, see [docs/user-guide.md](docs/user-guide.md).

## Configuration

Config is read from (in order):

- `~/.config/actui/config`
- `~/Library/Application Support/actui/config`

Writes go to `~/Library/Application Support/actui/config`.

Example TOML:

```toml
default_build_file = "Containerfile"
confirm_destructive_actions = true
theme_mode = "auto"
refresh_on_focus = false
log_retention_days = 7
```

Logs: `~/Library/Application Support/actui/command.log`

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

## Development Process

This project demonstrates a hybrid AI-assisted workflow:

- **Specification**: Artifacts (plan, spec, tasks, data model, contracts) authored with Claude Sonnet
- **Implementation**: Code written by Codex
- **Validation**: Final testing and oversight by the author

## Documentation

| Document | Description |
|---|---|
| [docs/user-guide.md](docs/user-guide.md) | Full workflow walkthroughs, key bindings, troubleshooting |
| [docs/binary-build-automation.md](docs/binary-build-automation.md) | CI build, SBOM generation, release automation |
| [docs/security-automation.md](docs/security-automation.md) | OSSF Scorecard, Dependabot, branch protection |
| [docs/ai-menu-map.md](docs/ai-menu-map.md) | AI agent navigation map and UI ownership guide |
| [docs/speckit-security-hardening-retrospective.md](docs/speckit-security-hardening-retrospective.md) | Retrospective on the security hardening sprint |

## Security

Repository security is enforced via OSSF Scorecard, Dependabot, and branch protection on `main`. Releases include SBOM (SPDX 2.3 JSON) and GitHub provenance attestations. See [docs/security-automation.md](docs/security-automation.md) and [SECURITY.md](SECURITY.md) for details.

## Releases & Binaries

Merging to `main` automatically triggers a build and publishes a versioned GitHub Release with the `actui-darwin-arm64` binary and its SBOM. See [docs/binary-build-automation.md](docs/binary-build-automation.md) for the full release pipeline.

## Contributing

All changes (human or AI-authored) follow the same workflow:

1. Create a branch from `main`
2. Commit focused changes with clear messages
3. Open a pull request to `main`
4. Wait for required checks to pass before merge

For AI agent guidelines, see [AGENTS.md](AGENTS.md) and [docs/ai-menu-map.md](docs/ai-menu-map.md).

Direct pushes to `main` are blocked.

## License

MIT — see [LICENSE](LICENSE) for details.
