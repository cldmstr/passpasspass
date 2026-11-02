# auction Specification

## Purpose
Run a complete bidding sequence for one deal: rotate through the seats, ask each seat's convention
sheet for a call, track what each player's bids reveal, and determine the final contract.

## Requirements

### Requirement: Dealer and rotation
Each auction run by a table SHALL deal four new hands, choose the dealer at random, and let the
dealer call first. Calls SHALL then rotate clockwise: North, East, South, West.

#### Scenario: East deals
- **WHEN** East is dealer
- **THEN** the calls are made in the order East, South, West, North, East, …

### Requirement: Bidding context
Before each call, the system SHALL give the seat's convention sheet a context holding the seat's
hand, its direction, the auction history so far, the inferred hand ranges of its partner, LHO and
RHO, and the partnership evaluation.

#### Scenario: Responder sees the opening
- **WHEN** North has opened `1NT` and South is about to call
- **THEN** South's context history contains North's `1NT` and East's call

### Requirement: Bid legality
If the call selected for a seat is not a Pass and does not outrank the current contract, the
system SHALL record a Pass instead.

#### Scenario: Lower bid replaced
- **WHEN** the contract is `2♠️` and a seat's sheet selects `1NT`
- **THEN** a Pass is recorded for that seat

### Requirement: Termination and contract
The auction SHALL end after four passes when nobody has bid, leaving no contract, or after three
consecutive passes following a bid. The contract SHALL be the last non-pass bid and the declarer
the seat that made it.

#### Scenario: Passed out
- **WHEN** all four seats pass
- **THEN** the auction ends with four calls and no contract

#### Scenario: Single opener
- **WHEN** North opens `1NT` and the other three seats pass
- **THEN** the contract is `1NT`, the declarer is North, and four calls were made

### Requirement: Bid inference
Every player SHALL start with an inferred hand range of 0–37 HCP and 0–13 cards in each suit.
After a non-pass call, the system SHALL narrow the caller's range by intersecting it with the bounds
of the first convention in the sheet whose bid equals the call. The bounds SHALL come from the
convention's conditions: `hcp` and `suit_length` set bounds, `all` intersects its children, `any`
takes the union of its children's bounds, and every other condition contributes nothing. A Pass
SHALL leave the range unchanged.

#### Scenario: Strong notrump opening
- **WHEN** North opens `1NT` with a sheet whose first `1NT` convention requires `balanced: true` and `hcp: [15, 17]`
- **THEN** North's inferred range is 15–17 HCP with 0–13 cards in every suit

#### Scenario: Union of alternatives
- **WHEN** a seat opens a Multi `2♦️` whose `any` covers 23–24, 6–10 with 6+ hearts, 6–10 with 6+ spades, and 27+
- **THEN** the seat's inferred range is 6–37 HCP with 0–13 cards in every suit

### Requirement: Partnership evaluation
The partnership evaluation SHALL combine the seat's own hand with partner's inferred range: the
combined HCP range is own HCP plus partner's minimum and maximum, and the fit in each suit is own
length plus partner's minimum and maximum length. A partnership never holds more than the deck
does, so both ends of the combined HCP range SHALL be capped at 40, and both ends of each suit's
fit at 13.

#### Scenario: Opposite a 1NT opening
- **WHEN** the seat holds 10 HCP and 3 hearts, and partner is inferred at 15–17 HCP with 0–13 hearts
- **THEN** the combined HCP range is 25–27 and the heart fit is 3–13

#### Scenario: Opposite a player who has shown nothing
- **WHEN** the seat holds 25 HCP and 7 hearts, and partner is inferred at 0–37 HCP with 0–13 hearts
- **THEN** the combined HCP range is 25–40 and the heart fit is 7–13

### Requirement: Auction display
An auction SHALL print each seat's hand with its HCP and distribution, the dealer, the calls in a
four-column grid headed North, East, South, West starting under the dealer's column, and then
either `Contract: <bid> by <declarer>` or `Passed out`.

#### Scenario: Passed out display
- **WHEN** the auction has no contract
- **THEN** the output ends with `Passed out`
