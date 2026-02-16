# Data Model: Binary Rename

**Feature**: 001-rename-binary-actui  
**Date**: 2026-02-16

## Overview

This feature involves renaming the binary executable from "apple-tui" to "actui". **No data model changes are required** as this is purely a naming/cosmetic change.

## Impact Analysis

### No New Entities

This feature does not introduce any new data entities or modify existing ones. The rename is limited to:
- Binary output filename
- File system paths (config and log directories)
- String references in code and documentation

### Existing Entities (Unchanged)

The following data structures from the existing codebase remain **unchanged**:

#### 1. Config Structure
- **Location**: `src/models/config.go`
- **Impact**: None - Structure unchanged, only directory name in path changes
- **Changes**: File system path strings updated in `src/services/config_manager.go`

#### 2. Container Models
- **Location**: `src/models/container.go`
- **Impact**: None - Container data structures completely unchanged

#### 3. Command Models
- **Location**: `src/models/command.go`
- **Impact**: None - Command structures completely unchanged

#### 4. Result Models
- **Location**: `src/models/result.go`
- **Impact**: None - Result data structures completely unchanged

### File System Paths

The only "data" changes are file system path strings:

| Current Path | New Path | Type |
|--------------|----------|------|
| `~/.config/apple-tui/config` | `~/.config/actui/config` | Config file |
| `~/Library/Application Support/apple-tui/config` | `~/Library/Application Support/actui/config` | Config file (macOS) |
| `~/Library/Application Support/apple-tui/command.log` | `~/Library/Application Support/actui/command.log` | Log file |

**These are string constant changes, not data model changes.**

## Migration Considerations

### User Data

**Config files**: 
- Existing config at `~/Library/Application Support/apple-tui/config` will NOT be automatically migrated
- New binary will create config at `~/Library/Application Support/actui/config`
- Users must manually copy config if desired

**Log files**:
- Existing logs at `~/Library/Application Support/apple-tui/command.log` will remain
- New binary will create new log at `~/Library/Application Support/actui/command.log`
- Old logs can be manually merged or archived if needed

### No Database Changes

This project uses local file storage only (TOML config files, JSONL logs). No database schema changes required.

## Architecture Impact

### No API Changes

This feature does not modify:
- CLI command structure
- Command flags or arguments
- TUI keyboard bindings
- Apple Container CLI integration
- Error handling or logging formats

The application architecture remains identical except for file paths.

## Summary

**Data Model Changes**: NONE

**Rationale**: This is a binary name and file path rename only. No data structures, APIs, or architectural components are modified. All existing data models in `src/models/` remain completely unchanged.

**Validation**: No data model validation required. Existing model tests continue to work without modification.
