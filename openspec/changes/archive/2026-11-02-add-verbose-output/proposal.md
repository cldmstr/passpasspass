# Change: Show what each player knew at every call

## Why
`ppp` prints the deal, the bidding and the contract, but not why a player made a call. Since
calls are read in context (`add-partner-aware-responses`, section 5), every player has a picture
of the others built from their calls alone. Seeing that picture at each step is the quickest way
to check a convention card and the inference behind it, and to learn from an auction.

## What Changes
- **Verbose flag**: `ppp -v` (or `--verbose`) prints, between the hands and the bidding grid, one
  entry per call: who called what, and the context their sheet saw: own HCP and shape, the calls
  so far, what partner, LHO and RHO had shown, and the partnership evaluation.
- **Readable ranges**: a player who has shown nothing is printed as `nothing shown`.
- **Auction steps**: `Auction.Run` records each call together with the context it was made in.
- Without the flag, the output is unchanged.

## Impact
- Affected specs: `cli`
- Affected code: `cmd/ppp/main.go` (flag parsing), `auction.go` (`Auction.Steps`, `String`
  methods for `AuctionContext`, `PartnershipEval` and `Direction`, `Auction.String` split into
  hands and bidding), `hand_model.go` (`HandRange.String`)
- `add-partner-aware-responses` changes `AuctionContext` to models with alternatives in task 6.3;
  its `String` follows then, printing each alternative on its own line.
