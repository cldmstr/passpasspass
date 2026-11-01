# Change: Parse conditions key by key into typed condition nodes

## Why
Conditions are decoded into `conditionNode`, a struct with one optional field per YAML key, and
`buildCondition` turns it into a closure with a `switch` that stops at the first field it finds.
This loses or misreads conditions without telling the author:

| YAML | Today |
|------|-------|
| `hcp: [5]` | panics with index out of range |
| `suit_length: {NT: [1, 2]}` | panics with index out of range |
| `balanced: false` (also `rule_of_20`, `no_prior_bids`) | behaves like `true` |
| `hcp` and `suit_length` in one mapping | `suit_length` is ignored |
| a misspelt key next to a valid one | ignored |
| `partner_bid: "1N"` | loads, never matches |
| `conditions: []` | loads, always matches |

Each condition also ends up as an opaque closure, so inference cannot ask a condition what it
shows. It walks the raw `conditionNode` again in `applyBounds`, a second `switch` with its own,
looser rules. `add-partner-aware-responses` would add a third walk and about eight new keys, each
needing a struct field and a case in every `switch`.

## What Changes
- **Key-by-key parsing**: the loader reads conditions from `yaml.Node`. A mapping is the
  conjunction of all its keys; each key is parsed by the parser registered for it, which decodes
  the value into the shape that key needs. A list is the conjunction of its items. `conditionNode`
  and `buildCondition` are removed.
- **Typed condition nodes**: `Condition` becomes an interface with `Eval(ctx) bool`. Each key has
  a small node type that implements it. The `New*Condition` constructors return these types, so
  callers keep working.
- **Strict loading**: unknown keys in a card, a sub-convention file, a convention entry or a
  condition are load errors. Errors name the file, the convention's bid and the line.
- **Validation per key**: ranges need exactly two values with `min ≤ max` inside the key's limits;
  `suit_length` accepts only the four suits; flag keys accept only `true`; seat bids are parsed at
  load time; empty conditions are load errors.
- **Inference on parsed nodes**: `Convention` keeps its parsed conditions instead of the raw YAML
  node, and `applyBounds` reads those until `add-partner-aware-responses` replaces it.

## Impact
- Affected specs: `convention-files`, `bidding-conditions`
- Affected code: `loader.go` (parsing, parser table, strict decoding, errors), `condition.go`
  (`Condition` interface, node types), `convention.go` (`Convention` holds `Condition`s, no raw
  node), `auction.go` (`applyBounds` over node types), `.golangci.yml` (drop the `buildCondition`
  exclusion)
- No change to the built-in data: every existing file loads and bids as before.
- Prerequisite for `add-partner-aware-responses`, which adds keys and a second method on the
  nodes to read what a call shows.
- Unchanged: `cards-and-hands`, `dealing`, `convention-evaluation`, `auction`,
  `standard-conventions`, `cli`
