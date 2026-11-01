# Project Context

## Purpose
passpasspass is a contract bridge bidding trainer. Learners practise bidding with a convention
card: either one of the built-in cards or their own, written in YAML by combining existing
sub-conventions or writing new ones from scratch.

The simulator deals four hands and runs an auction in which every seat bids by the same
convention card, then prints the deal, the bidding and the final contract.

## Tech Stack
- Go 1.25 (`go.mod`)
- `gopkg.in/yaml.v3` for convention files
- `github.com/matryer/is` for test assertions
- `embed` to ship the built-in convention data (`DefaultData`, rooted at `data/`)

## Project Conventions

### Code Style
- One package, `passpasspass`, at the repository root; the CLI lives in `cmd/ppp`.
- Exported identifiers carry doc comments.
- Per-suit arrays (`Distribution`, `DistributionRange`, `HandRange.SuitLength`) are always
  indexed in suit order: clubs, diamonds, hearts, spades.

### Architecture Patterns
- Dependency injection: constructor functions (`New*`) must not instantiate other structs
  internally. All dependencies are passed in as parameters.
- Bidding rules are data, not code. Go code provides condition primitives; YAML files compose
  them into conventions.
- Conventions are evaluated first-match in sheet order, so more specific conventions must come
  before more general ones.
- `Table` is a thin wiring layer; auction logic lives in `Auction`.

Main types and their files:

| Type | File | Role |
|------|------|------|
| `Card`, `Suit`, `Rank` | `card.go` | Card model and parsing |
| `Hand`, `Distribution` | `hand.go` | 13-card hand, HCP, suit lengths |
| `Dealer` | `dealer.go` | Shuffles and deals four hands |
| `Bid`, `Convention`, `ConventionSheet` | `convention.go` | Bids and first-match evaluation |
| `Condition` constructors | `condition.go` | Predicates over an `AuctionContext` |
| `Loader` | `loader.go` | YAML card and sub-convention loading |
| `Auction`, `AuctionContext`, `HandRange` | `auction.go` | Bidding loop and hand inference |
| `Table` | `table.go` | Wires dealer, sheet and auction |

### Testing Strategy
- All tests live in the external package `passpasspass_test`.
- Assertions use `github.com/matryer/is`.
- Loader tests write temporary YAML with `writeTempFile` into `t.TempDir()` and pass
  `os.DirFS(dir)` to `NewLoader`.
- Round-trip tests load the embedded cards and assert the bid chosen for a given hand. Comments
  on test hands show how the HCP total adds up, e.g. `// A♠️(4)+Q♠️(2) = 6 HCP`.
- Boundary values of every range (lowest, highest, one outside) get their own test case.

### Git Workflow
- Main branch: `main`.
- Do not add `Co-Authored-By: Claude` trailers to commit messages.

## Domain Context
- **Seats**: North, East, South, West, bidding clockwise. Partner sits opposite; LHO is the next
  seat to bid after you, RHO the seat that bid just before you.
- **HCP** (high-card points): A=4, K=3, Q=2, J=1. One hand holds at most 37; the deck holds 40.
- **Suit order**: Clubs(0) < Diamonds(1) < Hearts(2) < Spades(3) < NoTrump(4). Bids rank by
  level, then by suit.
- **Card strings**: `<rank><suit emoji>`, e.g. `A♣️`, `10❤️`. A suit emoji is 6 bytes, which
  `NewCardFromString` relies on when splitting rank from suit.
- **Suit strings**: emoji (`♣️ ♦️ ❤️ ♠️`), letter (`C D H S NT`) or name
  (`clubs diamonds hearts spades`).
- **Balanced**: no void or singleton, at most one doubleton, no suit of 6+ cards.
- **Rule of 20**: HCP plus the lengths of the two longest suits is at least 20.
- **Opening / response / rebid**: the first non-pass bid; partner's first call after it; the
  opener's second call.

## Important Constraints
- Every seat bids by the same convention sheet; there are no overcall conventions yet, so
  opponents pass after an opening.
- A convention card is a flat ordered list: includes first (in listed order), then the card's own
  conventions.

## External Dependencies
None at runtime. The built-in convention data is embedded in the binary.
