# Design: Condition nodes

## Context
`loader.go` decodes every condition into `conditionNode`, a struct with one optional field per
YAML key, and `buildCondition` turns it into a `Condition` closure. The struct suggests one key
per node, but a node is meant to hold any number of keys, all of which must hold. Because
closures cannot be inspected, `Convention` also keeps the raw `conditionNode` so that inference
(`applyBounds`) can walk it again.

## Goals / Non-Goals

**Goals**
- Read every key of a condition mapping, and reject anything the loader does not understand.
- Validate each key's value once, where the key is parsed, with errors that point at the line.
- Make conditions inspectable values, so later changes can ask a condition what it shows
  without a second walk over YAML.
- Make adding a condition key one parser, one node type and one table entry.

**Non-Goals**
- New condition keys, `not`, or new inference. Those belong to `add-partner-aware-responses`.
- Changing what any existing, valid condition means.

## Decisions

### D1: Conditions are parsed from `yaml.Node`, key by key
`conditions` and every nested condition are read as a `*yaml.Node`:

```go
// parseCondition turns a YAML condition into a Condition. A list is an `all` of its items;
// a mapping is an `all` of its keys, each parsed by the parser registered for that key.
func parseCondition(v *yaml.Node) (Condition, error)
```

- **Sequence**: every item is parsed with `parseCondition`; the result is the conjunction.
  A list with one item returns that item's condition directly.
- **Mapping**: every key/value pair is parsed by the parser registered for the key, in file
  order; the result is the conjunction. A mapping with one key returns that key's condition
  directly. A key that appears twice is a load error.
- **Anything else**, an empty sequence or an empty mapping: load error.

Each parser decodes its value with `v.Decode(&target)` into the shape its key needs (`[]int`,
`map[string][]int`, `bool`, `string`, or a nested condition). yaml.v3 reports type mismatches
itself, keeps keys in file order, and gives a line number for every key and value. It rejects
duplicate keys only when decoding into a Go map or struct, so `parseCondition` and the
`suit_length` parser, which walk mapping nodes themselves, reject them. `map[string]any` was rejected: it loses order and line numbers, and needs a type assertion
per level instead of one `Decode`.

`conditionNode`, its yaml tags, `buildCondition` and the hand-written `UnmarshalYAML` on
`yamlConvention` are removed.

### D2: A parser table is the single list of keys
```go
type conditionParser func(v *yaml.Node) (Condition, error)

var conditionParsers map[string]conditionParser // filled in init, see below
```

| Key | Value | Limits |
|-----|-------|--------|
| `hcp` | `[min, max]` | 0–40 |
| `suit_length` | `{suit: [min, max], ...}` | 0–13; clubs, diamonds, hearts or spades only |
| `balanced`, `rule_of_20`, `no_prior_bids` | `true` | only `true` |
| `all`, `any` | non-empty list of conditions | |
| `partner_bid` | a call, e.g. `"1NT"` or `Pass` | must parse as a bid |
| `partnership_hcp` | `[min, max]` | 0–40 |

HCP limits go to 40, the points in the deck, rather than 37, because the built-in data uses
`[8, 40]` and `[27, 40]` to mean "no upper limit".

`all` and `any` call `parseCondition`, which reads the table, so a package-level map literal would
form an initialisation cycle. The table is filled in `init`. Shared helpers keep each parser to a
few lines: `parseRange(v, ceiling)`, `parseFlag(v)` and `parseBid(v)`. Every range counts points
or cards, so its lower limit is always 0.

### D3: `Condition` is an interface; each key has a node type
```go
// Condition decides whether a convention applies in the given auction context.
type Condition interface {
    Eval(ctx *AuctionContext) bool
}

type hcpCondition struct{ lo, hi int }
type allCondition struct{ children []Condition }
// ... one type per key
```

The existing `New*Condition` constructors return these types as `Condition`, so code and tests that
build conditions directly keep compiling. `Convention.NextBid` calls `Eval`.

