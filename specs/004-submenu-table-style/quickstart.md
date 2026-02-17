# Quickstart: Consistent Table Styling for Submenus

**Feature**: 004-submenu-table-style  
**Date**: 2026-02-17  
**For**: Developers implementing visual consistency updates

## Prerequisites

- Go 1.21+ installed
- Feature 003-tui-table-format completed and merged (provides styling patterns)
- Bubbletea v1.2.4 and Lipgloss v1.0.0 in go.mod (no new dependencies)
- macOS 26.x on Apple Silicon with Terminal.app or iTerm2

## Quick Implementation Path (Priority Order)

### Step 1: Container Submenu Updates (P1 - 20 minutes)

**File**: `src/ui/container_submenu.go`

**Change 1**: Add width field to struct

```go
type ContainerSubmenuScreen struct {
    // ... existing fields
    width int  // NEW: Track terminal width
}
```

**Change 2**: Handle WindowSizeMsg

```go
func (m ContainerSubmenuScreen) Update(msg tea.Msg) (ContainerSubmenuScreen, tea.Cmd) {
    switch message := msg.(type) {
    case tea.WindowSizeMsg:  // NEW: Add this case first
        m.width = message.Width
        return m, nil
    // ... existing cases
    }
}
```

**Change 3**: Refactor View() method

```go
func (m ContainerSubmenuScreen) View() string {
    var builder strings.Builder
    
    // Title (existing)
    builder.WriteString(RenderTitle("Container Submenu") + "\n\n")
    
    // NEW: Bold details header
    headerStyle := lipgloss.NewStyle().Bold(true)
    builder.WriteString(headerStyle.Render("Container Details") + "\n")
    
    // Existing: Container info
    builder.WriteString(fmt.Sprintf("name: %s\n", m.container.Name))
    builder.WriteString(fmt.Sprintf("status: %s\n", m.container.Status))
    builder.WriteString(fmt.Sprintf("image: %s\n", m.container.Image))
    builder.WriteString("\n")
    
    // NEW: Horizontal separator
    width := m.width
    if width == 0 {
        width = 80
    }
    builder.WriteString(strings.Repeat("─", width) + "\n\n")
    
    // NEW: Bold actions header
    builder.WriteString(headerStyle.Render("Available Actions") + "\n")
    
    // NEW: Action items with inverse video selection (replace cursor prefix)
    actions := []string{"Stop container", "Tail container log", "Enter container", "Back"}
    normalStyle := lipgloss.NewStyle()
    selectedStyle := lipgloss.NewStyle().Reverse(true)
    
    for i, action := range actions {
        style := normalStyle
        if i == m.cursor {
            style = selectedStyle
        }
        builder.WriteString(style.Render(action) + "\n")
    }
    builder.WriteString("\n")
    
    // NEW: Horizontal separator
    builder.WriteString(strings.Repeat("─", width) + "\n\n")
    
    // Existing: Keyboard shortcuts
    builder.WriteString(RenderMuted("Keys: up/down=navigate, enter=select, esc=back") + "\n")
    
    return builder.String()
}
```

**Test**:
```bash
go run cmd/actui/main.go
# Press Enter on a container
# Verify: Bold headers, separators, inverse video selection (no cursor prefix)
```

---

### Step 2: Image Submenu Updates (P2 - 20 minutes)

**File**: `src/ui/image_submenu.go`

**Apply same pattern as Container Submenu**:

```go
// 1. Add width field to ImageSubmenuScreen struct
// 2. Add WindowSizeMsg handler
// 3. Update View() with:
//    - Bold "Image Details" header
//    - Horizontal separator
//    - Bold "Available Actions" header
//    - Inverse video selection for actions
//    - Horizontal separator
```

