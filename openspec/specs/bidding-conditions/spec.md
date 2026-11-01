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
between `min` and `max` cards inclusive. Unlisted suits SHALL be unconstrained. An unknown suit or
a range without exactly two values SHALL be a load error.

#### Scenario: Six-card heart suit
- **WHEN** the condition is `suit_length: {hearts: [6, 13]}` and the hand holds 6 hearts
- **THEN** the condition passes

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
`partner_bid`, `rho_bid` and `lho_bid` SHALL each pass when the most recent call by the bidder's
partner, right-hand opponent or left-hand opponent respectively equals the given bid. A seat that
has not yet called SHALL count as having passed.

#### Scenario: Partner opened 1NT
- **WHEN** the condition is `partner_bid: "1NT"` and partner's last call was `1NT`
- **THEN** the condition passes

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

### Requirement: Suit fit conditions
`heart_fit` and `spade_fit` SHALL pass when the combined card range in that suit (bidder's length
plus partner's inferred minimum and maximum length) lies entirely within `[min, max]`.
`major_fit` SHALL apply the same test to the better major, and `minor_fit` to the better minor,
where the better suit is the one with the higher combined maximum (spades or diamonds on a tie).

#### Scenario: Known eight-card heart fit
- **WHEN** the bidder holds 3 hearts, partner is inferred at 5–13 hearts, and the condition is `heart_fit: [8, 16]`
- **THEN** the condition passes

### Requirement: Opponent strength conditions
`lho_min_hcp: n` and `rho_min_hcp: n` SHALL pass when the left-hand or right-hand opponent's
inferred minimum HCP is at least `n`.

#### Scenario: Opponent opened
- **WHEN** RHO is inferred at 11–17 HCP and the condition is `rho_min_hcp: 11`
- **THEN** the condition passes
