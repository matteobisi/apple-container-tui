# API Contracts: Consistent Table Styling for Submenus

**Feature**: 004-submenu-table-style  
**Date**: 2026-02-17

## Overview

This feature has no external API contracts. All "contracts" are internal **styling patterns** and **View() method conventions** that submenu and help screens must follow for visual consistency.

---

## Styling Pattern Contracts

### Package
`container-tui/src/ui`

### Bold Header Pattern

**Purpose**: Render section headers with bold styling

**Pattern**:
```go
import "github.com/charmbracelet/lipgloss"

headerStyle := lipgloss.NewStyle().Bold(true)
builder.WriteString(headerStyle.Render("Section Header") + "\n")
```

**Usage Locations**:
- Container submenu: "Container Details", "Available Actions"
- Image submenu: "Image Details", "Available Actions"
- Help screen: "Navigation", "Container Actions", "Image Actions", "General"
- Daemon control: "Daemon Status", "Available Actions"

**Contract**:
- Header text MUST be concise (< 30 characters)
- Header MUST be followed by content or blank line (not directly by separator)
- Header MUST use `.Bold(true)` style, no additional color styling

---

### Horizontal Separator Pattern

**Purpose**: Render visual divider between sections

**Pattern**:
```go
import "strings"

// In screen struct:
type SomeScreen struct {
    // ... other fields
    width int  // Terminal width from WindowSizeMsg
}

// In View() method:
separatorWidth := m.width
if separatorWidth == 0 {
    separatorWidth = 80  // Fallback
}
builder.WriteString(strings.Repeat("─", separatorWidth) + "\n")
```

**Usage Locations**:
- After info section in container/image submenus
- After action list in container/image submenus
- Between help screen sections

**Contract**:
- Separator MUST use `─` character (U+2500)
- Separator MUST span full terminal width (stored in model)
- Separator MUST have blank line before it (after content)
- Separator MUST have blank line after it (before next section)

---

### Inverse Video Selection Pattern

**Purpose**: Highlight selected menu item with full-row inverse video

**Pattern**:
```go
import "github.com/charmbracelet/lipgloss"

normalStyle := lipgloss.NewStyle()
selectedStyle := lipgloss.NewStyle().Reverse(true)

for i, item := range items {
    style := normalStyle
    if i == m.cursor {
        style = selectedStyle
    }
    builder.WriteString(style.Render(item) + "\n")
}
```

**Usage Locations**:
- Container submenu action list
- Image submenu action list
- Daemon control action list
- File picker items (if applicable)

**Contract**:
- Selection MUST use `.Reverse(true)` style (no manual color swapping)
- Full item text MUST be rendered with selected style (not just part of line)
- NO cursor prefix (`>`) should be used (replaced by inverse video)
- Only ONE item can be selected at a time (enforced by cursor position)

---

### WindowSizeMsg Handler Pattern

**Purpose**: Track terminal width for dynamic separator rendering

**Pattern**:
```go
// In screen struct:
type SomeScreen struct {
    // ... other fields
    width int
}

// In Update() method:
func (m SomeScreen) Update(msg tea.Msg) (SomeScreen, tea.Cmd) {
    switch message := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = message.Width
        return m, nil
    // ... other cases
    }
    // ... rest of update logic
}
```

**Usage Locations**:
- Container submenu (ADD width field and handler)
- Image submenu (ADD width field and handler)
- Help screen (ADD width field if not present)
- Daemon control (ADD width field if not present)

**Contract**:
- Width field MUST be `int` type named `width`
- WindowSizeMsg MUST be first case in Update() switch (for consistency)
- Width MUST be stored in model, not recalculated on each render
- View() method MUST use stored width for separators

---

## View() Method Structure Contract

### Submenu View Structure

