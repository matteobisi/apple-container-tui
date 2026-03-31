# container-tui Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-03-31

## Active Technologies
- Local filesystem (config.toml, JSONL logs) (001-rename-binary-actui)
- Go 1.21 + Bubbletea v1.2.4 (TUI framework), Bubbles v0.20.0 (TUI components), Cobra v1.8.1 (CLI), Viper v1.19.0 (config) (002-refactor-menu-images)
- Local filesystem for logs and config (~/Library/Application Support/actui/), JSONL command logs (002-refactor-menu-images)
- Go 1.21 + Bubbletea v1.2.4 (TUI framework), Lipgloss v1.0.0 (styling), Bubbles v0.20.0 (components) (003-tui-table-format)
- N/A (display-only feature) (003-tui-table-format)
- Go 1.21+ (tested on Go 1.26.0) + Bubbletea v1.2.4 (TUI framework), Lipgloss v1.0.0 (styling), Bubbles v0.20.0 (004-submenu-table-style)
- Go 1.21+ (module target), validated in local Apple Container environment on macOS 26.x + Bubble Tea v1.2.4, Bubbles v0.20.0, Lipgloss v1.0.0, Cobra v1.8.1, Viper v1.19.0, Go standard library `encoding/json` (005-expand-container-workflows)
- Local filesystem only for exported OCI tar archives; no new persistent app-owned state required (005-expand-container-workflows)

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
- 005-expand-container-workflows: Added Go 1.21+ (module target), validated in local Apple Container environment on macOS 26.x + Bubble Tea v1.2.4, Bubbles v0.20.0, Lipgloss v1.0.0, Cobra v1.8.1, Viper v1.19.0, Go standard library `encoding/json`
- 004-submenu-table-style: Added Go 1.21+ (tested on Go 1.26.0) + Bubbletea v1.2.4 (TUI framework), Lipgloss v1.0.0 (styling), Bubbles v0.20.0
- 003-tui-table-format: Added Go 1.21 + Bubbletea v1.2.4 (TUI framework), Lipgloss v1.0.0 (styling), Bubbles v0.20.0 (components)


<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
