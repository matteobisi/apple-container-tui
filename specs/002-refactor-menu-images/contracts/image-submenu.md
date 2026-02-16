# Contract: Image Submenu

**Feature**: 002-refactor-menu-images  
**Component**: Image context menu for inspection and deletion operations  
**Date**: 2026-02-16

## Purpose

Define the user interactions and command contracts for the image submenu, which provides detailed operations for a selected image (inspect metadata, delete from repository).

---

## User Interface Contract

### Entry Point
- **Trigger**: User presses `Enter` on an image in the image list view
- **Precondition**: At least one image exists in the list
- **Effect**: Displays image submenu with action options

### Menu Display

```
Image: [image-name]:[tag]

> Inspect image
  Delete image
  Back

Keys: up/down=navigate, enter=select, esc=back
```

### Navigation
- **Arrow keys (up/down)**: Move highlight between menu options
- **Enter**: Execute highlighted action
- **Esc**: Return to image list (equivalent to selecting "Back")

### Menu Options

| Option | Visible When | Action | Maps To (Spec) |
|--------|-------------|--------|----------------|
| Inspect image | Always | Display formatted JSON metadata | FR-027, FR-028, FR-029 |
| Delete image | Always | Remove image from local registry | FR-030, FR-031 |
| Back | Always | Return to image list | FR-032 |

---

## Command Contracts

### 1. Inspect Image

**User Action**: Select "Inspect image"

**Command Generated**:
```bash
container image inspect <imageReference> | jq
```

**Parameters**:
- `<imageReference>`: Full image reference constructed as:
  - If tag exists: `NAME:TAG` (e.g., "ubuntu:latest", "ghcr.io/apple/container-builder-shim/builder:0.7.0")
  - If no tag (digest only): `NAME@DIGEST`
  - Constructed from selected Image entity (name, tag, digest fields from data-model.md)

**Expected Output**:
- Formatted JSON object with image metadata (from `jq`)
- Typical structure:
```json
{
  "Id": "sha256:32c70b3752ac28fd6f47c019...",
  "RepoTags": ["ubuntu:latest"],
  "RepoDigests": ["ubuntu@sha256:1a2b3c4d..."],
  "Created": "2026-01-15T10:30:45.123456789Z",
  "Size": 72800000,
  "VirtualSize": 72800000,
  "Config": {
    "Env": ["PATH=/usr/local/sbin:..."],
    "Cmd": ["/bin/bash"],
    ...
  },
  "Architecture": "arm64",
  "Os": "linux",
  ...
}
```

**Display Behavior**:
- Show JSON in scrollable viewport (using bubbles viewport component)
- Enable scrolling: arrow keys, Page Up/Down, Home/End
- Display syntax-highlighted if possible (optional enhancement)
- Show scroll position indicator (e.g., "Line 15/243")
- Esc key returns to image submenu (FR-029)

**Error Scenarios**:
- Image not found (deleted by external process): Display error, return to image submenu  
- `jq` not installed: Display raw JSON without formatting (graceful degradation)
- Inspect command fails: Display error message in viewport

**Success Criteria**:
- JSON fully displayed and scrollable (FR-028)
- Esc returns to submenu (FR-029)
- Large JSON (>1000 lines) scrolls smoothly
- If image deleted while viewing: error shown when Esc pressed (Clarification: edge case)

---

### 2. Delete Image

**User Action**: Select "Delete image"

**Confirmation Flow**:
1. Display type-to-confirm prompt: "Type 'delete' to remove image [image-name]:[tag]"
2. User must type exactly "delete" to confirm (matching container delete pattern)
3. On confirmation: Execute delete command
4. Show brief confirmation or error message
5. Return to image list view (not submenu) (FR-031 + spec scenario 4)
6. Image list automatically refreshes

**Command Generated** (after type-to-confirm):
```bash
container image rm <imageReference>
```

**Parameters**:
- `<imageReference>`: Full image reference constructed as:
  - If tag exists: `NAME:TAG` (e.g., "ubuntu:latest", "ghcr.io/owner/repo:v1.0")
  - If no tag (digest only): `NAME@DIGEST` (e.g., "ubuntu@sha256:abc123...")
  - Constructed from selected Image entity (name, tag, digest fields from data-model.md)

**Expected Output**:
```
Deleted: sha256:32c70b3752ac28fd6f47c019...
```
or
```
Untagged: ubuntu:latest
```

**Error Scenarios**:
- Image in use by container: Display error "Image is in use by container [container-name]", return to image submenu (FR-035)
- Image not found: Display error, return to image submenu
- Permission denied: Display error, return to image submenu
- User cancels type-to-confirm: Return to image submenu without deleting

