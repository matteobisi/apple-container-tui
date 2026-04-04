# Tasks: SBOM Generation for Binary Builds

**Input**: Design documents from `specs/011-sbom-binary-build/`  
**Prerequisites**: plan.md ✅ spec.md ✅ research.md ✅ data-model.md ✅ contracts/workflow-contract.md ✅ quickstart.md ✅

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story. No test tasks generated — none requested in the feature specification.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4)

---

## Phase 1: Setup

**Purpose**: Resolve the two SHA pins that are not pre-known and are required before any workflow file can be edited. All other SHAs are already listed in AGENTS.md.

- [x] T001 Resolve SHA pins for `actions/setup-go` (v5 series) and `anchore/sbom-action` (latest stable) by running the following commands and recording both the tag name and commit SHA for use in Phase 2–5 tasks:
  ```sh
  git ls-remote --tags https://github.com/actions/setup-go \
    | grep -E 'refs/tags/v5\.[0-9]+\.[0-9]+$' | sort -V | tail -1

  git ls-remote --tags https://github.com/anchore/sbom-action \
    | grep -E 'refs/tags/v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -1
  ```

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Pin the three existing bare-tag action references in `build-binary.yml` to immutable commit SHAs. This is a hard prerequisite for all subsequent tasks (same file) and directly improves the Scorecard `Pinned-Dependencies` score.

**⚠️ CRITICAL**: No user story tasks that touch `build-binary.yml` can begin until this phase is complete.

- [x] T002 Replace bare-tag action references in `.github/workflows/build-binary.yml` with immutable commit SHA pins and version comments:
  - `actions/checkout@v4` → `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2, Node 24 compatible`
  - `actions/setup-go@v5` → `actions/setup-go@<SHA-FROM-T001> # v5.X.Y, Node 24 compatible`
  - `actions/upload-artifact@v4` → `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02 # v4.6.2, Node 24 compatible`

**Checkpoint**: `build-binary.yml` has zero bare-tag action references. User story implementation in this file can now begin.

---

## Phase 3: User Story 1 — Generate SBOM Automatically During Binary Build (Priority: P1) 🎯 MVP

**Goal**: Every successful binary build run produces an SPDX 2.3 JSON SBOM file on the runner immediately after the binary is compiled.

**Independent Test**: Trigger the `Build Binary` workflow via `workflow_dispatch`. After it completes, confirm a file named `actui-linux-amd64.spdx.json` exists on the runner (visible in the step output), and that the SBOM step itself succeeded. The file content begins with `{"spdxVersion":"SPDX-2.3"` and contains at least one package entry.

### Implementation for User Story 1

- [x] T003 [US1] Add `anchore/sbom-action` SBOM generation step immediately after the `Build actui` step in `.github/workflows/build-binary.yml`, pinned to the SHA resolved in T001, producing `actui-linux-amd64.spdx.json` in SPDX JSON format:
  ```yaml
  - name: Generate SBOM
    uses: anchore/sbom-action@<SHA-FROM-T001> # vX.Y.Z, composite action
    with:
      artifact-name: actui-linux-amd64.spdx.json
      format: spdx-json
      output-file: actui-linux-amd64.spdx.json
      upload-artifact: false
  ```

**Checkpoint**: Running the workflow produces `actui-linux-amd64.spdx.json` on the runner after a successful build. User Story 1 is independently testable.

---

## Phase 4: User Story 2 — Publish SBOM as a Downloadable Build Artifact (Priority: P2)

**Goal**: The SBOM file produced in US1 is uploaded as a separately named workflow artifact (`actui-linux-amd64-sbom`) in the same run as the binary, downloadable from the GitHub Actions UI.

**Independent Test**: After a successful build run, open the run in the GitHub Actions UI → Artifacts section. Confirm `actui-linux-amd64-sbom` is listed alongside `actui-linux-amd64`. Download it and confirm it contains `actui-linux-amd64.spdx.json` with valid SPDX JSON content.

