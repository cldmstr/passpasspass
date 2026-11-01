## MODIFIED Requirements

### Requirement: Sub-convention files
A sub-convention file SHALL live at `conventions/<name>.yaml` and SHALL contain a `name` and a list
of `conventions`, each with a `bid` and `conditions`. A convention MAY also carry `establishes`, a
mapping with an optional `agreed` key (`clubs`, `diamonds`, `hearts`, `spades` or `notrump`) and
an optional `forcing` key (`one_round` or `game`), naming the facts its call establishes for the
caller's partnership. A convention without `establishes` SHALL establish nothing.

#### Scenario: Strong notrump sub-convention
- **WHEN** `conventions/strong_notrump.yaml` defines `1NT` with `balanced: true` and `hcp: [15, 17]`
- **THEN** including `strong_notrump` adds that convention to the sheet

#### Scenario: Declared facts
- **WHEN** a convention defines `bid: 2NT` with `establishes: {agreed: hearts, forcing: game}`
- **THEN** the loaded convention establishes hearts as the agreed strain and a game force

### Requirement: Load errors
Loading SHALL fail with an error for an unreadable or unparsable file, an invalid bid string, an
unknown suit name, a range without exactly two values, an empty condition node, an `auction`
element that is neither a bid nor a wildcard, a `not` without a child, or an `establishes` entry
with an unknown key, strain or forcing value.

#### Scenario: Invalid bid
- **WHEN** a convention has `bid: 9Z`
- **THEN** loading fails with an error naming the bid

#### Scenario: Unknown forcing value
- **WHEN** a convention has `establishes: {forcing: forever}`
- **THEN** loading fails with an error naming the value
