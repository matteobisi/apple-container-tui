# Feature Specification: Security Hardening Quick Wins (OSSF Scorecard 5.9 → 7.5+)

**Feature Branch**: `012-security-hardening-quick-wins`  
**Created**: 2026-04-04  
**Status**: Draft  
**Input**: User description: "Improve the OSSF Scorecard score from ~5.9 to ~7.5-8.0 by implementing two quick, low-effort security hardening items: Security Policy (SECURITY.md) and release build provenance attestation. The stale Dockerfile (from a different project) has been removed, which already fixes the Pinned-Dependencies score (7→10)."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Add Security Policy File (Priority: P1)

A security researcher or contributor finds the project on GitHub and wants to know how to responsibly disclose a vulnerability. They look for a `SECURITY.md` file in the repository root (the standard location), find clear instructions for reporting, and understand the expected response timeline before submitting a report. Separately, the OSSF Scorecard automated scanner detects the file and awards full marks for the Security-Policy check.

**Why this priority**: Zero effort to verify independently, highest impact-to-effort ratio (0→10 score jump on a MEDIUM-risk check), and directly serves human contributors — not just automated tooling.

**Independent Test**: Can be fully tested by verifying a `SECURITY.md` file exists at the repository root, contains vulnerability reporting instructions and a response timeline, and that the OSSF Scorecard Security-Policy check transitions from 0 to 10 on the next scan.

**Acceptance Scenarios**:

1. **Given** no `SECURITY.md` file exists, **When** the file is added to the repository root with reporting instructions and a response timeline, **Then** anyone visiting the repository can find and follow the vulnerability disclosure process without additional searching.
2. **Given** the `SECURITY.md` file is present, **When** the OSSF Scorecard workflow runs, **Then** the Security-Policy check score reports 10/10.
3. **Given** a vulnerability reporter follows the instructions in `SECURITY.md`, **When** they submit a report via the documented channel, **Then** they receive confirmation that the report has been received within the stated timeline.

---

### User Story 2 - Attest Build Provenance for Release Artifacts (Priority: P2)

A downstream user or auditor downloads a release binary from the GitHub Releases page and wants to verify that the binary was produced by the official CI pipeline and has not been tampered with. They can use the attestation attached to the release to verify the artifact's provenance — who built it, when, and from which source commit — without needing any external signing keys.

**Why this priority**: Highest security value (HIGH-risk Scorecard check, Signed-Releases 0→partial), but requires a new step in the release workflow. Changes are confined to a single workflow file and require no external services or keys.

**Independent Test**: Can be fully tested by triggering a release workflow on the branch, verifying that the provenance attestation is attached to the release artifacts on the GitHub Releases page, and confirming the attestation can be verified using the GitHub CLI or Sigstore tooling.

**Acceptance Scenarios**:

1. **Given** the release workflow runs and produces a binary and SBOM, **When** the provenance attestation step executes after artifact upload, **Then** a signed attestation is attached to all published release artifacts.
2. **Given** a signed release is published, **When** a user verifies the attestation on any release artifact, **Then** the attestation resolves to the correct repository, workflow run, and source commit without requiring any external private keys.
3. **Given** the attestation step is added to the workflow, **When** the OSSF Scorecard Signed-Releases check runs against the next release, **Then** the Signed-Releases score is non-zero (partial credit or full credit depending on Scorecard version).
4. **Given** the release workflow runs, **When** the attestation step requires GitHub's ID token, **Then** no additional secrets or credentials need to be configured — only the repository's default `GITHUB_TOKEN` permissions are used.

---

### Edge Cases

- What happens when the Scorecard scanner runs before a new release is published with attestation? The Signed-Releases score may not reflect the improvement until at least one attested release exists.
- What if the `SECURITY.md` file format is not recognized by Scorecard? Scorecard accepts any file named `SECURITY.md` at the repository root or in `.github/`; standard placement ensures detection.
- What if the provenance attestation step fails (e.g., token permissions insufficient)? The release workflow should fail clearly so the missing attestation is never silently skipped on a published release.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The repository MUST contain a `SECURITY.md` file at the repository root.
- **FR-002**: `SECURITY.md` MUST document at least one method for privately reporting security vulnerabilities (e.g., GitHub private vulnerability reporting, a dedicated email address, or GitHub Security Advisories).
- **FR-003**: `SECURITY.md` MUST state the maintainer's expected response timeline for vulnerability reports.
- **FR-004**: `SECURITY.md` MUST specify which versions of the project are currently supported and receive security fixes.
- **FR-005**: The release workflow MUST include a build provenance attestation step that runs after all release artifacts (binary and SBOM) have been uploaded.
- **FR-006**: The provenance attestation step MUST attest all release artifacts produced by the workflow in a single step invocation.
- **FR-007**: The provenance attestation MUST be signed automatically using GitHub's built-in identity infrastructure — no external private keys or third-party signing services are required.
- **FR-008**: All new or updated workflow action references MUST comply with the existing AGENTS.md policy: pinned to immutable commit SHAs with a version comment.
- **FR-009**: The provenance attestation step MUST require only the standard `GITHUB_TOKEN` with `id-token: write` and `attestations: write` permissions — no additional secrets.

### Key Entities

- **Security Policy**: The `SECURITY.md` document that describes the project's vulnerability disclosure process, supported versions, and response commitments.
- **Build Provenance Attestation**: A signed statement attached to a release artifact that cryptographically links the artifact to the source code, build workflow, and build environment that produced it.
- **Release Artifact**: Any file attached to a GitHub Release (binary, SBOM) for which provenance is attested.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The OSSF Scorecard overall score increases from 5.9 to at least 7.5 on the next automated scan after both changes are merged and a new release is published.
- **SC-002**: The Scorecard Security-Policy check score reaches 10/10.
- **SC-003**: The Scorecard Pinned-Dependencies check reaches 10/10 (already resolved by Dockerfile removal).
- **SC-004**: The Scorecard Signed-Releases check reports a non-zero score on the first release published after the provenance attestation step is merged.
- **SC-005**: Each of the two changes can be merged, verified, and demonstrated independently without requiring the other to be complete first.
- **SC-006**: A security researcher arriving at the repository can find the vulnerability reporting process within 30 seconds of landing on the repository home page.
- **SC-007**: A downstream user can verify the provenance of any release artifact using publicly available tooling without any private credentials.

## Assumptions

- The repository uses GitHub's native vulnerability reporting or a maintainer-controlled contact channel; no third-party bug bounty platform is assumed.
- The OSSF Scorecard GitHub Action is already configured and runs on a schedule; no changes to the Scorecard workflow itself are required.
- The `publish-release.yml` workflow already uploads the binary and SBOM as artifacts before the attestation step; the attestation step is appended after the existing upload steps.
- GitHub Actions `id-token: write` and `attestations: write` permissions can be granted to the release workflow without requiring organization-level approval.
- A single `SECURITY.md` at the repository root is sufficient for Scorecard detection; placement in `.github/` is an acceptable fallback but not required.
- The Scorecard Signed-Releases check awards partial credit for provenance attestation even if full SLSA signing is not implemented; full credit is not required to meet SC-001.
- Both changes together (plus the already-done Dockerfile removal) are sufficient to move the overall Scorecard score from 5.9 to the target range of 7.5–8.0; no other checks need to change to achieve this target.
