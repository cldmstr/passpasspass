<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# AGENTS.md

Guidance for AI coding agents (Claude Code, Codex, etc.) working in this repository.

- **Project context, conventions, domain glossary**: `openspec/project.md`
- **Current behaviour**, one spec per capability: `openspec/specs/<capability>/spec.md`
- **Planned work**: `openspec/changes/<change-id>/` (proposal, design, tasks, spec deltas)

## Commands

```bash
# Run all tests
go test ./...

# Run a single test
go test -run TestLoader_LoadCard_HCPCondition ./...

# Build and run the CLI (defaults to the embedded Colin card)
go run ./cmd/ppp/main.go

# Run the CLI with a custom convention card
go run ./cmd/ppp/main.go /path/to/data/dir cards/my_card.yaml

# Build the binary
go build -o ppp ./cmd/ppp
```

## Rules

- Constructor functions (`New*`) must not instantiate other structs internally. Pass all
  dependencies in as parameters.
- Trunk-based development: commit directly to `main`. Do not create branches before committing.
- Do not add `Co-Authored-By: Claude` trailers to commit messages.
- When behaviour changes, update the matching spec through a change's spec deltas.
