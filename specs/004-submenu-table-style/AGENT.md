# Feature 004: Submenu Table Styling - Agent Context

**Feature**: 004-submenu-table-style  
**Created**: February 17, 2026  
**Status**: Specification Complete, Ready for Planning  
**Depends On**: Feature 003-tui-table-format

## Quick Context

This feature extends the visual improvements from feature 003-tui-table-format to submenu screens and other UI elements that weren't covered in the original table formatting work.

### What Was Done in 003-tui-table-format

Feature 003 introduced professional table formatting for the main list views:
- **Container List**: Table with Name/State/Base Image columns
- **Image List**: Table with Name/Tag/Digest columns
- **Visual Elements**: Bold headers, horizontal separators (─), inverse video selection
- **New Component**: `src/ui/table.go` with TableColumn, TableRow, Table types
- **Styling**: Used Lipgloss `.Bold(true)` and `.Reverse(true)` for consistent appearance

### What's Missing (This Feature's Scope)

The 003 work only covered **list screens**. These screens still use the old styling:

#### Container Submenu (`src/ui/container_submenu.go`)
- **Current**: Plain text with `>` cursor prefix, basic color highlighting
- **Needs**: Bold headers, horizontal separators, inverse video selection

#### Image Submenu (`src/ui/image_submenu.go`)  
- **Current**: Plain text with `>` cursor prefix, basic color highlighting
- **Needs**: Bold headers, horizontal separators, inverse video selection

#### Help Screen (`src/ui/help.go`)
- **Current**: Plain text sections with no visual hierarchy
- **Needs**: Bold section headers, horizontal separators

#### Daemon Control Screen (`src/ui/daemon_control.go`)
- **Current**: Plain text menu with cursor prefix
- **Needs**: Consistent styling with bold headers and inverse video selection

## Key Implementation Files

### Source Files to Modify
- `src/ui/container_submenu.go` - Container action menu
- `src/ui/image_submenu.go` - Image action menu
- `src/ui/help.go` - Help/reference screen
- `src/ui/daemon_control.go` - Daemon management screen (optional)

### Reusable Components (Already Exist)
- `src/ui/table.go` - Has helper functions that might be reusable
  - Note: Full table rendering not needed for submenus (no columns)
  - Can reuse separator line logic: `strings.Repeat("─", width)`
- Lipgloss styles already available from 003:
  - `lipgloss.NewStyle().Bold(true)` - For headers
  - `lipgloss.NewStyle().Reverse(true)` - For selection

## Design Patterns from 003-tui-table-format

### Bold Headers
```go
headerStyle := lipgloss.NewStyle().Bold(true)
builder.WriteString(headerStyle.Render("Container Details") + "\n")
```

### Horizontal Separator
```go
builder.WriteString(strings.Repeat("─", width) + "\n")
```

### Inverse Video Selection
```go
normalStyle := lipgloss.NewStyle()
selectedStyle := lipgloss.NewStyle().Reverse(true)

for i, item := range items {
    style := normalStyle
    if i == cursor {
        style = selectedStyle
    }
    builder.WriteString(style.Render(item) + "\n")
}
```

## User Journey

1. **Problem**: After implementing 003, user noticed submenus still looked old/inconsistent
2. **Gap**: Table formatting improved lists but not submenus/help screens
3. **Solution**: This feature (004) applies the same visual language to all remaining screens

## Success Metrics

- Visual consistency: All screens use same bold/separator/inverse-video pattern
- User experience: No cognitive dissonance when navigating between list and submenu screens
- Professional appearance: Every screen follows the established design language

## Next Steps

1. Run `/speckit.plan` to create implementation plan
2. Focus on P1 (container submenu) first - most frequently used
3. Apply patterns consistently across all target screens
4. Manual testing to verify visual consistency with 003 work

## Constitutional Compliance

Same as 003-tui-table-format:
- **Principle I (Command-Safe TUI)**: ✅ Display-only changes, no command execution changes
- **Principle IV (Clear Observability)**: ✅ Enhanced - better visual hierarchy improves scannability
- **All Others**: ✅ Compliant - pure visual improvements

## References

- Feature 003 spec: `specs/003-tui-table-format/spec.md`
- Feature 003 implementation: `src/ui/table.go`, `src/ui/container_list.go`, `src/ui/image_list.go`
- Current submenus: `src/ui/container_submenu.go`, `src/ui/image_submenu.go`
