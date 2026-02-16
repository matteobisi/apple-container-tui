# Quickstart: Binary Rename Implementation

**Feature**: 001-rename-binary-actui  
**Purpose**: Step-by-step guide for implementing and verifying the binary rename  
**Date**: 2026-02-16

## Overview

This guide provides the implementation steps for renaming the binary from "apple-tui" to "actui", including build commands, testing procedures, and verification steps.

## Prerequisites

Before starting:
- [ ] Current directory: Repository root (`container-tui/`)
- [ ] Git branch: `001-rename-binary-actui` (already checked out)
- [ ] Go version: 1.21 or higher installed
- [ ] No uncommitted changes that conflict with rename

## Implementation Steps

### Phase 1: Directory Rename

Use `git mv` to preserve file history:

```bash
git mv cmd/apple-tui cmd/actui
```

**Verify**:
```bash
ls cmd/
# Should show: actui/
```

### Phase 2: Update .gitignore

Edit `.gitignore` to update binary name references:

**Find** (lines 23-24):
```
apple-tui
/cmd/apple-tui/apple-tui
```

**Replace with**:
```
actui
/cmd/actui/actui
```

### Phase 3: Update Source Code

#### File: `src/services/config_manager.go`

Update config paths (~line 25-28):

**Find**:
```go
configPaths := []string{
    filepath.Join(home, ".config", "apple-tui", "config"),
    filepath.Join(home, "Library", "Application Support", "apple-tui", "config"),
}
writePath := filepath.Join(home, "Library", "Application Support", "apple-tui", "config")
```

**Replace with**:
```go
configPaths := []string{
    filepath.Join(home, ".config", "actui", "config"),
    filepath.Join(home, "Library", "Application Support", "actui", "config"),
}
writePath := filepath.Join(home, "Library", "Application Support", "actui", "config")
```

#### File: `src/services/log_writer.go`

Update log path (~line 38):

**Find**:
```go
logPath := filepath.Join(home, "Library", "Application Support", "apple-tui", "command.log")
```

**Replace with**:
```go
logPath := filepath.Join(home, "Library", "Application Support", "actui", "command.log")
```

### Phase 4: Update Tests

#### File: `src/services/services_test.go`

Update test assertions:

**Line ~264**:
```go
configPath := filepath.Join(home, ".config", "actui", "config")
```

**Line ~305**:
```go
if _, err := os.Stat(filepath.Join(home, "Library", "Application Support", "actui", "command.log")); err != nil {
```

**Line ~317**:
```go
logPath := filepath.Join(home, "Library", "Application Support", "actui", "command.log")
```

### Phase 5: Update Documentation

#### File: `README.md`

Update all occurrences (lines 73, 81, 87, 101-102, 106, 120):

| Find | Replace |
|------|---------|
| `go build -o apple-tui ./cmd/apple-tui` | `go build -o actui ./cmd/actui` |
| `./apple-tui` | `./actui` |
| `./apple-tui --dry-run` | `./actui --dry-run` |
| `~/.config/apple-tui/config` | `~/.config/actui/config` |
| `~/Library/Application Support/apple-tui/config` | `~/Library/Application Support/actui/config` |
| `~/Library/Application Support/apple-tui/command.log` | `~/Library/Application Support/actui/command.log` |

#### File: `docs/user-guide.md`

Apply same replacements as README.md (lines 28, 34, 40, 54, 175-176, 180, 194)

## Build and Test

### Build the Binary

```bash
go build -o actui ./cmd/actui
```

**Expected output**: Binary file `actui` in current directory

**Verify**:
```bash
ls -lh actui
file actui
# Should show: actui: Mach-O 64-bit executable arm64
```

### Run the Binary

```bash
./actui --help
```

**Expected**: Help output with available commands

```bash
./actui --dry-run
```

**Expected**: TUI launches in dry-run mode

### Run Tests

```bash
go test ./...
```

**Expected**: All tests pass

**If tests fail**:
- Check that config paths in tests use "actui"
- Check that log paths in tests use "actui"
- Verify test fixtures updated

### Verify Config and Log Paths

**Test config creation**:
```bash
# Remove any existing config first
rm -rf ~/Library/Application\ Support/actui
rm -rf ~/.config/actui

# Run binary (should create default config)
./actui &
# Let it initialize, then quit (press 'q')

# Check config was created
ls ~/Library/Application\ Support/actui/config
# Should exist

# Check log file
ls ~/Library/Application\ Support/actui/command.log
# Should exist (may be empty if no commands run)
```

## Verification Checklist

After completing implementation:

- [ ] **Build Success**: `go build -o actui ./cmd/actui` completes without errors
- [ ] **Binary Name**: Output file is named `actui` (not `apple-tui`)
- [ ] **Directory**: `cmd/actui/` exists, `cmd/apple-tui/` does not
- [ ] **Tests Pass**: `go test ./...` shows all tests passing
- [ ] **Config Path**: Binary creates config at `~/Library/Application Support/actui/config`
- [ ] **Log Path**: Binary creates logs at `~/Library/Application Support/actui/command.log`
- [ ] **Documentation**: README shows `./actui` in all examples
- [ ] **No Old References**: `grep -r "apple-tui" src/ docs/ README.md .gitignore` returns no matches
- [ ] **Git Status**: All changed files staged for commit

## String Search Verification

Run comprehensive search to find any remaining old references:

```bash
# Search active code (should return 0 matches in src/)
grep -r "apple-tui" src/

# Search documentation (should return 0 matches)
grep -r "apple-tui" README.md docs/

# Search config files (should return 0 matches in .gitignore)
grep "apple-tui" .gitignore

# Historical specs are OK to still reference old name
# (specs/001-apple-container-tui/ is preserved as documentation)
```

**Success criteria**: Zero matches in:
- `src/` directory
- `README.md`
- `docs/` directory  
- `.gitignore`

## Commit Message

Suggested commit message after verification:

```
Rename binary from apple-tui to actui

- Rename cmd/apple-tui/ → cmd/actui/
- Update config paths to use actui directory
- Update log paths to use actui directory
- Update all documentation (README, user guide)
- Update .gitignore entries
- Update test assertions for new paths

Breaking change: Config location changes from
~/Library/Application Support/apple-tui/ to
~/Library/Application Support/actui/

Users must manually migrate config or start fresh.

Closes #001-rename-binary-actui
```

## Troubleshooting

### Build fails with "package not found"

**Issue**: Import paths may be incorrect after directory rename

**Solution**: 
```bash
go mod tidy
go clean -cache
```

### Tests fail with path not found

**Issue**: Test assertions still reference old paths

**Solution**: Double-check all path strings in test files use "actui"

### Binary still creates apple-tui directories

**Issue**: Hardcoded strings not updated in source

**Solution**: Search for "apple-tui" in all `.go` files:
```bash
find src/ -name "*.go" -exec grep -l "apple-tui" {} \;
```

## Next Steps

After successful verification:
1. Commit all changes
2. Push branch to remote
3. Run any CI/CD pipelines
4. Update release notes with migration instructions
5. Consider adding migration guide for users

## Reference

- Original spec: [spec.md](../spec.md)
- Change inventory: [contracts/rename-inventory.md](../contracts/rename-inventory.md)
- Research notes: [research.md](../research.md)
