# Tasks: add-partner-aware-responses

## 1. Auction history
- [x] 1.1 Add `type History []BidAction` with `BidsBy(d)`, `LastBidBy(d)`, `Opening()` and `Contract()`
- [x] 1.2 Change `AuctionContext.History` to `History`; replace `lastBidByDir` with `LastBidBy`
- [x] 1.3 Tests: each helper on an empty history, a passed-out history, and a history with leading passes before the opening

## 2. `auction` pattern condition
- [ ] 2.1 Add a pattern type whose elements are a call or one of the wildcards `Any`, `Any*` and `Pass*` (design D2), parsed from strings; any other element must parse as a bid
- [ ] 2.2 Add `NewAuctionPatternCondition(pattern)`: matches the whole history, from the dealer's first call, anchored at both ends; leading passes are part of the history
- [ ] 2.3 Add an `auction` parser to the condition parser table; an invalid element is a load error
- [ ] 2.4 Tests: `[]` matches only before the dealer's call; `[1H, Pass]` matches a dealer opening but not a third-seat opening; `[Pass*, 1H, Pass]` matches the opening in every seat; `[1H, Pass]` and `[Pass, 1H, Pass]` tell an unpassed responder from a passed one; `Any` matches exactly one call; `Any*` matches zero, one and several calls; `Pass*` does not match a bid; a longer history does not match; an opponent's bid in place of a Pass does not match; an invalid element fails to load

## 3. `not` node
- [ ] 3.1 Add `NewNotCondition(c)` and a `not` parser; its value is one condition (a mapping with one or more keys, or a list) negated as a whole; an empty value is a load error; the flag-key error message now suggests `not`
- [ ] 3.2 Tests: `not` flips a hand condition and an auction pattern; `not` over a mapping with two keys, over `any` and over `all`; nested inside `all`; empty value fails to load

## 4. Hand model with alternatives
- [ ] 4.1 Add `HandModel{Cases []HandRange}` with `Intersect`, `Tighten`, `Bounds`, `IsEmpty`, and a constructor for the unconstrained model
- [ ] 4.2 `Tighten`: 13-card rule per alternative (suit max ≤ 13 − Σ other mins, suit min ≥ 13 − Σ other maxes) and HCP ≤ 37; drop alternatives left empty
- [ ] 4.3 Tests: intersecting disjoint alternatives yields an empty model; 6+ hearts and 6+ spades tightens both minors to 0–1; `Bounds` covers all alternatives

## 5. Reading calls in context
- [ ] 5.1 Add `Shown` to the `Condition` interface and implement it on every node type: the `HandModel` a node allows for a given history prefix and the facts before the call, following the table in design D4 (three-valued guards, `any` → alternatives, `all` → intersection)
- [ ] 5.2 `balanced` contributes 2–5 cards in every suit
- [ ] 5.3 Translate `partnership_hcp` and `fit` into bounds on the caller's own hand, using partner's shown model at the time of the call
- [ ] 5.4 Implement `Shown` on the `not` node as in design D9: flip auction-only guards; complement `hcp` and `suit_length` boxes; push through `all` and `any` with De Morgan's laws; approximated nodes give one unconstrained alternative
- [ ] 5.5 Add the `Interpreter` interface and the default implementation: the union over every convention whose bid equals the call and whose model is not empty
- [ ] 5.6 Cap alternatives per convention (design risk 1, including alternatives added by `not`) and fail loading above the cap
- [ ] 5.7 Remove `inferAndUpdate`, `extractBounds`, `applyBounds`, `intersectRanges` and `unionRanges`
- [ ] 5.8 Tests on the Colin card: `1NT` opening shows 15–17 and 2–5 in every suit; a `2♦️` response to `1NT` shows 5+ hearts, not Multi; a Multi `2♦️` opening shows four alternatives; a Pass leaves the model unchanged
- [ ] 5.9 Tests for `not`: `not: {hcp: [0, 9]}` shows 10+; `not: {suit_length: {hearts: [4, 13]}}` shows 0–3 hearts; `not: {balanced: true}` leaves the model unconstrained; a negated auction guard that holds removes the convention

## 6. Auction state and per-player view
- [ ] 6.1 Add `AuctionState` (history, `[4]HandModel` shown models, facts per partnership, interpreter, sheets by seat) with `Record(d, bid)`; constructor takes all dependencies
- [ ] 6.2 Add `Perspective` (own hand, partner's model tightened by own cards, LHO/RHO passed through, own partnership's facts) and build `AuctionContext` and `PartnershipEval` from it
- [ ] 6.3 Change `HandRange` fields in `AuctionContext` to `HandModel`; `PartnershipEval` uses partner's `Bounds()`
- [ ] 6.4 Rewrite `Auction.Run` on top of `AuctionState` and remove its exclusion in `.golangci.yml`; `NewAuction` takes the state; `Table.RunAuction` wires dealer, sheet, interpreter and state
- [ ] 6.5 Tests: partner's suit maximum drops to 13 minus own length; partner's HCP maximum drops to 40 minus own HCP; existing auction tests still pass

