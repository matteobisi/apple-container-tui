# Research: SBOM Generation for Binary Builds

**Phase 0 output for** `011-sbom-binary-build`  
**Date**: 2026-04-04  
**Status**: Complete — all unknowns resolved

---

## Topic 1: OSSF Scorecard SBOM Check Criteria

**Decision**: Generate SBOM and attach it to every GitHub release as a release asset. No minimum score threshold is required; presence of a valid SBOM in the latest release awards the full Scorecard SBOM check score.

**Rationale**: The Scorecard `SBOM` check (introduced in Scorecard v5) scans the most recent repository release for a file whose name or content indicates a valid SBOM. Specifically it looks for:
- Release assets with `.spdx`, `.spdx.json`, or CycloneDX-compatible filenames
- The file must be parseable as a valid SPDX or CycloneDX document
- The SBOM must reference the repository and enumerate at least one package/component

The check does NOT require the SBOM to be signed or attested for the base score. Attaching the SBOM as a release asset is the primary qualifying action.

**Alternatives considered**:
- Committing SBOM to repository root: Not recognised by the Scorecard SBOM check; only release assets are evaluated.
- GitHub Dependency Graph only: The built-in dependency graph is not SBOM-compatible for Scorecard purposes; a standalone SBOM file in the release is required.

---

## Topic 2: SBOM Tool Selection for GitHub Actions + Go

**Decision**: Use `anchore/sbom-action` (powered by Syft) as the SBOM generation step.

**Rationale**: `anchore/sbom-action` is the most mature and widely adopted SBOM Action for Go projects. Syft natively resolves Go module dependency graphs from `go.mod`/`go.sum` and produces standards-compliant SPDX 2.3 JSON and CycloneDX outputs. It scans the repository source tree (not the binary) to build the component list, which for a Go project is comprehensive and accurate.

Key differentiators:
- Outputs SPDX 2.3 JSON, the format recommended by GitHub and OSSF
- Handles Go modules including indirect and transitive dependencies
- Can upload the SBOM directly as a workflow artifact in one step
- Actively maintained with Node 24 compatible runner support
- Industry-standard tool used in OSSF and GitHub reference implementations

**Alternatives considered**:
- `microsoft/sbom-tool`: Generates SPDX 2.2 JSON; less integrated with GitHub Actions; overkill for a pure Go CLI project.
- `actions/attest-sbom` (GitHub native attestation): Creates a signed attestation attached to GHCR or the GitHub Attestation Store, not a release asset file. Improves SLSA score but does NOT alone satisfy the Scorecard SBOM release-asset check. Can be added as a future enhancement.
- Manual `syft` CLI invocation: Same Syft engine without the Action wrapper; more complex to configure correctly and does not benefit from Action caching or input/output conventions.
- `go list -json`: Produces raw Go module metadata, not a valid SPDX/CycloneDX document. Requires significant post-processing to be Scorecard-compatible.

---

## Topic 3: SBOM Format and Naming Convention

**Decision**: Generate SPDX 2.3 JSON, named `actui-linux-amd64.spdx.json`, and attach it to releases under that filename.

**Rationale**: SPDX 2.3 JSON is the format used by GitHub's own dependency export API, referenced in the OSSF SBOM Everywhere guidance, and is the format most reliably detected by the Scorecard SBOM check. The `.spdx.json` extension is a recognised indicator in both Scorecard and SPDX tooling.

Naming convention follows the parallel with the binary artifact:
| Artifact type | Name |
|---|---|
| Binary | `actui-linux-amd64` |
| SBOM | `actui-linux-amd64.spdx.json` |

This naming makes the relationship between binary and SBOM unambiguous at the release asset level.

**Alternatives considered**:
- CycloneDX JSON (`actui-linux-amd64.cdx.json`): Also supported by Syft and compatible with Scorecard, but SPDX is the GitHub and OSSF reference standard and is preferred for maximum compatibility.
- SPDX tag-value format (`.spdx`): Text-based, harder to parse programmatically. JSON is the interoperable format.
- Generic name (`sbom.spdx.json`): Less clear when multiple artifacts exist in a release; the artifact-scoped name is clearer for audit purposes.

---

## Topic 4: anchore/sbom-action SHA Pin

**Decision**: Pin `anchore/sbom-action` to its latest stable commit SHA, to be verified at implementation time using `git ls-remote https://github.com/anchore/sbom-action`.

**Rationale**: The AGENTS.md policy requires all actions to be pinned to immutable commit SHAs. `anchore/sbom-action` is a composite action (no Node.js runtime dependency), so it is compatible with any runner. The SHA must be resolved at implementation time to guarantee both immutability and the latest security fixes.

