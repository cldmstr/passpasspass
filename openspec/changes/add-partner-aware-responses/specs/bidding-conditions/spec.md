## ADDED Requirements

### Requirement: Auction pattern condition
`auction: [<call>, ...]` SHALL pass when the calls made from the opening bid onwards, by all four
seats and including Passes, are exactly the listed calls in order. Passes made before the opening
bid SHALL be ignored. An empty list SHALL pass only when no bid has been made yet. Each listed call
SHALL be parsed as a bid, and an invalid call SHALL be a load error.

#### Scenario: Response to a one-level opening
- **WHEN** North opened `1♥️`, East passed, South is to call, and the condition is `auction: [1H, Pass]`
- **THEN** the condition passes

#### Scenario: Leading passes ignored
- **WHEN** North passed, East opened `1♥️`, South passed, West is to call, and the condition is `auction: [1H, Pass]`
- **THEN** the condition passes

#### Scenario: Later round
- **WHEN** the calls from the opening are `1♥️ Pass 2♥️ Pass` and the condition is `auction: [1H, Pass]`
- **THEN** the condition fails

#### Scenario: Contested auction
- **WHEN** North opened `1♥️`, East bid `1♠️`, and the condition is `auction: [1H, Pass]`
- **THEN** the condition fails

#### Scenario: Before the opening
- **WHEN** only passes have been made and the condition is `auction: []`
- **THEN** the condition passes

### Requirement: Partner HCP condition
`partner_hcp: [min, max]` SHALL pass when every alternative in partner's model, as seen by the
bidder, has its HCP range within `min`..`max`.

#### Scenario: Strong notrump partner
- **WHEN** partner has shown 15–17 and the condition is `partner_hcp: [15, 17]`
- **THEN** the condition passes

#### Scenario: Multi partner
- **WHEN** partner has shown a Multi `2♦️` and the condition is `partner_hcp: [6, 10]`
- **THEN** the condition fails because the strong alternatives fall outside the range

### Requirement: Partner suit length condition
`partner_suit_length: {<suit>: [min, max], ...}` SHALL pass when, for every listed suit, every
alternative in partner's model has its length range in that suit within `min`..`max`.

#### Scenario: Transfer accepted
- **WHEN** partner has shown 5+ hearts after a transfer and the condition is `partner_suit_length: {hearts: [5, 13]}`
- **THEN** the condition passes

### Requirement: Fit condition
`fit: {<suit>: [min, max], ...}` SHALL pass when, for every listed suit and every alternative in
partner's model, the bidder's length plus that alternative's length range lies within `min`..`max`.

#### Scenario: Eight-card fit
- **WHEN** the bidder holds 3 hearts, partner has shown 5–6 hearts, and the condition is `fit: {hearts: [8, 13]}`
- **THEN** the condition passes

#### Scenario: No known fit opposite Multi
- **WHEN** the bidder holds 3 hearts and 3 spades, partner has shown a Multi `2♦️`, and the condition is `fit: {hearts: [8, 13]}`
- **THEN** the condition fails

## MODIFIED Requirements

### Requirement: Partnership HCP condition
`partnership_hcp: [min, max]` SHALL pass when, for every alternative in partner's model as seen by
the bidder, the bidder's HCP plus that alternative's HCP range lies entirely within `min`..`max`.

#### Scenario: Game values
- **WHEN** the bidder has 10 HCP, partner has shown 15–17, and the condition is `partnership_hcp: [25, 40]`
- **THEN** the condition passes

#### Scenario: Range not contained
- **WHEN** the bidder has 8 HCP, partner has shown 15–17, and the condition is `partnership_hcp: [25, 40]`
- **THEN** the condition fails

### Requirement: Suit fit conditions
`heart_fit` and `spade_fit` SHALL pass when, for every alternative in partner's model as seen by
the bidder, the bidder's length in that suit plus the alternative's length range lies entirely
within `[min, max]`. `major_fit` SHALL apply the same test to the better major and `minor_fit` to
the better minor, choosing the better suit separately for each alternative as the one with the
higher combined maximum (spades or diamonds on a tie).

#### Scenario: Known eight-card heart fit
- **WHEN** the bidder holds 3 hearts, partner has shown 5–7 hearts, and the condition is `heart_fit: [8, 13]`
- **THEN** the condition passes

#### Scenario: One alternative without a fit
- **WHEN** the bidder holds 3 hearts, partner's model has one alternative with 6+ hearts and one with 0–5 hearts, and the condition is `heart_fit: [8, 13]`
- **THEN** the condition fails

### Requirement: Opponent strength conditions
`lho_min_hcp: n` and `rho_min_hcp: n` SHALL pass when every alternative in the left-hand or
right-hand opponent's shown model has a minimum HCP of at least `n`.

#### Scenario: Opponent opened
- **WHEN** RHO has shown 11–17 HCP and the condition is `rho_min_hcp: 11`
- **THEN** the condition passes
