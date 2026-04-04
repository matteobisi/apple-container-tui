# Apple Container TUI

## Overview

Apple Container TUI is a keyboard-first terminal UI for managing [Apple Container](https://github.com/apple/container)
operations on macOS. It lets you list containers, start/stop/delete them,
pull images, browse registries, export stopped containers, build from
Containerfiles, and manage the daemon with safe command previews and
confirmations.

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

### AI Contributor Notes

For AI-oriented navigation and UI ownership guidance, see [AGENTS.md](AGENTS.md) and the detailed map in [docs/ai-menu-map.md](docs/ai-menu-map.md).

## Features

- Container list with container action submenus (start/stop/logs/shell)
- Stopped-container export workflow with destination selection, command preview, and optional cleanup confirmation
- Safe delete with type-to-confirm
- Image management screen (`i`) with list/pull/build/prune and dedicated registries view
- Image submenu with inspect/delete
- Build form pull toggle enabled by default
- Daemon start/stop controls with structured status parsing and unknown fallback
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

Keys: up/down=navigate, enter=submenu, p=pull, g=registries, b=build, n=image-prune, r=refresh, esc=back
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
- `p`/`g`/`b`/`n` for pull/registries/build/prune inside image view
- `m` to manage the daemon

Additional workflows:

- Export is available from the submenu of stopped containers only, with an explicit prompt before removing the temporary export image
- Build previews show whether `--pull` will be applied
- Daemon status can render `running`, `stopped`, or `unknown`

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

## Repository Security Automation

This repository uses GitHub-native security automation:

- OSSF Scorecard workflow in [.github/workflows/scorecard.yml](.github/workflows/scorecard.yml)
- Dependabot configuration in [.github/dependabot.yml](.github/dependabot.yml)

Enforcement baseline on `main`:

- Required status check: `OSSF Scorecard`
- Branch protection enabled (force-push and deletion blocked)
- Dependabot security updates enabled

Operational guidance, branch-protection mapping, and troubleshooting are documented in [docs/security-automation.md](docs/security-automation.md).

## Binary Build and Release Automation

Build workflow operations, SBOM generation, retention policy, and troubleshooting are documented in [docs/binary-build-automation.md](docs/binary-build-automation.md).

Release publication is automated: a successful build on `main` automatically triggers the `Publish Release` workflow, which applies a deterministic semantic version tag (`v0.1.0`, `v0.1.1`, …) and publishes a GitHub Release with the `actui-darwin-arm64` macOS Apple Silicon binary and an `actui-darwin-arm64.spdx.json` SBOM (SPDX 2.3 JSON) attached. The full trigger chain, version-labeling policy, SBOM generation details, duplicate handling, and operator validation checklist are all in [docs/binary-build-automation.md](docs/binary-build-automation.md).

## Contribution Workflow

All changes (human or AI-authored) should follow the same repository workflow:

1. Create a new branch from `main`
2. Commit focused changes with clear messages
3. Open a pull request to `main`
4. Wait for required checks to pass before merge

Direct pushes to `main` are not part of the normal workflow.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details. You are free to use, modify, and distribute this software with
proper attribution to the original author.
