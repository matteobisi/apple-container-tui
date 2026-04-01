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