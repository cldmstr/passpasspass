# Design: Partner-aware responses

## Context
`Auction.Run` already keeps a `[4]HandRange` and narrows it after every call with
`inferAndUpdate`, and `AuctionContext` already exposes `Partner`, `LHO`, `RHO` and a
`PartnershipEval`. This change rebuilds that pipeline so that it can carry responses and rebids:
calls are read in context, bids with several meanings stay as separate alternatives, and each
player sees partner through their own cards.

```
Auction.Run
  └─ AuctionState          call history + what each player has shown (from the calls alone)
       ├─ History          []BidAction with per-player queries
       ├─ Interpreter      (sheet, history and facts before the call, caller, call)
       │                      → HandModel + facts the call establishes
       ├─ [4]HandModel     shown models, indexed by Direction
       └─ [2]Facts         agreed suit and force in effect, per partnership
  └─ Perspective           one player's view: own hand + partner's model tightened by own cards
       └─ AuctionContext → conditions → ConventionSheet.Bid
```

## Goals / Non-Goals

**Goals**
- Read every call in the auction context in which it was made.
- Keep a bid with several meanings (Multi) as alternatives, and narrow them as the auction goes on.
- Let responses and opener's rebids depend on partner's shown model and the player's own hand.
- Let a rule depend on what the auction has established (an agreed suit, a force) without listing
  every sequence that leads there.
- Ship standard responses and rebids for every built-in opening, and splinters as the first
  response file shared by several openings.

**Non-Goals**
- Using opponents' calls to choose a response or rebid. Opponents' models are still tracked and
  exposed unchanged, but no built-in convention uses them.
- Overcalls and contested auctions.
- Doubles and redoubles. `Bid` only holds a level and a strain, so contested auctions will need a
  new call type before `auction` patterns can describe them.
- Inference from calls that were *not* made: a Pass, or earlier conventions in the sheet that did
  not apply.
- Responder's second call and later rounds; opener's rebids after a 2/1 game-forcing response
  or after responses to a minor-suit opening; responses to Tartan or to the Colin `2NT`.
- Simulating possible hands by sampling. Ranges are cheap, behave predictably, and can be shown to
  a learner directly.

## Decisions

### D1: `History` as a named slice type
`type History []BidAction` with `BidsBy(d)`, `LastBidBy(d)` and `Opening()`.
`AuctionContext.History` changes to this type and `lastBidByDir` is folded into `LastBidBy`. The
auction as it stood when call *i* was made is `h[:i]`.

History is a plain record of calls and does not answer for the contract (user decision). The
auction already tracks its contract as calls are recorded and passes it on as
`AuctionContext.Contract`: the highest bid so far, or nil before the opening.

### D2: `auction` matches the whole call history (user decision)
`auction: [...]` matches when the pattern covers the entire history, from the dealer's first call
to the most recent one, by all four seats. Leading passes are part of the history and are never
dropped, so a pattern can tell which seat opened and whether the player to call has already
passed.

A pattern element is either a call (`1H`, `2NT`, `Pass`) or a wildcard:

| Element | Matches |
|---------|---------|
| `Any` | exactly one call of any kind |
| `Any*` | zero or more calls of any kind |
| `Pass*` | zero or more passes |

Matching is anchored at both ends, like a glob over calls. The wildcards are words rather than
`*` and `?` because those characters are YAML syntax inside a flow sequence and would have to be
quoted in every pattern. Any other element must parse as a bid; anything else is a load error.

| Pattern | Meaning |
|---------|---------|
| `[]` | the dealer is to call |
| `[Pass*]` | nobody has opened |
| `[Pass*, 1H, Pass]` | partner opened `1♥️` in any seat, RHO passed |
| `[1H, Pass]` or `[Pass, 1H, Pass]` | as above, and the player to call has not passed |
| `[Pass*, 1H, Pass, Any*]` | partner opened `1♥️` and RHO passed; any calls since |

Seat-independent rules use `Pass*`; rules that only apply to an unpassed hand (game-forcing
responses) list the two unpassed patterns under `any`.

The built-in data moves to `auction`, after which `no_prior_bids` and `partner_bid` are used by no
file and are removed (tasks 9.1 and 9.2). `rho_bid` and `lho_bid` were already removed by `refactor-condition-nodes`.

