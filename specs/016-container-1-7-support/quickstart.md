# Quickstart: Validate Container 1.7 Support

**Feature**: specs/016-container-1-7-support/spec.md | **Date**: 2026-11-09

## Prerequisites

- macOS 27.x on Apple Silicon with Apple Container 1.7.1: `container --version` shows 1.7.1; `container system start` running.
- Build actui: `go mod download && go build -o actui ./cmd/actui`.

## Validation scenarios

1. **No load-image anywhere**: Run `./actui`, press `k`, select a cluster with `enter`. Expect submenu actions refresh / write-config / delete / create / back — no load-image, no start. `grep -rn "load-image\|LoadImage\|ScreenKubernetesLoadImage" src/ --include='*.go'` returns nothing.
2. **Aggregate + per-node list**: With one multi-node cluster present, open `k`. Expect each cluster once with correct state. Open submenu detail and confirm node rows match `container cluster list --nodes` output.
3. **Shared-image-store + sized workers create**: From cluster list press `c`; fill workers (e.g. `count=2,cpus=2,memory=4G`) and enable shared image store. Expect preview `container cluster create ... --workers count=2,cpus=2,memory=4G --shared-image-store`. Approve with `y`; expect success or a readable CLI error, never a hang. Repeat with `--dry-run` and confirm nothing is created.
4. **Delete guard unchanged**: From submenu choose delete; expect type-to-confirm on the exact cluster name, preview `container cluster delete --name <name>`.
5. **Unchanged flows spot-check**: container commit preview, machine create with cpus/memory/home-mount, image pull — all previews identical to 1.6.1 behavior.
6. **Docs check**: README states 1.7.1 backend + macOS 27 prerequisite; `grep -rni "1\.6\.1\|macOS 26\|load-image" README.md docs/ .specify/memory/constitution.md` returns no stale user-facing references (historical spec dirs under `specs/015-*` excluded — they describe the prior release).

## Automated gate

`gofmt -l ./...` clean, `go vet ./...` clean, `go test ./...` all pass.
