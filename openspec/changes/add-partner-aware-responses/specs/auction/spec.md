## ADDED Requirements

### Requirement: Auction history
The auction SHALL keep every call in order together with the seat that made it, and SHALL answer
these queries: all calls by a given seat; a seat's most recent call (Pass if it has not called);
the opening bid and the seat that made it; and the current contract.

#### Scenario: Opening after leading passes
- **WHEN** North passes, East bids `1♥️`, South passes and West bids `2♥️`
- **THEN** the opening is `1♥️` by East, the contract is `2♥️`, and the calls by East are `[1♥️]`

#### Scenario: No opening yet
- **WHEN** only passes have been made
- **THEN** there is no opening and no contract

### Requirement: Hand model with alternatives
What a player has shown SHALL be a list of alternative ranges, each with HCP bounds and per-suit
length bounds. Every player SHALL start with one unconstrained alternative (0–37 HCP, 0–13 cards
per suit). Combining new information with a model SHALL keep every non-empty intersection of an
existing alternative with a new one. Every alternative SHALL respect the 13-card rule: no suit's
maximum exceeds 13 minus the sum of the other suits' minimums.

#### Scenario: Multi narrows after the correction
- **WHEN** West opens a Multi `2♦️`, East responds `2♥️`, and West corrects to `2♠️`
- **THEN** every alternative left in West's model has 6–10 HCP and 6+ spades

#### Scenario: Two long majors
- **WHEN** a model has 6+ hearts and 6+ spades
- **THEN** clubs and diamonds are limited to 0–1 cards each

### Requirement: Established facts
The auction SHALL keep, for each partnership, an agreed strain and the forces in effect. After each
call, the system SHALL record the facts the call establishes: those declared under `establishes`
by every convention that contributes to the call's meaning (see Bid inference), with the same
value. A fact declared by only some of those conventions, or with different values, SHALL NOT be
recorded. The facts SHALL apply to the caller's partnership and SHALL last as follows:

- `agreed`: until a later call by the same partnership establishes a different agreed strain.
- `forcing: one_round`: for partner's next call only.
- `forcing: game`: until the contract held by the partnership reaches game: `3NT`, `4♥️`, `4♠️`,
  `5♣️`, `5♦️` or higher.

Facts SHALL come from the calls alone, never from the caller's hand.

#### Scenario: Raise agrees the suit
- **WHEN** North opens `1♥️`, East passes, and South raises to `2♥️` with a convention that establishes `agreed: hearts`
- **THEN** North–South have agreed hearts and East–West have agreed nothing

#### Scenario: Later agreement replaces an earlier one
- **WHEN** a partnership has agreed hearts and one of its players then makes a call that establishes `agreed: notrump`
- **THEN** the partnership has agreed notrump

#### Scenario: Readings disagree
- **WHEN** a call is explained by two conventions, one establishing `agreed: hearts` and one establishing no agreed strain
- **THEN** no agreed strain is recorded for the call

#### Scenario: One-round force expires
- **WHEN** South's transfer establishes `forcing: one_round`, West passes, and North completes the transfer with a call that establishes no force
- **THEN** South is not under a force at South's next turn

#### Scenario: Game force lasts until game
- **WHEN** South's Jacoby `2NT` establishes `forcing: game` and North rebids `3♥️`
- **THEN** South is still under a force, and the force ends once the partnership bids `4♥️`

## MODIFIED Requirements

### Requirement: Bid inference
After each non-pass call, the system SHALL combine the caller's shown model with the meaning of the
call. The meaning SHALL be computed from the caller's convention sheet and the auction as it stood
before the call, without looking at the caller's hand. It SHALL be the union of the models of
every convention whose bid equals the call and whose conditions can hold in that auction, where a
convention's model is derived from its conditions:

- `hcp` and `suit_length` set bounds; `balanced` sets 2–5 cards in every suit; `rule_of_20` sets
  nothing.
- Auction conditions (`auction`, `no_prior_bids`, `partner_bid`, `rho_bid`, `lho_bid`, `agreed`,
  `forced`) are checked against the auction and the established facts before the call, and a
  convention whose auction conditions fail contributes nothing.
