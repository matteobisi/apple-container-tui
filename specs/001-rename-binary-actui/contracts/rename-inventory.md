# Binary Rename Inventory

**Feature**: 001-rename-binary-actui  
**Purpose**: Complete inventory of files requiring changes for apple-tui → actui rename  
**Date**: 2026-02-16

## Summary

This contract documents all files that must be changed to complete the binary rename from "apple-tui" to "actui". Changes are categorized by type and priority.

## Critical Changes (Must Complete)

### 1. Directory Rename

| Current Path | New Path | Action | Notes |
|--------------|----------|--------|-------|
| `cmd/apple-tui/` | `cmd/actui/` | `git mv` | Use git mv to preserve history |

**Impact**: Changes binary package name, affects build commands

### 2. Build Configuration

| File | Changes Required | Line References |
|------|-----------------|----------------|
| `.gitignore` | Update binary name and path entries | Lines 23-24 |
| `Dockerfile` | No changes needed | Dockerfile is for unrelated Python project |

### 3. Documentation Files

| File | Type | References | Priority |
|------|------|------------|----------|
| `README.md` | User docs | Lines 73, 81, 87, 101-102, 106, 120 | P1 - Critical |
| `docs/user-guide.md` | User guide | Lines 28, 34, 40, 54, 175-176, 180, 194 | P1 - Critical |

**Changes include**:
- Build commands: `go build -o actui ./cmd/actui`
- Execution examples: `./actui`, `./actui --dry-run`
- Config paths: `~/.config/actui/config`, `~/Library/Application Support/actui/config`
- Log paths: `~/Library/Application Support/actui/command.log`

### 4. Source Code

| File | Purpose | Changes | Impact |
|------|---------|---------|--------|
| `src/services/config_manager.go` | Config path resolution | Update directory names in paths (lines 25-26, 28) | User config location changes |
| `src/services/log_writer.go` | Command logging | Update log directory name (line 38) | Log file location changes |

**Important**: These changes affect where the application looks for configuration and writes logs. Users will need to migrate manually or start fresh.

### 5. Test Files

| File | Purpose | Changes | Impact |
|------|---------|---------|--------|
| `src/services/services_test.go` | Service tests | Update test paths (lines 264, 305, 317) | Test assertions must match new paths |

**Test validation**:
- Tests verify config path creation
- Tests verify log file creation
- All path assertions must use "actui" not "apple-tui"

## Out of Scope (Do Not Change)

| Path Pattern | Reason |
|--------------|--------|
| `specs/001-apple-container-tui/**` | Historical spec - preserved as documentation |
| `specs/001-rename-binary-actui/**` | Current feature spec - already uses correct terminology |
| Git commit history | Historical references remain unchanged |
| Module name `container-tui` in go.mod | Explicitly excluded per requirements |

## Change Matrix

### String Replacements

| Context | Find | Replace | Notes |
|---------|------|---------|-------|
| Binary name | `apple-tui` | `actui` | Build commands, .gitignore |
| Directory | `cmd/apple-tui` | `cmd/actui` | Build paths, documentation |
| Config path | `.config/apple-tui` | `.config/actui` | Code and docs |
| Config path | `Application Support/apple-tui` | `Application Support/actui` | Code and docs |
| Launch command | `./apple-tui` | `./actui` | All documentation |

### Files Explicitly Requiring Manual Review

1. **Dockerfile** (if exists) - Check for:
   - Build output name
   - COPY commands
   - Binary execution commands

2. **README.md** - Verify all example commands work after changes

3. **docs/user-guide.md** - Verify all workflows reference correct binary

4. **Config path code** - Ensure backward compatibility not assumed

## Verification Checklist

After implementation, verify:

- [ ] `go build ./cmd/actui` produces `actui` binary
- [ ] No "apple-tui" string in: .gitignore (active entries), README.md, docs/user-guide.md
- [ ] No "apple-tui" in src/ directory files
- [ ] No directory named `cmd/apple-tui/`
- [ ] Config manager creates `~/.config/actui/` or `~/Library/Application Support/actui/`
- [ ] Log writer creates `~/Library/Application Support/actui/command.log`
- [ ] All tests pass with new paths
- [ ] `grep -r "apple-tui" .` returns only historical refs in specs/001-apple-container-tui/

## Migration Notes

**User Impact**:
- Existing config files at `~/Library/Application Support/apple-tui/config` will NOT be automatically migrated
- Existing logs at `~/Library/Application Support/apple-tui/command.log` will NOT be automatically migrated
- Users must manually copy config if desired, or start with fresh config
- Old `apple-tui` binaries must be manually removed from PATH

**Documentation Required**:
- Release notes should mention the rename
- README should not reference old name except in "Migration" section if added
- User guide should use only new name

## Summary Statistics

| Category | Count |
|----------|-------|
| Files requiring changes | 8 (minimum, pending Dockerfile verification) |
| Directories to rename | 1 |
| Code files | 3 |
| Test files | 1 |
| Documentation files | 2 |
| Config files | 1 (.gitignore) |
| Total string occurrences (active files) | ~30+ |

**Estimated Complexity**: LOW - Mechanical search-and-replace with path updates

**Risk Level**: LOW - No functional changes, only naming
