# Feature Specification: Enhanced TUI Display with Table Layout

**Feature Branch**: `003-tui-table-format`  
**Created**: February 17, 2026  
**Status**: Draft  
**Input**: User description: "Improve the TUI making the command more readable and nice to see. First: i want that the tool is launched in a clear screen (like the btop command) and when close goes back to the previous terminal. Second, i want to increase the readability and clearness of the data for containers and images views using table formatting with columns and visual separators."

## Clarifications

### Session 2026-02-17

- Q: Digest Display Format - The image table includes a "Digest" column which can be quite long (SHA256 hashes are 64 characters). What is the standard truncation length for digest display? → A: 12 characters (matching Docker CLI standard: sha256: prefix removed, first 12 hex chars shown)
- Q: Container State Display Values - The container table shows a "State" column. What state values should be displayed? → A: Use STATE values from Apple Container's `container list` command output (e.g., running, stopped, etc.)
- Q: Column Width Priority - When terminal width is constrained and not all columns can display full content, which column should be prioritized? → A: Name column gets priority, then State/Tag (fixed width), then truncate Image/Digest
- Q: Row Selection Highlight Style - How should the selected row be highlighted when navigating? → A: Inverse video (swap foreground/background colors)
- Q: Table Header Visual Distinction - How should table headers be visually distinguished from data rows? → A: Bold text with separator line below

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Clean Screen Experience (Priority: P1)

When users launch the TUI tool, they want an uncluttered, isolated viewing experience similar to full-screen terminal applications like btop, htop, or vim. When they exit, they expect to return to their previous terminal state without the TUI output cluttering their terminal history.

**Why this priority**: This addresses user experience fundamentals - launching and exiting the application cleanly. Without this, users have a poor first impression and cluttered terminal history that interferes with other work. This is the foundation for professional TUI behavior.

**Independent Test**: Can be fully tested by launching the tool, verifying the screen clears and shows only TUI content, then exiting and confirming the previous terminal content is restored without TUI remnants in the scrollback. Delivers immediate professional UX improvement.

**Acceptance Scenarios**:

1. **Given** user has terminal open with command history visible, **When** user launches actui, **Then** screen clears completely and displays only the TUI interface without previous terminal content
2. **Given** TUI is running in full-screen mode, **When** user quits the application (via 'q' key), **Then** terminal returns to the exact previous state with command history visible and TUI output not added to scrollback
3. **Given** TUI is running, **When** user's terminal is resized, **Then** TUI adapts to new dimensions while maintaining clear screen mode

---

### User Story 2 - Tabular Container View (Priority: P2)

Users viewing the container list want to quickly scan and understand container information with clearly organized columns showing name, state, and base image. The data should be easy to read with aligned columns and clear headers, separated visually from the action menu.

**Why this priority**: This directly addresses the readability issue for the primary use case (viewing containers). Users currently struggle to parse the unstructured list format. Table formatting makes scanning and finding specific containers much faster.

**Independent Test**: Can be tested by launching the tool in container view and verifying that containers are displayed in a table with column headers (Name, State, Base Image), proper alignment, and a visual separator between the table and the command menu. Delivers immediate readability improvement for container management.

**Acceptance Scenarios**:

1. **Given** user has multiple containers, **When** viewing the container list, **Then** containers are displayed in a table with headers "Name", "State", and "Base Image" with columns properly aligned
2. **Given** container list table is displayed, **When** viewing the interface, **Then** a visual separator line appears between the container data and the keyboard shortcuts menu at the bottom
3. **Given** containers with varying name lengths, **When** displayed in table format, **Then** columns automatically adjust width to maintain alignment and readability
4. **Given** user selects a container, **When** navigating with arrow keys, **Then** the entire row is highlighted to show selection clearly

---

### User Story 3 - Tabular Image View (Priority: P3)

Users viewing the image list want to quickly scan and understand image information with clearly organized columns showing name, tag, and digest. The data should be easy to read with aligned columns and clear headers, separated visually from the action menu.

**Why this priority**: This extends the table formatting improvement to the images view. While important, it's lower priority than containers because users typically interact with containers more frequently than images.

**Independent Test**: Can be tested by navigating to the images view and verifying that images are displayed in a table with column headers (Name, Tag, Digest), proper alignment, and a visual separator between the table and the command menu. Delivers consistent table formatting across all list views.

