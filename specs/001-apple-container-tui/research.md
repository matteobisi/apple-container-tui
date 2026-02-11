# Research: Apple Container TUI

**Phase 0 Output** | **Date**: 2026-02-11  
**Purpose**: Resolve NEEDS CLARIFICATION items from Technical Context

## Technology Stack Decision

### Language and Framework

**Decision**: Go 1.21+ with Bubbletea/Lipgloss TUI framework

**Rationale**: Go with Bubbletea offers the optimal balance of developer productivity, performance, and deployment simplicity for this macOS container TUI project. The Elm-inspired architecture makes complex state management intuitive, automatic resize handling eliminates a major requirement concern, and compilation produces a single native ARM64 binary with zero runtime dependencies. Response times of 15-50ms are well under the 100ms requirement, and the gentle learning curve with fast compilation (1-3 seconds) significantly accelerates development velocity compared to Rust.

**Alternatives Considered**:
- **Rust + ratatui**: More performant (sub-10ms) but 2-3x longer development time due to ownership/borrowing complexity. Chosen only if ultra-low latency is critical or team already knows Rust.
- **Python + Textual**: Fastest prototyping with CSS-like styling and rich graphics, but binary distribution is complex on macOS (code signing issues with PyInstaller), slower startup (100-300ms), and GIL can cause input lag.
- **Swift**: Excellent macOS integration but no mature TUI framework exists - would require extensive custom work building on ncurses or Foundation APIs.

### Primary Dependencies

**Core TUI Stack**:
- `github.com/charmbracelet/bubbletea` v1.2.4 - TUI framework with Elm architecture
- `github.com/charmbracelet/lipgloss` v1.0.0 - Styling and layout
- `github.com/charmbracelet/bubbles` v0.20.0 - Pre-built components (viewport, list, spinner)

**Supporting Libraries**:
- `github.com/spf13/cobra` - CLI flag parsing and subcommands
- `github.com/spf13/viper` - Configuration file management (TOML/YAML/JSON)
- `golang.org/x/term` - Terminal capability detection

### Testing Framework

**Decision**: Go's built-in testing with `go test`

**Approach**:
- **Unit tests**: Test command builders, models, and business logic with standard Go tests
- **Contract tests**: Bubbletea's message-based architecture enables testing by sending messages and verifying model updates without rendering
- **Integration tests**: Optional tests that run against actual Apple Container CLI when available
- **Snapshot/golden file testing**: For UI component output verification

### Binary Distribution

**Approach**: Single native binary via `go build`

**Build command**:
```bash
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "apple-tui" ./cmd/apple-tui
```

**Distribution options**:
1. **Direct binary**: 8-15MB self-contained executable
2. **Homebrew tap**: For easy installation and updates
3. **GitHub releases**: With checksums and code signing

**Advantages**:
- No runtime dependencies
- Fast startup (<50ms)
- Cross-compilation trivial
- Code signing straightforward with Apple Developer ID

## Command-Safe Architecture Patterns

### Command Preview Implementation

**Pattern**: Multi-stage preview with visual hierarchy

**UI Placement**:
- **During configuration**: Live preview in bottom status bar updates as user selects options
- **Before execution**: Full command in confirmation modal with formatted output

**Example Layout**:
```
┌─────────────────────────────────────┐
│ Select Action:                      │
│ > Start Container                   │
│   Stop Container                    │
│   Delete Container                  │
├─────────────────────────────────────┤
│ Preview: acrun start --container app│
└─────────────────────────────────────┘
```

**Implementation**: Use Lipgloss for syntax highlighting in preview pane, distinguish executable from flags/arguments with color.

### Dry-Run Architecture

**Pattern**: Strategy pattern with executor interface

**Code Structure**:
```go
type CommandExecutor interface {
    Execute(cmd Command) (*Result, error)
}

type DryRunExecutor struct{}
func (d *DryRunExecutor) Execute(cmd Command) (*Result, error) {
    return &Result{
        DryRun: true,
        Command: cmd.String(),
        Output: fmt.Sprintf("Would execute: %s", cmd.String()),
    }, nil
}

type RealExecutor struct{}
func (r *RealExecutor) Execute(cmd Command) (*Result, error) {
    // Actual os/exec.Command execution
}
```

**Runtime Flag**: `--dry-run` CLI flag determines which executor instance to inject at startup

**Testing Strategy**:
- Unit tests default to `DryRunExecutor` for fast, safe testing
- Test command **building** separately from **execution**
- Integration tests use `RealExecutor` only when Apple Container CLI is available

### Destructive Action Safeguards

**Recommended Patterns**:

1. **Color-coded visual indicators**:
   - Red for destructive (delete, kill)
   - Yellow for impactful but reversible (stop, restart)
   - Green for safe operations (list, inspect, logs)

2. **Type-to-confirm for irreversible actions**:
```
┌─────────────────────────────────────────┐
│ ⚠️  Delete Container                    │
│                                         │
│ This will permanently delete:           │
│ • Container: production-db              │
│                                        │
│ Type "production-db" to confirm:        │
│ > _                                     │
└─────────────────────────────────────────┘
```

3. **Two-stage confirmation**: First y/N prompt, then type-to-confirm for critical operations

**UX Principles**:
- Only confirm truly destructive actions (data loss)
- Make safe actions fast (no confirmation for read-only)
- Always allow ESC or Ctrl+C to cancel
- Use symbols (⚠️, 🗑️, ✓) for instant action type recognition

### Command Composition Best Practices

**Validation**: Validate before building
```go
func (b *CommandBuilder) Validate() error {
    // Validate required flags present
    // Validate argument formats
    // Validate mutually exclusive options
    return nil
}

func (b *CommandBuilder) Build() (Command, error) {
    if err := b.Validate(); err != nil {
        return Command{}, err
    }
    return Command{...}, nil
}
```