### D3: `HandModel` holds a list of alternatives
```go
type HandModel struct{ Cases []HandRange }
```
- `Intersect` keeps every pairwise intersection that a hand can fit, tightened; an empty model
  means "impossible".
- `Tighten` applies the 13-card rule to each alternative: no suit's maximum exceeds 13 minus the
  sum of the other suits' minimums, and no suit's minimum is below 13 minus the sum of the other
  suits' maximums. Narrowing one suit can narrow another, so it repeats until nothing changes. It
  also caps HCP at 37, and drops alternatives no hand can fit.
- `Bounds` returns the smallest single range covering all alternatives, for display and for
  simple checks, and false for an empty model.

`HandModel`, `HandRange` and the range helpers `intersectRanges` and `unionRanges` live in
`hand_model.go`.

Ranges stay boxes. Shape facts that a box cannot express, such as "at most one doubleton" or the
Rule of 20, are approximated or dropped: `balanced` becomes 2–5 in every suit; `rule_of_20`
contributes nothing.

### D4: Calls are read from condition nodes with three-valued guards
Conditions are typed nodes (`refactor-condition-nodes`). The `Condition` interface gains a second
method next to `Eval`:

```go
// Shown returns the model of the hands for which the condition comes out as holds, for a call
// made in ctx.
Shown(ctx *CallContext, holds bool) HandModel
```

`CallContext` holds what the whole table knows when the call is made: the caller, the history
before the call and every player's shown model (later also the facts), never the caller's hand.
Asking for `holds == false` gives the negated model, so `not` only flips the flag for its child,
and De Morgan's laws (D9) follow from `all` and `any` answering both questions. Every node type
implements it:

| Node | Result |
|------|--------|
| `hcp`, `suit_length`, `balanced` | one alternative with those bounds |
| `rule_of_20` | one unconstrained alternative |
| auction-only (`auction`, `agreed`, `forced`) | evaluated on the history and facts before the call: true → one unconstrained alternative, false → empty |
| partner-model (`partner_hcp`, `partner_suit_length`) | evaluated on the shown models; true → unconstrained, false → empty. Without the caller's own cards these can only be checked approximately, so an inconclusive result counts as true |
| partnership (`partnership_hcp`, `fit`) | translated into bounds on the caller's own hand using partner's shown model, e.g. `partnership_hcp: [25, 40]` opposite 15–17 gives the caller 10+ HCP |
| `all` | intersection of the children |
| `any` | concatenation of the children's alternatives |
| `not` | see D9 |

The meaning of a call is the union of the models of **every** convention in the caller's sheet
whose bid equals the call and whose model is not empty. Taking every candidate, not just the
first, avoids claiming more than is known when an earlier convention's guards held but its hand
conditions might not have.

A third, unexported method, `maxCases(holds)`, bounds the alternatives `Shown` can return in any
auction: a sum for alternatives side by side, a product for intersections. The loader rejects a
convention above the cap of risk 1 (64), with the line of its `conditions`.

A call that no convention explains, or whose conventions all fail their guards, shows nothing:
the interpreter returns the unconstrained model rather than an empty one.

`extractBounds`, `applyBounds` and `inferAndUpdate` are replaced. Until `AuctionState` exists
(task 6.4), `Run` creates the interpreter itself and passes each model's `Bounds()` to the
`HandRange` fields of `AuctionContext`. The
interpreter is an interface injected into `AuctionState`, and takes the caller's sheet as an
argument so that partnerships can later play different cards.

### D5: Every condition on partner's model must hold for all alternatives (user decision)
`partner_hcp`, `partner_suit_length`, `fit` and `partnership_hcp` pass only when the test holds
for every remaining alternative. For a range-containment test this is the same as testing
`Bounds()`, which keeps the implementation simple. This errs on the side of not assuming a fit:
opposite Multi, neither major is a known fit, which is why the pass-or-correct responses exist.

A `fit` range lies within 8–13: a suit has 13 cards, so a partnership never holds more, and fewer
than 8 combined cards is not a fit. The absence of a fit is written `not: {fit: ...}`.