**Details Section Example**:
```go
builder.WriteString(headerStyle.Render("Image Details") + "\\n")
builder.WriteString(fmt.Sprintf("repository: %s\\n", m.image.Name))
builder.WriteString(fmt.Sprintf("tag: %s\\n", m.image.Tag))
builder.WriteString(fmt.Sprintf("digest: %s\\n", m.image.Digest))
```

**Actions Example**:
```go
actions := []string{"Inspect image", "Delete image", "Back"}
// ... same inverse video pattern as container submenu
```

**Test**:
```bash
go run cmd/actui/main.go
# Press 'i' to go to images, then Enter on an image
# Verify: Same visual style as container submenu
```

---

### Step 3: Help Screen Updates (P3 - 30 minutes)

**File**: `src/ui/help.go`

**Change 1**: Add width field (if not present)

```go
type HelpScreen struct {
    // ... existing fields
    width int  // NEW: Track terminal width
}
```

**Change 2**: Handle WindowSizeMsg (if not present)

```go
func (m HelpScreen) Update(msg tea.Msg) (HelpScreen, tea.Cmd) {
    switch message := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = message.Width
        return m, nil
    // ... existing cases
    }
}
```

**Change 3**: Reorganize View() with sections

```go
func (m HelpScreen) View() string {
    var builder strings.Builder
    
    builder.WriteString(RenderTitle("Help") + "\n\n")
    
    headerStyle := lipgloss.NewStyle().Bold(true)
    width := m.width
    if width == 0 {
        width = 80
    }
    
    // Section 1: Navigation
    builder.WriteString(headerStyle.Render("Navigation") + "\n")
    builder.WriteString("up/down, j/k       Navigate lists\n")
    builder.WriteString("enter              Open submenu/select\n")
    builder.WriteString("esc                Back/cancel\n")
    builder.WriteString("\n")
    builder.WriteString(strings.Repeat("─", width) + "\n\n")
    
    // Section 2: Container Actions
    builder.WriteString(headerStyle.Render("Container Actions") + "\n")
    builder.WriteString("s                  Start container\n")
    builder.WriteString("t                  Stop container\n")
    builder.WriteString("d                  Delete container\n")
    builder.WriteString("enter              Open container submenu\n")
    builder.WriteString("\n")
    builder.WriteString(strings.Repeat("─", width) + "\n\n")
    
    // Section 3: Image Actions
    builder.WriteString(headerStyle.Render("Image Actions") + "\n")
    builder.WriteString("i                  Switch to images view\n")
    builder.WriteString("p                  Pull image\n")
    builder.WriteString("b                  Build image\n")
    builder.WriteString("n                  Prune images\n")
    builder.WriteString("enter              Open image submenu\n")
    builder.WriteString("\n")
    builder.WriteString(strings.Repeat("─", width) + "\n\n")
    
    // Section 4: General
    builder.WriteString(headerStyle.Render("General") + "\n")
    builder.WriteString("r                  Refresh current view\n")
    builder.WriteString("m                  Manage daemon\n")
    builder.WriteString("?                  Show this help\n")
    builder.WriteString("q                  Quit application\n")
    builder.WriteString("\n")
    builder.WriteString(strings.Repeat("─", width) + "\n\n")
    
    builder.WriteString(RenderMuted("Press any key to return") + "\n")
    
    return builder.String()
}
```

**Test**:
```bash
go run cmd/actui/main.go
# Press '?'
# Verify: 4 sections with bold headers and separators
```

---

## Testing Checklist

### Visual Verification

**Container Submenu**:
- [ ] Bold "Container Details" header visible
- [ ] Horizontal separator after container info
- [ ] Bold "Available Actions" header visible
- [ ] Arrow keys navigate with inverse video selection (no cursor prefix)
- [ ] Horizontal separator after actions
- [ ] Separators span full terminal width

**Image Submenu**:
- [ ] Bold "Image Details" header visible
- [ ] Horizontal separator after image info
- [ ] Bold "Available Actions" header visible
- [ ] Arrow keys navigate with inverse video selection (no cursor prefix)
- [ ] Horizontal separator after actions
- [ ] Separators span full terminal width

