# Data Model: Automated Release Publishing

## Entity: QualifiedBuildRun

- Description: A completed build workflow run eligible for release publication.
- Fields:
  - `run_id` (string)
  - `workflow_name` (string): expected `Build Binary`
  - `conclusion` (enum): `success` | `failure` | `cancelled`
  - `head_sha` (string)
  - `event` (string): source trigger category
  - `artifact_name` (string): expected `actui-linux-amd64`
  - `completed_at` (datetime)
- Validation Rules:
  - `conclusion` must be `success` to qualify for release publication.
  - `artifact_name` must exist and match configured release asset source.
- Relationships:
  - One QualifiedBuildRun can produce at most one ReleasePublication.

## Entity: VersionLabelPolicy

- Description: Ruleset controlling release tag generation for automated publication.
- Fields:
  - `starting_version` (string): `0.1.0`
  - `prefix` (string): `v`
  - `increment_strategy` (enum): `patch`
  - `conflict_strategy` (enum): `skip_if_exists` | `find_next_patch`
- Validation Rules:
  - Version must follow semantic version format `MAJOR.MINOR.PATCH`.
  - Prefix policy must be consistent across all automated releases.
- Relationships:
  - One policy governs many ReleasePublication records.

## Entity: ReleasePublication

- Description: Published GitHub Release created from one qualified build.
- Fields:
  - `release_id` (string)
  - `tag_name` (string): e.g., `v0.1.3`
  - `title` (string)
  - `target_commitish` (string): source commit SHA
  - `published_at` (datetime)
  - `status` (enum): `published` | `skipped` | `failed`
  - `run_log_url` (string)
- Validation Rules:
  - `tag_name` must be unique in repository scope.
  - `status=published` requires attached release asset and commit reference.
- Relationships:
  - Each ReleasePublication references exactly one QualifiedBuildRun.
  - Each ReleasePublication has one or more ReleaseAsset entries.

## Entity: ReleaseAsset

- Description: Binary asset attached to a published release.
- Fields:
  - `asset_name` (string): e.g., `actui-linux-amd64`
  - `content_type` (string): binary media type
  - `size_bytes` (integer)
  - `download_url` (string)
  - `source_artifact_run_id` (string)
- Validation Rules:
  - `asset_name` must match artifact contract for the release run.
  - `size_bytes` must be greater than zero.
- Relationships:
  - Many ReleaseAsset records can belong to one ReleasePublication.

## Entity: ReleaseAutomationDocumentation

- Description: Maintainer runbook record describing operational release behavior.
- Fields:
  - `doc_path` (string): `docs/binary-build-automation.md`
  - `trigger_chain` (string): build workflow to release workflow linkage
  - `version_policy_summary` (string)
  - `failure_modes` (list of strings)
  - `last_reviewed_on` (date)
- Validation Rules:
  - Must include trigger behavior, versioning rules, idempotency behavior, and troubleshooting steps.
- Relationships:
  - References VersionLabelPolicy and ReleasePublication process semantics.
