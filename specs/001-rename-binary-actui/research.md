# Research: Binary Rename from apple-tui to actui

**Feature**: 001-rename-binary-actui  
**Phase**: 0 (Research & Discovery)  
**Date**: 2026-02-16

## Overview

Research findings for renaming the Go binary output from "apple-tui" to "actui". This covers Go build conventions, potential gotchas, and best practices for comprehensive renames.

## Research Areas

### 1. Go Binary Naming Conventions

**Decision**: Use "actui" as the binary name

**Rationale**: 
- Go binaries are typically lowercase, single-word names (e.g., `kubectl`, `docker`, `git`)
- Short names are easier to type and remember
- "actui" clearly indicates "Apple Container TUI" while being concise
- No special characters or spaces align with Unix/POSIX conventions

**Alternatives considered**: 
- "ac-tui" (rejected: hyphen less common in binary names)
- "apple-container-tui" (rejected: too verbose)
- "acontainer" (rejected: less clear about TUI nature)

### 2. Go Module vs Binary Name

**Decision**: Keep module name as `container-tui`, change only binary output name

**Rationale**:
- Module name affects import paths across the codebase
- Changing module name would require updating all internal imports
- Module name and binary name can differ without issues
- User requirement explicitly states "i don't need to update the module name focus on binary"

**Implementation**:
- Binary name is controlled by the directory under `cmd/` (currently `cmd/apple-tui/`)
- Go build uses package name from directory: `go build -o <output> ./cmd/<dirname>`
- Renaming `cmd/apple-tui/` → `cmd/actui/` changes default binary name
- Build commands and Makefile/Dockerfile need explicit `-o actui` or path updates

### 3. Directory Rename Impact

**Decision**: Rename `cmd/apple-tui/` to `cmd/actui/` 

**Rationale**:
- Consistency: directory name should match binary name
- Developer experience: reduces confusion when navigating code
- Build simplification: default package name matches desired output

**Potential Issues**:
- Git history: `git mv` preserves history better than delete+add
- IDE references: most IDEs handle directory renames gracefully
- Import paths: `cmd/` packages are typically not imported, so no import path changes

### 4. Build Configuration Updates

**Decision**: Update all build configurations to reference "actui"

**Files requiring updates**:
- `Dockerfile`: Look for binary name in COPY or RUN commands
- Makefiles (if present): Update binary output names
- CI/CD configs (if present): Update artifact names
- Build scripts: Update any hardcoded binary names

**Best practices**:
- Use variables for binary name in build configs (easier future changes)
- Ensure multi-platform builds (if any) all use consistent naming
- Verify output path doesn't hardcode old name

### 5. Documentation and Test Updates

**Decision**: Comprehensive search-and-replace with manual verification

**Strategy**:
1. **Automated search**: `grep -r "apple-tui" .` to find all references
2. **Categorize findings**:
   - Code comments: Safe to replace
   - Documentation: Must replace (README, user guides, etc.)
   - Test assertions: Must replace if testing binary name
   - Git history/commit messages: Leave unchanged
3. **Manual verification**: Review each change to avoid false positives

**Potential gotchas**:
- String literals in tests that check binary name
- Error messages that reference the binary name
- Example commands in documentation
- Log messages that include binary name

### 6. Testing Strategy

**Decision**: Validate rename through build and test execution

**Test plan**:
1. **Build test**: `go build -o actui ./cmd/actui/` should succeed
2. **Binary verification**: Check output file is named "actui"
3. **Execution test**: `./actui --help` should run without errors
4. **Test suite**: All existing tests should pass unchanged (unless they check binary name)
5. **Documentation check**: Search for remaining "apple-tui" references

**Success criteria** (from spec):
- Build produces `actui` binary
- Zero occurrences of "apple-tui" in active code/docs/tests
- All tests pass

### 7. User Migration Path

**Decision**: No backward compatibility needed (clean break)

**Rationale**:
- Project is proof-of-concept stage
- No package manager distribution (manual binary replacement)
- Users explicitly confirmed clean break is acceptable

**User impact**:
- Existing `apple-tui` binaries must be manually removed
- No automatic migration or symlinks
- README should note the rename in release notes/changelog

## Implementation Checklist

Based on research, the implementation requires:

- [ ] Rename directory: `cmd/apple-tui/` → `cmd/actui/`
- [ ] Update Dockerfile (if binary name referenced)
- [ ] Update README.md (all examples and installation instructions)
- [ ] Update docs/user-guide.md (if exists)
- [ ] Search and replace "apple-tui" → "actui" in:
  - [ ] All code comments
  - [ ] All documentation files
  - [ ] All test files
  - [ ] All configuration files
- [ ] Verify build produces `actui` binary
- [ ] Run full test suite
- [ ] Update any CI/CD configurations

## Open Questions

**NONE** - All clarifications obtained during spec phase:
- ✅ Package manager distribution: No
- ✅ Backward compatibility: No
- ✅ Directory rename: Yes
- ✅ Root-level files: Remove/rename
- ✅ Test updates: Yes

## References

- Go Command Documentation: https://golang.org/cmd/go/
- Go Project Layout: https://github.com/golang-standards/project-layout
- Git best practices for file renames: `git mv` preserves history