### Implementation for User Story 2

- [x] T004 [US2] Add `actions/upload-artifact` upload step for the SBOM artifact immediately after the SBOM generation step (T003) in `.github/workflows/build-binary.yml`:
  ```yaml
  - name: Upload SBOM artifact
    uses: actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02 # v4.6.2, Node 24 compatible
    with:
      name: actui-linux-amd64-sbom
      path: actui-linux-amd64.spdx.json
      if-no-files-found: error
      retention-days: 14
  ```

**Checkpoint**: Workflow run artifacts include both `actui-linux-amd64` and `actui-linux-amd64-sbom`. User Story 2 is independently testable.

---

## Phase 5: User Story 3 — Include SBOM in Published Releases (Priority: P3)

**Goal**: Every published GitHub release includes `actui-linux-amd64.spdx.json` as a downloadable release asset alongside the binary, satisfying the OSSF Scorecard `SBOM` check.

**Independent Test**: Trigger a full build-to-release cycle. On the repository's Releases page, confirm the new release lists both `actui-linux-amd64` and `actui-linux-amd64.spdx.json` under Assets. Download the SBOM asset and confirm it is valid SPDX 2.3 JSON.

### Implementation for User Story 3

- [x] T005 [US3] Add `actions/download-artifact` step after the existing binary download step in `.github/workflows/publish-release.yml` to download the `actui-linux-amd64-sbom` artifact from the triggering build run into `release-assets/`:
  ```yaml
  - name: Download artifact actui-linux-amd64-sbom
    uses: actions/download-artifact@cc203385981b70ca67e1cc392babf9cc229d5806 # v4.1.9, Node 24 compatible
    with:
      name: actui-linux-amd64-sbom
      path: release-assets/
      run-id: ${{ github.event.workflow_run.id }}
      github-token: ${{ secrets.GITHUB_TOKEN }}
  ```

- [x] T006 [US3] Add SBOM verification step after the existing binary verification step in `.github/workflows/publish-release.yml` to fail the workflow with a clear error if the SBOM file is missing:
  ```yaml
  - name: Verify SBOM artifact
    run: |
      if [ ! -f "release-assets/actui-linux-amd64.spdx.json" ]; then
        echo "::error::SBOM file 'actui-linux-amd64.spdx.json' not found in release-assets/. Cannot publish release."
        ls -la release-assets/ || true
        exit 1
      fi
      echo "SBOM verified: $(ls -lh release-assets/actui-linux-amd64.spdx.json)"
  ```

- [x] T007 [US3] Update the `gh release create` command in the `Publish release` step of `.github/workflows/publish-release.yml` to attach the SBOM as a second release asset alongside the binary:
  ```sh
  gh release create "$TAG" \
    release-assets/actui-linux-amd64 \
    release-assets/actui-linux-amd64.spdx.json \
    --title "$TITLE" \
    --generate-notes \
    --target "$BUILD_SHA"
  ```
  Also update the echo statement for published assets to mention both files.

**Checkpoint**: Published releases include both `actui-linux-amd64` and `actui-linux-amd64.spdx.json`. OSSF Scorecard `SBOM` check should score 10/10 on the next evaluation after a release is published.

---

## Phase 6: User Story 4 — Document the SBOM Generation Process (Priority: P4)

**Goal**: `docs/binary-build-automation.md` contains a dedicated SBOM section that allows a new contributor to understand what is generated, where to find it, and how to verify it without reading the workflow YAML.

**Independent Test**: Read the updated `docs/binary-build-automation.md`. Without looking at any workflow file, a reader must be able to answer: (1) what SBOM format is produced, (2) when it is generated, (3) what the workflow artifact is named, (4) what the release asset is named, and (5) how to do a quick format verification.

### Implementation for User Story 4

