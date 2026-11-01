## ADDED Requirements

### Requirement: Responses to one of a major
The `five_card_major` sub-convention SHALL define responses after `auction: [1H, Pass]` and
`auction: [1S, Pass]`, where M is the opened major, evaluated in this order:

| Order | Bid | Shows |
|-------|-----|-------|
| 1 | `2NT` (Jacoby) | 4+ cards in M, 13+ HCP |
| 2 | `4M` | 5+ cards in M, 0–9 HCP |
| 3 | `3M` (limit raise) | 4+ cards in M, 10–12 HCP |
| 4 | `2M` | 3+ cards in M, 6–10 HCP |
| 5 | `1♠️` (over `1♥️` only) | 4+ spades, 6+ HCP |
| 6 | `2♥️` (over `1♠️` only) | 5+ hearts, 13+ HCP |
| 7 | `2♦️` | 5+ diamonds, 13+ HCP |
| 8 | `2♣️` | 3+ clubs, 13+ HCP |
| 9 | `2♦️` | 3+ diamonds, 13+ HCP |
| 10 | `1NT` (forcing) | 6–12 HCP |

Any other hand SHALL pass.

#### Scenario: Simple raise
- **WHEN** partner opened `1♥️` and responder holds 8 HCP with 3 hearts
- **THEN** responder bids `2♥️`

#### Scenario: Limit raise boundary
- **WHEN** partner opened `1♠️` and responder holds 10 HCP with 4 spades
- **THEN** responder bids `3♠️`

#### Scenario: Game-forcing raise
- **WHEN** partner opened `1♥️` and responder holds 14 HCP with 4 hearts and no singleton or void
- **THEN** responder bids `2NT`

#### Scenario: Spades over hearts
- **WHEN** partner opened `1♥️` and responder holds 9 HCP with 2 hearts and 4 spades
- **THEN** responder bids `1♠️`

#### Scenario: Forcing notrump
- **WHEN** partner opened `1♠️` and responder holds 7 HCP with 2 spades and no 4-card heart suit
- **THEN** responder bids `1NT`

#### Scenario: Too weak
- **WHEN** partner opened `1♥️` and responder holds 5 HCP with 2 hearts
- **THEN** responder passes

### Requirement: Responses to one of a minor
The `five_card_major` sub-convention SHALL define responses after `auction: [1C, Pass]` and
`auction: [1D, Pass]`, where m is the opened minor, evaluated in this order:

| Order | Bid | Shows |
|-------|-----|-------|
| 1 | `1♠️` | 5+ spades, 6+ HCP |
| 2 | `1♥️` | 4+ hearts, 6+ HCP |
| 3 | `1♠️` | 4+ spades, 6+ HCP |
| 4 | `1♦️` (over `1♣️` only) | 4+ diamonds, 6+ HCP |
| 5 | `3NT` | balanced, 13–15 HCP |
| 6 | `2NT` | balanced, 11–12 HCP |
| 7 | `2♣️` (over `1♦️` only) | 4+ clubs, 13+ HCP |
| 8 | `3m` (limit raise) | 10–12 HCP; 5+ clubs over `1♣️`, 4+ diamonds over `1♦️` |
| 9 | `2m` | 6–10 HCP; same support as the limit raise |
| 10 | `1NT` | 6–10 HCP |

Any other hand SHALL pass.

#### Scenario: Hearts up the line
- **WHEN** partner opened `1♦️` and responder holds 8 HCP with 4 hearts and 4 spades
- **THEN** responder bids `1♥️`

#### Scenario: Five spades first
- **WHEN** partner opened `1♣️` and responder holds 8 HCP with 4 hearts and 5 spades
- **THEN** responder bids `1♠️`

#### Scenario: Invitational notrump
- **WHEN** partner opened `1♣️` and responder holds 11 HCP, balanced, with 3 diamonds and no 4-card major
- **THEN** responder bids `2NT`

### Requirement: Opener's rebids after a major opening
The `five_card_major` sub-convention SHALL define opener's rebids, evaluated in the order listed
within each auction:

| Auction | Rebid | Shows |
|---------|-------|-------|
| `[1M, Pass, 2M, Pass]` | `4M` | 17 HCP |
| | `3M` | 15–16 HCP |
| | Pass | otherwise |
| `[1M, Pass, 3M, Pass]` | `4M` | 14+ HCP |
| | Pass | otherwise |
| `[1M, Pass, 2NT, Pass]` | `3` of another suit | 0–1 cards in that suit |
| | `3M` | 15+ HCP |
| | `4M` | otherwise |
| `[1H, Pass, 1S, Pass]` | `4♠️` | 4+ spades, 17 HCP |
| | `3♠️` | 4+ spades, 15–16 HCP |
| | `2♠️` | 4+ spades |
| | `2♥️` | 6+ hearts |
| | `1NT` | balanced |
| | `2♣️` | 4+ clubs |
| | `2♦️` | 4+ diamonds |
| | `2♥️` | otherwise |
| `[1M, Pass, 1NT, Pass]` | `2M` | 6+ cards in M |
| | `2♥️` (after `1♠️` only) | 4+ hearts |
| | `2♦️` | 4+ diamonds |
| | `2♣️` | 3+ clubs |
| | `2♦️` | 3+ diamonds |
| | `2M` | otherwise |

