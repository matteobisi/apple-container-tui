# Research: Kubernetes Workflow

## Decision: Implement the locally available Kubernetes command set exactly

The Kubernetes workflow will expose `list`, `create`, `start`, `delete`,
`load-image`, and `write-config`, which are the subcommands reported by the
installed Apple Container 1.4.1 CLI.

**Rationale**: `container k8s --help` and detailed help were run locally on
2026-09-14. The commands and syntax are documented in
[contracts/kubernetes-commands.md](contracts/kubernetes-commands.md). The CLI
does not expose `stop` or `inspect`, so they are excluded rather than simulated
through undocumented behavior.

**Alternatives considered**: Retaining the original broad lifecycle language
would promise operations the local CLI cannot perform. Adding workload or
remote-cluster management would exceed the requested local Container scope.

## Decision: Use a cluster list plus a per-cluster submenu

Lowercase `k` from the container list opens a Kubernetes cluster list. Enter on
a cluster opens its submenu. Create, load-image, and write-config use focused
forms; start uses a command preview; delete uses the existing type-to-confirm
modal.

**Rationale**: This matches the existing `MachineList -> MachineSubmenu ->
form/action` model and preserves keyboard navigation, preview, dry-run, result,
and back-navigation behavior.

**Alternatives considered**: A single menu with no list would force repeated
manual cluster-name entry and hide node state. New modal types are unnecessary
because established preview and destructive confirmation components satisfy
the safety requirements.

## Decision: Parse the Kubernetes list as a tolerant terminal table

The parser will consume `container k8s list` output as a whitespace table,
skip header and blank lines, carry the preceding cluster name into continuation
rows, and retain each cluster's nodes. Unknown or malformed rows are skipped
without inventing state.

**Rationale**: The locally observed output has no JSON/format option, wraps the
`PORTS` header in the current terminal, and leaves the cluster column blank for
following node rows. A tolerant parser and fixtures covering wrapped headers,
multiple nodes, blank cluster cells, empty output, and malformed rows produces
safe display data.

**Alternatives considered**: Passing an unsupported `--format json` flag would
break on the installed CLI. Reliance on exact visual spacing would be brittle.

## Decision: Inject the release version at build time

The build workflow computes the next semantic tag before compiling, passes it
through `go build -ldflags "-X main.version=<tag-without-v>"`, and uploads a
small release-version artifact. The publish workflow downloads that value and
uses it as the release tag rather than recomputing it.

**Rationale**: The existing release workflow currently chooses a tag only after
the binary has already been built. Computing once before compilation makes the
binary and GitHub release share a single source of truth. A default value such
as `dev` identifies local builds without claiming a release version.

**Alternatives considered**: Rebuilding in the publish job duplicates the
build. A repository version file makes manual updates a source of drift.
Requiring pre-merge tags changes the established release workflow more broadly.

## Decision: Amend the constitution's Kubernetes scope boundary

The constitution will permit local Apple Container Kubernetes cluster lifecycle
operations while keeping remote management, cloud integration, and workload
orchestration out of scope.

**Rationale**: The present boundary expressly excludes Kubernetes. A governed
amendment keeps the baseline accurate and makes the feature compliant.

**Alternatives considered**: Treating the feature as an exception would leave a
permanent contradiction between the product scope and governance.