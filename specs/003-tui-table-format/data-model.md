# Data Model: Enhanced TUI Display with Table Layout

**Feature**: 003-tui-table-format  
**Date**: 2026-02-17

## Overview

This feature is primarily a display enhancement and does not introduce new domain entities or persistent data structures. It defines rendering data structures for table visualization of existing container and image models.

## Display Entities

### TableColumn

**Purpose**: Defines a column in a table with rendering metadata

**Fields**:
- `Header` (string): Column title displayed in header row
- `Width` (int): Calculated width for this column in characters
- `MinWidth` (int): Minimum width required (for fixed-width columns like State/Tag)
- `Priority` (int): Layout priority (1=highest, truncated last; 3=lowest, truncated first)
- `Align` (string): Text alignment ("left", "right", "center")

**Relationships**: 
- One TableColumn per table column (3 columns for container table, 3 for image table)

**Validation Rules**:
- Width must be >= MinWidth
- Priority must be 1-3
- Align must be one of: "left", "right", "center"

---

### TableRow

**Purpose**: Represents one row of table data with cell values

**Fields**:
- `Cells` ([]string): Cell values in column order
- `Selected` (bool): Whether this row is currently selected
- `RawData` (interface{}): Reference to source entity (Container or Image) for event handling

**Relationships**:
- Multiple TableRows belong to one Table
- Each Cell corresponds to one TableColumn by index

**Validation Rules**:
- len(Cells) must equal number of columns in table
- Empty cells rendered as "-" placeholder

---

### Table

**Purpose**: Complete table structure with columns, rows, and rendering state

**Fields**:
- `Columns` ([]TableColumn): Column definitions
- `Rows` ([]TableRow): Data rows
- `Width` (int): Total available width for table rendering
- `Cursor` (int): Index of selected row (0-based)

**Relationships**:
- One Table contains multiple Columns and Rows
- Table owned by ContainerListScreen or ImageListScreen

**State Transitions**:
1. **Initialize**: Create Table with column definitions
2. **Populate**: Add rows from container/image data
3. **Render**: Calculate column widths, apply styles, generate output string
4. **Update**: Cursor changes on navigation, rows replaced on data refresh

**Validation Rules**:
- At least one column required
- Rows can be empty (show "No items found")
- Cursor must be within valid row range (0 to len(Rows)-1)

---

## Rendering Data Flow

```
Container/Image Models  →  TableRow Conversion  →  Table Structure  →  Rendered String
      (existing)              (new mapping)          (new entity)        (View output)

Example for Container:
models.Container{
  ID: "17007fa5...",
  Name: "ubuntu-test",
  Status: "running",
  Image: "docker.io/library/ubuntu:latest"
}
          ↓
TableRow{
  Cells: ["ubuntu-test", "running", "docker.io/library/ubuntu:latest"],
  Selected: true,
  RawData: &container
}
          ↓
Table.Render(width=120)
          ↓
"| ubuntu-test              | running | docker.io/library/ubuntu:latest    |"
```

---

## Column Width Calculation Algorithm

**Goal**: Distribute available width across columns based on priority and content

**Inputs**:
- Terminal width (from WindowSizeMsg)
- Column definitions with MinWidth and Priority
- Row data for content measurement

**Algorithm**:

```
1. Calculate content width for each column:
   contentWidth[i] = max(len(column[i].Header), max(len(row[j].Cells[i]) for all rows))

2. Allocate fixed-width columns (Priority=2, State/Tag):
   allocated[i] = max(contentWidth[i], column[i].MinWidth) + 2 (padding)
   remainingWidth -= allocated[i]

3. Allocate high-priority column (Priority=1, Name):
   allocated[name] = min(contentWidth[name] + 2, remainingWidth * 0.4)
   remainingWidth -= allocated[name]

4. Allocate low-priority column (Priority=3, Image/Digest):
   allocated[image] = remainingWidth

5. If any allocated[i] < MinWidth, truncate lower priority columns:
   - Reduce Priority=3 columns first
   - Then reduce Priority=1 if necessary
   - Priority=2 never truncated below MinWidth
   - Apply ellipsis (...) to truncated values
```

**Edge Cases**:
- Width < 80 chars: Truncate Image/Digest aggressively, preserve Name at minimum 20 chars, State/Tag at content width
- Empty table (no rows): Headers still rendered with minimum column widths
- Very long names (>50 chars): Truncate with ellipsis to maintain table structure

---

## Styling Specifications

### Header Row
- **Bold**: Enabled via Lipgloss `.Bold(true)`
- **Separator**: Unicode `─` (U+2500) repeated to table width
- **Padding**: 1 space on each side of text

### Data Rows
- **Normal**: Default terminal foreground/background
- **Selected**: Inverse video via Lipgloss `.Reverse(true)`
- **Padding**: 1 space on each side of text
- **Alignment**: Left-aligned for Name/Image, Center for State/Tag

### Visual Separator (Menu Divider)
- **Character**: Unicode `─` (U+2500) or `_` (fallback)
- **Width**: Full terminal width
- **Position**: Between last data row and keyboard shortcuts menu

---

## Data Mapping

### Container to TableRow

```go
func containerToTableRow(c models.Container, selected bool) TableRow {
    return TableRow{
        Cells: []string{
            c.Name,           // Column 1: Name
            c.Status,         // Column 2: State (from Apple Container STATE field)
            c.Image,          // Column 3: Base Image
        },
        Selected: selected,
        RawData: c,
    }
}
```

### Image to TableRow

```go
func imageToTableRow(img models.Image, selected bool) TableRow {
    digest := TruncateDigest(img.Digest) // sha256: prefix removed, first 12 chars
    return TableRow{
        Cells: []string{
            img.Name,         // Column 1: Name (repository)
            img.Tag,          // Column 2: Tag
            digest,           // Column 3: Digest (truncated)
        },
        Selected: selected,
        RawData: img,
    }
}
```

---

## Non-Functional Considerations

### Performance
- Table rendering must complete in <50ms for up to 100 rows
- Column width calculation is O(n*m) where n=rows, m=columns (acceptable for typical sizes)
- No caching of rendered output (stateless rendering on each frame)

### Memory
- TableRow stores minimal data (3 strings + bool + pointer)
- Full table for 100 containers ~10KB memory overhead (negligible)

### Accessibility
- Screen readers: Table structure not accessible in terminal context (limitation of medium)
- Keyboard navigation: Fully accessible (arrow keys, existing patterns maintained)
- Color blindness: Inverse video works independent of color (no color-based differentiation)

---

## Unchanged Entities

The following existing models are used but not modified:

- **models.Container**: Source data for container table (fields: ID, Name, Status, Image, etc.)
- **models.Image**: Source data for image table (fields: ID, Name, Tag, Digest, etc.)
- **ContainerListScreen**: Owns container table, handles navigation and commands
- **ImageListScreen**: Owns image table, handles navigation and commands

These remain as-is; only their `View()` rendering methods change to use table component.