- Conditions on the partnership (`partnership_hcp`, `fit`, `heart_fit`, `spade_fit`, `major_fit`,
  `minor_fit`) become bounds on the caller's own hand, given partner's shown model at the time of
  the call.
- `all` intersects its children; `any` keeps each child as a separate alternative.
- `not` flips an auction condition; turns `hcp` and `suit_length` into their complement, with one
  alternative for each bound that can be broken; is pushed through `all` and `any` by De Morgan's
  laws; and leaves the model unconstrained for `balanced`, `rule_of_20`, conditions on other
  players' models, and partnership conditions.

A Pass SHALL leave the model unchanged.

#### Scenario: Strong notrump opening
- **WHEN** North opens `1NT` with a sheet whose `1NT` opening requires `balanced: true` and `hcp: [15, 17]`
- **THEN** North's model is one alternative of 15–17 HCP with 2–5 cards in every suit

#### Scenario: Same bid, different context
- **WHEN** North opens `1NT`, East passes, and South bids `2♦️` on the Colin card
- **THEN** South's model shows 5+ hearts and does not include the Multi opening's alternatives

#### Scenario: Multi opening keeps alternatives
- **WHEN** a seat opens a Multi `2♦️`
- **THEN** its model holds four alternatives: 23–24 balanced, 6–10 with 6+ hearts, 6–10 with 6+ spades, and 27+

#### Scenario: Negated HCP range
- **WHEN** the only convention for a call requires `not: {hcp: [0, 9]}`
- **THEN** the caller's model shows 10–37 HCP

#### Scenario: Negated approximation
- **WHEN** the only convention for a call requires `not: {balanced: true}`
- **THEN** the call leaves the caller's model unchanged

#### Scenario: Partnership condition becomes own bounds
- **WHEN** partner has shown 15–17 HCP and the caller's convention requires `partnership_hcp: [25, 40]`
- **THEN** the caller's model shows at least 10 HCP

### Requirement: Bidding context
Before each call, the system SHALL give the seat's convention sheet a context built from that
seat's own view of the auction: the seat's hand, its direction, the auction history so far,
partner's shown model tightened by the seat's own cards, the LHO's and RHO's shown models
unchanged, the facts established for the seat's partnership, and the partnership evaluation.
Tightening SHALL limit partner to at most 13 minus the seat's own length in each suit and at most
40 minus the seat's own HCP.

#### Scenario: Responder sees the opening
- **WHEN** North has opened `1NT` and South is about to call
- **THEN** South's context history contains North's `1NT` and East's call

#### Scenario: Own cards limit partner
- **WHEN** partner has shown 5+ spades and the seat holds 6 spades
- **THEN** partner's spade range in the seat's context is 5–7

### Requirement: Bid legality
The system SHALL only record legal calls. A seat's call SHALL be the first convention in its sheet
that matches and whose bid is a Pass or outranks the current contract, as defined by
convention-evaluation. When no such convention exists, the seat SHALL pass.

#### Scenario: Illegal match falls through
- **WHEN** the contract is `2♠️`, a seat's first matching convention bids `1NT`, and a later matching convention bids `3♥️`
- **THEN** the seat bids `3♥️`

### Requirement: Partnership evaluation
The partnership evaluation SHALL combine the seat's own hand with partner's model from the seat's
view: the combined HCP range is own HCP plus the lowest and highest HCP across partner's
alternatives, and the fit in each suit is own length plus the lowest and highest length across
partner's alternatives. The better major SHALL be the major with the higher combined maximum
(spades on a tie), and likewise the better minor (diamonds on a tie).

#### Scenario: Opposite a 1NT opening
- **WHEN** the seat holds 10 HCP and 3 hearts, and partner has shown 15–17 HCP and 2–5 hearts
- **THEN** the combined HCP range is 25–27 and the heart fit is 5–8

#### Scenario: Opposite Multi
- **WHEN** the seat holds 2 hearts and partner has shown a Multi `2♦️`
- **THEN** the heart fit spans from the lowest to the highest heart length across the four alternatives plus 2
