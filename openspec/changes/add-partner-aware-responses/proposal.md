# Change: Add partner-aware responses and opener's rebids

## Why
The trainer only practises opening bids. To practise responses and rebids, each player needs to
know what partner's bids have shown and to choose a call from that information. The current
inference cannot support this:

- It reads a bid without context, taking the *first* convention with a matching bid anywhere in
  the sheet. A Jacoby `2♦️` response to `1NT` is read as a Multi `2♦️` opening.
- It reduces a bid with several meanings to one merged range. Multi `2♦️` becomes "6–37 HCP, any
  shape", so responder learns nothing.
- It ignores `balanced`, the fact that a hand has 13 cards, and the cards the player holds.
- `partner_bid` only checks partner's latest call, so a response rule also fires in later rounds
  and in contested auctions.

## What Changes
- **Auction history**: per-player queries over the call record (calls by a player, last call,
  opening bid, current contract).
- **Hand model with alternatives**: what a player has shown becomes a list of alternative
  HCP and suit-length ranges instead of a single range.
- **Bids read in context**: a call's meaning comes from the caller's sheet and the auction as it
  stood when the call was made, never from the caller's hand.
- **Per-player view**: each player sees partner's model tightened by their own cards. Responses
  and rebids in this change use partner's information only.
- **Auction patterns**: `auction` matches the whole call history from the dealer's first call,
  including leading passes and opponents' calls. Wildcards `Any` (one call), `Any*` (zero or more
  calls) and `Pass*` (zero or more passes) let one rule cover every seat, while an exact pattern
  can still tell a passed hand from an unpassed one.
- **Negation**: a `not` node passes when its child fails, so a rule can exclude auctions or hands
  without listing every alternative.
- **Established facts**: a convention can declare what its call establishes for the partnership:
  an agreed suit and a one-round or game force. The `agreed` and `forced` conditions read these
  facts, so a rule can apply after a suit is agreed however the auction got there.
- **New conditions**: `partner_hcp`, `partner_suit_length`, `fit`, `agreed`, `forced`. Every
  condition on partner's model holds only when it is true for every remaining alternative.
- **Legal bid selection**: the sheet skips conventions whose bid does not outrank the current
  contract and tries the next one. Previously an illegal bid was replaced by Pass.
- **Built-in responses and rebids**: standard responses and opener's rebids are added to the file
  that holds the opening (`five_card_major`, `strong_notrump`, `strong_two_clubs`,
  `multi_two_diamond`). Responses that apply to several openings get their own file; the first is
  `splinters`.
- **BREAKING (data)**: openings switch from `no_prior_bids: true` to `auction: [Pass*]`, and the
  1NT response files (`jacoby_transfer`, the 1NT responses in `multi_two_diamond`, `baron`,
  `colin_stayman`) switch from `partner_bid` to `auction` patterns.
  `baron` and `colin_stayman` move ahead of `strong_notrump` in the Colin card so Colin's scheme
  beats standard Stayman, and the Two Over One card gains `jacoby_transfer`.

## Impact
- Affected specs: `auction`, `bidding-conditions`, `convention-evaluation`, `convention-files`,
  `standard-conventions`
- Affected code: `auction.go` (history, hand model, interpretation, per-player view,
  established facts), `condition.go` (new conditions, `not`, pattern matching, all-alternatives
  semantics), `convention.go` (legal selection, declared facts), `loader.go` (new keys,
  `establishes`, compiling condition trees into hand models), `table.go` (wiring),
  `data/conventions/*.yaml`, new `data/conventions/splinters.yaml`, `data/cards/*.yaml`
- Unchanged: `cards-and-hands`, `dealing`, `cli`
