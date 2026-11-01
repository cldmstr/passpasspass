# dealing Specification

## Purpose
Produce a standard 52-card deck and deal it into four 13-card hands.

## Requirements

### Requirement: Standard deck
The system SHALL provide a standard deck of 52 distinct cards: thirteen ranks (2 to ace) in each
of the four suits.

#### Scenario: Deck contents
- **WHEN** a standard deck is created
- **THEN** it holds 52 cards with no duplicates and 13 cards of each suit

### Requirement: Random deal
The random dealer SHALL shuffle one standard deck and deal four hands of 13 cards each, so that
every card appears in exactly one hand. Two dealers created with the same seed SHALL produce the
same deal.

#### Scenario: Every card dealt once
- **WHEN** the random dealer deals
- **THEN** it returns four hands of 13 cards that together contain all 52 cards

#### Scenario: Seeded determinism
- **WHEN** two random dealers are created with the same seed and each deals once
- **THEN** both deals are identical