**Required Pattern**:
```go
func (m SomeSubmenuScreen) View() string {
    var builder strings.Builder
    
    // 1. Title (existing pattern)
    builder.WriteString(RenderTitle("Screen Title") + "\n\n")
    
    // 2. Details Section
    headerStyle := lipgloss.NewStyle().Bold(true)
    builder.WriteString(headerStyle.Render("Details Header") + "\n")
    // ... info fields
    builder.WriteString("\n")
    
    // 3. Separator
    width := m.width
    if width == 0 {
        width = 80
    }
    builder.WriteString(strings.Repeat("─", width) + "\n\n")
    
    // 4. Actions Section
    builder.WriteString(headerStyle.Render("Available Actions") + "\n")
    // ... action items with inverse video
    builder.WriteString("\n")
    
    // 5. Separator
    builder.WriteString(strings.Repeat("─", width) + "\n\n")
    
    // 6. Keyboard shortcuts (existing pattern)
    builder.WriteString(RenderMuted("Keys: ...") + "\n")
    
    return builder.String()
}
```

**Contract Requirements**:
- MUST have blank line after title
- Details header MUST be bold
- MUST have blank line before separator
- Separator MUST be followed by blank line
- Actions header MUST be bold
- Action items MUST use inverse video for selection (NO cursor prefix)
- MUST have blank line before final separator
- Keyboard shortcuts MUST use RenderMuted()

---

### Help Screen View Structure

**Required Pattern**:
```go
func (m HelpScreen) View() string {
    var builder strings.Builder
    
    builder.WriteString(RenderTitle("Help") + "\n\n")
    
    headerStyle := lipgloss.NewStyle().Bold(true)
    width := m.width
    if width == 0 {
        width = 80
    }
    
    // For each section:
    builder.WriteString(headerStyle.Render("Section Name") + "\n")
    // ... key bindings
    builder.WriteString("\n")
    builder.WriteString(strings.Repeat("─", width) + "\n\n")
    
    // Footer
    builder.WriteString(RenderMuted("Press any key to return") + "\n")
    
    return builder.String()
}
```

**Contract Requirements**:
- MUST have at least 3 sections (Navigation, Actions, General)
- Each section header MUST be bold
- Each section MUST end with separator
- Blank line MUST separate sections

---

## Style Consistency Checklist

When implementing this feature, verify:

- [ ] All section headers use `lipgloss.NewStyle().Bold(true)`
- [ ] All separators use `─` character (not `-`, `=`, or other)
- [ ] All separators span full terminal width (`m.width`)
- [ ] All menu selections use `.Reverse(true)` (not color accent)
- [ ] NO cursor prefix (`>`) appears before menu items
- [ ] All screens have `width int` field in struct
- [ ] All screens handle `tea.WindowSizeMsg` to update width
- [ ] Blank lines properly separate sections (before/after separators)
- [ ] Keyboard shortcuts use existing `RenderMuted()` helper

---

## Integration with Existing Code

### Existing Helpers (Reused)

These functions already exist in `src/ui/theme.go` and should continue to be used:

```go
func RenderTitle(text string) string
    // Returns: Bold blue text
    // Use for: Screen titles

func RenderMuted(text string) string
    // Returns: Gray text
    // Use for: Keyboard shortcuts, help text

func RenderError(text string) string
    // Returns: Red text
    // Use for: Error messages

func RenderAccent(text string) string
    // Returns: Accent color text
    // DEPRECATED for selection: Use inverse video instead
    // Still valid for: Emphasis in non-selection contexts
```

### New Patterns (This Feature)

```go
// Bold section headers (NEW)
headerStyle := lipgloss.NewStyle().Bold(true)
headerStyle.Render("Header Text")

// Horizontal separators (NEW)
strings.Repeat("─", width)

// Inverse video selection (NEW)
selectedStyle := lipgloss.NewStyle().Reverse(true)
selectedStyle.Render(itemText)
```

---

## Testing Expectations

### Visual Verification

- Launch app and navigate to each submenu - verify bold headers visible
- Verify horizontal separators span full terminal width
- Use arrow keys to navigate - verify inverse video selection works
- Resize terminal - verify separators adjust to new width
- Press '?' for help - verify sections are separated and headers are bold

### Regression Prevention

- Verify existing functionality unchanged (navigation, command execution, confirmations)
- Verify dry-run mode still works
- Verify command preview modals still display correctly
- Verify type-to-confirm destructive actions still work

---

## Example Implementation Reference

For complete implementation example, see:
- Feature 003: `src/ui/table.go` (for bold and reverse styles)
- Feature 003: `src/ui/container_list.go` (for width handling and WindowSizeMsg)

The patterns in this feature apply the **same styling techniques** to **different screen types** (submenus instead of lists).
