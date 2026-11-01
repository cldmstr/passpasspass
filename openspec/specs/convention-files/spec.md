# convention-files Specification

## Purpose
Load convention cards and sub-conventions from YAML into a ready-to-use convention sheet.

## Requirements

### Requirement: Sub-convention files
A sub-convention file SHALL live at `conventions/<name>.yaml` and SHALL contain a `name` and a list
of `conventions`, each with a `bid` and `conditions`.

#### Scenario: Strong notrump sub-convention
- **WHEN** `conventions/strong_notrump.yaml` defines `1NT` with `balanced: true` and `hcp: [15, 17]`
- **THEN** including `strong_notrump` adds that convention to the sheet

### Requirement: Convention card files
A card file SHALL contain a `name`, an optional list of `includes` naming sub-conventions, and an
optional list of card-level `conventions`. The loader SHALL add the conventions of each include in
the listed order, followed by the card-level conventions.

#### Scenario: Include order
- **WHEN** a card includes `a` then `b` and defines one card-level convention
- **THEN** the sheet holds `a`'s conventions, then `b`'s, then the card-level one

### Requirement: Include resolution
Includes SHALL resolve to `conventions/<name>.yaml` relative to the root of the loader's file
system. In production the root is the embedded data sub-rooted at `data/`; tests may supply any
`fs.FS`.

#### Scenario: Missing include
- **WHEN** a card includes a sub-convention that does not exist
- **THEN** loading fails with an error naming the sub-convention

### Requirement: Conditions syntax
`conditions` SHALL accept either a sequence, treated as an implicit `all`, or a mapping holding a
single condition key such as `all` or `any`.

#### Scenario: Sequence form
- **WHEN** a convention lists `- hcp: [11, 17]` and `- rule_of_20: true` as a sequence
- **THEN** both conditions must pass for the convention to match

### Requirement: Load errors
Loading SHALL fail with an error for an unreadable or unparsable file, an invalid bid string, an
unknown suit name, a range without exactly two values, or an empty condition node.

#### Scenario: Invalid bid
- **WHEN** a convention has `bid: 9Z`
- **THEN** loading fails with an error naming the bid
