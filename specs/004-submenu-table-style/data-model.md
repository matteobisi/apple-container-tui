# Data Model: Consistent Table Styling for Submenus

**Feature**: 004-submenu-table-style  
**Date**: 2026-02-17

## Overview

This feature is purely a visual enhancement with no new data structures or domain entities. It defines **display patterns and styling conventions** for existing submenu and help screens.

## Display Patterns

### MenuSection

**Purpose**: Represents a logical section in a submenu or help screen with visual styling

**Display Structure**:
- `HeaderText` (string): Section label text (e.g., "Container Details", "Available Actions")
- `HeaderStyle` (Lipgloss style): Bold styling applied to header
- `SeparatorLine` (string): Horizontal line of ─ characters spanning terminal width
- `ContentLines` ([]string): Text lines within this section

**Rendering Rules**:
1. Header rendered with bold style: `lipgloss.NewStyle().Bold(true).Render(HeaderText)`
2. Content lines rendered with normal style
3. Separator line appended after all content: `strings.Repeat("─", width)`

**Example Rendering**:
```
Container Details           ← Bold header
name: ubuntu-test          ← Content line 1
status: running            ← Content line 2
image: ubuntu:latest       ← Content line 3
──────────────────────     ← Separator line
```

---

### MenuItem

**Purpose**: Represents one actionable item in a submenu with selection state

**Display Structure**:
- `Text` (string): Action description (e.g., "Stop container", "Inspect image")
- `Selected` (bool): Whether this item is currently selected by cursor
- `SelectionStyle` (Lipgloss style): Inverse video style for selected state
- `NormalStyle` (Lipgloss style): Default style for unselected state

**Rendering Rules**:
1. If Selected == true: Apply `lipgloss.NewStyle().Reverse(true)`
2. If Selected == false: Apply `lipgloss.NewStyle()` (default)
3. Render full line with selected style for complete row highlighting

**State Transitions**:
```
[Created] → Selected = false (normal style)
    ↓
[User presses arrow key] → Cursor moves
    ↓
[This item == cursor position] → Selected = true (inverse video)
    ↓
[User presses arrow key] → Cursor moves away
    ↓
[This item != cursor position] → Selected = false (normal style)
```

**Example Rendering**:
```
Normal item               ← NormalStyle (default colors)
Selected item             ← SelectionStyle (inverse video - white on black)
Another normal item       ← NormalStyle
```

---

### SubmenuLayout

**Purpose**: Complete visual structure for container or image submenu screens

**Structure**:
```
[Screen Title]                    ← RenderTitle() (already bold)
  (blank line)
[Details Section]                 ← MenuSection with HeaderText = "Container Details" or "Image Details"
  - Info fields (name, status, etc.)
  - Horizontal separator
  (blank line)
[Actions Section]                 ← MenuSection with HeaderText = "Available Actions"
  - MenuItem[] (action list)
  - Horizontal separator
  (blank line)
[Keyboard Shortcuts]              ← RenderMuted() (already muted)
```

**Width Responsiveness**:
- Terminal width stored in screen model: `m.width int`
- Updated via WindowSizeMsg handler
- Separators dynamically sized: `strings.Repeat("─", m.width)`
- Fallback: If width == 0, use 80 characters

---

### HelpScreenLayout

**Purpose**: Organized help reference with multiple functional sections

**Structure**:
```
Help                                   ← RenderTitle()
  (blank line)
[Navigation Section]                   ← MenuSection with HeaderText = "Navigation"
  - Key bindings for navigation
  - Horizontal separator
  (blank line)
[Container Actions Section]            ← MenuSection with HeaderText = "Container Actions"
  - Key bindings for container operations
  - Horizontal separator
  (blank line)
[Image Actions Section]                ← MenuSection with HeaderText = "Image Actions"
  - Key bindings for image operations
  - Horizontal separator
  (blank line)
[General Section]                      ← MenuSection with HeaderText = "General"
  - Key bindings for app-wide actions
  - Horizontal separator
  (blank line)
[Footer]                               ← RenderMuted()
```

**Section Contents**:
- **Navigation**: up/down, enter, esc
- **Container Actions**: s (start), t (stop), d (delete), enter (submenu)
- **Image Actions**: p (pull), b (build), n (prune), enter (submenu)
- **General**: r (refresh), m (manage daemon), ? (help), q (quit)

---

## Rendering Data Flow

```
Screen Model (with width) → View() method → String Builder
                                ↓
                        Section Rendering:
                        1. Bold header
                        2. Content lines
                        3. Horizontal separator
                                ↓
                        Menu Item Rendering:
                        1. Check cursor position
                        2. Apply appropriate style
                        3. Full-row rendering
                                ↓
                        Final String Output → Terminal Display
```

---

## Styling Constants

**From Feature 003** (already established):
- Bold Header: `lipgloss.NewStyle().Bold(true)`
- Inverse Video Selection: `lipgloss.NewStyle().Reverse(true)`
- Normal Style: `lipgloss.NewStyle()`
- Separator Character: `─` (U+2500 Box Drawing Light Horizontal)

**Existing Helpers** (reused):
- `RenderTitle(text string) string` - Bold blue title
- `RenderMuted(text string) string` - Gray muted text
- `RenderError(text string) string` - Red error text
- `RenderAccent(text string) string` - Accent color (to be phased out for selection)

---

## Implementation Impact

### Screens Requiring Updates

| Screen | File | Changes Required |
|--------|------|------------------|
| Container Submenu | `src/ui/container_submenu.go` | Add width field, bold headers, separators, inverse video selection |
| Image Submenu | `src/ui/image_submenu.go` | Add width field, bold headers, separators, inverse video selection |
| Help Screen | `src/ui/help.go` | Add width field (if not present), bold section headers, separators |
| Daemon Control | `src/ui/daemon_control.go` | Apply same pattern for consistency |

### No New Data Structures

This feature does NOT introduce:
- New domain models (Container, Image remain unchanged)
- New service layers
- New command builders
- New configuration options

All changes are **View() method modifications** using existing Lipgloss styling capabilities.

---

## Validation Rules

- **Header Text**: Must not be empty, should be concise (< 30 chars)
- **Separator Width**: Must match terminal width if known, default to 80 if unknown
- **Menu Items**: At least one item required, maximum reasonable count ~10 items
- **Selection State**: Exactly one item selected at a time (enforced by cursor position)

---

## Example: Container Submenu Before/After

### Before (Current):
```
Container: ubuntu-test (running)

> Stop container
  Tail container log
  Enter container
  Back
```

### After (This Feature):
```
Container Submenu

Container Details
name: ubuntu-test
status: running
image: ubuntu:latest
────────────────────────────────

Available Actions
Stop container              ← Inverse video when selected
Tail container log
Enter container
Back
────────────────────────────────

Keys: up/down=navigate, enter=select, esc=back
```

**Visual Improvements**:
1. Bold "Container Details" and "Available Actions" headers → Clear section identification
2. Horizontal separators → Visual section boundaries
3. Inverse video selection → Stronger selection feedback (no cursor prefix)
4. Structured layout → Improved scannability
