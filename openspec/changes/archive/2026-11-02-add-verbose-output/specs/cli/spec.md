## ADDED Requirements

### Requirement: Verbose output
Run with `-v` or `--verbose`, `ppp` SHALL also print, between the hands and the bidding grid, one
entry per call in the order the calls were made. Each entry SHALL give the call's number, the
caller and the call, and the context the caller's convention sheet saw: the caller's HCP and
distribution, the calls before it, the ranges partner, LHO and RHO had shown, and the partnership
evaluation. Suit ranges SHALL be listed from spades down to clubs, like the hands. A player who
has shown nothing SHALL be printed as `nothing shown`. Without the flag, the output SHALL be
unchanged.

#### Scenario: Verbose run
- **WHEN** the user runs `ppp -v` and the auction has six calls
- **THEN** six entries numbered 1 to 6 are printed between the hands and the bidding grid

#### Scenario: Partner's opening in the context
- **WHEN** North opens `1NT` on the Colin card and South is the next player of North's partnership to call
- **THEN** South's entry shows partner at 15–17 HCP with 2–5 cards in every suit

#### Scenario: Nothing shown
- **WHEN** a player has made no call yet, or only passes
- **THEN** their range is printed as `nothing shown`

## MODIFIED Requirements

### Requirement: Card selection
Run without arguments, `ppp` SHALL use the embedded `cards/colin.yaml`. Run with two arguments,
it SHALL load the card at the path given by the second argument, relative to the data directory
given by the first. Options such as `-v` SHALL come before the two arguments.

#### Scenario: Custom card
- **WHEN** the user runs `ppp ./mydata cards/my_card.yaml`
- **THEN** the auction uses the card loaded from `./mydata/cards/my_card.yaml`

#### Scenario: Custom card with options
- **WHEN** the user runs `ppp -v ./mydata cards/my_card.yaml`
- **THEN** the auction uses the card loaded from `./mydata/cards/my_card.yaml` and prints the verbose output
