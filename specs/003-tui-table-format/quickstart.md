# Quickstart: Enhanced TUI Display with Table Layout

**Feature**: 003-tui-table-format  
**Date**: 2026-02-17  
**For**: Developers implementing table rendering and alternate screen mode

## Prerequisites

- Go 1.21+ installed
- Existing container-tui project cloned
- Bubbletea v1.2.4 and Lipgloss v1.0.0 already in go.mod (no new dependencies)
- macOS 26.x on Apple Silicon with Terminal.app or iTerm2

## Quick Implementation Path (Priority Order)

### Step 1: Enable Alternate Screen Mode (P1 - 5 minutes)

**File**: `cmd/actui/main.go`

**Change**: Add `tea.WithAltScreen()` option to tea.NewProgram

```go
// Before
p := tea.NewProgram(model)

// After
p := tea.NewProgram(model, tea.WithAltScreen())
```

**Test**: 
```bash
go run cmd/actui/main.go
# Screen should clear completely
# Press 'q' to quit - terminal should restore previous state
```

**Success Criteria**: Screen clears on launch, previous terminal content returns on exit

---

### Step 2: Create Table Component (P2 - 30 minutes)

**File**: `src/ui/table.go` (new file)

**Core Structure**:
```go
package ui

import (
    "strings"
    "github.com/charmbracelet/lipgloss"
)

type TableColumn struct {
    Header   string
    MinWidth int
    Priority int
    Align    string
}

type TableRow struct {
    Cells    []string
    Selected bool
    Data     interface{}
}

type Table struct {
    columns []TableColumn
    rows    []TableRow
}

func NewTable(columns []TableColumn) *Table {
    return &Table{columns: columns}
}

func (t *Table) SetRows(rows []TableRow) {
    t.rows = rows
}

func (t *Table) Render(width int, cursor int) string {
    // TODO: Implement rendering logic
    // 1. Calculate column widths
    // 2. Render header with bold
    // 3. Render separator line
    // 4. Render rows with selection highlighting
    return ""
}
```

**Test**:
```bash
go test ./src/ui -run TestTable
```

---

### Step 3: Integrate Container Table (P2 - 20 minutes)

**File**: `src/ui/container_list.go`

**Modify**: `View()` method

```go
func (m ContainerListScreen) View() string {
    // Existing loading/error handling...
    
    // NEW: Create table
    table := NewTable([]TableColumn{
        {Header: "Name", MinWidth: 10, Priority: 1, Align: "left"},
        {Header: "State", MinWidth: 8, Priority: 2, Align: "center"},
        {Header: "Base Image", MinWidth: 15, Priority: 3, Align: "left"},
    })
    
    // NEW: Populate rows
    rows := make([]TableRow, len(m.containers))
    for i, c := range m.containers {
        rows[i] = TableRow{
            Cells: []string{c.Name, c.Status, c.Image},
            Selected: i == m.cursor,
            Data: c,
        }
    }
    table.SetRows(rows)
    
    // NEW: Render table
    output := table.Render(m.width, m.cursor)
    output += "\n" + strings.Repeat("─", m.width) + "\n"
    
    // Existing status bar rendering...
    return output
}
```

**Test**:
```bash
go run cmd/actui/main.go
# Navigate to container list
# Verify table format with columns and headers
```

---

### Step 4: Integrate Image Table (P3 - 20 minutes)

**File**: `src/ui/image_list.go`

**Modify**: `View()` method (similar to container_list.go)

```go
func (m ImageListScreen) View() string {
    // Existing loading/error handling...
    
    // NEW: Create table
    table := NewTable([]TableColumn{
        {Header: "Name", MinWidth: 10, Priority: 1, Align: "left"},
        {Header: "Tag", MinWidth: 10, Priority: 2, Align: "center"},
        {Header: "Digest", MinWidth: 12, Priority: 3, Align: "left"},
    })
    
    // NEW: Populate rows with truncated digests
    rows := make([]TableRow, len(m.images))
    for i, img := range m.images {
        digest := TruncateDigest(img.Digest) // Helper function
        rows[i] = TableRow{
            Cells: []string{img.Name, img.Tag, digest},
            Selected: i == m.cursor,
            Data: img,
        }
    }
    table.SetRows(rows)
    
    // NEW: Render table
    output := table.Render(m.width, m.cursor)
    output += "\n" + strings.Repeat("─", m.width) + "\n"
    
    // Existing status bar rendering...
    return output
}
```

**Helper Function** (add to `table.go`):
```go
func TruncateDigest(digest string) string {
    clean := strings.TrimPrefix(digest, "sha256:")
    if len(clean) > 12 {
        return clean[:12]
    }
    return clean
}
```

**Test**:
```bash
go run cmd/actui/main.go
# Press 'i' to switch to images
# Verify table format with digests truncated to 12 chars
```

---

## Table Rendering Implementation Guide

### Column Width Algorithm (for table.go)

```go
func (t *Table) calculateColumnWidths(totalWidth int) []int {
    widths := make([]int, len(t.columns))
    remaining := totalWidth
    
    // 1. Measure content widths
    for colIdx, col := range t.columns {
        contentWidth := len(col.Header)
        for _, row := range t.rows {
            if len(row.Cells[colIdx]) > contentWidth {
                contentWidth = len(row.Cells[colIdx])
            }
        }
        widths[colIdx] = contentWidth + 2 // padding
    }
    
    // 2. Allocate by priority
    // Priority 2 (State/Tag) - fixed width
    // Priority 1 (Name) - flexible, gets preference
    // Priority 3 (Image/Digest) - gets remainder, truncated first
    
    // Simplified allocation (full logic in data-model.md)
    for i, col := range t.columns {
        if col.Priority == 2 {
            remaining -= widths[i]
        }
    }
    
    // Distribute remaining to Priority 1 and 3...
    // (See data-model.md for complete algorithm)
    
    return widths
}
```

