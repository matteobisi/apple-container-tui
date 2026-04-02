# Research: Automated Release Publishing After Binary Build

## Decision 1: Trigger release publication only after successful binary build completion

- Decision: Release publication automation will run only when the binary build workflow completes successfully and produces the expected artifact.
- Rationale: Prevents publishing broken or incomplete releases and keeps release entries aligned with validated build output.
- Alternatives considered:
  - Trigger directly on `push` to `main`: rejected because release flow would be decoupled from actual build success.
  - Manual-only release publication: rejected because it does not satisfy automation goals and introduces operator drift.

## Decision 2: Use workflow chaining via `workflow_run` for build-to-release orchestration

- Decision: Implement a release workflow that is triggered by `workflow_run` completion of `Build Binary` and gated on `conclusion == success`.
- Rationale: Cleanly models dependency between workflows and avoids race conditions around artifact availability.
- Alternatives considered:
  - Add release steps in the existing build workflow: viable, but rejected for this feature because separation improves observability and rollback isolation.
  - External webhook/service orchestration: rejected as out of scope and adds unnecessary operational complexity.

## Decision 3: Adopt deterministic semantic version labels starting at `0.1.0`

- Decision: First automated release label is `v0.1.0`; subsequent automated releases increment patch (`v0.1.1`, `v0.1.2`, ...).
- Rationale: Matches user requirement for predictable progression and keeps initial policy simple, audit-friendly, and deterministic.
- Alternatives considered:
  - Date-based versioning: rejected due to weaker semantic meaning for compatibility.
  - Commit-SHA-only tagging: rejected because it is less human-readable for release consumers.

## Decision 4: Enforce idempotency to avoid duplicate releases per source commit

- Decision: Release automation must check if a release already exists for the candidate tag/source reference and skip publication on duplicates.
- Rationale: Prevents reruns from creating duplicate release records and satisfies FR-009.
- Alternatives considered:
  - Allow duplicate prereleases on rerun: rejected because it introduces ambiguity and violates measurable outcomes.

## Decision 5: Keep permissions explicit and minimal for release publication

- Decision: Release workflow permissions include only required scopes (`contents: write` for release creation; read-only scopes otherwise).
- Rationale: Aligns with repository security posture and least-privilege principle.
- Alternatives considered:
  - Default broad token permissions: rejected due to avoidable security exposure.

## Decision 6: Keep artifact contract stable and attach build artifact unchanged

- Decision: Release asset will reuse existing artifact naming convention (`actui-linux-amd64`) and include provenance metadata in release notes.
- Rationale: Preserves compatibility with existing build contract and documentation.
- Alternatives considered:
  - Renaming artifacts only for release stage: rejected to avoid confusion between CI artifacts and release assets.

## Decision 7: Document the end-to-end build-to-release process in existing runbook

- Decision: Extend `docs/binary-build-automation.md` with release flow, version-labeling policy, failure triage, and operator checklist, and add/adjust top-level references where needed.
- Rationale: Keeps one canonical operations document and satisfies requested docs updates.
- Alternatives considered:
  - New dedicated docs file: rejected because existing runbook already owns build automation lifecycle.

## Decision 8: Require Node 24-compatible GitHub Actions in release workflow

- Decision: Release automation must use actions verified as Node 24 compatible and avoid introducing Node 20-only actions.
- Rationale: Repository operations already encountered Node 20 deprecation warnings, and future runner defaults/removals require Node 24 readiness.
- Alternatives considered:
  - Keep existing Node 20-only actions until hard cutoff: rejected because it creates avoidable migration risk and warning noise.