### D6: `AuctionState` and `Perspective` are injected
`NewAuctionState(interp Interpreter, sheets [4]*ConventionSheet, shown [4]HandModel)` and
`NewAuction(dealer, hands, state)`. `Table.RunAuction`, a method rather than a constructor,
builds the per-deal state, which follows the AGENTS.md dependency-injection rule.

`Perspective` combines the player's own hand with partner's shown model, tightened by the player's
cards: partner holds at most `13 − own length` in each suit and at most `40 − own HCP` points.
LHO and RHO models are passed through untightened. The current contract, which `AuctionState`
tracks as calls are recorded, reaches the context as `AuctionContext.Contract`. Facts (D10) join
`AuctionState` and `Perspective` in task 7.3.

Two conflicts are resolved in favour of what was said most recently (user decision):

- **A call that contradicts the caller's earlier calls**: when the meaning of a call has no hand
  in common with what the caller had shown, the caller's model becomes the meaning of the latest
  call. Keeping the empty intersection would leave every later reader with nothing.
- **Partner's model against the player's own cards**: when the player's cards leave no hand
  partner could hold (partner showed 6+ spades, the player holds 8), the context keeps partner's
  model as partner showed it and sets `HandModel.Conflict`. The flag tells later rules that
  partner has misbid or is lying; handling that is left to a later change.

### D7: The sheet only returns legal bids
`ConventionSheet.Bid` skips a convention whose bid is neither Pass nor higher than
`ctx.Contract`, and tries the next. The replace-with-Pass step in `Run` is removed. A
convention may bid `Pass` explicitly, which stops a weaker rule further down the sheet from
firing (e.g. opener passing `2♥️` after Multi with weak hearts).

### D8: Where responses live (user decision)
- Standard responses and opener's rebids go into the file that holds the opening.
- A response that applies to several openings gets its own file, included **before** the opening
  files it extends so that first-match picks it over the more general responses. `splinters.yaml`
  is the first such file.
- The 1NT responses are split across `jacoby_transfer`, `multi_two_diamond` (`3NT`, the `2NT`
  weakness transfer and the `3♣️` transfer), `baron` and `colin_stayman`. Both cards include
  `jacoby_transfer`. The Colin card includes `baron` and `colin_stayman` before `strong_notrump`,
  so Colin's scheme wins and the standard responses in `strong_notrump` only catch what it leaves
  over.

Every response and rebid is guarded by an `auction` pattern, so responses never collide with
openings (guarded by `auction: [Pass*]`) whatever the include order.

### D9: `not` negates one condition
`not: <condition>` passes when its condition fails. The value is one condition: a mapping, which
may hold several keys, or a list. Like everywhere else, several keys or items form a conjunction,
so `not: {hcp: [0, 9], suit_length: {hearts: [4, 13]}}` means "not both", and is read as `not`
over an `all`.

When a call is read (D4), the negated child is handled by its kind:

| Child | Result |
|-------|--------|
| auction-only, including `agreed` and `forced` | the child's truth value is flipped: true → empty, false → one unconstrained alternative |
| `hcp`, `suit_length` | the complement of the child's box: one alternative for each bound that can be broken (below the minimum, above the maximum), per dimension |
| `all`, `any` | pushed down with De Morgan's laws, then handled as above |
| `balanced`, `rule_of_20`, partner-model, partnership | one unconstrained alternative |

The last row is the safe choice: these nodes are already approximations, and the complement of an
approximation is not an approximation of the complement. Treating them as unconstrained can only
widen the shown model, never make it wrong. Complements add alternatives, so they count toward the
cap in risk 1.

### D10: Conventions declare the facts their call establishes
A convention entry may carry `establishes` next to `bid` and `conditions`:

```yaml
- bid: 2H
  establishes: {agreed: hearts}
  conditions: ...
- bid: 2NT
  establishes: {agreed: hearts, forcing: game}
  conditions: ...
```

| Fact | Values | Lasts |
|------|--------|-------|
| `agreed` | `clubs`, `diamonds`, `hearts`, `spades`, `notrump` | until a later call by the same partnership agrees a different strain |
| `forcing` | `one_round` | for partner's next call only |
| | `game` | until the partnership's contract reaches game (`3NT`, `4♥️`, `4♠️`, `5♣️`, `5♦️` or higher) |