`add-partner-aware-responses` adds a second method to the interface for reading what a call shows
(its `Shown`, design D4 there). Every node type then implements it next to `Eval`, with no
second `switch`.

### D4: Flag keys accept only `true`
`balanced: false` today means "balanced". Reading it as "not balanced" would give a second way to
write negation once `not` exists, and would make three keys behave differently from every other
key, so `false` and any non-boolean value are load errors that state that only `true` is
accepted. When `add-partner-aware-responses` adds `not`, the message also suggests `not`.

### D5: Strict loading at every level
Card files, sub-convention files, convention entries and conditions reject unknown and repeated
keys. One helper, `checkKeys`, checks every level against its list of keys: the parser table for
conditions, fixed lists for the rest. `yaml.Decoder.KnownFields(true)` was rejected: its errors
span several lines and name Go struct types, and `yaml.Node.Decode` cannot use it, so convention
entries and conditions would need their own check anyway. A convention entry needs `bid` and a
non-empty `conditions`. Every load error names the file (card or sub-convention), the convention's bid
where there is one, and the line of the offending key. It starts with `file:line`, as compiler
errors do, so editors and terminals can jump to it; the bid and the keys leading to the offending
one follow, the keys joined with dots so they stand apart from the message:

```
conventions/baron.yaml:8: bid 2S: suit_lenght: unknown condition key (known keys: all, any, balanced, ...)
conventions/baron.yaml:12: bid 2S: any.all.suit_length.hearts: [0, 14] is outside 0 to 13
cards/colin.yaml: include barn: conventions/barn.yaml: file does not exist
```

Errors inside nested conditions keep the line of the innermost offending key or value. yaml.v3
syntax errors are rewritten into the same form.

### D6: Inference reads parsed nodes until it is replaced
`Convention` drops its `node conditionNode` field and keeps its parsed `Condition`s. `applyBounds`
switches on node types (`hcpCondition`, `suitLengthCondition`, `allCondition`, `anyCondition`)
instead of `conditionNode` fields, with the same results as today. `add-partner-aware-responses`
removes it in its task 5.7.

### D7: Unused keys are removed first
`rho_bid`, `lho_bid`, `lho_min_hcp`, `rho_min_hcp`, `heart_fit`, `spade_fit`, `major_fit` and
`minor_fit` are used by no convention file; only their own unit tests call them. They are removed,
together with `PartnershipEval.BestMajorFit` and `BestMinorFit`, before the parser table exists,
so no parser, validation rule or inference rule is written for them.

| Removed | Replacement |
|---------|-------------|
| `rho_bid`, `lho_bid` | `auction` patterns (`add-partner-aware-responses`), which can name any seat's call |
| `heart_fit`, `spade_fit` | `fit: {<suit>: [min, max]}` (`add-partner-aware-responses`), for any suit |
| `major_fit`, `minor_fit` | `any` of two `fit` conditions |
| `lho_min_hcp`, `rho_min_hcp` | none; opponents' calls are out of scope until overcalls are planned |

`partner_bid` and `no_prior_bids` stay: the built-in data uses them until
`add-partner-aware-responses` moves it to `auction` patterns. `partnership_hcp` stays because that
change uses it. `PartnershipEval.SuitFit` stays for the planned `fit`.

## Risks / Trade-offs
- **Custom cards using a removed key stop loading**: the error names the key and its line. No
  built-in file is affected.
- **No single struct lists every key with its type**: the parser table takes that role, and each
  parser's doc comment states the value's shape.
- **Stricter loading may reject custom cards that load today**: every rejected case is one that is
  silently wrong today, and the error names the line. The built-in data loads unchanged.
- **Type switch in `applyBounds`** (D6): a stopgap that keeps today's inference working; it is
  removed with `applyBounds` by the next change.

## Migration Plan
1. Add the node types and the `Condition` interface; keep the constructors' signatures.
2. Add `parseCondition`, the parser table and strict file decoding; switch the loader to them.
3. Move `applyBounds` to node types and delete `conditionNode` and `buildCondition`.
4. Existing tests pass unchanged after each step; new tests cover every case in the proposal.
