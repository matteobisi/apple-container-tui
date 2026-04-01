# Data Model: Automated Binary Build Workflow

## Entity: BuildTrigger

- Description: Represents the repository event that starts a binary build workflow run.
- Fields:
  - `trigger_type` (enum): `push_main` | `manual_dispatch`
  - `source_ref` (string): branch or ref associated with the run
  - `commit_sha` (string): revision under build
  - `initiator` (string): actor that caused the run
- Validation Rules:
  - `commit_sha` must be present for all runs.
  - `trigger_type` must map to approved workflow triggers.
- Relationships:
  - One BuildTrigger creates one or more BuildRun records over time.

## Entity: BuildRun

- Description: Execution record for a single automated binary build.
- Fields:
  - `run_id` (string)
  - `started_at` (datetime)
  - `completed_at` (datetime)
  - `status` (enum): `success` | `failure` | `cancelled`
  - `log_url` (string)
  - `summary` (string)
- Validation Rules:
  - `completed_at` must be greater than or equal to `started_at` when present.
  - `status` must be one of the defined enum values.
- Relationships:
  - Each BuildRun is associated to exactly one BuildTrigger.
  - Each BuildRun may produce one or more BuildArtifact records.

## Entity: BuildArtifact

- Description: Packaged binary output published from a successful build run.
- Fields:
  - `artifact_name` (string)
  - `platform` (string): e.g., `darwin-arm64`
  - `produced_from_run_id` (string)
  - `retention_window_days` (integer)
  - `download_reference` (string)
- Validation Rules:
  - `artifact_name` and `platform` are required.
  - Artifact records should only exist for successful runs.
- Relationships:
  - Many BuildArtifact records can belong to one BuildRun.

## Entity: BuildAutomationDoc

- Description: Maintainer-facing operational document for the build workflow.
- Fields:
  - `doc_path` (string)
  - `workflow_name` (string)
  - `trigger_description` (string)
  - `validation_steps` (list of strings)
  - `troubleshooting_steps` (list of strings)
  - `last_reviewed_on` (date)
- Validation Rules:
  - Must include trigger behavior, artifact expectations, and acceptance guidance.
  - Must include reference validation environment details.
- Relationships:
  - References ValidationEnvironmentProfile as contextual metadata.

## Entity: ValidationEnvironmentProfile

- Description: Reference machine configuration used for build verification context.
- Fields:
  - `machine_model` (string): `Macbook M4`
  - `os_version` (string): `macOS 26.4`
  - `memory` (string): `32GB`
  - `apple_container_version` (string): `0.10.0`
- Validation Rules:
  - All fields are mandatory in documentation.
- Relationships:
  - Linked from BuildAutomationDoc as the canonical validation baseline.