- [x] T008 [P] [US4] Add an `## SBOM Generation` section to `docs/binary-build-automation.md` immediately after the `## Artifact Contract` section, covering: format (SPDX 2.3 JSON), generation trigger (same conditions as binary build), workflow artifact name (`actui-linux-amd64-sbom`), artifact file inside (`actui-linux-amd64.spdx.json`), retention (14 days, matches binary), release asset name (`actui-linux-amd64.spdx.json`), quick verification command (`cat actui-linux-amd64.spdx.json | jq '.spdxVersion, (.packages | length)'`), and Scorecard relevance (`.spdx.json` extension triggers Scorecard SBOM check).

**Checkpoint**: User Story 4 is independently testable by reading the documentation section.

---

## Phase 7: Polish & End-to-End Validation

**Purpose**: Confirm the full pipeline works together from push to release.

- [x] T009 Trigger the `Build Binary` workflow via `workflow_dispatch` and execute all verification steps from `specs/011-sbom-binary-build/quickstart.md` in sequence: (1) confirm both artifacts in the workflow run, (2) download and inspect `actui-linux-amd64-sbom`, (3) verify SPDX JSON content with `jq`, (4) after the triggered `Publish Release` run completes, confirm both release assets on the Releases page, (5) note whether Scorecard `SBOM` and `Pinned-Dependencies` checks improved.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on T001 (needs `setup-go` SHA) — BLOCKS T003, T004
- **US1 (Phase 3)**: Depends on T001 (needs `sbom-action` SHA) and T002 (same file) — start after T002
- **US2 (Phase 4)**: Depends on T003 (same file, SBOM file must exist before upload step) — start after T003
- **US3 (Phase 5)**: Logically depends on T003+T004 (artifact must be produced to be consumed); code can be written independently in a different file after T001 is complete
- **US4 (Phase 6)**: No blocking dependencies — `[P]` relative to all other phases
- **Polish (Phase 7)**: Depends on all phases 1–6 completing

### User Story Dependencies

| Story | Depends on | File touched | Can start after |
|---|---|---|---|
| US1 (T003) | T001 + T002 | `.github/workflows/build-binary.yml` | T002 |
| US2 (T004) | T003 | `.github/workflows/build-binary.yml` | T003 |
| US3 (T005–T007) | T001 (logically T003+T004) | `.github/workflows/publish-release.yml` | T001 (code); T004 (full pipeline) |
| US4 (T008) | None | `docs/binary-build-automation.md` | Immediately |

### Within Each User Story

- US1: Single task (T003) — no internal sequencing needed
- US2: Single task (T004) — depends on T003 (same file)
- US3: T005 → T006 → T007 (same file, sequential)
- US4: Single task (T008) — fully independent

### Parallel Opportunities

- **T008 [US4]** can run in parallel with any other task — different file, no dependencies
- **T005** (US3 — different file from build-binary.yml tasks) can be coded in parallel with T003/T004 if the SBOM artifact name and path from the contract are sufficient for the implementor
- All three `publish-release.yml` changes (T005, T006, T007) must be applied sequentially in the same file

---

## Parallel Example: US1 + US4

While the implementor edits `build-binary.yml` for T003 (US1), a second implementor can work on T008 (US4) documentation in parallel:

```
Implementor A:  T001 → T002 → T003 → T004 → T005 → T006 → T007 → T009
                                                           ↑
Implementor B:                               T008 ─────────┘ (merge before T009)
```

---

## Implementation Strategy

**MVP scope (minimum to demonstrate value)**: Complete T001 through T004 (Phases 1–4). This delivers a fully automated SBOM generation and artifact upload on every build — the core deliverable of US1 and US2. The SBOM is accessible per workflow run.

**Full scope**: T005–T007 (US3) attaches the SBOM to releases, which is what actually triggers the Scorecard improvement. T008 (US4) documents the process. T009 validates end-to-end.

**Recommended delivery order**: Implement all tasks in a single PR since the total change surface is small (two workflow files, one docs file). The Scorecard benefit only materialises after the first release with the SBOM attached, so T005–T007 should not be deferred.
