# Quickstart: Validate Container 1.6 Support

**Feature**: specs/015-container-1-6-support/spec.md | **Date**: 2026-11-02

## Prerequisites

- macOS Apple Silicon with Apple Container 1.6.1: `container --version` shows 1.6.1; `container system start` running.
- Build actui: `go mod download && go build -o actui ./cmd/actui`.

## Validation scenarios

1. **Kubernetes rename**: Run `./actui`, press `k`. Expect the cluster list to load via `container cluster list` (no "unknown command" error). Select a cluster with `enter`; expect submenu actions refresh / load-image / write-config / delete / create / back — no start action.
2. **Multi-worker create**: From cluster list press `c`; fill workers (e.g. `2`) and CNI (`NONE` or blank). Expect the preview to contain `--workers 2` (and `--cni NONE` when set). Approve with `y`; expect success or a readable CLI error, never a hang.
3. **Node-targeted load-image**: From a cluster submenu choose load image; enter image plus node names. Expect preview `container cluster load-image --name <c> <img> [--platform ...] [--node <n> ...]`. Approve and check the reported result.
4. **Commit container**: From the container list select a container, open submenu with `enter`, choose commit, enter a reference like `demo:actui-test`. Expect preview `container commit <id> demo:actui-test`. Approve; verify via image list or readable error. Repeat with `--dry-run` and confirm nothing is created.
5. **Machine options**: Open machine list (`M`), press `c`, fill image plus CPUs/memory/home-mount. Expect preview `container machine create <img> --name ... --cpus ... --memory ... --home-mount ...`. Invalid CPUs (e.g. `0`) must show a validation message before preview.
6. **Docs check**: `grep -rn "k8s" README.md docs/ src/ --include='*.go' --include='*.md'` returns no command references (parser comments updated); README states 1.6.1 backend and documents commit + cluster workflows.

## Automated gate

`gofmt -l ./...` clean, `go vet ./...` clean, `go test ./...` all pass.
