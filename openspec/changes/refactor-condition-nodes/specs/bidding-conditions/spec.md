## ADDED Requirements

### Requirement: Range values
Every condition whose value is a range SHALL take exactly two integers `[min, max]` with
`min ≤ max`, both within the limits of what the condition measures:

- HCP of one hand or of a partnership (`hcp`, `partnership_hcp`): 0–40;
- length of one suit in one hand (`suit_length`): 0–13;
- a fit, the combined length of one suit in a partnership (`heart_fit`, `spade_fit`,
  `major_fit`, `minor_fit`): 8–13. A suit has 13 cards, and a fit is at least 8 of them.

A single threshold (`lho_min_hcp`, `rho_min_hcp`) SHALL be an integer from 0 to 40. Any other
value SHALL be a load error.

#### Scenario: One value
- **WHEN** a condition is `hcp: [5]`
- **THEN** loading fails

#### Scenario: Reversed range
- **WHEN** a condition is `hcp: [12, 10]`
- **THEN** loading fails

#### Scenario: Out of limits
- **WHEN** a condition is `suit_length: {hearts: [0, 14]}`
- **THEN** loading fails

#### Scenario: Fit beyond the suit
- **WHEN** a condition is `heart_fit: [8, 14]`
- **THEN** loading fails, because the partnership cannot hold more than 13 hearts

#### Scenario: Fewer than eight cards is not a fit
- **WHEN** a condition is `spade_fit: [7, 13]`
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

### Requirement: Suit fit conditions
`heart_fit` and `spade_fit` SHALL pass when the combined card range in that suit (bidder's length
plus partner's inferred minimum and maximum length) lies entirely within `[min, max]`.
`major_fit` SHALL apply the same test to the better major, and `minor_fit` to the better minor,
where the better suit is the one with the higher combined maximum (spades or diamonds on a tie).
Ranges SHALL follow the range values requirement for a fit (8–13).

#### Scenario: Known eight-card heart fit
- **WHEN** the bidder holds 3 hearts, partner is inferred at 5–10 hearts, and the condition is `heart_fit: [8, 13]`
- **THEN** the condition passes

### Requirement: Seat bid conditions
`partner_bid`, `rho_bid` and `lho_bid` SHALL each pass when the most recent call by the bidder's
partner, right-hand opponent or left-hand opponent respectively equals the given bid. A seat that
has not yet called SHALL count as having passed. The given bid SHALL be parsed when the file is
loaded, and an invalid bid SHALL be a load error.

#### Scenario: Partner opened 1NT
- **WHEN** the condition is `partner_bid: "1NT"` and partner's last call was `1NT`
- **THEN** the condition passes

#### Scenario: Invalid bid
- **WHEN** the condition is `partner_bid: "1N"`
- **THEN** loading fails with an error naming `1N`
