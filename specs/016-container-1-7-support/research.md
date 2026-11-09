# Research: Container 1.7 Support

**Feature**: specs/016-container-1-7-support/spec.md | **Date**: 2026-11-09
**Method**: installed `container 1.7.1` CLI help output + apple/container GitHub release 1.7.1 (user also asked to map `container k8s --help` directly)

## Verified CLI facts (1.7.1 installed, author machine macOS 27.2)

- `container --version` reports `container CLI version 1.7.1 (build: release, commit: a15b460)`; `sw_vers` reports `ProductVersion 27.2`.
- `container k8s --help` fails with `Error: unknown command 'k8s'` (exit 64) — the legacy alias is fully gone, not just deprecated. actui has no `k8s` builders left (015 already renamed), so no code references break; only docs/tests mentioning `k8s start` in forbidden-lists remain relevant.
- `container cluster --help` subcommands: `create`, `delete/rm`, `list/ls`, `write-config`. **No `load-image`, no `run`, no `start`.** Note: the 1.7.1 release notes say "`--shared-image-store` option on `container cluster run`" (PR #2437) but the installed 1.7.1 CLI exposes no `cluster run` subcommand — help output is authoritative here, and the examples block shows `--shared-image-store` belongs to `cluster create` (`container cluster create --name my-cluster --shared-image-store`). The plan follows installed-CLI help, not the release-note wording.
- `container cluster list` prints one row per cluster by default (`NAME NODES K8S VERSION STATE` — verified against empty output `[]` JSON and header-only table). `--nodes` restores the old per-node table (`CLUSTER NODE ROLE STATE CPUS MEMORY ADDR PORTS` — header verified live). JSON list output is `[]` when empty; parser must handle both shapes.
- `cluster create --help`: `--name`, `--rm`, `--cpus`, `--memory`, `--scheme`, `--max-concurrent-downloads`, `--node-image`, `--cni` (manifest path or `NONE`), repeatable `--workers` in form `count=<n>,cpus=<n>,memory=<size>` (a bare number also accepted for count-only), `--shared-image-store` (no value). Control-plane `--cpus/--memory` semantics unchanged from 1.6.1.
- `machine create` / `machine set` help output is identical in flags to the 1.6.1 baseline (cpus, memory, home-mount ro/rw/none, virtualization, kernel). No machine changes needed.
- `container build --help` shows no syntax change (same options as 1.6.1 baseline). BuildKit fixes are internal; no actui change.
- `container commit <id> <reference>` unchanged. `stats` gains `cpuUtilization` in machine-readable output but actui has no stats screen — no change.
- `system status/dns/property/kernel/df`, `network`, `volume`, `image` families: no new subcommands relevant to actui surfaces (DNS defaults now configurable via daemon `config.toml` — daemon-side, no TUI surface).

## Release-note deltas (1.6.1 -> 1.7.1)

Source: https://github.com/apple/container/releases/tag/1.7.1 (Highlights + What's Changed).

1. Kubernetes (all breaking, all in scope): `load-image` removed, replaced by `--shared-image-store` on create (PR #2437); `cluster list` one-row-per-cluster, `--nodes` for per-node view (PR #2451); lifecycle commands (`start/stop/delete/...`) no longer operate on cluster nodes (PR #2462); worker sizing `--workers count=<n>,cpus=<n>,memory=<size>` (PR #2448).
2. Core (out of scope, no actui surface): `stats` `cpuUtilization`, Time Machine exclusions, Kata kernel 4.3.1, Containerization 0.54.0.
3. Networking (out of scope): default DNS via `config.toml` — daemon config file, not a CLI flag actui previews.
4. Build (out of scope): BuildKit + build-context fixes, same `build` syntax.
5. Infra (noted for docs): upstream CI moved to macOS 27 runners (PR #2436) — corroborates the user's macOS 27 prerequisite request.

## Decisions

- Decision: Remove the load-image workflow entirely (builder + screen + submenu option + app wiring + tests), replace with a `--shared-image-store` toggle on the create form.
  - Rationale: the subcommand does not exist on 1.7.1; every load-image preview actui builds today fails. The toggle is the upstream-blessed replacement and fits the existing create-form pattern.
  - Alternatives considered: keep builder for older CLIs — rejected per single-version target (spec assumption); keep screen but show "unavailable" — rejected as dead-code surface violating command-safe principle.
- Decision: Keep the existing cluster-list screen and extend the parser to accept BOTH the one-row-per-cluster shape and the per-node `--nodes` shape; run list with `--nodes` when node detail is needed (submenu detail already shows per-node rows from stored cluster state).
  - Rationale: verified live that default and `--nodes` outputs have different headers (`NAME NODES K8S VERSION STATE` vs `CLUSTER NODE ROLE ...`); a shape-detecting parser covers both without a second screen.
  - Alternatives considered: separate nodes view screen — rejected; existing submenu already renders node detail, and list screen shows aggregates.
- Decision: Worker sizing is a single free-text field accepting the upstream spec (`<n>` or `count=<n>,cpus=<n>,memory=<size>`, repeatable via comma/space separation passed through as repeated `--workers` flags), validated leniently (non-empty token, must parse count >= 0), matching how actui already treats Workers/CNI as pass-through strings.
  - Rationale: mirrors the 015 approach (workers as validated string, not structured form); avoids building a multi-group sub-form for a rarely repeated flag.
  - Alternatives considered: structured count/cpus/memory sub-fields — rejected as form bloat for a single CLI flag.
- Decision: No changes to containers/images/machines/registries/export/commit/daemon flows; machine `create`/`set` and `build` help verified unchanged.
  - Rationale: 1.7.1 help output shows identical flags to the 1.6.1 baseline for every builder actui owns outside `cluster`.
- Decision: Docs sweep covers README.md, docs/user-guide.md, docs/ai-menu-map.md, docs/binary-build-automation.md (26.x verification references), .github/agents/copilot-instructions.md (26.x validation note), tests/integration/documentation_test.go (required/forbidden strings), help screen version text if it names a backend version, plus the constitution amendment macOS 26.x -> 27.x (FR-011, governance rule).
  - Rationale: user explicitly requested README-first full doc update plus the macOS 27 prerequisite; grep shows stale version strings in exactly these files.
- Decision: Node-identifier lifecycle guard is a docs + error-readability concern only, not a new validation layer.
  - Rationale: actui never constructs node IDs as command targets (list shows them, actions target cluster names); upstream refusal surfaces through existing error paths. No code change needed beyond keeping previews accurate.

## Risks

- The release notes mention `container cluster run` but installed 1.7.1 has no such subcommand; if a later 1.7.x adds it, the create-only `--shared-image-store` mapping may need revisiting. Plan pins to installed-CLI help output (authoritative) and notes the discrepancy here.
- `cluster list --format json` shape on 1.7.1 for populated clusters is unverified (no clusters on author machine at plan time); parser stays whitespace-table-first with JSON tolerance, and quickstart requires live verification with a real cluster before merge.
- Partial cluster-create failures (port conflicts) surface as CLI errors; no actui-side retry logic — documented as read-through behavior (same as 015).