**Escaping**: Use `exec.Command` properly
```go
// ✅ CORRECT - exec.Command handles escaping automatically
cmd := exec.Command("acrun", "start", "--container", userInput)

// ❌ WRONG - shell injection vulnerability
cmd := exec.Command("sh", "-c", fmt.Sprintf("acrun start %s", userInput))
```

**Key Principles**:
- Never use shell execution (`sh -c`) unless absolutely necessary
- Go's `os/exec` automatically handles argument escaping
- Separate command building from display formatting
- For preview display, use proper quoting without actual shell execution

### Error Display Patterns

**stdout/stderr Handling**: Separate streams with visual distinction

**Buffering vs Streaming**:
- **Buffering** (default): Collect full output, display in formatted panel - better for short commands (<10s)
- **Streaming**: For long operations (logs --follow, build) - use Bubbles viewport with auto-scroll

**UI Components**:
- Bubbles `viewport` for scrollable output
- Bubbles `list` for structured output (container listings)
- Status bar for quick success/failure indication
- Modal/dialog for detailed error messages

**Error Display Format**:
```go
type OutputDisplay struct {
    Status   Status        // success, error, warning
    Title    string        // "Container Started" or "Error: Failed to start"
    Stdout   string        
    Stderr   string        // displayed in red/yellow
    ExitCode int          
    Command  string        // show what command failed
}
```

**Always show for errors**:
1. What operation failed
2. The actual command that was run
3. Full error output (stderr)
4. Exit code
5. Suggested next steps or troubleshooting hints

## macOS Apple Silicon Considerations

### Go Binary Compatibility

**Status**: First-class support since Go 1.16

**Build Configuration**:
- Target: `GOOS=darwin GOARCH=arm64`
- No special flags required for M-series CPUs
- Universal binaries possible via `lipo` if Intel support needed (out of scope)

### Terminal Framework Compatibility

**Crossterm** (Bubbletea's default backend):
- Handles macOS terminal nuances for Terminal.app, iTerm2, Alacritty
- Window resize events work seamlessly
- Keyboard input properly mapped on macOS

### Configuration Paths

**Dual-path Support** (per FR-011):
- Read both: `~/.config/apple-tui/config` and `~/Library/Application Support/apple-tui/config`
- Write to: `~/Library/Application Support/apple-tui/config` (macOS standard)
- Implementation: Check both paths in order, use first found; create macOS standard path for writes

### Apple Container CLI Detection

**Startup Check** (per FR-014):
```go
func checkAppleContainerCLI() error {
    path, err := exec.LookPath("acrun")
    if err != nil {
        return fmt.Errorf("Apple Container CLI not found in PATH. Please install from https://github.com/apple/container")
    }
    // Optionally verify version compatibility
    return nil
}
```

**Error on Missing CLI**: Exit immediately with clear installation instructions

## Performance Analysis

### Response Time Goals

**Target**: <100ms for UI updates

**Expected Performance**:
- Go TUI rendering: 15-50ms per frame
- Command execution: Depends on Apple Container CLI (1-5 seconds typical)
- Configuration load: <10ms
- Container list parsing: <50ms for 100 containers

**Optimization Strategies**:
- Cache container list until manual refresh (per FR-015)
- Lazy load detailed container info
- Use goroutines for concurrent CLI calls when safe
- Debounce rapid user input events

### Memory Footprint

**Expected**: 20-40MB RSS for running TUI

**Considerations**:
- Go runtime overhead: ~10MB
- Bubbletea framework: ~5MB
- Container list in memory: Minimal (<1MB for 100 containers)
- Output buffering: Configurable limit (default 1MB per command output)

## Development Workflow

### Project Setup

**Initialize**:
```bash
go mod init github.com/user/apple-tui
go get github.com/charmbracelet/bubbletea@v1.2.4
go get github.com/charmbracelet/lipgloss@v1.0.0
go get github.com/charmbracelet/bubbles@v0.20.0
```

**Hot Reload**: Use `entr` or `watchexec` for auto-recompile during development
```bash
ls **/*.go | entr -r go run ./cmd/apple-tui
```

### Testing Approach

**Unit Tests**: Fast, no external dependencies
```bash
go test ./...
```

**Contract Tests**: Verify command building
```bash
go test -v ./internal/commands
```

**Integration Tests** (optional): Require Apple Container CLI
```bash
go test -tags=integration ./...
```

### Build and Release

**Development Build**:
```bash
go build -o apple-tui ./cmd/apple-tui
```

**Release Build** (optimized):
```bash
GOOS=darwin GOARCH=arm64 go build \
  -ldflags="-s -w -X main.version=1.0.0" \
  -o apple-tui ./cmd/apple-tui
```

**Code Signing** (for distribution):
```bash
codesign --sign "Developer ID Application: Your Name" apple-tui
```

## Summary

All NEEDS CLARIFICATION items from Technical Context have been resolved:

✅ **Language/Version**: Go 1.21+ chosen for optimal balance of productivity, performance, and binary distribution  
✅ **Primary Dependencies**: Bubbletea v1.2.4, Lipgloss v1.0.0, Bubbles v0.20.0 for TUI framework  
✅ **Testing**: Go's built-in `go test` with strategy pattern for dry-run/execute separation  
✅ **Command Safety**: Executor interface, command preview patterns, type-to-confirm for destructive actions  
✅ **macOS Compatibility**: First-class Go ARM64 support, dual config path handling, startup CLI check  
✅ **Performance**: 15-50ms UI response well under 100ms target, manual refresh per FR-015

**Next Phase**: Generate data-model.md, contracts/, and quickstart.md based on these technology decisions.
