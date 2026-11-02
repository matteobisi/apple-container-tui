# Research: Container 1.6 Support

**Feature**: specs/015-container-1-6-support/spec.md
**Date**: 2026-11-02
**Method**: installed `container 1.6.1` CLI help output + apple/container GitHub releases 1.6.0/1.6.1

## Verified CLI facts (1.6.1 installed)

- `container --version` reports `container CLI version 1.6.1 (build: release, commit: aa8d76b)`.
- Top-level commands include new `commit` ("Create a new image from a container's filesystem": `container commit <id> <reference>`) and `clean`; `cluster` replaces `k8s`.
- `container k8s --help` fails with `Error: unknown command 'k8s'` — actui's current `k8s ...` builders are broken on 1.6.1.
- `container cluster` subcommands: `create`, `delete/rm`, `list/ls`, `load-image`, `write-config`. No `start` subcommand exists (upstream removed `k8s start` in 1.5.0; 1.6.1 `cluster` has no start either).
- `container cluster list` prints the same whitespace table actui already parses (`CLUSTER NODE ROLE STATE CPUS MEMORY ADDR PORTS`); the existing parser works unchanged against live output verified on this machine.
- `cluster create --help`: `--name`, `--rm`, `--cpus`, `--memory`, `--scheme`, `--max-concurrent-downloads`, `--node-image`, `--cni` (manifest path or `NONE`), `--workers` (default 0). So the create form gains exactly two fields: Workers, CNI.
- `cluster load-image --help`: `--name`, `--platform` (default `linux/arm64`), repeatable `--node` for targeting specific nodes. So the load-image form gains one field: Nodes (comma-separated, split into repeated `--node` flags).
- `container machine create --help`: image + `--name`, `--set-default`, `--no-boot`, `--cpus`, `--memory`, `--home-mount (ro|rw|none)`, `--virtualization`, `--kernel`, arch/os/platform, registry/progress/fetch options. The pragmatic actui addition is CPUs, Memory, HomeMount (the three also supported by `machine set`, matching the existing edit-resources screen).
- `container machine set --help`: `cpus=<N> memory=<size> home-mount=<ro|rw|none> virtualization=<bool> kernel=<path>`. Existing `MachineSetBuilder` already covers cpus/memory/home-mount, so no edit-screen change is needed beyond keeping parity.
- `container machine run` (used by the "start machine" action) still exists, so the machine start/stop flow is unaffected.
- `container build`/`container builder` syntax is unchanged (plugin repackaging only) — no actui build-flow change.

## Release-note deltas (1.5.0 -> 1.6.0 -> 1.6.1)

1.6.0 highlights: `container commit` command; DNS aliases for containers (network-attach flag, no actui surface); `container build`/`builder` become CLI plugins (same syntax); multi-worker `container k8s` clusters (`--workers`, `--node` load targeting, `--cni`); security fixes (digest pinning, build symlink). 1.6.1: rename `container k8s` to `container cluster`; signal forwarding fixes; cluster-create port-conflict fixes. Out of scope for actui: digest verification internals, build symlink hardening, DNS alias flags (no run/create screen), network/volume screens (constitution out of scope).

## Decisions

- Decision: Target 1.6.1 only; no dual `k8s`/`cluster` fallback.
  - Rationale: spec assumption; fallback adds detection complexity for a CLI the author has already upgraded, and failures surface readably through existing error paths.
  - Alternatives considered: version probe + fallback to `k8s` — rejected as speculative complexity.
- Decision: Remove the cluster start action entirely (builder + submenu option) rather than hiding it conditionally.
  - Rationale: no start subcommand exists on 1.6.1; dead code paths violate the command-safe constitution principle.
  - Alternatives considered: keep builder for older CLIs — rejected per single-version target.
- Decision: Machine creation gains CPUs/Memory/HomeMount only (not virtualization/kernel/no-boot/set-default).
  - Rationale: matches the three values `machine set` supports and the existing edit-resources screen; keeps the form small and consistent.
  - Alternatives considered: full flag parity — rejected as form bloat for rarely used exotic flags.
- Decision: Commit is a direct preview action from the container submenu (stopped or running), not a dedicated screen.
  - Rationale: single input (image reference) fits the existing modal input pattern used elsewhere; a full screen is unnecessary.
  - Alternatives considered: dedicated commit screen — rejected as over-engineering for one field.
- Decision: `cluster` table parsing is unchanged; only the invoked command name changes.
  - Rationale: live `container cluster list` output verified identical in shape to the `k8s list` output the parser already handles.

## Risks

- Upstream may rename `cluster` output columns later; parser stays whitespace-based and tolerant, covered by existing parser tests.
- Partial cluster-create failures (port conflicts) surface as CLI errors; no actui-side retry logic — documented as read-through behavior.
