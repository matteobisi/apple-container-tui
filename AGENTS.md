# AGENTS

This repository keeps AI-focused navigation and change guidance in [docs/ai-menu-map.md](docs/ai-menu-map.md).

Use that file before editing UI flows, adding screens, or wiring new commands.

Priority sources of truth:

- Screen ids and screen names: [src/ui/messages.go](src/ui/messages.go)
- Screen registration and navigation flow: [src/ui/app.go](src/ui/app.go)
- AI menu map and change recipes: [docs/ai-menu-map.md](docs/ai-menu-map.md)

Working rules for agents:

- Do not add a new screen without updating both [src/ui/messages.go](src/ui/messages.go) and [src/ui/app.go](src/ui/app.go)
- Prefer extending an existing workflow when the new behavior clearly belongs to an existing menu branch
- Keep service builders/parsers in [src/services](src/services) aligned with the UI flow that calls them
- Update [docs/ai-menu-map.md](docs/ai-menu-map.md) whenever screens, entry points, or major UI ownership change
- Agents MUST always work on a non-main branch, commit their changes, and open a pull request to `main` for approval
- Agents MUST NOT push directly to `main`; pull-request review plus required checks is the only enabled workflow
## GitHub Actions Node version policy

- Use only Node 24 compatible action versions, pinned by commit SHA with a version comment: `uses: owner/action@<sha> # vX.Y.Z, Node 24 compatible`
- Do NOT use bare version tags (e.g. `@v4`) or any action version that runs on Node 20
- Node 20 support ends June 2, 2026; Node 20 is removed from runners September 16, 2026
- Known Node 24 compatible SHA pins (update as new versions are verified):
  - `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683` (v4.2.2)
  - `actions/upload-artifact@ea165f8d65b6e75b540449e92b4886f43607fa02` (v4.6.2)
  - `actions/download-artifact@cc203385981b70ca67e1cc392babf9cc229d5806` (v4.1.9)
  - `ossf/scorecard-action@4eaacf0543bb3f2c246792bd56e8cdeffafb205a` (v2.4.3)
  - `github/codeql-action/upload-sarif@5c8a8a642e79153f5d047b10ec1cba1d1cc65699` (v3.35.1)
