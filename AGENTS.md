# AGENTS.md

This file provides guidance to AI coding agents (Claude Code, Codex, etc.) when working with code in this repository.

## Project Goal

passpasspass is a contract bridge bidding application for learners to practice correct opening bids across different bidding systems. Players can use standard built-in conventions or define their own by authoring a custom convention card in YAML — selecting from existing sub-conventions or writing new ones from scratch.

## Commands

```bash
# Run all tests
go test ./...

# Run a single test
go test -run TestLoader_LoadCard_HCPCondition ./...

# Build and run the CLI
go run ./cmd/ppp/main.go

# Run CLI with a custom convention card
go run ./cmd/ppp/main.go /path/to/data/dir cards/my_card.yaml

# Build the binary
go build -o ppp ./cmd/ppp
```

## Architecture

The simulator models a single round of contract bridge opening bids. The core flow is:

1. A `Dealer` shuffles and deals four `Hand`s.
2. A `ConventionSheet` evaluates each `Hand` and returns the opening `Bid`.
3. A `Table` binds a `ConventionSheet` to four `Hand`s and produces the display output.

### Convention evaluation

`ConventionSheet` holds an ordered slice of `Convention`s. On `Bid(hand)`, it walks the slice and returns the first convention whose conditions all pass. This makes ordering significant — more specific conventions must come before more general ones.

Each `Convention` wraps a target `Bid` and a slice of `Condition` functions (`func(*Hand) bool`). The conditions are ANDed implicitly. `condition.go` provides the built-in condition constructors:

- `NewHCPIntervalCondition(min, max)` — high-card-point range
- `NewMinMaxDistributionCondition(DistributionRange)` — per-suit length bounds; `DistributionRange` is `[4]SuitRange` indexed `[clubs, diamonds, hearts, spades]`
- `NewBalancedDistributionCondition()` — no void/singleton, at most one doubleton, no 6+ suit
- `NewRuleOfTwentyCondition()` — HCP + two longest suits ≥ 20
- `NewMultiCondition(...)` — OR combinator
- `NewAllCondition(...)` — AND combinator

### YAML convention files

Conventions are defined in YAML and loaded via `Loader`. Two file types:

**Sub-convention file** (`data/conventions/<name>.yaml`):
```yaml
name: strong_notrump
conventions:
  - bid: 1NT
    conditions:
      all:
        - balanced: true
        - hcp: [15, 17]
```

**Convention card file** (`data/cards/<name>.yaml`):
```yaml
name: Two Over One
includes:
  - strong_notrump   # loaded first, in order
conventions:
  - bid: 1C          # card-level conventions appended after includes
    conditions:
      - hcp: [11, 17]
```

The `conditions` field accepts either a sequence (implicit `all`) or a mapping with `all`/`any` keys for explicit logic. Supported condition keys: `hcp: [min, max]`, `balanced: true`, `rule_of_20: true`, `suit_length: {suit: [min, max]}`.

`Loader` is constructed with an `fs.FS` — in production the embedded `DefaultData` FS is sub-rooted at `data/`; in tests a `t.TempDir()` FS is used. Include paths are always resolved as `conventions/<name>.yaml` within the loader's FS root.

### Suit and rank conventions

Suits are ordered numerically throughout: `Clubs(0) < Diamonds(1) < Hearts(2) < Spades(3) < NoTrump(4)`. `DistributionRange` and `Distribution` arrays are always indexed in this order. Suit strings accepted: emoji (`♣️ ♦️ ❤️ ♠️`), letter (`C D H S NT`), or full name (`clubs diamonds hearts spades`).

Cards are represented as `<rank><suit-emoji>` strings (e.g. `"A♣️"`, `"10❤️"`). The unicode suit emoji occupies 6 bytes, which `NewCardFromString` accounts for when splitting rank from suit.

### Testing

Uses `github.com/matryer/is` for assertions. All tests live in `package passpasspass_test` (external test package). Loader tests write temp YAML files via `writeTempFile` and pass `os.DirFS(dir)` to `NewLoader`.

## Design Conventions

Constructor functions (`New*`) must not instantiate other structs internally. All dependencies must be passed in as parameters — use dependency injection.

## Git

Do not add `Co-Authored-By: Claude` trailers to commit messages.