`AuctionState` keeps the facts for each partnership. When the interpreter reads a call, it
returns the facts as well as the hand model. A fact is recorded only when **every** convention
that contributes to the call's meaning (D4) declares it with the same value, the same caution as
D5: if one reading of the call agrees hearts and another does not, hearts are not agreed. Facts
come from calls alone, so both partnerships see the same facts.

Two conditions read the facts of the bidder's partnership:

- `agreed: <strain>` passes when that strain is the agreed one.
- `forced: true` passes when partner's most recent call established `forcing: one_round`, or a
  `forcing: game` is in effect and the contract is below game. Like the other flag keys it
  accepts only `true`; the opposite is `not: {forced: true}`.

This answers "after a suit has been agreed" in one condition, however the auction got there:
`agreed: hearts` holds after `1♥️–2♥️`, `1♥️–3♥️` and `1♥️–2NT` alike, but not after `1♥️–1♠️`.

## Risks / Trade-offs
- **Alternatives multiply**: `all` over several `any` nodes forms a cross product. Mitigation:
  drop empty intersections straight away, and fail loading above a fixed cap (e.g. 64
  alternatives per convention).
- **Approximate public evaluation** (D4): conditions on partner's model are checked without the
  caller's cards, so the interpreter may consider a convention the caller could not actually have
  used. This only widens the shown model and never makes it wrong.
- **Box approximation**: a "balanced" alternative still allows 5-4-2-2 hands. This is acceptable
  for range-based training feedback.
- **Ranges are a judgement call**: the response and rebid ranges below are standard 2/1 values
  but debatable. They live only in YAML and in the `standard-conventions` spec, so tuning them
  needs no code change.
- **Gaps until the next change**: auctions can stop early where no rebid is defined (e.g. opener
  passes a 2/1 response). This is listed under Non-Goals and covered by a follow-up change.
- **Longer patterns**: because leading passes stay in the history (D2), seat-independent rules
  need a `Pass*` prefix, and rules for unpassed hands need two patterns under `any`.
- **Simplified fact lifetimes** (D10): a force is either one round or until game, and an agreed
  strain can only be replaced, not withdrawn. Auctions that need finer rules (a force that ends
  at 3NT but not at 4♣️, or a later suit preference) are not covered yet.
- **Facts must be declared**: a raise that forgets `establishes` leaves the suit unagreed. The
  round-trip tests check the declared facts for every built-in raise and forcing response.

## Migration Plan
1. Phases 1–8 in `tasks.md` are refactors plus new primitives; existing round-trip tests must pass
   unchanged after each phase.
2. Phase 9 moves the openings to `auction: [Pass*]` and the 1NT response files to
   `auction: [Pass*, 1NT, Pass]`, and reorders the Colin card. Existing Colin 1NT response tests
   switch to full auction histories instead of a lone partner call.
3. Phases 9–10 add data and tests only.

## Open Questions
- The Two Over One card includes `multi_two_diamond`, so it inherits the Colin-style `2NT`
  weakness transfer, which comes before the standard `2NT` invitation in `strong_notrump`. With
  6–9 HCP and 6+ spades responder bids `2NT` as a transfer; with 8–9 HCP and no long spades the
  same bid is an invitation. The interpreter will read it as either. Should the 1NT responses in
  `multi_two_diamond` move to a Colin-only file?
- Splinter range: 12–15 HCP is proposed. Some partnerships use 10–13 or 13–16.
- Opener's 18–19 point unbalanced hands have no opening on either card (one-level openings stop
  at 17, `2♣️` starts at 20). This is out of scope here but affects rebid ranges later.
- Slam bidding has no conventions yet: this change adds responses to `2♣️` and opener's rebid
  (task 10.4), but nothing that asks for aces or bids a slam. `TestAuction_GrandSlam` in
  `slam_test.go` holds a deal on which East–West should reach a grand slam (West deals and opens
  `2♣️` with 21 HCP and seven solid hearts; East has six solid clubs) and is skipped until slam
  bidding exists. Should slam bidding (for example key-card Blackwood after an agreed suit, using
  the `agreed` fact) be a follow-up change, and which card should carry it?
