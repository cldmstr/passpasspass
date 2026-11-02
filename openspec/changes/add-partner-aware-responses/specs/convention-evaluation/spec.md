## MODIFIED Requirements

### Requirement: First-match selection
A convention sheet SHALL hold conventions in order and SHALL return the bid of the first convention
that matches the context and whose bid is legal: either Pass or a bid that outranks the current
contract given in the context. A convention whose bid is not legal SHALL be skipped even if it
matches. A convention MAY bid Pass explicitly, which ends the search. When no legal convention
matches, the sheet SHALL return Pass.

#### Scenario: Earlier convention wins
- **WHEN** two conventions both match, both are legal, and the first one bids `2♦️`
- **THEN** the sheet returns `2♦️`

#### Scenario: Nothing matches
- **WHEN** no convention matches
- **THEN** the sheet returns Pass

#### Scenario: Illegal match skipped
- **WHEN** the contract is `2♠️`, the first matching convention bids `2♥️`, and the next matching convention bids `3♥️`
- **THEN** the sheet returns `3♥️`

#### Scenario: Explicit pass
- **WHEN** the first matching convention bids `Pass` and a later matching convention bids `2♠️`
- **THEN** the sheet returns Pass
