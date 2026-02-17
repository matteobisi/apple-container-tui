package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TableColumn defines a column in a table with rendering metadata
type TableColumn struct {
	Header   string // Column title displayed in header row
	MinWidth int    // Minimum width required for this column
	Priority int    // Layout priority (1=highest, truncated last; 3=lowest, truncated first)
	Align    string // Text alignment: "left", "right", "center"
}

// TableRow represents one row of table data with cell values
type TableRow struct {
	Cells    []string    // Cell values in column order
	Selected bool        // Whether this row is currently selected
	Data     interface{} // Reference to source entity for event handling
}

// Table encapsulates table structure and rendering logic
type Table struct {
	columns []TableColumn
	rows    []TableRow
}

// NewTable creates a new table with column definitions
func NewTable(columns []TableColumn) *Table {
	return &Table{
		columns: columns,
		rows:    []TableRow{},
	}
}

// SetRows updates the table with new data rows
func (t *Table) SetRows(rows []TableRow) {
	t.rows = rows
}

// Render generates the formatted table string for display
func (t *Table) Render(width int, cursor int) string {
	if len(t.columns) == 0 {
		return ""
	}

	// Calculate column widths based on available space
	columnWidths := t.calculateColumnWidths(width)

	var sb strings.Builder

	// Render header row with bold styling
	headerStyle := lipgloss.NewStyle().Bold(true)
	headerParts := make([]string, len(t.columns))
	for i, col := range t.columns {
		headerParts[i] = t.padOrTruncate(col.Header, columnWidths[i], col.Align)
	}
	sb.WriteString(headerStyle.Render(strings.Join(headerParts, " ")))
	sb.WriteString("\n")

	// Render separator line below header
	sb.WriteString(strings.Repeat("─", width))
	sb.WriteString("\n")

	// Render data rows
	if len(t.rows) == 0 {
		// Empty state
		emptyMsg := "No items found"
		padding := (width - len(emptyMsg)) / 2
		if padding < 0 {
			padding = 0
		}
		sb.WriteString(strings.Repeat(" ", padding))
		sb.WriteString(emptyMsg)
		sb.WriteString("\n")
	} else {
		normalStyle := lipgloss.NewStyle()
		selectedStyle := lipgloss.NewStyle().Reverse(true)

		for i, row := range t.rows {
			// Determine if this row is selected based on cursor position
			isSelected := i == cursor || row.Selected

			// Build row string
			cellParts := make([]string, len(t.columns))
			for j, cell := range row.Cells {
				if j < len(t.columns) {
					cellParts[j] = t.padOrTruncate(cell, columnWidths[j], t.columns[j].Align)
				}
			}
			rowStr := strings.Join(cellParts, " ")

			// Apply styling based on selection
			if isSelected {
				sb.WriteString(selectedStyle.Render(rowStr))
			} else {
				sb.WriteString(normalStyle.Render(rowStr))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// calculateColumnWidths determines the width for each column based on priority and available space
func (t *Table) calculateColumnWidths(totalWidth int) []int {
	if len(t.columns) == 0 {
		return []int{}
	}

	widths := make([]int, len(t.columns))

	// Account for spacing between columns (1 space between each column)
	spacingWidth := len(t.columns) - 1
	availableWidth := totalWidth - spacingWidth

	if availableWidth < 0 {
		availableWidth = 0
	}

	// Start with minimum widths
	remainingWidth := availableWidth
	for i, col := range t.columns {
		widths[i] = col.MinWidth
		remainingWidth -= col.MinWidth
	}

	if remainingWidth <= 0 {
		return widths
	}

	// Measure actual content width for each column
	contentWidths := make([]int, len(t.columns))
	for i, col := range t.columns {
		maxWidth := len(col.Header)
		for _, row := range t.rows {
			if i < len(row.Cells) {
				cellWidth := len(row.Cells[i])
				if cellWidth > maxWidth {
					maxWidth = cellWidth
				}
			}
		}
		contentWidths[i] = maxWidth
	}

	// Allocate remaining width based on priority (1=high, 3=low)
	// Priority 1 columns get expanded first, then 2, then 3
	for priority := 1; priority <= 3; priority++ {
		if remainingWidth <= 0 {
			break
		}

		// Find columns with this priority
		priorityCols := []int{}
		totalNeeded := 0
		for i, col := range t.columns {
			if col.Priority == priority {
				needed := contentWidths[i] - widths[i]
				if needed > 0 {
					priorityCols = append(priorityCols, i)
					totalNeeded += needed
				}
			}
		}

		if len(priorityCols) == 0 {
			continue
		}

		// Distribute available width among priority columns
		if totalNeeded <= remainingWidth {
			// All priority columns can get their full content width
			for _, i := range priorityCols {
				widths[i] = contentWidths[i]
			}
			remainingWidth -= totalNeeded
		} else {
			// Distribute proportionally
			for _, i := range priorityCols {
				needed := contentWidths[i] - widths[i]
				allocation := (needed * remainingWidth) / totalNeeded
				widths[i] += allocation
			}
			remainingWidth = 0
		}
	}

	return widths
}

// padOrTruncate formats a cell value to fit the specified width with alignment
func (t *Table) padOrTruncate(value string, width int, align string) string {
	if width <= 0 {
		return ""
	}

	valueLen := len(value)

	// Truncate if too long
	if valueLen > width {
		if width <= 3 {
			return strings.Repeat(".", width)
		}
		return value[:width-3] + "..."
	}

	// Pad if too short
	if valueLen < width {
		padding := width - valueLen
		switch align {
		case "right":
			return strings.Repeat(" ", padding) + value
		case "center":
			leftPad := padding / 2
			rightPad := padding - leftPad
			return strings.Repeat(" ", leftPad) + value + strings.Repeat(" ", rightPad)
		default: // "left"
			return value + strings.Repeat(" ", padding)
		}
	}

	return value
}

// TruncateDigest truncates a digest string to 12 characters, removing the "sha256:" prefix if present
func TruncateDigest(digest string) string {
	// Remove "sha256:" prefix if present
	digest = strings.TrimPrefix(digest, "sha256:")

	// Truncate to 12 characters
	if len(digest) > 12 {
		return digest[:12]
	}

	return digest
}
