# Feature Specification: SBOM Generation for Binary Builds

**Feature Branch**: `011-sbom-binary-build`  
**Created**: 2026-04-04  
**Status**: Draft  
**Input**: User description: "I want to increase my security posture adding the SBOM generation when a new binary will be built"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Generate SBOM Automatically During Binary Build (Priority: P1)

As a maintainer, I want a Software Bill of Materials (SBOM) to be automatically generated every time the binary build workflow produces a new artifact, so that I have a complete, machine-readable inventory of all components included in each build without any manual effort.

**Why this priority**: SBOM generation is the core deliverable. Without it, no downstream security, compliance, or audit workflows can function. This is the foundation of the entire feature.

**Independent Test**: Can be fully tested by triggering a binary build, confirming it completes successfully, and verifying that an SBOM file is produced alongside the binary artifact.

**Acceptance Scenarios**:

1. **Given** a push to the default branch triggers the binary build, **When** the build completes successfully, **Then** an SBOM document is generated that lists all dependencies included in the binary.
2. **Given** a manual dispatch triggers the binary build, **When** the build completes successfully, **Then** an SBOM document is generated with the same completeness as an automated run.
3. **Given** the binary build fails before producing an artifact, **When** maintainers review the run, **Then** no partial or empty SBOM is published.

---

### User Story 2 - Publish SBOM as a Downloadable Build Artifact (Priority: P2)

As a maintainer, I want the generated SBOM to be uploaded as a build artifact alongside the binary, so that anyone reviewing or consuming a build can access the SBOM from the same workflow run.

**Why this priority**: An SBOM that is not accessible has no practical value. Publishing it as a retrievable artifact enables downstream consumption by auditors, security tools, and release workflows.

**Independent Test**: Can be fully tested by triggering a build, navigating to the workflow run in the repository, and confirming the SBOM artifact is listed and downloadable separately from the binary artifact.

**Acceptance Scenarios**:

1. **Given** a successful binary build that produces an SBOM, **When** a maintainer views the workflow run artifacts, **Then** the SBOM is available as a separately named, downloadable artifact.
2. **Given** the SBOM artifact is downloaded, **When** a maintainer opens the file, **Then** it contains a valid, well-formed SBOM document that can be parsed by standard tooling.

---

### User Story 3 - Include SBOM in Published Releases (Priority: P3)

As a maintainer, I want the SBOM to be attached to published releases alongside the binary, so that end users downloading a release can also obtain the corresponding SBOM for their own compliance and security review.

**Why this priority**: Distributing the SBOM with the release extends transparency to end users and supports supply-chain security best practices at the distribution point.

**Independent Test**: Can be fully tested by triggering a full build-to-release cycle and confirming the published release includes both the binary and the SBOM as separate downloadable assets.

**Acceptance Scenarios**:

1. **Given** a successful build triggers the release workflow, **When** the release is published, **Then** both the binary and the SBOM document are attached as release assets.
2. **Given** a user visits the release page, **When** they inspect the assets list, **Then** the SBOM asset is clearly identified and downloadable independently from the binary.

---

### User Story 4 - Document the SBOM Generation Process (Priority: P4)

As a maintainer, I want the SBOM generation process documented in the existing build automation documentation, so that future contributors understand what is generated, where to find it, and how to verify it.

**Why this priority**: Documentation preserves operational knowledge and ensures the SBOM workflow remains maintainable as the project evolves.

**Independent Test**: Can be fully tested by reading the updated documentation and confirming it describes the SBOM format, generation trigger, artifact location, and verification steps without requiring code inspection.

**Acceptance Scenarios**:

1. **Given** a new contributor reads the build documentation, **When** they look for SBOM information, **Then** they find clear descriptions of what the SBOM contains, when it is generated, and where it is stored.
2. **Given** the SBOM generation process changes in the future, **When** maintainers update the documentation, **Then** the documentation structure supports updating SBOM-specific details without restructuring other sections.

