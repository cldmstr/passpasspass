## ADDED Requirements

### Requirement: Range values
Every condition whose value is a range SHALL take exactly two integers `[min, max]` with
`min ≤ max`, both within the limits of what the condition measures:

- HCP of one hand or of a partnership (`hcp`, `partnership_hcp`): 0–40;
- length of one suit in one hand (`suit_length`): 0–13.

Any other value SHALL be a load error.

#### Scenario: One value
- **WHEN** a condition is `hcp: [5]`
- **THEN** loading fails

#### Scenario: Reversed range
- **WHEN** a condition is `hcp: [12, 10]`
- **THEN** loading fails

#### Scenario: Out of limits
- **WHEN** a condition is `suit_length: {hearts: [0, 14]}`
- **THEN** loading fails

#### Scenario: No upper limit
- **WHEN** a condition is `hcp: [27, 40]`
- **THEN** loading succeeds and the condition passes for a hand with 30 HCP

### Requirement: Flag values
`balanced`, `rule_of_20` and `no_prior_bids` SHALL accept only the value `true`. Any other value,
including `false`, SHALL be a load error stating that only `true` is accepted.

#### Scenario: False flag
- **WHEN** a condition is `balanced: false`
- **THEN** loading fails with an error stating that only `true` is accepted

## MODIFIED Requirements

### Requirement: Suit length condition
`suit_length: {<suit>: [min, max], ...}` SHALL pass when, for every listed suit, the bidder holds
between `min` and `max` cards inclusive. Unlisted suits SHALL be unconstrained. A suit SHALL be
clubs, diamonds, hearts or spades, in any spelling accepted for suits; notrump or an unknown suit
SHALL be a load error. Ranges SHALL follow the range values requirement.

#### Scenario: Six-card heart suit
- **WHEN** the condition is `suit_length: {hearts: [6, 13]}` and the hand holds 6 hearts
- **THEN** the condition passes

#### Scenario: Notrump is not a suit
- **WHEN** the condition is `suit_length: {NT: [1, 2]}`
- **THEN** loading fails

### Requirement: Seat bid conditions
`partner_bid` SHALL pass when the most recent call by the bidder's partner equals the given bid.
A partner who has not yet called SHALL count as having passed. The given bid SHALL be parsed when the file is
loaded, and an invalid bid SHALL be a load error.

#### Scenario: Partner opened 1NT
- **WHEN** the condition is `partner_bid: "1NT"` and partner's last call was `1NT`
- **THEN** the condition passes

#### Scenario: Invalid bid
- **WHEN** the condition is `partner_bid: "1N"`
- **THEN** loading fails with an error naming `1N`

## REMOVED Requirements

### Requirement: Suit fit conditions
**Reason**: `heart_fit`, `spade_fit`, `major_fit` and `minor_fit` are used by no convention file.
**Migration**: none needed for the built-in data. `add-partner-aware-responses` adds
`fit: {<suit>: [min, max]}` for any suit; a major or minor fit is `any` of two `fit` conditions.

### Requirement: Opponent strength conditions
**Reason**: `lho_min_hcp` and `rho_min_hcp` are used by no convention file, and conventions that
depend on opponents' calls are out of scope until overcalls are planned.
**Migration**: none needed for the built-in data.
