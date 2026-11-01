## ADDED Requirements

### Requirement: Auction pattern condition
`auction: [<element>, ...]` SHALL pass when the pattern matches the whole auction so far: every
call from the dealer's first call to the most recent one, by all four seats, including Passes
made before the opening. The match SHALL be anchored at both ends. Each element SHALL be either a
call, which matches exactly that call, or one of these wildcards:

- `Any` matches exactly one call of any kind.
- `Any*` matches zero or more calls of any kind.
- `Pass*` matches zero or more Passes.

An empty list SHALL pass only before the dealer has called. Any element that is not a wildcard
SHALL be parsed as a bid, and an element that is neither SHALL be a load error.

#### Scenario: Response to a dealer opening
- **WHEN** North dealt and opened `1♥️`, East passed, South is to call, and the condition is `auction: [1H, Pass]`
- **THEN** the condition passes

#### Scenario: Leading passes are part of the auction
- **WHEN** North dealt and passed, East opened `1♥️`, South passed, West is to call, and the condition is `auction: [1H, Pass]`
- **THEN** the condition fails

#### Scenario: Opening in any seat
- **WHEN** North dealt and passed, East opened `1♥️`, South passed, West is to call, and the condition is `auction: [Pass*, 1H, Pass]`
- **THEN** the condition passes

#### Scenario: Passed responder
- **WHEN** North dealt and passed, East passed, South opened `1♥️`, West passed, North is to call, and the condition is `any` of `auction: [1H, Pass]` and `auction: [Pass, 1H, Pass]`
- **THEN** the condition fails, because North has already passed

#### Scenario: Pass wildcard does not match a bid
- **WHEN** the calls are `1♣️ Pass 1♥️ Pass` and the condition is `auction: [Pass*, 1H, Pass]`
- **THEN** the condition fails

#### Scenario: Single-call wildcard
- **WHEN** the calls are `1♥️ Pass 2♥️ Pass` and the condition is `auction: [1H, Pass, Any, Pass]`
- **THEN** the condition passes

#### Scenario: Any number of calls
- **WHEN** the calls are `1♥️ Pass 2♥️ Pass 3♥️ Pass` and the condition is `auction: [Pass*, 1H, Pass, Any*]`
- **THEN** the condition passes

#### Scenario: Later round
- **WHEN** the calls are `1♥️ Pass 2♥️ Pass` and the condition is `auction: [Pass*, 1H, Pass]`
- **THEN** the condition fails

#### Scenario: Contested auction
- **WHEN** North opened `1♥️`, East bid `1♠️`, and the condition is `auction: [Pass*, 1H, Pass]`
- **THEN** the condition fails

#### Scenario: Before the opening
- **WHEN** only passes have been made and the condition is `auction: [Pass*]`
- **THEN** the condition passes

#### Scenario: Dealer to call
- **WHEN** no call has been made and the condition is `auction: []`
- **THEN** the condition passes, and it fails once the dealer has passed

#### Scenario: Invalid element
- **WHEN** a convention file contains `auction: [1H, Foo]`
- **THEN** loading fails with an error naming `Foo`

### Requirement: Agreed suit condition
`agreed: <strain>` SHALL pass when the bidder's partnership has an agreed strain, as established
by the auction, and it is the given strain. Strains SHALL be `clubs`, `diamonds`, `hearts`,
`spades` or `notrump`, and an unknown strain SHALL be a load error.

#### Scenario: Agreed by a raise
- **WHEN** the bidder opened `1♥️` and partner raised to `2♥️`, a call that establishes `agreed: hearts`, and the condition is `agreed: hearts`
- **THEN** the condition passes

#### Scenario: Agreed by a different route
- **WHEN** the bidder opened `1♥️` and partner bid Jacoby `2NT`, a call that establishes `agreed: hearts`, and the condition is `agreed: hearts`
- **THEN** the condition passes

#### Scenario: Nothing agreed
- **WHEN** the bidder opened `1♥️` and partner responded `1♠️`, a call that establishes no agreed strain, and the condition is `agreed: hearts`
- **THEN** the condition fails

#### Scenario: Opponents' agreement
- **WHEN** the opponents have agreed hearts and the bidder's partnership has agreed nothing, and the condition is `agreed: hearts`
- **THEN** the condition fails

### Requirement: Forcing condition
`forced: true` SHALL pass when the bidder is under a force from partner: partner's most recent call
established `forcing: one_round`, or a `forcing: game` established by the bidder's partnership is
in effect and the current contract is below game. `forced: false` SHALL pass exactly when
`forced: true` fails.

#### Scenario: One-round force
- **WHEN** partner's last call was a transfer that establishes `forcing: one_round`, and the condition is `forced: true`
- **THEN** the condition passes

#### Scenario: One-round force has ended
- **WHEN** partner's call two rounds ago established `forcing: one_round` and partner's last call established no force, and the condition is `forced: true`
- **THEN** the condition fails

#### Scenario: Game force below game
- **WHEN** partner's Jacoby `2NT` established `forcing: game` and the contract is `3♥️`, and the condition is `forced: true`
- **THEN** the condition passes

#### Scenario: Game reached
- **WHEN** the partnership is under a game force and the contract is `4♥️`, and the condition is `forced: true`
- **THEN** the condition fails

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

### Requirement: Combinators
`all: [...]` SHALL pass when every child passes. `any: [...]` SHALL pass when at least one child
passes. `not: <condition>` SHALL take exactly one child and SHALL pass when that child fails; a
`not` without a child SHALL be a load error. Combinators SHALL nest.

#### Scenario: Any of two shapes
- **WHEN** the condition is `any` of 6+ hearts or 6+ spades, and the hand holds 6 spades and 2 hearts
- **THEN** the condition passes

#### Scenario: Excluding a shape
- **WHEN** the condition is `not: {suit_length: {hearts: [4, 13]}}` and the hand holds 3 hearts
- **THEN** the condition passes

#### Scenario: Excluding an auction
- **WHEN** the condition is `not: {auction: [Pass*, 1NT, Pass]}` and the calls so far are `1NT Pass`
- **THEN** the condition fails

#### Scenario: Negating a combinator
- **WHEN** the condition is `not: {any: [{hcp: [0, 5]}, {hcp: [16, 37]}]}` and the hand holds 10 HCP
- **THEN** the condition passes


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
