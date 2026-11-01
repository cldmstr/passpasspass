# Tasks: refactor-condition-nodes

## 1. Condition nodes
- [ ] 1.1 Change `Condition` to an interface with `Eval(ctx *AuctionContext) bool`
- [ ] 1.2 Add one node type per existing key (design D2) and make each `New*Condition` constructor return it; `NewAllCondition` and `NewMultiCondition` return `allCondition` and `anyCondition`
- [ ] 1.3 `Convention.NextBid` calls `Eval`
- [ ] 1.4 Existing condition, convention and auction tests pass unchanged

## 2. Parsing conditions key by key
- [ ] 2.1 Add `parseCondition(*yaml.Node)`: a sequence or a mapping becomes the conjunction of its items or keys (design D1); a mapping with one key returns that key's condition
- [ ] 2.2 Add the parser table, filled in `init`, with a parser for every existing key and the shared helpers `parseRange`, `parseFlag` and `parseBid` (design D2)
- [ ] 2.3 Unknown condition keys, empty mappings, empty lists, and values that are neither a mapping nor a list are load errors
- [ ] 2.4 Tests: two keys in one mapping both apply; keys in a list and in a mapping give the same result; nested `all` and `any` with several keys per item; unknown key; empty `conditions`, `all: []` and `any: []`

## 3. Validation
- [ ] 3.1 Ranges: exactly two integers, `min ≤ max`, within the limits in design D2
- [ ] 3.2 `suit_length`: only clubs, diamonds, hearts and spades in any accepted spelling; `NT` is a load error
- [ ] 3.3 Flag keys (`balanced`, `rule_of_20`, `no_prior_bids`) accept only `true` (design D4)
- [ ] 3.4 `partner_bid`, `rho_bid` and `lho_bid` parse their bid at load time
- [ ] 3.5 Tests, one per case in the proposal table and per limit: `hcp: [5]`, `hcp: [12, 10]`, `hcp: [0, 41]`, `suit_length: {NT: [1, 2]}`, `suit_length: {hearts: [0, 14]}`, `balanced: false`, `balanced: yes please`, `partner_bid: "1N"`, `heart_fit: [8, 14]`, `spade_fit: [7, 13]`; boundary values that still load (`hcp: [0, 40]`, `suit_length: {spades: [0, 13]}`, `heart_fit: [8, 13]`)

## 4. Strict files and error messages
- [ ] 4.1 Decode card and sub-convention files with unknown keys rejected; a convention entry needs `bid` and `conditions`
- [ ] 4.2 Load errors name the file, the convention's bid and the line of the offending key (design D5)
- [ ] 4.3 Tests: misspelt `condtions`, `bid` missing, unknown top-level key in a card; error text contains the file, bid and line

## 5. Switching over and cleaning up
- [ ] 5.1 `Convention` keeps its parsed `Condition`s and drops `node`; `buildConventions` uses `parseCondition`
- [ ] 5.2 `applyBounds` switches on node types instead of `conditionNode` fields (design D6)
- [ ] 5.3 Remove `conditionNode`, `buildCondition` and `yamlConvention.UnmarshalYAML`
- [ ] 5.4 Remove the `buildCondition` exclusion from `.golangci.yml`; `golangci-lint run ./...` is clean
- [ ] 5.5 All built-in data loads and every existing round-trip test passes unchanged

## 6. Documentation and archive
- [ ] 6.1 Update the condition syntax and load errors in `openspec/project.md`, if described there
- [ ] 6.2 After merge, apply the spec deltas to `openspec/specs/` and move this change to `openspec/changes/archive/`