#### Scenario: Invitation after a simple raise
- **WHEN** the auction is `1♥️ Pass 2♥️ Pass` and opener holds 15 HCP
- **THEN** opener bids `3♥️`

#### Scenario: Minimum after a simple raise
- **WHEN** the auction is `1♠️ Pass 2♠️ Pass` and opener holds 12 HCP
- **THEN** opener passes

#### Scenario: Shortness after Jacoby 2NT
- **WHEN** the auction is `1♥️ Pass 2NT Pass` and opener holds a singleton club
- **THEN** opener bids `3♣️`

#### Scenario: Raising responder's spades
- **WHEN** the auction is `1♥️ Pass 1♠️ Pass` and opener holds 12 HCP with 4 spades
- **THEN** opener bids `2♠️`

### Requirement: Standard responses to 1NT
The `strong_notrump` sub-convention SHALL define responses after `auction: [1NT, Pass]`, in this
order: `2♣️` (Stayman) with 8+ HCP and 4+ hearts or 4+ spades; `2NT` with 8–9 HCP; `3NT` with
10–15 HCP. Any other hand SHALL pass. Transfers come from `jacoby_transfer`.

It SHALL define opener's replies: after `[1NT, Pass, 2C, Pass]` bid `2♥️` with 4+ hearts, else
`2♠️` with 4+ spades, else `2♦️`; after `[1NT, Pass, 2NT, Pass]` bid `3NT` with 16–17 HCP, else
pass.

#### Scenario: Stayman answer
- **WHEN** the auction is `1NT Pass 2♣️ Pass` and opener holds 4 spades and 3 hearts
- **THEN** opener bids `2♠️`

#### Scenario: Invitation declined
- **WHEN** the auction is `1NT Pass 2NT Pass` and opener holds 15 HCP
- **THEN** opener passes

### Requirement: Responses to a strong 2♣
The `strong_two_clubs` sub-convention SHALL define responses after `auction: [2C, Pass]`, in this
order: `2♥️` with 5+ hearts and 8+ HCP; `2♠️` with 5+ spades and 8+ HCP; `3♣️` with 5+ clubs and
8+ HCP; `3♦️` with 5+ diamonds and 8+ HCP; `2NT` balanced with 8+ HCP; otherwise `2♦️` (waiting).

It SHALL define opener's rebids after `[2C, Pass, 2D, Pass]`, in this order: `2NT` balanced with
20–22 HCP; `3NT` balanced with 25–26 HCP; `2♠️` with 5+ spades; `2♥️` with 5+ hearts; `3♦️` with
5+ diamonds; `3♣️` with 5+ clubs; otherwise `2NT`.

#### Scenario: Waiting response
- **WHEN** partner opened `2♣️` and responder holds 4 HCP
- **THEN** responder bids `2♦️`

#### Scenario: Positive response
- **WHEN** partner opened `2♣️` and responder holds 9 HCP with 5 spades and no other 5-card suit
- **THEN** responder bids `2♠️`

#### Scenario: Balanced rebid
- **WHEN** the auction is `2♣️ Pass 2♦️ Pass` and opener holds 21 HCP, balanced
- **THEN** opener bids `2NT`

### Requirement: Responses to Multi 2♦
The `multi_two_diamond` sub-convention SHALL define responses after `auction: [2D, Pass]`, in this
order: `2NT` (enquiry) with 16+ HCP; `2♠️` with 3+ hearts and 10–15 HCP (pass with spades, or
correct to `3♥️` with hearts); otherwise `2♥️` (pass or correct).

It SHALL define opener's rebids, in the order listed within each auction:

| Auction | Rebid | Shows |
|---------|-------|-------|
| `[2D, Pass, 2H, Pass]` | `2NT` | balanced, 23–24 HCP |
| | `3NT` | 27+ HCP |
| | Pass | 6–10 HCP, 6+ hearts |
| | `2♠️` | 6–10 HCP, 6+ spades |
| `[2D, Pass, 2S, Pass]` | `2NT` | balanced, 23–24 HCP |
| | `3NT` | 27+ HCP |
| | Pass | 6–10 HCP, 6+ spades |
| | `3♥️` | 6–10 HCP, 6+ hearts |
| `[2D, Pass, 2NT, Pass]` | `3NT` | balanced, 23–24 HCP |
| | `4NT` | 27+ HCP |
| | `3♥️` | 6–10 HCP, 6+ hearts |
| | `3♠️` | 6–10 HCP, 6+ spades |

#### Scenario: Pass or correct
- **WHEN** partner opened Multi `2♦️` and responder holds 7 HCP with 2 hearts and 3 spades
- **THEN** responder bids `2♥️`

#### Scenario: Correction to spades
- **WHEN** the auction is `2♦️ Pass 2♥️ Pass` and opener holds 8 HCP with 6 spades and 2 hearts
- **THEN** opener bids `2♠️`

