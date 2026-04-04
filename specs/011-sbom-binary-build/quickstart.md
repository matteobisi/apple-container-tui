# Quickstart: SBOM Generation for Binary Builds

**Feature**: `011-sbom-binary-build`  
**Date**: 2026-04-04  
**Audience**: Maintainers implementing or verifying this feature

---

## What This Feature Does

Every time the `Build Binary` workflow runs successfully, it now generates a Software Bill of Materials (SBOM) alongside the binary and uploads both as separate workflow artifacts. When the `Publish Release` workflow publishes a release, both the binary and the SBOM are attached as release assets. This makes the SBOM discoverable by the OSSF Scorecard tool and downloadable by anyone reviewing a release.

---

## Implementation Checklist

### Step 1: Resolve SHA pins

Before editing any workflow file, resolve the commit SHAs for the two unpinned actions:

```sh
# Resolve setup-go latest v5 tag SHA
git ls-remote --tags https://github.com/actions/setup-go \
  | grep -E 'refs/tags/v5\.[0-9]+\.[0-9]+$' | sort -V | tail -1

# Resolve anchore/sbom-action latest tag SHA
git ls-remote --tags https://github.com/anchore/sbom-action \
  | grep -E 'refs/tags/v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -1
```

Record the tag name and SHA. Use the pattern:
```yaml
uses: actions/setup-go@<SHA> # v5.X.Y, Node 24 compatible
uses: anchore/sbom-action@<SHA> # vX.Y.Z, composite action
```

### Step 2: Update build-binary.yml

**Pin all three existing bare-tag references to commit SHAs** (quick win for Scorecard `Pinned-Dependencies`):

```yaml
# Before                              # After
actions/checkout@v4                   actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2
actions/setup-go@v5                   actions/setup-go@<RESOLVED-SHA> # v5.X.Y
actions/upload-artifact@v4            actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02 # v4.6.2
```

**Add SBOM generation step** after the `Build actui` step and before the `Upload artifact` step:

```yaml
- name: Generate SBOM
  uses: anchore/sbom-action@<RESOLVED-SHA> # vX.Y.Z, composite action
  with:
    artifact-name: actui-linux-amd64.spdx.json
    format: spdx-json
    output-file: actui-linux-amd64.spdx.json
    upload-artifact: false  # we upload manually below for explicit control

- name: Upload SBOM artifact
  uses: actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02 # v4.6.2, Node 24 compatible
  with:
    name: actui-linux-amd64-sbom
    path: actui-linux-amd64.spdx.json
    if-no-files-found: error
    retention-days: 14
```

### Step 3: Update publish-release.yml

**Add SBOM download step** after the existing binary download step:

```yaml
- name: Download artifact actui-linux-amd64-sbom
  uses: actions/download-artifact@cc203385981b70ca67e1cc392babf9cc229d5806 # v4.1.9, Node 24 compatible
  with:
    name: actui-linux-amd64-sbom
    path: release-assets/
    run-id: ${{ github.event.workflow_run.id }}
    github-token: ${{ secrets.GITHUB_TOKEN }}
```

**Add SBOM verification step** after the binary verification step:

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

**Extend the publish release step** to include the SBOM:

```yaml
gh release create "$TAG" \
  release-assets/actui-linux-amd64 \
  release-assets/actui-linux-amd64.spdx.json \
  --title "$TITLE" \
  --generate-notes \
  --target "$BUILD_SHA"
```

### Step 4: Update documentation

Update `docs/binary-build-automation.md` to add an SBOM section describing:
- When the SBOM is generated
- The format (SPDX 2.3 JSON)
- Where to find it (workflow artifact and release asset)
- How to verify it

---

## Verification After Deployment

### Verify SBOM artifact in workflow run

1. Navigate to the repository → Actions → Build Binary
2. Open any successful run after this change
3. Under **Artifacts**, confirm `actui-linux-amd64-sbom` is listed alongside `actui-linux-amd64`
4. Download `actui-linux-amd64-sbom` and confirm the file `actui-linux-amd64.spdx.json` is inside
5. Open the file and confirm it starts with `{"spdxVersion":"SPDX-2.3"` and contains package entries

### Verify SBOM is attached to release

1. Navigate to the repository → Releases
2. Open the latest release
3. Under **Assets**, confirm `actui-linux-amd64.spdx.json` appears alongside `actui-linux-amd64`
4. Download `actui-linux-amd64.spdx.json` and verify it is valid SPDX JSON

### Verify SBOM content (optional but recommended)

Install [Syft CLI](https://github.com/anchore/syft) or [spdx-tools](https://github.com/spdx/tools-golang) and validate:

```sh
# Quick format check with jq
cat actui-linux-amd64.spdx.json | jq '.spdxVersion, (.packages | length)'
# Expect: "SPDX-2.3" and a count > 0

# Full SPDX validation (requires spdx-tools)
pyspdxtools validate actui-linux-amd64.spdx.json
```

### Verify Scorecard improvement

After the first release with the SBOM attached:

1. Navigate to the repository's OSSF Scorecard output (GitHub Security tab → Supply-chain security, or Scorecard API)
2. Confirm the `SBOM` check now shows a passing score (10/10)
3. Confirm the `Pinned-Dependencies` check score has improved or remains at its current level

---

## Troubleshooting

| Symptom | Likely cause | Fix |
|---|---|---|
| SBOM artifact missing from workflow run | `anchore/sbom-action` step failed; check `if-no-files-found: error` | Review step logs; check `go.mod` is present at repository root |
| SBOM artifact empty or malformed | Syft could not resolve Go module graph | Ensure `go.sum` is committed and `go.mod` is valid |
| Publish fails with SBOM not found | Download artifact step uses wrong `name` or `path` | Verify `name: actui-linux-amd64-sbom` matches exactly what build workflow uploads |
| Scorecard SBOM check still 0 after release | SBOM filename does not end in `.spdx.json` | Confirm release asset is named `actui-linux-amd64.spdx.json` |
| SHA pin resolution fails | Tag does not exist or repo is unreachable | Use `curl -s https://api.github.com/repos/anchore/sbom-action/releases/latest` to find the latest tag |