**Verification procedure** (for the implementing task):
```sh
git ls-remote https://github.com/anchore/sbom-action refs/tags/v0.17.0
# or query the latest tag first:
git ls-remote --tags https://github.com/anchore/sbom-action | grep -E 'refs/tags/v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -1
```

**Alternatives considered**:
- Using the bare tag (`anchore/sbom-action@v0.17.0`): Prohibited by the repo's action pinning policy. Tags are mutable and can be moved.
- Using `@main`: Prohibited; refers to a moving target.

---

## Topic 5: Workflow Permission Requirements

**Decision**: SBOM generation without GitHub attestation requires only `contents: read` at the job level (the default set in `build-binary.yml`) with no additional permissions beyond what is already present.

**Rationale**:  
- Generating the SBOM file and uploading it as a workflow artifact requires only read access to the repository contents.
- `actions/upload-artifact` does not require write permissions at the workflow level.
- The existing `build-binary.yml` already has `permissions: contents: read` which covers the SBOM generation step.

If GitHub signed attestation is added in a future iteration:
| Extra capability | Required permission |
|---|---|
| `actions/attest-sbom` | `id-token: write`, `attestations: write` |

The `publish-release.yml` already has `contents: write` (required to create releases) and `actions: read` (required to download cross-run artifacts), which is sufficient to download the SBOM artifact and attach it to the release.

**Alternatives considered**:
- Granting `write-all` for simplicity: Violates least-privilege principle and would lower the Scorecard `Token-Permissions` score, counter-productive to the goal.

---

## Topic 6: Integration Architecture

**Decision**: Add SBOM generation as a new step in `build-binary.yml` immediately after the binary build step. Upload the SBOM as a separate named artifact. Update `publish-release.yml` to download the SBOM artifact and attach it to the release alongside the binary.

**Rationale**:  
- Keeping SBOM generation inside `build-binary.yml` ensures the SBOM is generated from the same source state that produced the binary; they are co-located by workflow run.
- A separate workflow would require coordinating on the same workflow run ID, adding complexity.
- Using the existing `workflow_run` downstream mechanism in `publish-release.yml` for the SBOM ensures the SBOM is always attached to every published release automatically.

**Artifact flow**:
```
build-binary.yml (push to main / workflow_dispatch)
  ├── Build binary                → actui
  ├── Generate SBOM (Syft)        → actui-linux-amd64.spdx.json
  ├── Upload binary artifact      → artifact: actui-linux-amd64
  └── Upload SBOM artifact        → artifact: actui-linux-amd64-sbom

publish-release.yml (triggered by workflow_run on "Build Binary")
  ├── Download binary artifact    ← actui-linux-amd64
  ├── Download SBOM artifact      ← actui-linux-amd64-sbom
  ├── Verify & rename binary      → release-assets/actui-linux-amd64
  ├── Verify & rename SBOM        → release-assets/actui-linux-amd64.spdx.json
  └── Publish release             → attaches both assets
```

**Alternatives considered**:
- Separate SBOM workflow: Adds unnecessary coordination complexity (same run ID linking needed) and makes the build pipeline harder to reason about.
- Generate SBOM in `publish-release.yml` at publish time: The SBOM would reflect the release workflow's source state, not necessarily the exact source used for the build. Coupling SBOM to the build run is more accurate and reproducible.

---

## Topic 7: Additional Scorecard Quick Wins

**Decision**: Pin the three unpinned action references in `build-binary.yml` to their commit SHAs as part of this change. This directly improves the Scorecard `Pinned-Dependencies` check.

**Rationale**: `build-binary.yml` currently uses bare major-version tags:
```yaml
actions/checkout@v4         # should be SHA-pinned
actions/setup-go@v5         # should be SHA-pinned
actions/upload-artifact@v4  # should be SHA-pinned
```
The AGENTS.md lists commit SHAs for `actions/checkout` and `actions/upload-artifact`. The `actions/setup-go@v5` SHA must be verified at implementation time. Fixing these in the same PR that adds SBOM maximises the Scorecard improvement from a single change.

**Verified SHA pins from AGENTS.md**:
| Action | Tag | Commit SHA |
|---|---|---|
| `actions/checkout` | v4.2.2 | `11bd71901bbe5b1630ceea73d27597364c9af683` |
| `actions/upload-artifact` | v4.6.2 | `ea165f8d65b6e75b540449e92b4886f43607fa02` |
| `actions/setup-go` | (verify at impl) | (resolve via `git ls-remote`) |
| `anchore/sbom-action` | (verify at impl) | (resolve via `git ls-remote`) |

**Alternatives considered**:
- Deferring the SHA pinning fix to a separate PR: The Scorecard benefit is incremental; doing it in the same PR maximises the lift for a single reviewer cycle.
