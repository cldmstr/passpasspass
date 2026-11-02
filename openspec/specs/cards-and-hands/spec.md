# cards-and-hands Specification

## Purpose
Represent suits, ranks, cards and 13-card hands, and evaluate a hand's high-card points and suit
distribution.

## Requirements

### Requirement: Suit order
The system SHALL order suits as Clubs (0) < Diamonds (1) < Hearts (2) < Spades (3) < NoTrump (4),
and every per-suit array SHALL be indexed in this order.

#### Scenario: Distribution indexing
- **WHEN** a hand holds 1 club, 2 diamonds, 4 hearts and 6 spades
- **THEN** its distribution, indexed by suit, is `[1,2,4,6]`

### Requirement: Suit parsing
The system SHALL parse a suit from its emoji (`♣️ ♦️ ❤️ ♠️`), its letter (`C D H S NT`) or its
lowercase name (`clubs diamonds hearts spades`), and SHALL return an error for any other string.

#### Scenario: Equivalent spellings
- **WHEN** parsing `♠️`, `S` or `spades`
- **THEN** each yields Spades

#### Scenario: Unknown suit
- **WHEN** parsing `X`
- **THEN** an error is returned

### Requirement: Card notation
A card SHALL be written as `<rank><suit emoji>`, where rank is one of `2`–`10`, `J`, `Q`, `K`,
`A`. Parsing an unknown rank or suit SHALL return an error.

#### Scenario: Ten of hearts
- **WHEN** parsing `10❤️`
- **THEN** the card is the ten of hearts and prints as `10❤️`

#### Scenario: Unknown rank
- **WHEN** parsing `1♣️`
- **THEN** an error is returned

### Requirement: High-card points
A card SHALL be worth 4 HCP for an ace, 3 for a king, 2 for a queen, 1 for a jack and 0
otherwise. A hand's HCP SHALL be the sum over its cards.

#### Scenario: One of each honour
- **WHEN** a hand holds A♣️, K♦️, Q❤️, J♠️ and nine spot cards
- **THEN** its HCP is 10

### Requirement: Hand validation
A hand built from card strings SHALL contain exactly 13 cards with no duplicates; otherwise
construction SHALL fail with an error.

#### Scenario: Too few cards
- **WHEN** a hand is built from 12 cards
- **THEN** an error is returned

#### Scenario: Duplicate card
- **WHEN** a hand lists `A♣️` twice
- **THEN** an error is returned

### Requirement: Hand display
A hand SHALL print one line per suit from spades down to clubs, with the cards in ascending rank
separated by commas and an empty line for a void. Wherever a hand's shape is printed, its suit
lengths SHALL be listed in the same order, spades first.

#### Scenario: Spades first
- **WHEN** a hand with cards in every suit is printed
- **THEN** the lines list spades, hearts, diamonds and clubs, in that order

#### Scenario: Void in clubs
- **WHEN** a hand with no clubs is printed
- **THEN** the first three lines list spades, hearts and diamonds, and the last line is empty

#### Scenario: Shape
- **WHEN** a hand holds 1 club, 2 diamonds, 4 hearts and 6 spades
- **THEN** its shape prints as `[6,4,2,1]`
