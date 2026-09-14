# Quickstart: Kubernetes Workflow Validation

## Prerequisites

- macOS 26.x on Apple Silicon
- Apple Container 1.4.1 installed and available on `PATH`
- `container system start` has completed
- Go 1.24.2 for source validation

Run focused automated validation:

```bash
go test ./src/services ./src/ui ./cmd/actui
go test ./...
```

## Keyboard Navigation

1. Build and run actui.
2. From the container list, press `k`.
3. Confirm the Kubernetes cluster list appears, including readable empty or
unavailable states when applicable.
4. Select a listed cluster with `enter`; confirm `esc` returns first to the
cluster list, then to the container list.

## Command Safety

1. Use the Kubernetes submenu to open create, load-image, and write-config.
2. Enter valid values and confirm each shows its command preview before `y` or
`enter` can execute it.
3. Run `./actui --dry-run`, repeat those operations, and confirm displayed
commands do not change the local system.
4. Select delete and confirm the exact cluster name must be typed. Decline the
action and verify the cluster remains listed.

The expected command formats and inputs are in
[contracts/kubernetes-commands.md](contracts/kubernetes-commands.md). The
expected parsed cluster/node fields are in [data-model.md](data-model.md).

## Local CLI Compatibility

Run:

```bash
container k8s --help
container k8s list
```

Confirm that the installed CLI reports the six contracted commands and that the
actui list presents available cluster/node state. Test an unavailable Kubernetes
installation or simulated command error and verify the UI provides readable
diagnostics plus a back-navigation route.

## Release Version Alignment

For a local build, run:

```bash
go build -o actui ./cmd/actui
./actui --version
```

Expected output identifies a non-release build (for example, `actui version
dev`). For a published artifact, run the downloaded binary's `--version` and
compare it with the associated GitHub release tag after omitting the tag's `v`
prefix. They must match exactly.

## Manual Validation Record

On macOS 26.x Apple Silicon with Apple Container 1.4.1, `container k8s --help`
confirmed the contracted `create`, `delete`, `list`, `load-image`, `start`, and
`write-config` commands. The non-destructive `container k8s list` check returned
a node record with an empty cluster column, so no cluster could be safely selected
for the remaining preview, dry-run, and delete-cancellation checks. The parser
ignores that malformed row rather than fabricating a cluster name. Published
release-version comparison remains pending the next release artifact.