### Row Rendering with Highlighting

```go
func (t *Table) Render(width int, cursor int) string {
    widths := t.calculateColumnWidths(width)
    var output strings.Builder
    
    // Header row (bold)
    headerStyle := lipgloss.NewStyle().Bold(true)
    headerRow := ""
    for i, col := range t.columns {
        cell := padOrTruncate(col.Header, widths[i])
        headerRow += cell + " "
    }
    output.WriteString(headerStyle.Render(headerRow) + "\n")
    
    // Separator line
    output.WriteString(strings.Repeat("─", width) + "\n")
    
    // Data rows
    normalStyle := lipgloss.NewStyle()
    selectedStyle := lipgloss.NewStyle().Reverse(true)
    
    for i, row := range t.rows {
        style := normalStyle
        if i == cursor {
            style = selectedStyle
        }
        
        rowStr := ""
        for j, cell := range row.Cells {
            cellStr := padOrTruncate(cell, widths[j])
            rowStr += cellStr + " "
        }
        output.WriteString(style.Render(rowStr) + "\n")
    }
    
    return output.String()
}

func padOrTruncate(text string, width int) string {
    if len(text) > width {
        if width < 3 {
            return text[:width]
        }
        return text[:width-3] + "..."
    }
    return text + strings.Repeat(" ", width-len(text))
}
```

---

## Testing Checklist

### Manual Testing

1. **Alternate Screen**:
   - [ ] Launch app - screen clears completely
   - [ ] Quit app - previous terminal content restored
   - [ ] No TUI output in scrollback after quit

2. **Container Table**:
   - [ ] Column headers visible and bold
   - [ ] Separator line below headers
   - [ ] 3 columns aligned: Name, State, Base Image
   - [ ] Arrow keys navigate - selected row inverts colors
   - [ ] Long container names truncate with "..."

3. **Image Table**:
   - [ ] Press 'i' to switch to images view
   - [ ] 3 columns aligned: Name, Tag, Digest
   - [ ] Digests show 12 characters only
   - [ ] Long repository names truncate with "..."

4. **Terminal Resize**:
   - [ ] Resize terminal window
   - [ ] Tables re-render with adjusted column widths
   - [ ] No visual artifacts or misalignment

5. **Edge Cases**:
   - [ ] No containers - shows "No items found"
   - [ ] No images - shows "No items found"
   - [ ] Very narrow terminal (80 chars) - still readable

### Automated Testing

```bash
# Unit tests
go test ./src/ui -v

# Integration tests
go test ./tests/integration -v

# Full test suite
go test ./... -v
```

**Key Tests to Add**:
- `TestTableRender` in `src/ui/table_test.go`
- `TestContainerTableView` in `tests/integration/`
- `TestImageTableView` in `tests/integration/`

---

## Troubleshooting

### Alternate Screen Not Working

**Symptom**: TUI output mixed with terminal history  
**Solution**: Verify `tea.WithAltScreen()` is passed to `tea.NewProgram()`  
**Check**: Terminal supports ANSI alternate screen (all modern macOS terminals do)

### Table Columns Misaligned

**Symptom**: Data doesn't line up under headers  
**Solution**: Ensure padding is consistent (same space count in header and data rows)  
**Debug**: Print column widths to stderr: `fmt.Fprintf(os.Stderr, "widths: %v\n", widths)`

### Selection Highlight Not Visible

**Symptom**: Can't tell which row is selected  
**Solution**: Verify `lipgloss.NewStyle().Reverse(true)` is applied to selected row  
**Check**: Terminal supports ANSI inverse video (standard feature)

### Digest Not Truncating

**Symptom**: Full SHA256 digest shown  
**Solution**: Call `TruncateDigest()` before adding to TableRow.Cells  
**Verify**: Check digest field in Apple Container output includes "sha256:" prefix

---

## Performance Targets

- **Table Render**: <1ms for 100 rows (no caching needed)
- **Terminal Resize**: <50ms to re-render with new widths
- **Row Selection**: <16ms to update highlight (60fps target)

**Profiling** (if needed):
```bash
go test -cpuprofile=cpu.prof -bench=BenchmarkTableRender ./src/ui
go tool pprof cpu.prof
```

---

## Next Steps After Implementation

1. **Manual Verification**: Test on macOS Terminal.app and iTerm2
2. **Update Tests**: Add unit tests for table component
3. **Documentation**: Update README with new table view screenshots
4. **Constitution Re-check**: Verify all principles still satisfied (should be ✅)

**Estimated Total Time**: 1.5-2 hours for full implementation and testing

---

## Reference Materials

- **Bubbletea Alternate Screen**: https://github.com/charmbracelet/bubbletea/tree/master/examples/altscreen
- **Lipgloss Styling**: https://github.com/charmbracelet/lipgloss#quick-start
- **Box Drawing Characters**: U+2500 (horizontal line: ─), U+2502 (vertical line: │)
- **ANSI Reverse Video**: ESC[7m (handled by Lipgloss .Reverse(true))

**Full Design Details**: See [data-model.md](data-model.md) and [contracts/table-interface.md](contracts/table-interface.md)
