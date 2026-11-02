# Contracts: Container 1.6 Command Mapping

**Feature**: specs/015-container-1-6-support/spec.md | **Date**: 2026-11-02
**Source of truth**: installed `container 1.6.1` `--help` output, verified 2026-11-02.

Every state-changing action shows this exact command in a preview and requires approval. Removals require type-to-confirm.

## Kubernetes (all renamed `k8s` -> `cluster`, start removed)

| Action | Command |
|---|---|
| List clusters | `container cluster list` |
| Create cluster | `container cluster create [--name <name>] [--cpus <n>] [--memory <size>] [--rm] [--node-image <ref>] [--workers <n>] [--cni <path\|NONE>]` (empty optionals omitted) |
| Delete cluster | `container cluster delete --name <name>` (type-to-confirm) |
| Load image | `container cluster load-image --name <name> <image> [--platform <p>] [--node <n> ...]` (repeatable `--node`) |
| Write kubeconfig | `container cluster write-config --name <name> [--kubeconfig <path>]` |
| Start cluster | REMOVED — no 1.6.1 subcommand |

## Containers (new)

| Action | Command |
|---|---|
| Commit container to image | `container commit <id> <reference>` (preview + approve, dry-run safe) |

## Machines (extended create)

| Action | Command |
|---|---|
| Create machine | `container machine create <image> [--name <name>] [--cpus <n>] [--memory <size>] [--home-mount <ro\|rw\|none>]` (empty optionals omitted) |
| Edit resources | `container machine set -n <id> cpus=<n> memory=<m> home-mount=<mode>` (unchanged) |
| Start machine | `container machine run -n <id>` (unchanged) |

## Docs contract

- README states Apple Container 1.6.1 backend, `container cluster` family, no start action, commit action, `--workers`/`--cni`/`--node` options, machine cpus/memory/home-mount options.
- `docs/user-guide.md`, `docs/ai-menu-map.md` updated to match; no `k8s` references remain in user-facing instructions or code.
