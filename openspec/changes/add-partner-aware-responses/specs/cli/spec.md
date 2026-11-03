## MODIFIED Requirements

### Requirement: Verbose output
Run with `-v` or `--verbose`, `ppp` SHALL also print, between the hands and the bidding grid, one
entry per call in the order the calls were made. Each entry SHALL give the call's number, the
caller and the call, and the context the caller's convention sheet saw: the caller's HCP and
distribution, the calls before it, the models partner, LHO and RHO had shown, and the partnership
evaluation. A model with several alternatives SHALL print one alternative per line, and a partner
model that conflicts with the caller's hand SHALL be marked as such. Suit ranges SHALL be listed
from spades down to clubs, like the hands. A player who has shown nothing SHALL be printed as
`nothing shown`. Without the flag, the output SHALL be unchanged.

#### Scenario: Verbose run
- **WHEN** the user runs `ppp -v` and the auction has six calls
- **THEN** six entries numbered 1 to 6 are printed between the hands and the bidding grid

#### Scenario: Partner's opening in the context
- **WHEN** North opens `1NT` on the Colin card and South is the next player of North's partnership to call
- **THEN** South's entry shows partner at 15–17 HCP with 2–5 cards in every suit

#### Scenario: Several alternatives
- **WHEN** partner has opened a Multi `2♦️` on the Colin card
- **THEN** the next entry of partner's partnership shows partner's four alternatives, one per line

#### Scenario: Nothing shown
- **WHEN** a player has made no call yet, or only passes
- **THEN** their range is printed as `nothing shown`