**Success Criteria**:
- Type-to-confirm safety guard active (FR-030)
- Error for in-use image clearly states which container(s) using it (FR-035)
- On success: Returns to image list, list refreshed to exclude deleted image
- On error: Returns to image submenu with error message

---

### 3. Back

**User Action**: Select "Back" or press Esc

**Behavior**:
- Navigates back to image list view
- No command executed
- No confirmation required

---

## State Management

**Navigation State Changes**:
```
Initial: ImageList view
User presses Enter on image
→ NavigationState.selectedImage = selected image
→ NavigationState.currentView = ImageSubmenu
→ Push ImageList to navigationStack

User selects "Inspect image"
→ NavigationState.currentView = ImageInspect
→ Push ImageSubmenu to navigationStack
→ Execute `container image inspect <image> | jq`

User presses Esc in ImageInspect
→ Pop from navigationStack (return to ImageSubmenu)
→ NavigationState.currentView = ImageSubmenu

User selects "Delete image"
→ Show type-to-confirm
→ On confirmation: Execute delete
→ Pop to ImageList (skip ImageSubmenu)
→ NavigationState.currentView = ImageList
→ NavigationState.selectedImage = null
→ Refresh image list

User presses Esc in ImageSubmenu
→ Pop from navigationStack (return to ImageList)
→ NavigationState.currentView = ImageList
→ NavigationState.selectedImage = null
```

---

## Image Inspect View Details

### Viewport Configuration
- **Initial scroll**: Top of JSON (line 0)
- **Width**: Terminal width - 4 (for margins)
- **Height**: Terminal height - 6 (for header, help text)
- **Scroll keys**:
  - Up/Down: Move 1 line
  - Page Up/Down: Move 10 lines
  - Home: Jump to top
  - End: Jump to bottom
  - Esc: Return to submenu

### Header Display
```
Image Inspection: [image-name]:[tag]

[JSON content in viewport]

Keys: up/down/pgup/pgdn=scroll, home/end=jump, esc=back | Line 15/243
```

### Performance
- Large JSON (>10,000 lines): Should still scroll smoothly
- Parsing happens once (when command executes)
- Viewport only renders visible portion

---

## Testing Requirements

### Unit Tests
- [ ] Image submenu option generation (always 3 options: Inspect, Delete, Back)
- [ ] Type-to-confirm validation for delete ("delete" keyword exact match)

### Contract Tests
- [ ] `container image inspect <name> | jq` command generated correctly with full image reference
- [ ] `container image rm <name>` command generated correctly with full image reference
- [ ] Delete command only executed after type-to-confirm matches

### Integration Tests
- [ ] Full flow: Enter submenu → Select "Inspect" → View JSON → Scroll → Esc → Return to submenu
- [ ] Full flow: Enter submenu → Select "Delete" → Type confirm → Delete success → Return to list → List refreshed
- [ ] Error flow: Select "Delete" → Image in use → Error displayed → Stay in submenu
- [ ] Error flow: Select "Inspect" → Image not found (deleted externally) → Error displayed
- [ ] Escape flow: Enter submenu → Select "Inspect" → Press Esc → Return to submenu → Press Esc → Return to list
- [ ] Large JSON (simulated >1000 lines) scrolls smoothly in inspect view

---

## Acceptance Criteria (from Spec)

Links to User Story 3 acceptance scenarios:

1. ✅ **Given** image list with image selected, **When** press Enter, **Then** submenu shows "Inspect image", "Delete image", "Back"
2. ✅ **Given** image submenu, **When** select "Inspect image", **Then** see formatted JSON with scrolling and Esc to exit
3. ✅ **Given** image submenu viewing inspection, **When** press Esc, **Then** return to image submenu
4. ✅ **Given** image submenu, **When** select "Delete image", **Then** type-to-confirm prompt, then image removed, return to image list
5. ✅ **Given** image submenu, **When** select "Back" or Esc, **Then** return to image list view

---

## Error Message Catalog

| Scenario | Error Message | Action |
|----------|---------------|--------|
| Image in use | "Cannot delete image: in use by container(s) [names]" | Stay in submenu |
| Image not found | "Image not found (may have been deleted externally)" | Return to submenu |
| Inspect failed | "Failed to inspect image: [error details]" | Display in viewport |
| jq not available | [Display raw JSON without formatting] | Viewport still works |
| Delete permission denied | "Permission denied: cannot delete image" | Stay in submenu |

**Contract Complete**: 2026-02-16
