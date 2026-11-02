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
Loading SHALL fail with an error for:

- an unreadable, unparsable or empty file;
- an unknown key in a card file, a sub-convention file, a convention entry, an `establishes`
  mapping or a condition;
- a key that appears twice in the same mapping, including a suit in `suit_length`;
- a convention entry without `bid` or without `conditions`;
- an invalid bid string, in `bid` or in a condition;
- an `auction` element that is neither a bid nor a wildcard;
- an empty condition: an empty mapping, an empty sequence, an empty `all` or `any`, or a `not`
  without a value;
- an unknown strain or forcing value in `establishes`;
- a condition value of the wrong shape or outside its limits, as defined by bidding-conditions.

Every error SHALL start with the path of the file (card or sub-convention) and the line of the
offending key or value, written `file:line:`. An error inside a convention SHALL then name the
convention's bid, and an error inside nested conditions SHALL name the keys leading to the
offending one, joined with dots.

#### Scenario: Invalid bid
- **WHEN** a convention has `bid: 9Z`
- **THEN** loading fails with an error naming the bid

#### Scenario: Misspelt condition key
- **WHEN** a condition mapping holds `hcp: [0, 10]` and `suit_lenght: {spades: [6, 13]}`
- **THEN** loading fails with an error naming `suit_lenght` and its line

#### Scenario: Misspelt convention key
- **WHEN** a convention entry has `condtions:` instead of `conditions:`
- **THEN** loading fails with an error naming `condtions`

#### Scenario: Repeated key
- **WHEN** a condition mapping holds `hcp: [10, 12]` and, below it, `hcp: [15, 17]`
- **THEN** loading fails with an error naming `hcp` and the line of the second key

#### Scenario: Empty conditions
- **WHEN** a convention has `conditions: []`
- **THEN** loading fails

#### Scenario: Error location
- **WHEN** the sub-convention `baron` has `suit_length: {hearts: [0, 14]}` inside `all` and `any` in its `2S` convention, with `hearts` on line 8
- **THEN** the error reads `conventions/baron.yaml:8: bid 2S: all.any.suit_length.hearts: [0, 14] is outside 0 to 13`

#### Scenario: Unknown forcing value
- **WHEN** a convention has `establishes: {forcing: forever}`
- **THEN** loading fails with an error naming the value
