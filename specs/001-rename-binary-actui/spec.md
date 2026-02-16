# Feature Specification: Rename Binary from apple-tui to actui

**Feature Branch**: `001-rename-binary-actui`  
**Created**: February 16, 2026  
**Status**: Draft  
**Input**: User description: "currently this repository is set to generate a go binary named apple-tui i want to change it to actui, making the according changes to the code and also to the README. and everything else needed. i don't need to update the module name focus on binary"

## Clarifications

### Session 2026-02-16

- Q: Is this project distributed through package managers (like Homebrew, apt, yum, etc.)? → A: No, not distributed via package managers - only direct binary distribution
- Q: Should the rename maintain backward compatibility for users who have "apple-tui" installed? → A: No backward compatibility needed - clean break to new name
- Q: Should the source directory `cmd/apple-tui/` be renamed to `cmd/actui/` to match the new binary name? → A: Yes, rename directory to cmd/actui/ for consistency
- Q: Should the root-level file "apple-tui" (if it exists as a binary or script) be renamed or removed? → A: Remove or rename to actui
- Q: Should test files and test documentation be updated to reference "actui" instead of "apple-tui"? → A: Yes, update all test files and test documentation

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Build Produces Correctly Named Binary (Priority: P1)

Developers need to build the project and receive a binary file named "actui" instead of "apple-tui". This ensures the product is correctly identified by its intended name.

**Why this priority**: This is the core requirement - without the binary being correctly named, the entire rename effort fails. All downstream usage depends on this.

**Independent Test**: Can be fully tested by running the build process and verifying the output binary filename is "actui".

**Acceptance Scenarios**:

1. **Given** the project repository with build configuration, **When** a developer runs the build command, **Then** the output binary is named "actui"
2. **Given** a clean build environment, **When** the build process completes successfully, **Then** no binary named "apple-tui" is generated
3. **Given** the binary is built, **When** checking the executable name, **Then** it matches "actui" exactly

---

### User Story 2 - Documentation Reflects Correct Binary Name (Priority: P2)

Users and developers reading the README and other documentation need to see accurate references to the "actui" binary name for installation, usage, and troubleshooting instructions.

**Why this priority**: Documentation accuracy is critical for user onboarding and support. Incorrect documentation creates confusion and frustration.

**Independent Test**: Can be fully tested by searching all documentation files for "apple-tui" vs "actui" and verifying correct usage throughout.

**Acceptance Scenarios**:

1. **Given** a user reads the README, **When** they follow installation instructions, **Then** all references use "actui" as the binary name
2. **Given** a developer reviews usage examples, **When** they see command-line examples, **Then** the binary is invoked as "actui"
3. **Given** troubleshooting documentation exists, **When** users reference it, **Then** error messages and examples reference "actui"

---

### User Story 3 - Clean Removal of Old Binary Name References (Priority: P3)

The codebase should have no lingering references to "apple-tui" that could cause confusion or inconsistency in future development.

**Why this priority**: While not immediately blocking functionality, residual references can cause confusion during maintenance and future development.

**Independent Test**: Can be fully tested by performing a comprehensive search across all project files for "apple-tui" and verifying no inappropriate references remain.

**Acceptance Scenarios**:

1. **Given** all source code files, **When** searching for "apple-tui", **Then** only historical references (if any) in comments or documentation remain
2. **Given** configuration files and scripts, **When** examining build-related files, **Then** all active configuration uses "actui"
3. **Given** the entire repository, **When** a new developer joins, **Then** they see consistent use of "actui" throughout

### Edge Cases

- What happens if there are existing "apple-tui" binaries in user's PATH? (Users need to manually remove old binaries and replace with "actui" - no automatic migration)
- What happens if documentation is cached or distributed elsewhere? (External references like blog posts or third-party sites may need separate updates)
- What happens to git history references to "apple-tui"? (Historical references remain intact - only current/active code is updated)
- What happens if CI/CD pipelines reference the old binary name? (Build pipelines must be updated to reference "actui")
- What about package manager configurations? (Not applicable - project uses direct binary distribution only, no package managers)

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

- **FR-001**: Build system MUST generate a binary executable named "actui"
- **FR-002**: README documentation MUST reference "actui" in all installation and usage instructions
- **FR-003**: All user-facing documentation MUST use "actui" as the binary name consistently
- **FR-004**: Build configuration files MUST specify "actui" as the output binary name
- **FR-005**: Module name MUST remain unchanged (only binary name changes)
- **FR-006**: All code comments or documentation referencing the binary name MUST use "actui"
- **FR-007**: Command-line examples in documentation MUST show "actui" invocation
- **FR-008**: Build system MUST NOT generate any binary named "apple-tui"
- **FR-009**: Project files MUST NOT contain active references to "apple-tui" binary name (git history excluded)
- **FR-010**: Package manager configurations are out of scope (project uses direct binary distribution only)
- **FR-011**: No backward compatibility or migration tooling required (clean break from "apple-tui" to "actui")
- **FR-012**: Source directory MUST be renamed from "cmd/apple-tui/" to "cmd/actui/" for consistency
- **FR-013**: Root-level "apple-tui" file (if present) MUST be removed or renamed to "actui"
- **FR-014**: Test files and test documentation MUST be updated to reference "actui" instead of "apple-tui"

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: Build process successfully produces binary named "actui" on all platforms
- **SC-002**: Zero occurrences of "apple-tui" remain in active code, configuration, documentation, or test files
- **SC-003**: README and all documentation files reference "actui" consistently (100% accuracy)
- **SC-004**: Build completes successfully without errors or warnings related to binary naming
- **SC-005**: Users following README instructions use "actui" command without confusion
- **SC-006**: All test suites pass with updated "actui" references
