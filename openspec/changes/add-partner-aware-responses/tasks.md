# Tasks: add-partner-aware-responses

## 1. Auction history
- [ ] 1.1 Add `type History []BidAction` with `BidsBy(d)`, `LastBidBy(d)`, `Opening()` and `Contract()`
- [ ] 1.2 Change `AuctionContext.History` to `History`; replace `lastBidByDir` with `LastBidBy`
- [ ] 1.3 Tests: each helper on an empty history, a passed-out history, and a history with leading passes before the opening

## 2. `auction` pattern condition
- [ ] 2.1 Add `NewAuctionPatternCondition(calls []*Bid)`: matches the history from the opening bid onwards against all listed calls; an empty list matches when no bid has been made
- [ ] 2.2 Add the `auction: [<call>, ...]` YAML key to `conditionNode` and `buildCondition`; an invalid call is a load error
- [ ] 2.3 Tests: exact match; leading passes ignored; a longer history does not match; an opponent's bid in place of a Pass does not match; an empty pattern before and after the opening; invalid call fails to load

## 3. Hand model with alternatives
- [ ] 3.1 Add `HandModel{Cases []HandRange}` with `Intersect`, `Tighten`, `Bounds`, `IsEmpty`, and a constructor for the unconstrained model
- [ ] 3.2 `Tighten`: 13-card rule per alternative (suit max ≤ 13 − Σ other mins, suit min ≥ 13 − Σ other maxes) and HCP ≤ 37; drop alternatives left empty
- [ ] 3.3 Tests: intersecting disjoint alternatives yields an empty model; 6+ hearts and 6+ spades tightens both minors to 0–1; `Bounds` covers all alternatives

## 4. Reading calls in context
- [ ] 4.1 Add a function that compiles a `conditionNode` into a `HandModel` for a given history prefix, following the table in design D4 (three-valued guards, `any` → alternatives, `all` → intersection)
- [ ] 4.2 `balanced` contributes 2–5 cards in every suit
- [ ] 4.3 Translate `partnership_hcp`, `fit` and the `*_fit` keys into bounds on the caller's own hand, using partner's shown model at the time of the call
- [ ] 4.4 Add the `Interpreter` interface and the default implementation: the union over every convention whose bid equals the call and whose model is not empty
- [ ] 4.5 Cap alternatives per convention (design risk 1) and fail loading above the cap
- [ ] 4.6 Remove `inferAndUpdate`, `extractBounds`, `applyBounds` and `intersectRanges`
- [ ] 4.7 Tests on the Colin card: `1NT` opening shows 15–17 and 2–5 in every suit; a `2♦️` response to `1NT` shows 5+ hearts, not Multi; a Multi `2♦️` opening shows four alternatives; a Pass leaves the model unchanged

## 5. Auction state and per-player view
- [ ] 5.1 Add `AuctionState` (history, `[4]HandModel` shown models, interpreter, sheets by seat) with `Record(d, bid)`; constructor takes all dependencies
- [ ] 5.2 Add `Perspective` (own hand, partner's model tightened by own cards, LHO/RHO passed through) and build `AuctionContext` and `PartnershipEval` from it
- [ ] 5.3 Change `HandRange` fields in `AuctionContext` to `HandModel`; `PartnershipEval` uses partner's `Bounds()`
- [ ] 5.4 Rewrite `Auction.Run` on top of `AuctionState`; `NewAuction` takes the state; `Table.RunAuction` wires dealer, sheet, interpreter and state
- [ ] 5.5 Tests: partner's suit maximum drops to 13 minus own length; partner's HCP maximum drops to 40 minus own HCP; existing auction tests still pass

## 6. Conditions on partner's model and legal selection
- [ ] 6.1 Add `partner_hcp: [min, max]`, `partner_suit_length: {suit: [min, max]}` and `fit: {suit: [min, max]}` (conditions and YAML keys)
- [ ] 6.2 Apply the all-alternatives rule (design D5) to the new keys and to `partnership_hcp`, `major_fit`, `minor_fit`, `heart_fit` and `spade_fit`
- [ ] 6.3 `ConventionSheet.Bid` skips conventions whose bid is neither Pass nor higher than the current contract; remove the replace-with-Pass step from `Run`
- [ ] 6.4 Tests: each new key at its boundaries; opposite a Multi `2♦️` neither `fit: {hearts: [8, 13]}` nor `fit: {spades: [8, 13]}` holds; an illegal first match falls through to the next convention; an explicit `Pass` convention stops the search

## 7. Migrating existing data
- [ ] 7.1 Replace `no_prior_bids: true` with `auction: []` in all opening conventions
- [ ] 7.2 Replace `partner_bid: "1NT"` with `auction: [1NT, Pass]` in `jacoby_transfer.yaml`, `multi_two_diamond.yaml`, `baron.yaml` and `colin_stayman.yaml`
- [ ] 7.3 Add opener's replies: transfer completion in `jacoby_transfer.yaml`; `3♠️` after `2NT` and `3♦️` after `3♣️` in `multi_two_diamond.yaml`; Baron replies after `2♠️` in `baron.yaml`
- [ ] 7.4 Reorder the Colin card includes: `jacoby_transfer`, `multi_two_diamond`, `strong_two_clubs`, `baron`, `colin_stayman`, `strong_notrump`, `tartan`, `five_card_major`; add `jacoby_transfer` first in the Two Over One card
- [ ] 7.5 Update the Colin 1NT response tests to use full auction histories; all existing round-trip tests pass

## 8. Built-in responses and rebids
Every response and rebid is guarded by an `auction` pattern. Ranges are listed in the
`standard-conventions` spec delta.
- [ ] 8.1 `five_card_major.yaml`: responses to `1♥️`/`1♠️`, opener's rebids after raises, after Jacoby `2NT`, after `1♥️–1♠️` and after `1M–1NT`
- [ ] 8.2 `five_card_major.yaml`: responses to `1♣️`/`1♦️`
- [ ] 8.3 `strong_notrump.yaml`: standard responses to `1NT` (Stayman, `2NT`, `3NT`) and opener's replies
- [ ] 8.4 `strong_two_clubs.yaml`: responses to `2♣️` and opener's rebid after `2♦️`
- [ ] 8.5 `multi_two_diamond.yaml`: responses to Multi `2♦️` and opener's rebids, including explicit `Pass`
- [ ] 8.6 New `splinters.yaml`: splinter responses to `1♥️`/`1♠️` and opener's `4M` rebid
- [ ] 8.7 Include `splinters` in both cards, directly before `five_card_major`
- [ ] 8.8 Round-trip tests per sub-convention: one hand per response and rebid, plus boundary hands for every range (lowest, highest, one outside on each side)
- [ ] 8.9 Full-auction tests: fixed deals that run `1♥️–2♥️–Pass`, `1NT–2♦️–2♥️`, `2♦️–2♥️–2♠️` and a splinter sequence end to end
- [ ] 8.10 Random-deal property test: with 6+ HCP, responder never passes partner's opening of one of a suit (Two Over One card)

## 9. Documentation and archive
- [ ] 9.1 Update the condition keys and YAML examples in `openspec/project.md`, if any have changed
- [ ] 9.2 After merge, apply the spec deltas to `openspec/specs/` and move this change to `openspec/changes/archive/`
