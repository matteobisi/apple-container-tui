# Research: Enhanced TUI Display with Table Layout

**Feature**: 003-tui-table-format  
**Date**: 2026-02-17  
**Status**: Complete

## Research Tasks

### 1. Alternate Screen Buffer in Bubbletea

**Question**: How to enable alternate screen mode in Bubbletea applications to achieve btop-like full-screen behavior?

**Decision**: Use `tea.WithAltScreen()` program option when creating the tea.Program in main.go

**Rationale**: 
- Bubbletea provides native alternate screen buffer support through the `tea.WithAltScreen()` option
- This is the standard approach used by all full-screen Bubbletea applications
- Handles terminal state save/restore automatically on program exit
- Works reliably across Terminal.app, iTerm2, and other ANSI-compatible terminals
- No additional dependencies or manual ANSI escape sequence management required

**Implementation Notes**:
- Modify `cmd/actui/main.go` to add `tea.WithAltScreen()` to tea.NewProgram() call
- No changes needed in screen components - Bubbletea handles buffer switching transparently
- Terminal dimensions are still available through tea.WindowSizeMsg as normal

**Alternatives Considered**:
- Manual ANSI escape sequences (`\033[?1049h` for enter, `\033[?1049l` for exit) - Rejected because Bubbletea's built-in support is more robust and handles edge cases (signals, crashes) better
- Third-party terminal libraries - Rejected as unnecessary given native Bubbletea support

---

### 2. Table Rendering with Dynamic Column Sizing

**Question**: What's the best approach for rendering aligned tables with dynamic column widths in a terminal UI?

**Decision**: Create custom table component using width calculation algorithm with priority-based truncation

