# container-tui Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-04-04

## Active Technologies
- Local filesystem (config.toml, JSONL logs) (001-rename-binary-actui)
- Go 1.21 + Bubbletea v1.2.4 (TUI framework), Bubbles v0.20.0 (TUI components), Cobra v1.8.1 (CLI), Viper v1.19.0 (config) (002-refactor-menu-images)
- Local filesystem for logs and config (~/Library/Application Support/actui/), JSONL command logs (002-refactor-menu-images)
- Go 1.21 + Bubbletea v1.2.4 (TUI framework), Lipgloss v1.0.0 (styling), Bubbles v0.20.0 (components) (003-tui-table-format)
- N/A (display-only feature) (003-tui-table-format)
- Go 1.21+ (tested on Go 1.26.0) + Bubbletea v1.2.4 (TUI framework), Lipgloss v1.0.0 (styling), Bubbles v0.20.0 (004-submenu-table-style)
- Go 1.21+ (module target), validated in local Apple Container environment on macOS 26.x + Bubble Tea v1.2.4, Bubbles v0.20.0, Lipgloss v1.0.0, Cobra v1.8.1, Viper v1.19.0, Go standard library `encoding/json` (005-expand-container-workflows)
- Local filesystem only for exported OCI tar archives; no new persistent app-owned state required (005-expand-container-workflows)
- YAML (GitHub Actions workflow and Dependabot configuration), Go 1.21 module context for dependency ecosystem detection + GitHub Actions runner (`ubuntu-latest`), OSSF Scorecard GitHub Action, Dependabot version updates engine (006-repo-security-hardening)
- Git repository configuration files only; no runtime data storage (006-repo-security-hardening)
- Go 1.21+ for project build, GitHub Actions workflow YAML + `actions/checkout`, `actions/setup-go`, `actions/upload-artifact`, Go toolchain from `go.mod` (007-build-binary-action)
- GitHub Actions artifact storage and repository documentation files (007-build-binary-action)
- Go 1.21+ for build command; GitHub Actions YAML for automation + `actions/checkout`, `actions/setup-go`, `actions/upload-artifact`, Go toolchain from `go.mod` (007-build-binary-action)
- GitHub Actions artifact storage and repository Markdown docs in `docs/` (007-build-binary-action)
- GitHub Actions YAML workflows plus shell scripting on Ubuntu runners + `actions/checkout@v4`, `actions/setup-go@v5`, `actions/upload-artifact@v4`, GitHub CLI/API release actions (`actions/download-artifact`, `softprops/action-gh-release` or equivalent) (010-auto-release-publish)
- GitHub Actions artifact storage and GitHub Releases assets; repository Markdown docs in `docs/` (010-auto-release-publish)
- Go 1.24.2 (module `container-tui`) + `anchore/sbom-action` (Syft, composite GHA); existing: `actions/checkout`, `actions/setup-go`, `actions/upload-artifact`, `actions/download-artifact` (011-sbom-binary-build)
- GitHub Actions workflow artifacts (14-day retention) + GitHub Release assets (permanent) (011-sbom-binary-build)
- N/A (no Go/Python source changes) + GitHub Actions — `actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` (v4.1.0, Node 24 compatible); Docker Hub — `python:3.13-slim-bullseye` (012-security-hardening-quick-wins)
- Go 1.24 (TUI application); YAML (GitHub Actions workflows); Markdown (SECURITY.md) + `actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` (v4.1.0, Node 24 compatible) (012-security-hardening-quick-wins)
- N/A — no runtime storage changes (012-security-hardening-quick-wins)

- Go 1.21+ (chosen for optimal balance of productivity, performance, binary distribution, and TUI library maturity) + Bubbletea v1.2.4 (TUI framework), Lipgloss v1.0.0 (styling), Bubbles v0.20.0 (UI components), Cobra (CLI), Viper (config management) (001-apple-container-tui)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.21+ (chosen for optimal balance of productivity, performance, binary distribution, and TUI library maturity)

## Code Style

Go 1.21+ (chosen for optimal balance of productivity, performance, binary distribution, and TUI library maturity): Follow standard conventions

## Recent Changes
- 012-security-hardening-quick-wins: Added Go 1.24 (TUI application); YAML (GitHub Actions workflows); Markdown (SECURITY.md) + `actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` (v4.1.0, Node 24 compatible)
- 012-security-hardening-quick-wins: Added Go 1.24 (TUI application); YAML (GitHub Actions workflows); Markdown (SECURITY.md) + `actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` (v4.1.0, Node 24 compatible)
- 012-security-hardening-quick-wins: Added N/A (no Go/Python source changes) + GitHub Actions — `actions/attest-build-provenance@a2bbfa25375fe432b6a289bc6b6cd05ecd0c4c32` (v4.1.0, Node 24 compatible); Docker Hub — `python:3.13-slim-bullseye`


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
