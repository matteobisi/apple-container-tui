# Contract: Image List View

**Feature**: 002-refactor-menu-images  
**Component**: Image management screen with list display and quick actions  
**Date**: 2026-02-16

## Purpose

Define the user interactions and command contracts for the image list view, which displays all local container images and provides access to pull, build, prune, and per-image operations.

---

## User Interface Contract

### Entry Point
- **Trigger**: User presses `i` key in the main container list view
- **Precondition**: None (works with empty or populated image list)
- **Effect**: Displays image list view with all local images

### Display Format

**With Images**:
```
Images

NAME                                          TAG                  DIGEST
> ghcr.io/apple/container-builder-shim/builder  0.7.0               sha256:32c70b3752ac28fd...
  markitdown                                     latest              sha256:7f8a9d3e4b2c1a5f...
  docker.io/library/ubuntu                       latest              sha256:1a2b3c4d5e6f7g8h...

Keys: up/down=navigate, enter=submenu, p=pull, b=build, n=image-prune, esc=back to main menu
```

**Empty List**:
```
Images

No images found. Press 'p' to pull an image or 'b' to build from Containerfile.

Keys: p=pull, b=build, esc=back to main menu
```

### Navigation
- **Arrow keys (up/down)**: Move selection between images
- **Enter**: Open image submenu for selected image
- **p**: Trigger image pull workflow (delegates to existing feature)
- **b**: Trigger image build workflow (delegates to existing feature)
- **n**: Trigger image prune operation
- **r**: Refresh image list (re-execute `container image list`)
- **Esc**: Return to main container list view

### Column Display Rules
- **NAME**: Full image repository name, left-aligned
- **TAG**: Image tag, left-aligned
- **DIGEST**: First 24 characters of digest (e.g., "sha256:32c70b3752ac28fd..."), left-aligned
- If content exceeds terminal width: truncate with "..." (FR-037)
- Empty digest: display as "(none)" or empty cell

---

## Command Contracts

### 1. List Images

**Auto-executed**: On view entry and after pull/build/prune operations

**Command Generated**:
```bash
container image list
```

**Expected Output Format**:
```
NAME                                          TAG                  DIGEST
ghcr.io/apple/container-builder-shim/builder  0.7.0                sha256:32c70b3752ac28fd6f47c019e4b234a1...
markitdown                                     latest               sha256:7f8a9d3e4b2c1a5f8e9d0c3b6a7f4e2d...
docker.io/library/ubuntu                       latest               sha256:1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p...
```

**Parsing Logic** (from research.md):
1. Split output by newlines
2. First line contains headers (NAME, TAG, DIGEST, ...)
3. Detect column indices from headers using `strings.Fields()`
4. Parse each subsequent line into Image entity using column indices
5. Handle edge cases: missing digest, long names, empty output

**Error Scenarios**:
- Container daemon not running: Display "Daemon not running. Press 'm' to start daemon."
- Command execution failure: Display error message
- Empty output: Display "No images found" message (not an error)
- Malformed output: Display "Failed to parse image list"

**Success Criteria**:
- Displays all local images in <2 seconds for up to 100 images (SC-006)
- Correctly parses NAME, TAG, DIGEST columns (FR-016)
- Handles empty list gracefully (FR-036)

---

### 2. Pull Image

**User Action**: Press 'p' in image list view

**Command Generated**: Delegates to existing image pull workflow

**Behavior**:
- Navigates to existing pull screen (already implemented)
- User enters image reference (name:tag)
- Shows progress during pull
- On completion: Returns to image list view (not main container list)
- Image list automatically refreshes to show new image (FR-024)

**Error Scenarios**:
- Handled by existing pull workflow
- Network errors, invalid image reference, etc.

**Success Criteria**:
- Consistent with existing pull UX
- Automatic return to image list + refresh (FR-019, FR-024)

---

### 3. Build Image

**User Action**: Press 'b' in image list view

**Command Generated**: Delegates to existing image build workflow

**Behavior**:
- Navigates to existing build screen (already implemented)
- User selects Containerfile/Dockerfile
- User names the image
- Shows progress during build
- On completion: Returns to image list view (not main container list)
- Image list automatically refreshes to show new image (FR-024)

**Error Scenarios**:
- Handled by existing build workflow
- Build failures, missing Containerfile, etc.

**Success Criteria**:
- Consistent with existing build UX
- Automatic return to image list + refresh (FR-020, FR-024)

---

### 4. Prune Images

**User Action**: Press 'n' in image list view

