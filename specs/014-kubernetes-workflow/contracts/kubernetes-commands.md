# Kubernetes Command Contract

Source: local `container k8s --help` and `container k8s help <subcommand>` on
Apple Container 1.4.1, captured 2026-09-14.

| UI action | Command | Input and validation | Safety behavior |
|-----------|---------|----------------------|-----------------|
| List clusters | `container k8s list` | None | Read-only; execute directly |
| Create cluster | `container k8s create [--name <name>] [--rm] [--cpus <cpus>] [--memory <memory>] [--node-image <image>]` | Name token optional; CPUs positive; memory uses CLI-supported units; image non-blank when present | Preview, then `y`/`enter` approval; dry-run supported |
| Start cluster | `container k8s start [--name <name>]` | Selected cluster name required | Preview, then `y`/`enter` approval; dry-run supported |
| Delete cluster | `container k8s delete [--name <name>]` | Selected cluster name required | Type cluster name to confirm; dry-run supported |
| Load image | `container k8s load-image [--name <name>] <image> [--platform <platform>]` | Image required; name and platform optional | Preview, then `y`/`enter` approval; dry-run supported |
| Write kubeconfig | `container k8s write-config [--name <name>] [--kubeconfig <path>]` | Selected cluster name required; alternate path optional | Preview, then `y`/`enter` approval; dry-run supported |

The CLI also accepts create-only `--scheme` and `--max-concurrent-downloads`.
They remain at the CLI defaults for this feature. `delete` has alias `rm` and
`list` has alias `ls`; actui uses the canonical command spellings above.

`container k8s list` has no documented JSON or formatting option. Its output
is a table containing `CLUSTER`, `NODE`, `ROLE`, `STATE`, `CPUS`, `MEMORY`,
`ADDR`, and `PORTS`; rows belonging to the same cluster can omit repeated
cluster names. The UI parser must retain cluster context across those rows.