**Help Screen**:
- [ ] 4 bold section headers (Navigation, Container Actions, Image Actions, General)
- [ ] Horizontal separators between all sections
- [ ] Content properly organized by function
- [ ] Separators span full terminal width

### Responsive Behavior

- [ ] Resize terminal - separators adjust to new width
- [ ] Narrow terminal (80 chars) - still readable
- [ ] Wide terminal (150+ chars) - separators span full width

### Functional Regression

- [ ] Container submenu actions still work (stop, logs, shell, back)
- [ ] Image submenu actions still work (inspect, delete, back)
- [ ] Help screen dismisses on any key press
- [ ] Navigation between screens unaffected

---

## Common Issues & Solutions

### Issue: Separators Not Full Width

**Symptom**: Separators only 80 characters even in wider terminal

**Solution**: Verify WindowSizeMsg handler is present and correctly updates `m.width`

```go
case tea.WindowSizeMsg:
    m.width = message.Width  // Make sure this line exists
    return m, nil
```

### Issue: Selection Not Visible

**Symptom**: Can't tell which menu item is selected

**Solution**: Ensure `.Reverse(true)` is applied to selected item

```go
selectedStyle := lipgloss.NewStyle().Reverse(true)
if i == m.cursor {
    style = selectedStyle  // Make sure selectedStyle is used
}
```

### Issue: Headers Not Bold

**Symptom**: Headers look same as regular text

**Solution**: Create headerStyle once and reuse, don't forget `.Render()`

```go
headerStyle := lipgloss.NewStyle().Bold(true)
builder.WriteString(headerStyle.Render("Header Text") + "\n")
//                   ^^^ Don't forget .Render()
```

### Issue: Cursor Prefix Still Showing

**Symptom**: Still seeing `>` before selected items

**Solution**: Remove old cursor prefix logic completely

```go
// OLD (remove this):
cursor := " "
if i == m.cursor {
    cursor = ">"
}
line := cursor + " " + action

// NEW (replace with this):
style := normalStyle
if i == m.cursor {
    style = selectedStyle
}
builder.WriteString(style.Render(action) + "\n")
```

---

## Performance Notes

- Styling operations are lightweight (<1ms per screen)
- No caching needed - recalculate on every render
- Inverse video is native ANSI, no performance penalty
- Terminal width updates are event-driven, not polled

---

## Code Patterns Reference

### From Feature 003 (Reuse These)

```go
// Bold styling
headerStyle := lipgloss.NewStyle().Bold(true)

// Inverse video selection
selectedStyle := lipgloss.NewStyle().Reverse(true)

// Horizontal separator
strings.Repeat("─", width)

// Width handling
case tea.WindowSizeMsg:
    m.width = message.Width
```

### Existing UI Helpers (Continue Using)

```go
RenderTitle(text)   // Screen titles (already bold)
RenderMuted(text)   // Keyboard shortcuts
RenderError(text)   // Error messages
```

---

## Validation

After implementing all changes:

1. **Build**: `go build -o actui cmd/actui/main.go`
2. **Run**: `./actui`
3. **Navigate**: Test all updated screens (container submenu, image submenu, help)
4. **Resize**: Verify separators respond to terminal resize
5. **Compare**: Ensure visual consistency with feature 003 table views

**Expected Outcome**: All screens follow same design language - bold headers, horizontal separators, inverse video selection.

---

## Estimated Time

- Container Submenu: 20 minutes
- Image Submenu: 20 minutes
- Help Screen: 30 minutes
- Testing & Validation: 15 minutes

**Total**: ~1.5 hours (similar to feature 003 complexity)

---

## Next Steps After Implementation

1. Manual verification on macOS Terminal.app and iTerm2
2. Update constitution if needed (likely no changes required - display-only)
3. Commit changes with descriptive message
4. Merge to main branch
