# convention-evaluation Specification

## Purpose
Represent bids, and choose a seat's call by evaluating an ordered list of conventions against the
current auction context.

## Requirements

### Requirement: Bid notation
A bid SHALL be written as a level from 1 to 7 followed by a suit in any accepted suit spelling
(e.g. `1NT`, `2D`, `2♦️`), or as `Pass` / `pass`. Any other string SHALL be rejected with an
error. A bid SHALL print as the level followed by the suit emoji or `NT`; Pass SHALL print as
`Pass`. Every Pass SHALL be the same shared value.

#### Scenario: Letter and emoji spellings
- **WHEN** parsing `2D` and `2♦️`
- **THEN** both yield the same bid, printed as `2♦️`

#### Scenario: Level out of range
- **WHEN** parsing `8C`
- **THEN** an error is returned

### Requirement: Bid ranking
Bids SHALL rank by level first and by suit order within a level, so that `1NT` outranks `1♠️`
and `2♣️` outranks `1NT`.

#### Scenario: Notrump outranks spades
- **WHEN** comparing `1NT` with `1♠️`
- **THEN** `1NT` ranks higher

### Requirement: Convention matching
A convention SHALL pair a bid with a list of conditions, and SHALL match a context only when every
condition passes.

#### Scenario: One condition fails
- **WHEN** a convention requires 15–17 HCP and balanced shape, and the hand is balanced with 14 HCP
- **THEN** the convention does not match

### Requirement: First-match selection
A convention sheet SHALL hold conventions in order and SHALL return the bid of the first
convention that matches the context. When no convention matches, it SHALL return Pass.

#### Scenario: Earlier convention wins
- **WHEN** two conventions both match and the first one bids `2♦️`
- **THEN** the sheet returns `2♦️`

#### Scenario: Nothing matches
- **WHEN** no convention matches
- **THEN** the sheet returns Pass
