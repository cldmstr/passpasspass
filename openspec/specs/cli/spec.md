# cli Specification

## Purpose
Provide the `ppp` command that deals one hand, runs an auction with a convention card, and prints
the result.

## Requirements

### Requirement: Card selection
Run without arguments, `ppp` SHALL use the embedded `cards/colin.yaml`. Run with two arguments,
it SHALL load the card at the path given by the second argument, relative to the data directory
given by the first.

#### Scenario: Custom card
- **WHEN** the user runs `ppp ./mydata cards/my_card.yaml`
- **THEN** the auction uses the card loaded from `./mydata/cards/my_card.yaml`

### Requirement: Load failure
If the card cannot be loaded, `ppp` SHALL print the error to standard error and exit with status 1.

#### Scenario: Missing card
- **WHEN** the card path does not exist
- **THEN** an error is printed to standard error and the exit status is 1

### Requirement: Single auction output
`ppp` SHALL deal with a time-seeded random dealer, run one auction, and print it to standard
output.

#### Scenario: Default run
- **WHEN** the user runs `ppp`
- **THEN** four hands, the bidding grid and the contract (or `Passed out`) are printed