#### Scenario: Weak hearts accepted
- **WHEN** the auction is `2♦️ Pass 2♥️ Pass` and opener holds 7 HCP with 6 hearts
- **THEN** opener passes

#### Scenario: Strong enquiry
- **WHEN** partner opened Multi `2♦️` and responder holds 17 HCP
- **THEN** responder bids `2NT`

### Requirement: Splinter responses
The `splinters` sub-convention SHALL define splinter responses that show 4+ cards in opener's
major, 12–15 HCP, and 0–1 cards in the bid suit: after `auction: [1H, Pass]` the bids `3♠️`,
`4♣️` and `4♦️`; after `auction: [1S, Pass]` the bids `4♣️`, `4♦️` and `4♥️`. After any splinter,
opener SHALL bid four of the opened major.

#### Scenario: Club splinter over hearts
- **WHEN** partner opened `1♥️` and responder holds 13 HCP with 4 hearts and a singleton club
- **THEN** responder bids `4♣️` rather than `2NT`

#### Scenario: Opener signs off
- **WHEN** the auction is `1♠️ Pass 4♦️ Pass`
- **THEN** opener bids `4♠️`

## MODIFIED Requirements

### Requirement: Jacoby transfers
The `jacoby_transfer` sub-convention SHALL respond after `auction: [1NT, Pass]` in this order:
`2♦️` with 5+ hearts (any strength); `2♥️` with exactly 5 spades. It SHALL define opener's
completion: after `[1NT, Pass, 2D, Pass]` bid `2♥️`; after `[1NT, Pass, 2H, Pass]` bid `2♠️`.

#### Scenario: Transfer to hearts
- **WHEN** partner opened `1NT` and responder holds 3 HCP with 5 hearts
- **THEN** responder bids `2♦️` and opener then bids `2♥️`

#### Scenario: Not a later round
- **WHEN** the calls from the opening are `1NT Pass 2♣️ Pass 2♥️ Pass`
- **THEN** neither transfer applies

### Requirement: 1NT responses in the Multi file
The `multi_two_diamond` sub-convention SHALL also respond after `auction: [1NT, Pass]`, after its
opening, in this order: `3NT` with 10–14 HCP, balanced and no 4-card major; `2NT` weakness
transfer with 6–9 HCP and 6+ spades; `3♣️` transfer with 6–9 HCP and 6+ diamonds. It SHALL define
opener's replies: after `[1NT, Pass, 2NT, Pass]` bid `3♠️`; after `[1NT, Pass, 3C, Pass]` bid
`3♦️`.

#### Scenario: Weakness transfer completed
- **WHEN** the auction is `1NT Pass 2NT Pass`
- **THEN** opener bids `3♠️`

### Requirement: Baron
The `baron` sub-convention SHALL respond `2♠️` after `auction: [1NT, Pass]` with 8+ HCP. It SHALL
define opener's reply after `[1NT, Pass, 2S, Pass]`: `2NT` with 3–4 cards in every suit,
otherwise the lowest 4-card suit at the three level.

#### Scenario: Baron with a four-card major
- **WHEN** responder holds 10 HCP, 4 hearts and a singleton, after `1NT Pass` on the Colin card
- **THEN** they respond `2♠️`

#### Scenario: Baron reply
- **WHEN** the auction is `1NT Pass 2♠️ Pass` and opener holds 4 clubs and 4 hearts
- **THEN** opener bids `3♣️`

### Requirement: Colin Stayman
The `colin_stayman` sub-convention SHALL respond `2♣️` after `auction: [1NT, Pass]` with 4+ hearts
or 4+ spades. Opener's replies come from `strong_notrump`.

#### Scenario: Weak hand with a four-card major
- **WHEN** partner opened `1NT` and responder holds 4 HCP with 4 hearts, on the Colin card
- **THEN** they respond `2♣️`

### Requirement: Two Over One card
The `two_over_one` card SHALL include, in order, `jacoby_transfer`, `multi_two_diamond`,
`strong_two_clubs`, `strong_notrump`, `splinters` and `five_card_major`, with no card-level
conventions.

#### Scenario: Too few points
- **WHEN** a hand holds 0 HCP before anyone has opened
- **THEN** the card selects Pass

#### Scenario: Standard 1NT responses available
- **WHEN** partner opened `1NT` and responder holds 10 HCP, balanced, with no 4-card major
- **THEN** the card selects `3NT`

### Requirement: Colin card
The `colin` card SHALL include, in order, `jacoby_transfer`, `multi_two_diamond`,
`strong_two_clubs`, `baron`, `colin_stayman`, `strong_notrump`, `tartan`, `splinters` and
`five_card_major`, followed by a
card-level `2NT` opening before anyone has opened, with 10–12 HCP, 5+ clubs and 5+ diamonds.

#### Scenario: Minor two-suiter
- **WHEN** an opener holds 10 HCP, 5 clubs and 5 diamonds
- **THEN** they open `2NT`

#### Scenario: Colin 1NT scheme takes precedence
- **WHEN** partner opened `1NT` and responder holds 10 HCP with 4 hearts and a singleton
- **THEN** the card selects Baron `2♠️`, not standard Stayman `2♣️`