**Rationale**:
- Bubbles (Bubbletea's component library) includes a `table` component, but it's designed for interactive selection and doesn't match our existing list navigation patterns
- Custom implementation provides full control over column prioritization (Name > State/Tag > Image/Digest)
- Can integrate seamlessly with existing cursor-based selection and keyboard handling
- Allows proper handling of minimum widths and ellipsis truncation per spec requirements

**Algorithm**:
```
1. Measure content width for each column across all rows
2. Calculate fixed widths for State/Tag columns (longest value + 2 padding)
3. Allocate remaining width to Name column (priority)
4. If space remains, allocate to Image/Digest column
5. If width < 80 chars, truncate Image/Digest first, then Name if necessary
6. Apply ellipsis (...) to truncated values
```

**Implementation Notes**:
- Create `src/ui/table.go` with reusable TableRenderer struct
- Methods: `NewTableRenderer()`, `SetColumns()`, `SetRows()`, `Render(width int) string`
- Use `lipgloss.NewStyle()` for column alignment, padding, and bold headers
- Unicode box-drawing characters for separators: `─` (U+2500)

**Alternatives Considered**:
- Bubbles table component - Rejected because it enforces its own selection model incompatible with our screen navigation
- Fixed-width columns - Rejected because it wastes space for short values and breaks for long repository names
- Printf-style formatting with %*s - Rejected because it doesn't handle overflow gracefully

---

### 3. Row Highlighting with Inverse Video

**Question**: How to implement inverse video (swapped foreground/background) for selected row in Lipgloss?

**Decision**: Use Lipgloss `Reverse(true)` style attribute for selected row

**Rationale**:
- Lipgloss provides `.Reverse(true)` method that applies ANSI inverse video escape sequences
- Works universally across all color schemes (light/dark terminals)
- More reliable than explicit foreground/background color swaps which break in custom terminal themes
- Matches standard TUI conventions (vim, htop, less)
- Zero configuration - works immediately without theme detection

**Implementation**:
```go
selectedStyle := lipgloss.NewStyle().Reverse(true)
normalStyle := lipgloss.NewStyle()

for i, row := range rows {
    style := normalStyle
    if i == cursor {
        style = selectedStyle
    }
    renderedRow := style.Render(row)
}
```

**Alternatives Considered**:
- Explicit color swap (e.g., white-on-black → black-on-white) - Rejected because it requires detecting current theme and breaks with custom color schemes
- Arrow prefix (▶) instead of highlighting - Rejected because spec explicitly requires row highlighting
- Background color only - Rejected because inverse video provides stronger visual distinction

---

### 4. Terminal Resize Handling

**Question**: How to efficiently re-render tables when terminal is resized?

**Decision**: Leverage Bubbletea's tea.WindowSizeMsg to trigger table re-render with new dimensions

**Rationale**:
- Bubbletea automatically sends tea.WindowSizeMsg on SIGWINCH (terminal resize signal)
- All screens already handle this message for basic layout adjustments
- Table component can recalculate column widths on each render based on passed width
- No caching needed - rendering is fast enough (<1ms for typical table sizes)
- Bubbletea's message-based architecture ensures consistent state updates

**Implementation Notes**:
- Container_list.go and image_list.go already handle WindowSizeMsg
- Store width/height in screen structs
- Pass current width to table.Render(width) on each View() call
- Table component recalculates layout on every render (stateless)

**Alternatives Considered**:
- Cache rendered output and invalidate on resize - Rejected as premature optimization; adds complexity without measurable benefit for typical list sizes
- Manual signal handling - Rejected because Bubbletea abstracts this properly

---

### 5. Digest Truncation Strategy

**Question**: How to truncate SHA256 digests to 12 characters while handling various input formats?

**Decision**: Strip `sha256:` prefix if present, then take first 12 hexadecimal characters

**Rationale**:
- Matches Docker CLI convention for digest display
- Apple Container includes `sha256:` prefix in digest field (seen in image metadata)
- 12 characters provide sufficient uniqueness for human identification (2^48 combinations)
- Fixed width allows table column sizing to be predictable

**Implementation**:
```go
func TruncateDigest(digest string) string {
    // Remove sha256: prefix if present
    clean := strings.TrimPrefix(digest, "sha256:")
    // Take first 12 chars
    if len(clean) > 12 {
        return clean[:12]
    }
    return clean
}
```

**Edge Cases**:
- Empty digest: return empty string (displayed as "-" in table)
- Short digest (<12 chars): return as-is without padding
- Non-hex characters: display as-is (let Apple Container validation handle invalid formats)

---

### 6. Table Header Styling

**Question**: How to style table headers with bold text and separator line?

**Decision**: Use Lipgloss `.Bold(true)` for header text and render Unicode horizontal line below

**Rationale**:
- Lipgloss `.Bold(true)` applies ANSI bold escape sequence (works universally)
- Unicode box-drawing character `─` (U+2500) for horizontal lines is widely supported in modern terminal fonts
- Two-technique approach (bold + line) provides clear visual hierarchy
- Separator line naturally divides header from data without additional spacing

**Implementation**:
```go
headerStyle := lipgloss.NewStyle().Bold(true)
separator := strings.Repeat("─", totalWidth)

output := headerStyle.Render(headerRow) + "\n" + separator + "\n"
```

**Fallback**: If Unicode box-drawing not supported, fallback to ASCII dash `-` (detection via render test not needed - U+2500 is standard in macOS Terminal.app)

---

## Technology Stack Summary

| Component | Technology | Version | Purpose |
|-----------|------------|---------|---------|
| TUI Framework | Bubbletea | v1.2.4 | Event loop, screen management, alternate screen |
| Styling | Lipgloss | v1.0.0 | Text formatting, colors, alignment, inverse video |
| Language | Go | 1.21 | Implementation language |
| Terminal | ANSI-compatible | - | Target (Terminal.app, iTerm2 on macOS) |

## Best Practices Applied

1. **Stateless Rendering**: Table component recalculates layout on each render rather than caching, simplifying state management
2. **Priority-Based Truncation**: Name column prioritized over Image/Digest follows information hierarchy principles
3. **Universal Styling**: Inverse video and bold text work across all terminal color schemes
4. **Native Integration**: Use Bubbletea/Lipgloss primitives rather than manual ANSI sequences for maintainability
5. **Graceful Degradation**: Tables remain readable even at 80-character minimum width

## Open Questions

None - all technical decisions resolved with reasonable defaults aligned with industry standards.
