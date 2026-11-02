## MODIFIED Requirements

### Requirement: Suit order
The system SHALL order suits as Clubs (0) < Diamonds (1) < Hearts (2) < Spades (3) < NoTrump (4),
and every per-suit array SHALL be indexed in this order.

#### Scenario: Distribution indexing
- **WHEN** a hand holds 1 club, 2 diamonds, 4 hearts and 6 spades
- **THEN** its distribution, indexed by suit, is `[1,2,4,6]`

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
