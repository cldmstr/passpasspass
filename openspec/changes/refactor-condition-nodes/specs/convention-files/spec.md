## MODIFIED Requirements

### Requirement: Conditions syntax
`conditions` SHALL accept either a sequence or a mapping, and so SHALL every nested condition. A
sequence SHALL pass when every item passes. A mapping SHALL pass when the condition of every key in
it passes; each key SHALL be read as its own condition, whatever other keys the mapping holds. An
empty sequence or an empty mapping SHALL be a load error.

#### Scenario: Sequence form
- **WHEN** a convention lists `- hcp: [11, 17]` and `- rule_of_20: true` as a sequence
- **THEN** both conditions must pass for the convention to match

#### Scenario: Several keys in one mapping
- **WHEN** a convention's conditions are the mapping `{hcp: [0, 10], suit_length: {spades: [6, 13]}}` and the hand holds 4 HCP and 3 spades
- **THEN** the convention does not match, because the spade condition fails

#### Scenario: Several keys inside a combinator
- **WHEN** a condition is `any` of `{hcp: [6, 10], suit_length: {hearts: [6, 13]}}` and `{hcp: [6, 10], suit_length: {spades: [6, 13]}}`, and the hand holds 8 HCP and 6 spades
- **THEN** the condition passes

### Requirement: Load errors
Loading SHALL fail with an error for:

- an unreadable or unparsable file;
- an unknown key in a card file, a sub-convention file, a convention entry or a condition;
- a convention entry without `bid` or without `conditions`;
- an invalid bid string, in `bid` or in a condition;
- an empty condition: an empty mapping, an empty sequence, or an empty `all` or `any`;
- a condition value of the wrong shape or outside its limits, as defined by bidding-conditions.

Every error SHALL name the file (card or sub-convention), the convention's bid when the error is
inside a convention, and the line of the offending key or value.

#### Scenario: Invalid bid
- **WHEN** a convention has `bid: 9Z`
- **THEN** loading fails with an error naming the bid

#### Scenario: Misspelt condition key
- **WHEN** a condition mapping holds `hcp: [0, 10]` and `suit_lenght: {spades: [6, 13]}`
- **THEN** loading fails with an error naming `suit_lenght` and its line

#### Scenario: Misspelt convention key
- **WHEN** a convention entry has `condtions:` instead of `conditions:`
- **THEN** loading fails with an error naming `condtions`

#### Scenario: Empty conditions
- **WHEN** a convention has `conditions: []`
- **THEN** loading fails

#### Scenario: Error location
- **WHEN** the sub-convention `baron` has an invalid condition in its `2S` convention on line 8
- **THEN** the error names `baron`, `2S` and line 8
