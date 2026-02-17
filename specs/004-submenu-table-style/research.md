# Research: Consistent Table Styling for Submenus

**Feature**: 004-submenu-table-style  
**Date**: 2026-02-17  
**Status**: Complete

## Research Tasks

### 1. Reusable Styling Components from Feature 003

**Question**: What styling components and patterns from 003-tui-table-format can be directly reused for submenu formatting?

**Decision**: Use Lipgloss styling patterns directly; do NOT reuse full Table component

**Rationale**:
- Feature 003 created `src/ui/table.go` with TableColumn, TableRow, Table types for **columnar data**
- Submenus don't need columns - just single-line menu items with selection highlighting
- The core Lipgloss styles are what we need:
  - `lipgloss.NewStyle().Bold(true)` for headers
  - `lipgloss.NewStyle().Reverse(true)` for selection
  - `strings.Repeat("─", width)` for separators
- Using full Table component would be over-engineering (introduces unnecessary column width calculations)

**Implementation Notes**:
- Copy styling patterns, not components
- Apply bold/reverse/separator directly in submenu View() methods
- Keep implementation simple and focused on visual consistency

**Alternatives Considered**:
- Create MenuRenderer component similar to Table - Rejected because submenus vary in structure (some show info, some don't), making a generic component complex
- Extend Table component for single-column use - Rejected because Table is optimized for multi-column data with priority-based width allocation, unnecessary for simple menu lists

---

### 2. Menu Item Selection Highlighting Pattern

**Question**: How should menu items transition from cursor-prefix (`>`) to inverse video selection?

**Decision**: Replace cursor prefix with full-row inverse video highlighting

**Rationale**:
- Matches user expectation set by 003-tui-table-format where rows highlight on selection
- Inverse video provides stronger visual feedback than color accent alone
- More accessible (works in all terminal color schemes without configuration)
- Industry standard pattern (vim, htop, less all use inverse video)

**Implementation Pattern**:
```go
// Before (current submenu pattern):
cursor := " "
if i == m.cursor {
    cursor = ">"
}
line := cursor + " " + action
if i == m.cursor {
    line = RenderAccent(line)  // Just color change
}

// After (new pattern matching 003):
normalStyle := lipgloss.NewStyle()
selectedStyle := lipgloss.NewStyle().Reverse(true)

for i, action := range actions {
    style := normalStyle
    if i == m.cursor {
        style = selectedStyle
    }
    builder.WriteString(style.Render(action) + "\n")
}
```

**Alternatives Considered**:
- Keep cursor prefix + add inverse video - Rejected because redundant visual cue clutters interface
- Use background color only (no reverse) - Rejected because requires theme-specific color choices, less universal

---

### 3. Section Header and Separator Placement

**Question**: Where should bold headers and horizontal separators be placed in submenu screens?

**Decision**: Use 3-section layout with headers and separators between each section

**Layout Pattern**:
```
[Title (already bold from RenderTitle())]

[Bold] Container Details  ← NEW: Bold section header
name: container-name
status: running
image: docker.io/library/ubuntu:latest
─────────────────────────  ← NEW: Horizontal separator

[Bold] Available Actions   ← NEW: Bold section header
Stop container             ← Inverse video when selected
Tail container log
Enter container
Back
─────────────────────────  ← NEW: Horizontal separator

Keys: [keyboard shortcuts]
```

**Rationale**:
- Clear visual hierarchy matches 003 table headers
- Separators divide functional sections (info vs actions vs help)
- Consistent with table layout: headers are bold, data is separated from controls
- Easy to scan - users can jump directly to action section

**Implementation Notes**:
- Use same `RenderTitle()` for screen title (already exists and is bold)
- Add bold style wrapper for section headers
- Insert separator after each major section
- Keep keyboard shortcuts section as-is (already muted styling)

**Alternatives Considered**:
- Single separator between info and actions only - Rejected because doesn't separate actions from keyboard help
- No section headers, just separators - Rejected because doesn't communicate purpose of each section
- Different separator character (═, ━) - Rejected to maintain consistency with 003 which uses ─

---

### 4. Width Handling for Separators

**Question**: How should horizontal separators handle variable terminal width?

**Decision**: Store terminal width in screen model and render separators dynamically

**Rationale**:
- Feature 003 already established pattern: pass `width` through WindowSizeMsg
- Container and image submenus need to add width field to model (currently missing)
- Help screen and daemon control screen also need width tracking
- Separators use `strings.Repeat("─", width)` to span full terminal width

**Implementation**:
```go
// Add to struct:
type ContainerSubmenuScreen struct {
    // ... existing fields
    width int
}

// Add to Update():
case tea.WindowSizeMsg:
    m.width = message.Width
    return m, nil

// Use in View():
builder.WriteString(strings.Repeat("─", m.width) + "\n")
```

**Fallback**: If width is 0 (not yet set), default to 80 characters (same as 003)

**Alternatives Considered**:
- Fixed-width separators (80 chars) - Rejected because doesn't respond to terminal resize
- Query terminal width on every render - Rejected because Bubbletea provides width through WindowSizeMsg, no need for external queries

---

### 5. Help Screen Section Organization

**Question**: How should the help screen be restructured with bold headers and separators?

**Decision**: Organize into 4 sections with bold headers: Navigation, Container Actions, Image Actions, General

**Rationale**:
- Current help screen is flat text making it hard to scan for specific commands
- Logical groupings match user mental model (what can I do with containers? images? navigation?)
- Bold headers make scanning faster
- Separators create visual breaks reducing cognitive load

**Section Structure**:
```
Help

[Bold] Navigation
up/down, j/k       Navigate lists
enter              Open submenu/select
esc                Back/cancel
─────────────────

[Bold] Container Actions
s                  Start container
t                  Stop container
d                  Delete container
... (etc)
─────────────────

[Bold] Image Actions
p                  Pull image
b                  Build image
... (etc)
─────────────────

[Bold] General
r                  Refresh
m                  Manage daemon
?                  Help
q                  Quit
─────────────────

[keyboard shortcuts muted]
```

**Alternatives Considered**:
- Alphabetical listing - Rejected because doesn't group related functionality
- Table format with columns - Rejected because help text varies in length, single column is clearer
- Collapsible sections - Rejected because adds interaction complexity for reference screen

---

## Summary of Design Decisions

| Decision | Approach | Rationale |
|----------|----------|-----------|
| Reusable components | Lipgloss styles only, not Table component | Submenus don't need columnar layout |
| Selection highlighting | Inverse video (.Reverse(true)) | Matches 003, universal across themes |
| Layout structure | 3-section with headers + separators | Clear hierarchy, easy scanning |
| Width handling | Store in model, use WindowSizeMsg | Consistent with 003 pattern |
| Help organization | 4 sections grouped by function | Logical structure matches user tasks |

All decisions prioritize **visual consistency** with feature 003 while keeping implementation **simple and focused** on the specific needs of submenu/help screens.