**Acceptance Scenarios**:

1. **Given** user has multiple images, **When** viewing the image list, **Then** images are displayed in a table with headers "Name", "Tag", and "Digest" with columns properly aligned
2. **Given** image list table is displayed, **When** viewing the interface, **Then** a visual separator line appears between the image data and the keyboard shortcuts menu at the bottom
3. **Given** images with long repository names, **When** displayed in table format, **Then** long names are truncated with ellipsis (...) while maintaining digest visibility
4. **Given** user selects an image, **When** navigating with arrow keys, **Then** the entire row is highlighted to show selection clearly

---

### Edge Cases

- What happens when terminal width is too narrow to display all table columns? (Columns should wrap or truncate gracefully with ellipsis, prioritizing critical information)
- What happens when there are no containers or images to display? (Show table headers with "No items found" message in a centered row)
- What happens when a container name or image repository name is extremely long? (Truncate with ellipsis to prevent table misalignment)
- What happens if the user resizes the terminal while viewing a table? (Re-render the table with adjusted column widths to fit new dimensions)
- What happens if terminal doesn't support alternate screen buffer? (Fall back gracefully to standard output without clear screen, maintaining table formatting)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST launch in full-screen mode using alternate screen buffer to isolate TUI display from terminal scrollback
- **FR-002**: System MUST restore previous terminal state when user exits, without leaving TUI output in scrollback history
- **FR-003**: Container list MUST display data in table format with three columns: Name, State, and Base Image
- **FR-004**: Image list MUST display data in table format with three columns: Name, Tag, and Digest
- **FR-005**: Both container and image tables MUST include column header rows with clearly labeled column names
- **FR-006**: System MUST display a visual separator (horizontal line) between the data table and the keyboard shortcuts menu
- **FR-007**: Table columns MUST be automatically sized to accommodate content while maintaining alignment; when space is constrained, Name column gets display priority, State/Tag columns use fixed width, and Base Image/Digest columns are truncated first
- **FR-008**: System MUST highlight the currently selected row using inverse video (swapped foreground/background colors) when user navigates with arrow keys
- **FR-009**: Long text values in table cells MUST be truncated with ellipsis (...) when they exceed available column width
- **FR-010**: System MUST handle terminal resize events and re-render tables to fit new dimensions
- **FR-011**: Tables MUST maintain readable formatting even when terminal width is constrained (minimum viable width: 80 characters)
- **FR-012**: Image digest values MUST be displayed truncated to 12 characters (sha256: prefix removed, first 12 hexadecimal characters shown) following Docker CLI convention

### Key Entities

- **Container Display Row**: Represents one container entry with three data points: container name/ID, current state (as returned by Apple Container's STATE field), and base image reference. Must maintain visual alignment across all rows.

- **Image Display Row**: Represents one image entry with three data points: repository name, tag identifier, and digest hash (truncated to 12 characters matching Docker CLI standard). Must maintain visual alignment across all rows.

- **Table Header**: Column labels that identify the meaning of each data column. Must be displayed in bold text with a horizontal separator line below to distinguish from data rows.

- **Visual Separator**: Horizontal line element that divides the data display area from the command menu area. Must span the full width of the terminal.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can launch the tool and see only TUI content without any previous terminal history visible on screen
- **SC-002**: Users can exit the tool and return to their previous terminal state with command history intact and no TUI scrollback visible
- **SC-003**: Container data is displayed in aligned columns with headers, making individual container attributes scannable in under 2 seconds per row
- **SC-004**: Image data is displayed in aligned columns with headers, making individual image attributes scannable in under 2 seconds per row
- **SC-005**: Visual separation between data and menu is clearly visible, improving interface structure recognition
- **SC-006**: Tables remain readable and properly formatted when terminal is resized to any width above 80 characters
- **SC-007**: Users can identify selected container or image instantly through row highlighting

## Assumptions

- Terminal emulator supports ANSI escape sequences for alternate screen buffer (standard for modern terminals)
- Minimum terminal width of 80 characters is acceptable for table display (industry standard)
- Horizontal line characters (─) or similar Unicode box-drawing characters are available in the terminal font
- Users understand standard TUI navigation patterns (arrow keys, q to quit)
- Container names and image repository names don't contain special characters that would break table alignment (or will be sanitized)
