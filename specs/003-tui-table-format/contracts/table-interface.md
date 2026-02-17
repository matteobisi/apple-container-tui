# API Contracts: Enhanced TUI Display with Table Layout

**Feature**: 003-tui-table-format  
**Date**: 2026-02-17

## Overview

This feature is a pure UI enhancement with no external API contracts. All "contracts" are internal Go interfaces for the table rendering component. This document defines the programmatic interface that screen components use to render tables.

---

## Table Renderer Interface

### Package
`container-tui/src/ui`

### Public Types

#### TableColumn
```go
type TableColumn struct {
    Header   string  // Column title
    MinWidth int     // Minimum width in characters
    Priority int     // Layout priority (1=high, 3=low)
    Align    string  // "left", "center", "right"
}
```

**Usage**: Define columns when creating a table

**Validation**:
- Priority must be 1-3
- Align must be "left", "center", or "right"
- MinWidth must be >= 0

---

#### TableRow
```go
type TableRow struct {
    Cells    []string
    Selected bool
    Data     interface{} // Reference to source entity
}
```

**Usage**: Represent one row of data in the table

**Validation**:
- len(Cells) must match number of columns in table
- nil Data is allowed (for empty states)

---

#### Table
```go
type Table struct {
    columns []TableColumn
    rows    []TableRow
}
```

**Usage**: Encapsulates table structure and rendering logic

---

### Public Functions

#### NewTable
```go
func NewTable(columns []TableColumn) *Table
```

**Purpose**: Create a new table with column definitions

**Parameters**:
- `columns`: Slice of TableColumn defining table structure

**Returns**: Pointer to Table instance

**Preconditions**:
- len(columns) > 0
- All columns have valid Priority (1-3) and Align values

**Example**:
```go
table := NewTable([]TableColumn{
    {Header: "Name", MinWidth: 10, Priority: 1, Align: "left"},
    {Header: "State", MinWidth: 8, Priority: 2, Align: "center"},
    {Header: "Base Image", MinWidth: 15, Priority: 3, Align: "left"},
})
```

---

#### SetRows
```go
func (t *Table) SetRows(rows []TableRow)
```

**Purpose**: Update table with new row data

**Parameters**:
- `rows`: Slice of TableRow (can be empty)

**Effects**: Replaces all existing rows

**Validation**: All rows must have len(Cells) == len(t.columns)

**Example**:
```go
table.SetRows([]TableRow{
    {Cells: []string{"ubuntu-test", "running", "ubuntu:latest"}, Selected: true},
    {Cells: []string{"nginx-web", "stopped", "nginx:alpine"}, Selected: false},
})
```

---

#### Render
```go
func (t *Table) Render(width int, cursor int) string
```

**Purpose**: Generate formatted table output string

**Parameters**:
- `width`: Available terminal width in characters
- `cursor`: Index of selected row (0-based, -1 for none)

**Returns**: Multi-line string with ANSI escape codes for styling

**Behavior**:
1. Calculate column widths based on available width and priority
2. Render header row with bold styling
3. Render separator line below header
4. Render data rows, applying inverse video to selected row
5. Apply truncation with ellipsis (...) to cells exceeding column width

**Example Output** (simplified, actual includes ANSI codes):
```
Name                 | State   | Base Image              
─────────────────────────────────────────────────────────
ubuntu-test          | running | ubuntu:latest           
nginx-web            | stopped | nginx:alpine            
```

**Edge Cases**:
- width < 80: Truncate low-priority columns, maintain minimum readability
- len(rows) == 0: Render headers with "No items found" centered below
- cursor out of bounds: Treat as -1 (no selection)

---

### Helper Functions

#### TruncateDigest
```go
func TruncateDigest(digest string) string
```

**Purpose**: Truncate SHA256 digest to 12 characters following Docker convention

**Parameters**:
- `digest`: Full digest string (may include "sha256:" prefix)

**Returns**: First 12 hex characters of digest

**Example**:
```go
TruncateDigest("sha256:341bf0f3ce6c5277d6002cf6...")
// Returns: "341bf0f3ce6c"

TruncateDigest("341bf0f3ce6c5277d6002cf6...")
// Returns: "341bf0f3ce6c"
```

---

#### TruncateWithEllipsis
```go
func TruncateWithEllipsis(text string, maxWidth int) string
```

**Purpose**: Truncate text to fit width, adding ellipsis if truncated

**Parameters**:
- `text`: String to truncate
- `maxWidth`: Maximum width in characters

**Returns**: Truncated string with "..." at end if truncated

**Behavior**:
- If len(text) <= maxWidth: return text unchanged
- If maxWidth < 3: return text[:maxWidth] (no ellipsis)
- Otherwise: return text[:maxWidth-3] + "..."

**Example**:
```go
TruncateWithEllipsis("very-long-repository-name", 15)
// Returns: "very-long-re..."
```

---

## Integration Contract

### ContainerListScreen Integration

**Modified Method**: `View() string`

**Before**: Manual string formatting with list items
**After**: Use Table component for rendering