**Command Generated** (after type-to-confirm):
```bash
container image prune
```

**Confirmation Flow**:
1. Display type-to-confirm prompt: "Type 'prune' to remove all unused images"
2. User must type exactly "prune" to confirm
3. On confirmation: Execute prune command
4. Show progress/spinner during execution
5. Display result: "Pruned X images, freed Y MB" (from command output)
6. Automatically refresh image list (FR-024)

**Expected Output**:
```
Deleted Images:
sha256:abc123...
sha256:def456...

Total reclaimed space: 1.2 GB
```

**Error Scenarios**:
- No unused images: Display "No unused images to prune"
- Command fails: Display error message
- Prune interrupted: Display error, refresh list to show current state
- User cancels type-to-confirm: Return to image list without executing

**Success Criteria**:
- Type-to-confirm safety guard active (FR-022 + Clarification Q5)
- Result summary displayed clearly
- List refreshed to reflect removed images (FR-024)

---

### 5. Refresh Image List

**User Action**: Press 'r' in image list view

**Command Generated**:
```bash
container image list
```

**Behavior**:
- Re-executes list command
- Updates displayed images
- Maintains selection cursor position if possible
- Shows brief "Refreshing..." indicator

**Error Scenarios**:
- Same as "List Images" above

**Success Criteria**:
- Quick refresh (<2 seconds)
- Selection preserved when possible

---

## State Management

**Navigation State Changes**:
```
Initial: ContainerList view
User presses 'i'
→ NavigationState.currentView = ImageList
→ Push ContainerList to navigationStack
→ Execute `container image list`

User presses 'p' (pull)
→ NavigationState.currentView = ImagePull
→ Push ImageList to navigationStack
→ [pull workflow]
→ On completion: Pop to ImageList, refresh list

User presses Enter on image
→ NavigationState.selectedImage = selected image
→ NavigationState.currentView = ImageSubmenu
→ Push ImageList to navigationStack

User presses Esc in ImageList
→ Pop from navigationStack (return to ContainerList)
→ NavigationState.currentView = ContainerList
```

---

## Display Rules

### Column Widths (Dynamic)
- Calculate based on terminal width
- Priority order: NAME (flexible), TAG (min 12 chars), DIGEST (24 chars)
- If terminal too narrow: truncate NAME with "..."

### Truncation Example
```
Terminal width: 80 columns
NAME (max): 45 chars
TAG: 12 chars
DIGEST: 24 chars (sha256:... format)
Spacing: ~4 chars

Example:
ghcr.io/apple/container-builder-shim/buil...  0.7.0        sha256:32c70b3752ac28fd...
```

### Empty State
```
Images

No images found. Press 'p' to pull an image or 'b' to build from Containerfile.

Keys: p=pull, b=build, esc=back to main menu
```

---

## Testing Requirements

### Unit Tests
- [ ] Image list parser handles well-formed output
- [ ] Image list parser handles empty output
- [ ] Image list parser handles missing digest column
- [ ] Image list parser handles long names (truncation)
- [ ] Column width calculation for various terminal sizes

### Contract Tests
- [ ] `container image list` command generated correctly
- [ ] `container image prune` command generated correctly
- [ ] Pull workflow returns to image list (not container list)
- [ ] Build workflow returns to image list (not container list)

### Integration Tests
- [ ] Full flow: Press 'i' → View images → Press 'p' → Pull image → Return to image list → List refreshed
- [ ] Full flow: Press 'i' → View images → Press 'b' → Build image → Return to image list → List refreshed
- [ ] Full flow: Press 'i' → View images → Press 'n' → Type confirm → Prune → List refreshed
- [ ] Empty image list displays correct message and key hints
- [ ] Image list with 50+ images displays correctly, scrolls smoothly

---

## Acceptance Criteria (from Spec)

Links to User Story 2 acceptance scenarios:

1. ✅ **Given** main container list, **When** press 'i', **Then** see image list with NAME, TAG, DIGEST columns
2. ✅ **Given** image list, **When** press 'p', **Then** taken to existing image pull workflow
3. ✅ **Given** image list, **When** press 'b', **Then** taken to existing image build workflow
4. ✅ **Given** image list, **When** press 'n', **Then** prompted with type-to-confirm, then prune executes
5. ✅ **Given** image list, **When** press Esc, **Then** return to main container list view
6. ✅ **Given** image list, **When** images pulled/built/pruned, **Then** list automatically refreshes

**Contract Complete**: 2026-02-16
