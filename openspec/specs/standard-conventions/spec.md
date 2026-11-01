# standard-conventions Specification

## Purpose
Define the bidding conventions and convention cards shipped in the embedded data directory
(`data/conventions/`, `data/cards/`).

## Requirements

### Requirement: Multi 2♦ opening
The `multi_two_diamond` sub-convention SHALL open `2♦️` when no prior bid has been made and the
hand is one of: balanced with 23–24 HCP; 6–10 HCP with 6+ hearts; 6–10 HCP with 6+ spades; or
27+ HCP.

#### Scenario: Weak hearts upper bound
- **WHEN** an opener holds 10 HCP and 6 hearts
- **THEN** they open `2♦️`

#### Scenario: Too strong for the weak range
- **WHEN** an opener holds 11 HCP and 6 hearts
- **THEN** they do not open `2♦️`

### Requirement: Strong 2♣ opening
The `strong_two_clubs` sub-convention SHALL open `2♣️` when no prior bid has been made and the
hand holds 20–22 or 25–26 HCP.

#### Scenario: Twenty points
- **WHEN** an opener holds 20 HCP that no earlier convention covers
- **THEN** they open `2♣️`

### Requirement: Strong 1NT opening
The `strong_notrump` sub-convention SHALL open `1NT` when no prior bid has been made and the hand
is balanced with 15–17 HCP.

#### Scenario: Sixteen balanced
- **WHEN** an opener holds 16 HCP and 4-3-3-3 shape
- **THEN** they open `1NT`

### Requirement: Five-card major openings
The `five_card_major` sub-convention SHALL, when no prior bid has been made, the hand meets the
Rule of 20, and it holds 11–17 HCP, open in this order: `1♠️` with 5–6 spades, `1♥️` with 5–6
hearts, `1♦️` with 5–6 diamonds, `1♣️` with 2–6 clubs.

#### Scenario: Six spades
- **WHEN** an opener holds 14 HCP and 6 spades
- **THEN** they open `1♠️`

### Requirement: Tartan two-bids
The `tartan` sub-convention SHALL open, when no prior bid has been made and the hand holds 11–15
HCP, `2♥️` with 5+ hearts and another 5+ card suit, or `2♠️` with 5+ spades and another 5+ card
suit.

#### Scenario: Hearts and clubs
- **WHEN** an opener holds 13 HCP, 5 hearts and 5 clubs
- **THEN** they open `2♥️`

### Requirement: Jacoby transfers
The `jacoby_transfer` sub-convention SHALL respond to partner's `1NT` in this order: `2♦️` with
5+ hearts (any strength); `2♥️` with exactly 5 spades.

#### Scenario: Weak hand with five hearts
- **WHEN** partner opened `1NT` and responder holds 3 HCP with 5 hearts
- **THEN** they respond `2♦️`

### Requirement: 1NT responses in the Multi file
The `multi_two_diamond` sub-convention SHALL also respond to partner's `1NT`, after its opening, in
this order: `3NT` with 10–14 HCP, balanced and no 4-card major; `2NT` weakness transfer with 6–9
HCP and 6+ spades; `3♣️` transfer with 6–9 HCP and 6+ diamonds.

#### Scenario: Direct game
- **WHEN** partner opened `1NT` and responder holds 10 HCP, balanced, with no 4-card major
- **THEN** they respond `3NT`

#### Scenario: Weakness transfer to spades
- **WHEN** partner opened `1NT` and responder holds 6 HCP with 6 spades
- **THEN** they respond `2NT`

### Requirement: Baron
The `baron` sub-convention SHALL respond `2♠️` to partner's `1NT` with 8+ HCP.

#### Scenario: Baron with a four-card major
- **WHEN** partner opened `1NT` and responder holds 10 HCP, 4 hearts and a singleton, on the Colin card
- **THEN** they respond `2♠️`

### Requirement: Colin Stayman
The `colin_stayman` sub-convention SHALL respond `2♣️` to partner's `1NT` with 4+ hearts or 4+
spades. On the Colin card it is included after `baron`, so it only applies below 8 HCP.

#### Scenario: Weak hand with a four-card major
- **WHEN** partner opened `1NT` and responder holds 4 HCP with 4 hearts, on the Colin card
- **THEN** they respond `2♣️`

### Requirement: Two Over One card
The `two_over_one` card SHALL include, in order, `multi_two_diamond`, `strong_two_clubs`,
`strong_notrump` and `five_card_major`, with no card-level conventions.

#### Scenario: Too few points
- **WHEN** a hand holds 0 HCP
- **THEN** the card selects Pass

#### Scenario: 1NT responses come with the Multi file
- **WHEN** partner opened `1NT` and responder holds 10 HCP, balanced, with no 4-card major
- **THEN** the card selects `3NT`

### Requirement: Colin card
The `colin` card SHALL include, in order, `jacoby_transfer`, `multi_two_diamond`,
`strong_two_clubs`, `strong_notrump`, `tartan`, `five_card_major`, `baron` and `colin_stayman`,
followed by a card-level `2NT`
opening with no prior bids, 10–12 HCP, 5+ clubs and 5+ diamonds.

#### Scenario: Minor two-suiter
- **WHEN** an opener holds 10 HCP, 5 clubs and 5 diamonds
- **THEN** they open `2NT`
