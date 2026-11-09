# Contracts: Container 1.7 Command Mapping

**Feature**: specs/016-container-1-7-support/spec.md | **Date**: 2026-11-09
**Source of truth**: installed `container 1.7.1` `--help` output (verified 2026-11-09 on macOS 27.2) + https://github.com/apple/container/releases/tag/1.7.1. Installed-CLI help is authoritative where it conflicts with release-note wording (see research.md: no `cluster run` subcommand exists on 1.7.1; `--shared-image-store` is a `cluster create` flag).

Every state-changing action shows this exact command in a preview and requires approval. Removals require type-to-confirm.

## Kubernetes (`cluster` family — load-image removed, list reshaped, create extended)

| Action | Command |
|---|---|
| List clusters (default, one row per cluster) | `container cluster list` |
| List clusters (per-node detail) | `container cluster list --nodes` |
| Create cluster | `container cluster create [--name <name>] [--cpus <n>] [--memory <size>] [--rm] [--node-image <ref>] [--workers <spec> ...] [--cni <path\|NONE>] [--shared-image-store]` (empty optionals omitted; `--workers` repeatable; `<spec>` is `<n>` or `count=<n>,cpus=<n>,memory=<size>`) |
| Delete cluster | `container cluster delete --name <name>` (type-to-confirm) |
| Write kubeconfig | `container cluster write-config --name <name> [--kubeconfig <path>]` |
| Load image | REMOVED — no 1.7.1 subcommand; replaced at create time by `--shared-image-store` |
| Start cluster | REMOVED — no 1.7.1 subcommand (unchanged from 015) |

## Unchanged surfaces (help output identical to 1.6.1 baseline)

Containers (`commit`, list/start/stop/delete/logs/shell/export), images (list/pull/build/prune/inspect/delete), registries, machines (create with cpus/memory/home-mount, set, run, inspect, logs, delete, set-default), daemon (start/stop/status) — no contract change.

## Docs contract

- README states Apple Container 1.7.1 backend, macOS 27.x on Apple Silicon prerequisite, `container cluster` family, no load-image/start actions, `--shared-image-store` replacement, one-row-per-cluster list with `--nodes` detail, worker sizing, unchanged machine options.
- `docs/user-guide.md`, `docs/ai-menu-map.md`, `docs/binary-build-automation.md` (26.x verification references), `.github/agents/copilot-instructions.md` (26.x validation note), `tests/integration/documentation_test.go` (required/forbidden strings) updated to match; no stale `load-image`/`1.6.1`/`macOS 26` user-facing instructions remain.
- Constitution `II. macOS 26.x + Apple Silicon Target` amended to 27.x with rationale and migration notes (FR-011).