## 7. Established facts
- [ ] 7.1 Add the optional `establishes: {agreed: <strain>, forcing: one_round | game}` key to convention entries, decoded with unknown keys rejected; store it on `Convention`; an unknown fact, strain or forcing value is a load error
- [ ] 7.2 The interpreter returns the facts a call establishes: only facts that every contributing convention declares with the same value (design D10)
- [ ] 7.3 `AuctionState.Record` updates the caller's partnership: a new `agreed` replaces the old one; `one_round` applies to partner's next call only; `game` lasts until the partnership's contract reaches game
- [ ] 7.4 Add `agreed: <strain>` and `forced: true` (conditions and parsers; `forced` accepts only `true`), reading the bidder's partnership facts from `AuctionContext`; the interpreter treats both as auction-only guards
- [ ] 7.5 Tests with a temporary card: a rule guarded by `agreed: hearts` fires after `1♥️–2♥️`, `1♥️–3♥️` and `1♥️–2NT` but not after `1♥️–1♠️`; a later agreement replaces an earlier one; facts of one partnership are not seen as the other's; `forced` holds for partner's next call after `one_round` and not the call after; `forced` holds after a game force until game is bid; two conventions that explain the same call with different facts record neither

## 8. Conditions on partner's model and legal selection
- [ ] 8.1 Add `partner_hcp: [min, max]`, `partner_suit_length: {suit: [min, max]}` and `fit: {suit: [min, max]}` (conditions and parsers; ranges follow the range-value limits)
- [ ] 8.2 Apply the all-alternatives rule (design D5) to the new keys and to `partnership_hcp`
- [ ] 8.3 `ConventionSheet.Bid` skips conventions whose bid is neither Pass nor higher than the current contract; remove the replace-with-Pass step from `Run`
- [ ] 8.4 Tests: each new key at its boundaries; `fit: {hearts: [8, 14]}` and `fit: {spades: [7, 13]}` fail to load; opposite a Multi `2♦️` neither `fit: {hearts: [8, 13]}` nor `fit: {spades: [8, 13]}` holds; an illegal first match falls through to the next convention; an explicit `Pass` convention stops the search

## 9. Migrating existing data
- [ ] 9.1 Replace `no_prior_bids: true` with `auction: [Pass*]` in all opening conventions, including the Colin card's `2NT`; then remove the `no_prior_bids` condition, its parser and its tests
- [ ] 9.2 Replace `partner_bid: "1NT"` with `auction: [Pass*, 1NT, Pass]` in `jacoby_transfer.yaml`, `multi_two_diamond.yaml`, `baron.yaml` and `colin_stayman.yaml`; then remove the `partner_bid` condition, its parser and its tests
- [ ] 9.3 Add opener's replies: transfer completion in `jacoby_transfer.yaml`; `3♠️` after `2NT` and `3♦️` after `3♣️` in `multi_two_diamond.yaml`; Baron replies after `2♠️` in `baron.yaml`
- [ ] 9.4 Declare `forcing: one_round` on the strong `2♣️` and Multi `2♦️` openings and on the 1NT responses that ask opener to bid (transfers, the `2NT` and `3♣️` transfers, Baron, Stayman); declare `agreed: notrump` on the Multi file's `3NT` response
- [ ] 9.5 Reorder the Colin card includes: `jacoby_transfer`, `multi_two_diamond`, `strong_two_clubs`, `baron`, `colin_stayman`, `strong_notrump`, `tartan`, `five_card_major`; add `jacoby_transfer` first in the Two Over One card
- [ ] 9.6 Update the Colin 1NT response tests to use full auction histories; all existing round-trip tests pass

## 10. Built-in responses and rebids
Every response and rebid is guarded by an `auction` pattern. Ranges and declared facts are listed
in the `standard-conventions` spec delta.
- [ ] 10.1 `five_card_major.yaml`: responses to `1♥️`/`1♠️`, opener's rebids after raises, after Jacoby `2NT`, after `1♥️–1♠️` and after `1M–1NT`
- [ ] 10.2 `five_card_major.yaml`: responses to `1♣️`/`1♦️`
- [ ] 10.3 `strong_notrump.yaml`: standard responses to `1NT` (Stayman, `2NT`, `3NT`) and opener's replies
- [ ] 10.4 `strong_two_clubs.yaml`: responses to `2♣️` and opener's rebid after `2♦️`
- [ ] 10.5 `multi_two_diamond.yaml`: responses to Multi `2♦️` and opener's rebids, including explicit `Pass`
- [ ] 10.6 New `splinters.yaml`: splinter responses to `1♥️`/`1♠️` and opener's `4M` rebid
- [ ] 10.7 Include `splinters` in both cards, directly before `five_card_major`
- [ ] 10.8 Declare `establishes` on every built-in response and rebid as listed in the spec delta
- [ ] 10.9 Round-trip tests per sub-convention: one hand per response and rebid, plus boundary hands for every range (lowest, highest, one outside on each side)
- [ ] 10.10 Unpassed-hand tests: Jacoby `2NT`, the game-forcing 2/1 responses and splinters do not fire for a passed responder
- [ ] 10.11 Fact tests: after each built-in raise the partnership has agreed the major; after Jacoby `2NT`, a splinter and a 2/1 response a game force is in effect
- [ ] 10.12 Full-auction tests: fixed deals that run `1♥️–2♥️–Pass`, `1NT–2♦️–2♥️`, `2♦️–2♥️–2♠️` and a splinter sequence end to end, with leading passes in at least one deal
- [ ] 10.13 Random-deal property test: with 6+ HCP, responder never passes partner's opening of one of a suit (Two Over One card)

## 11. Documentation and archive
- [ ] 11.1 Update the condition keys, the `establishes` key and YAML examples in `openspec/project.md`
- [ ] 11.2 After merge, apply the spec deltas to `openspec/specs/` and move this change to `openspec/changes/archive/`
