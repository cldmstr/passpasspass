# Tasks: add-verbose-output

## 1. Library
- [x] 1.1 `Auction.Steps` records, for every call, the `AuctionContext` the caller's sheet saw and the call made
- [x] 1.2 `String` for `HandRange` (`nothing shown` when unconstrained), `PartnershipEval`, `AuctionContext` and `Direction`
- [x] 1.3 Split `Auction.String` into `HandsString` and `BiddingString`; `String` prints both, unchanged
- [x] 1.4 Tests: steps of a known auction, including partner's range after an opening; the text of a context

## 2. CLI
- [x] 2.1 `-v` and `--verbose` with the `flag` package; options come before `<data dir> <card path>`
- [x] 2.2 `main` delegates to `run(args, stdout, stderr) int` so the command can be tested
- [x] 2.3 Tests: verbose output lists one entry per call; plain output has none; a bad card exits 1

## 3. Documentation and archive
- [x] 3.1 Add the verbose command to the commands in `AGENTS.md`
- [x] 3.2 Apply the spec delta to `openspec/specs/` and move this change to `openspec/changes/archive/`
