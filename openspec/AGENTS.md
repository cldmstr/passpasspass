# OpenSpec Instructions

Instructions for AI coding assistants using OpenSpec for spec-driven development in this
repository.

## Layout

```
openspec/
├── project.md                    # project context and conventions
├── specs/                        # current truth: what IS built
│   └── <capability>/spec.md
└── changes/                      # proposals: what SHOULD change
    ├── <change-id>/
    │   ├── proposal.md           # why, what, impact
    │   ├── design.md             # technical decisions (when needed)
    │   ├── tasks.md              # implementation checklist
    │   └── specs/<capability>/spec.md   # deltas against specs/
    └── archive/YYYY-MM-DD-<change-id>/  # completed changes
```

Current capabilities: `cards-and-hands`, `dealing`, `convention-evaluation`,
`bidding-conditions`, `convention-files`, `auction`, `standard-conventions`, `cli`.

## Three stages

### 1. Creating a change
Create a change for new features, breaking changes, architecture shifts, or new conventions in
the built-in data. Skip it for bug fixes that restore specified behaviour, typos, formatting,
and test-only edits.

1. Read `project.md`, the relevant `specs/`, and any active `changes/` to avoid conflicts.
2. Pick a kebab-case, verb-led change id (`add-`, `update-`, `remove-`, `refactor-`).
3. Write `proposal.md`, `tasks.md`, `design.md` when there are cross-cutting technical
   decisions, and one delta per affected capability under `specs/`.
4. Do not start implementation until the proposal is approved.

### 2. Implementing a change
1. Read `proposal.md`, `design.md` and `tasks.md`.
2. Work through the tasks in order.
3. Tick a task (`- [x]`) only once it is implemented and its tests pass.

### 3. Archiving a change
After the change is merged, apply its deltas to `specs/` and move the folder to
`changes/archive/YYYY-MM-DD-<change-id>/`.

## Spec format

```markdown
# <capability> Specification

## Purpose
<one or two sentences>

## Requirements

### Requirement: <name>
The system SHALL <behaviour>.

#### Scenario: <name>
- **WHEN** <condition>
- **THEN** <outcome>
```

- Every requirement has at least one scenario.
- Scenario headers use exactly four `#`.
- Normative statements use SHALL / MUST.

## Delta format

Deltas in `changes/<id>/specs/<capability>/spec.md` use these sections:

- `## ADDED Requirements`: new, standalone requirements.
- `## MODIFIED Requirements`: the **complete** updated requirement, with its header copied
  exactly from the current spec. Partial text loses detail at archive time.
- `## REMOVED Requirements`: header plus `**Reason**` and `**Migration**`.
- `## RENAMED Requirements`: `- FROM: ### Requirement: Old` / `- TO: ### Requirement: New`.

## CLI

If the `openspec` CLI is installed, `openspec list`, `openspec show <id>`,
`openspec validate <id> --strict` and `openspec archive <id>` automate the steps above.
Without it, follow the same layout by hand.