**Implementation**:
```go
func (m ContainerListScreen) View() string {
    if !m.hasLoaded {
        return "Loading containers..."
    }
    
    table := NewTable([]TableColumn{
        {Header: "Name", MinWidth: 10, Priority: 1, Align: "left"},
        {Header: "State", MinWidth: 8, Priority: 2, Align: "center"},
        {Header: "Base Image", MinWidth: 15, Priority: 3, Align: "left"},
    })
    
    rows := make([]TableRow, len(m.containers))
    for i, c := range m.containers {
        rows[i] = TableRow{
            Cells: []string{c.Name, c.Status, c.Image},
            Selected: i == m.cursor,
            Data: c,
        }
    }
    table.SetRows(rows)
    
    return table.Render(m.width, m.cursor) + "\n" + 
           renderSeparator(m.width) + "\n" +
           renderStatusBar(m)
}
```

**Contract Guarantees**:
- Width passed to table is from latest WindowSizeMsg
- Cursor index matches m.cursor (navigation state)
- Row order matches m.containers slice order

---

### ImageListScreen Integration

**Modified Method**: `View() string`

**Before**: Manual string formatting with list items
**After**: Use Table component for rendering

**Implementation**:
```go
func (m ImageListScreen) View() string {
    if !m.hasLoaded {
        return "Loading images..."
    }
    
    table := NewTable([]TableColumn{
        {Header: "Name", MinWidth: 10, Priority: 1, Align: "left"},
        {Header: "Tag", MinWidth: 10, Priority: 2, Align: "center"},
        {Header: "Digest", MinWidth: 12, Priority: 3, Align: "left"},
    })
    
    rows := make([]TableRow, len(m.images))
    for i, img := range m.images {
        rows[i] = TableRow{
            Cells: []string{img.Name, img.Tag, TruncateDigest(img.Digest)},
            Selected: i == m.cursor,
            Data: img,
        }
    }
    table.SetRows(rows)
    
    return table.Render(m.width, m.cursor) + "\n" + 
           renderSeparator(m.width) + "\n" +
           renderStatusBar(m)
}
```

**Contract Guarantees**:
- Digest truncated to 12 characters before passing to table
- Width passed to table is from latest WindowSizeMsg
- Cursor index matches m.cursor (navigation state)

---

## Alternate Screen Buffer Contract

### Entry Point Modification

**File**: `cmd/actui/main.go`

**Modified Function**: `func main()`

**Before**:
```go
p := tea.NewProgram(model)
if err := p.Start(); err != nil {
    fmt.Printf("Error: %v\n", err)
    os.Exit(1)
}
```

**After**:
```go
p := tea.NewProgram(model, tea.WithAltScreen())
if err := p.Start(); err != nil {
    fmt.Printf("Error: %v\n", err)
    os.Exit(1)
}
```

**Contract**:
- `tea.WithAltScreen()` enables alternate screen buffer mode
- Bubbletea handles terminal state save/restore on program start/exit
- All screen rendering happens in alternate buffer (isolated from scrollback)
- Terminal returns to previous state on program exit (normal or crash)

**Observable Behavior**:
- Launch: Screen clears, TUI takes over full terminal
- Exit: Previous terminal content restored, TUI output not in scrollback
- Resize: TUI adapts, alternate buffer maintained

---

## Testing Contracts

### Unit Test Requirements

**File**: `tests/unit/table_test.go`

**Test Cases**:

1. **TestTableColumnWidthCalculation**
   - Input: Table with 3 columns, 5 rows, width=100
   - Expected: Columns sized according to priority, total width <= 100, no truncation

2. **TestTableTruncationLowPriority**
   - Input: Table with 3 columns, width=60 (constrained)
   - Expected: Priority=3 column truncated with ellipsis, Priority=1 preserved

3. **TestTableEmptyState**
   - Input: Table with 0 rows
   - Expected: Headers rendered, "No items found" message, no panic

4. **TestTableInverseVideoSelection**
   - Input: Table with cursor=1
   - Expected: Row 1 includes ANSI reverse video codes, other rows normal

5. **TestTruncateDigest**
   - Input: "sha256:341bf0f3ce6c5277d6002cf6..."
   - Expected: "341bf0f3ce6c"

6. **TestTruncateDigestNoPrefix**
   - Input: "341bf0f3ce6c5277d6002cf6..."
   - Expected: "341bf0f3ce6c"

### Integration Test Requirements

**File**: `tests/integration/table_display_test.go`

**Test Cases**:

1. **TestContainerTableRendering**
   - Setup: Mock 3 containers
   - Action: Render ContainerListScreen
   - Assert: Output contains table headers, 3 data rows, separator

2. **TestImageTableRendering**
   - Setup: Mock 3 images with long names
   - Action: Render ImageListScreen with width=80
   - Assert: Long names truncated, table maintains alignment

3. **TestNavigationPreservesTableFormat**
   - Setup: ContainerListScreen with 5 containers
   - Action: Press down arrow 3 times
   - Assert: Each render maintains table structure, selection moves correctly

---

## Backward Compatibility

**Breaking Changes**: None

**Existing Behavior Preserved**:
- Keyboard navigation (arrow keys, enter, etc.) unchanged
- Command execution flow unchanged
- Modal dialogs (preview, confirm) unchanged
- All existing tests remain valid

**New Behavior**:
- Lists now render as tables (visual change only)
- Alternate screen mode enabled (screen management change)
- Row highlighting uses inverse video (visual change only)

**Migration**: No migration needed - purely additive UI enhancement
