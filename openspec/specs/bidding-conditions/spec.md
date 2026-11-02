# bidding-conditions Specification

## Purpose
Define the predicates a convention can require, each with its YAML key and meaning. Conditions read
the bidder's own hand, the auction history, or the hand ranges inferred from other players' bids.

## Requirements

### Requirement: HCP condition
`hcp: [min, max]` SHALL pass when the bidder's HCP lies within `min`..`max` inclusive.

#### Scenario: Boundary values
- **WHEN** the condition is `hcp: [6, 10]`
- **THEN** it passes for 6 and 10 HCP and fails for 5 and 11 HCP

### Requirement: Balanced condition
`balanced: true` SHALL pass when the bidder's hand has no void or singleton, at most one doubleton,
and no suit of six or more cards.

#### Scenario: Two doubletons
- **WHEN** the hand is 5-4-2-2
- **THEN** the condition fails

#### Scenario: Semi-balanced
- **WHEN** the hand is 5-3-3-2
- **THEN** the condition passes

### Requirement: Rule of 20 condition
`rule_of_20: true` SHALL pass when the bidder's HCP plus the lengths of their two longest suits is
at least 20.

#### Scenario: Exactly twenty
- **WHEN** the hand has 11 HCP and 5-4-2-2 shape
- **THEN** the condition passes

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

### Requirement: Combinators
`all: [...]` SHALL pass when every child passes. `any: [...]` SHALL pass when at least one child
passes. Combinators SHALL nest.

#### Scenario: Any of two shapes
- **WHEN** the condition is `any` of 6+ hearts or 6+ spades, and the hand holds 6 spades and 2 hearts
- **THEN** the condition passes

### Requirement: No prior bids condition
`no_prior_bids: true` SHALL pass when every call made so far in the auction is a Pass.

#### Scenario: Third seat after two passes
- **WHEN** the two calls so far are both Pass
- **THEN** the condition passes

#### Scenario: After an opening
- **WHEN** an earlier seat has bid `1♣️`
- **THEN** the condition fails

### Requirement: Seat bid conditions
`partner_bid` SHALL pass when the most recent call by the bidder's partner equals the given bid.
A partner who has not yet called SHALL count as having passed. The given bid SHALL be parsed when
the file is loaded, and an invalid bid SHALL be a load error.

#### Scenario: Partner opened 1NT
- **WHEN** the condition is `partner_bid: "1NT"` and partner's last call was `1NT`
- **THEN** the condition passes

#### Scenario: Invalid bid
- **WHEN** the condition is `partner_bid: "1N"`
- **THEN** loading fails with an error naming `1N`

### Requirement: Partnership HCP condition
`partnership_hcp: [min, max]` SHALL pass when the inferred combined HCP range of the bidder and
their partner (bidder's actual HCP plus partner's inferred minimum and maximum) lies entirely
within `min`..`max`.

#### Scenario: Game values
- **WHEN** the bidder has 10 HCP, partner is inferred at 15–17, and the condition is `partnership_hcp: [25, 40]`
- **THEN** the condition passes

#### Scenario: Range not contained
- **WHEN** the bidder has 8 HCP, partner is inferred at 15–17, and the condition is `partnership_hcp: [25, 40]`
- **THEN** the condition fails