---

### Edge Cases

- What happens when the SBOM generation step fails but the binary build itself succeeds? The build should still produce the binary artifact, but the SBOM failure must be clearly reported, and the overall workflow should be marked as failed to signal incomplete security artifacts.
- What happens when the project has no external dependencies? The SBOM should still be generated, listing only the standard library and the project module itself.
- How does the SBOM reflect vendored or replaced dependencies? The SBOM must accurately represent the actual dependency graph resolved at build time, including any replacements or vendored modules.
- What happens when the SBOM artifact exceeds retention limits? The SBOM follows the same retention policy as the binary artifact; expired artifacts are handled by the repository's standard artifact lifecycle.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The binary build workflow MUST generate an SBOM document for every successful build that produces a binary artifact.
- **FR-002**: The SBOM MUST use an industry-standard format (SPDX or CycloneDX) to ensure compatibility with common security scanning and compliance tools.
- **FR-003**: The SBOM MUST enumerate all direct and transitive dependencies included in the built binary, including module names, versions, and license information where available.
- **FR-004**: The SBOM MUST be uploaded as a separate, named build artifact in the same workflow run as the binary artifact.
- **FR-005**: The SBOM artifact MUST follow the same retention policy as the binary artifact.
- **FR-006**: The release workflow MUST download and attach the SBOM as an additional release asset when publishing a new release.
- **FR-007**: The SBOM generation MUST NOT alter or delay the binary build step itself; it runs as a post-build step using the already-produced binary or build metadata.
- **FR-008**: If the SBOM generation step fails, the overall workflow run MUST be marked as failed, even if the binary was built successfully.
- **FR-009**: The existing build automation documentation MUST be updated to describe the SBOM generation process, including format, location, and verification steps.
- **FR-010**: The SBOM generation MUST use automation actions pinned by commit SHA with version comments, consistent with the repository's existing action pinning policy.

### Key Entities

- **SBOM Document**: A machine-readable inventory of all software components included in the built binary. Key attributes: format standard (SPDX or CycloneDX), list of components with names and versions, license identifiers, relationships between components, and the build artifact it describes.
- **Binary Artifact**: The compiled executable produced by the build workflow. The SBOM is generated from and associated with a specific binary artifact.
- **Release Asset**: A file attached to a published repository release. Both the binary and the SBOM are published as separate release assets.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of successful binary builds produce a corresponding SBOM artifact in the same workflow run.
- **SC-002**: The generated SBOM passes validation by at least one standard SBOM validation tool without errors.
- **SC-003**: Every published release includes both the binary and the SBOM as separately downloadable assets.
- **SC-004**: Maintainers can locate SBOM documentation and understand the generation process within 5 minutes of reading the build documentation.
- **SC-005**: The SBOM accurately lists all dependencies present in the built binary, with zero missing direct or transitive dependencies.
- **SC-006**: The SBOM generation adds no more than 2 minutes to the total build workflow duration.

## Assumptions

- The existing binary build workflow (`.github/workflows/build-binary.yml`) and release workflow (`.github/workflows/publish-release.yml`) are operational and serve as the integration points for SBOM generation.
- The repository uses Go modules (`go.mod`) as the dependency management system, which provides the dependency graph needed for SBOM generation.
- SPDX is the preferred SBOM format unless tooling constraints favor CycloneDX, as SPDX is an ISO/IEC standard and widely adopted in open-source ecosystems.
- SBOM generation tooling is available as a GitHub Action or command-line tool that can run on the existing `ubuntu-latest` runner environment.
- The SBOM artifact naming convention follows the existing pattern: `actui-linux-amd64-sbom` to parallel the binary artifact name `actui-linux-amd64`.
- Documentation updates target the existing `docs/binary-build-automation.md` file to keep build-related information consolidated.
- No signature or attestation of the SBOM is in scope for this feature; signing and provenance may be addressed in a future security enhancement.
