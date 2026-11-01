# Tasks: refactor-condition-nodes

## 1. Condition nodes
- [x] 1.1 Change `Condition` to an interface with `Eval(ctx *AuctionContext) bool`
- [x] 1.2 Add one node type per existing key (design D2) and make each `New*Condition` constructor return it; `NewAllCondition` and `NewMultiCondition` return `allCondition` and `anyCondition`
- [x] 1.3 `Convention.NextBid` calls `Eval`
- [x] 1.4 Existing condition, convention and auction tests pass unchanged

## 2. Remove unused condition keys
- [x] 2.1 Remove `rho_bid`, `lho_bid`, `lho_min_hcp`, `rho_min_hcp`, `heart_fit`, `spade_fit`, `major_fit` and `minor_fit`: constructors, condition types, `conditionNode` fields, `buildCondition` cases and their tests (design D7)
- [x] 2.2 Remove `PartnershipEval.BestMajorFit` and `BestMinorFit` and their tests
- [x] 2.3 All built-in data loads and bids exactly as before

## 3. Parsing conditions key by key
- [x] 3.1 Add `parseCondition(*yaml.Node)`: a sequence or a mapping becomes the conjunction of its items or keys (design D1); a mapping with one key returns that key's condition
- [x] 3.2 Add the parser table, filled in `init`, with a parser for every existing key and the shared helpers `parseRange`, `parseFlag` and `parseBid` (design D2)
- [x] 3.3 Unknown condition keys, empty mappings, empty lists, and values that are neither a mapping nor a list are load errors
- [x] 3.4 Tests: two keys in one mapping both apply; keys in a list and in a mapping give the same result; nested `all` and `any` with several keys per item; unknown key; empty `conditions`, `all: []` and `any: []`

## 4. Validation
- [ ] 4.1 Ranges: exactly two integers, `min ≤ max`, within the limits in design D2
- [ ] 4.2 `suit_length`: only clubs, diamonds, hearts and spades in any accepted spelling; `NT` is a load error
- [ ] 4.3 Flag keys (`balanced`, `rule_of_20`, `no_prior_bids`) accept only `true` (design D4)
- [ ] 4.4 `partner_bid` parses its bid at load time
- [ ] 4.5 Tests, one per case in the proposal table and per limit: `hcp: [5]`, `hcp: [12, 10]`, `hcp: [0, 41]`, `suit_length: {NT: [1, 2]}`, `suit_length: {hearts: [0, 14]}`, `balanced: false`, `balanced: yes please`, `partner_bid: "1N"`; boundary values that still load (`hcp: [0, 40]`, `suit_length: {spades: [0, 13]}`)

## 5. Strict files and error messages
- [ ] 5.1 Decode card and sub-convention files with unknown keys rejected; a convention entry needs `bid` and `conditions`
- [ ] 5.2 Load errors name the file, the convention's bid and the line of the offending key (design D5)
- [ ] 5.3 Tests: misspelt `condtions`, `bid` missing, unknown top-level key in a card; error text contains the file, bid and line

## 6. Switching over and cleaning up
- [ ] 6.1 `Convention` keeps its parsed `Condition`s and drops `node`; `buildConventions` uses `parseCondition`
- [ ] 6.2 `applyBounds` switches on node types instead of `conditionNode` fields (design D6)
- [ ] 6.3 Remove `conditionNode`, `buildCondition` and `yamlConvention.UnmarshalYAML`
- [ ] 6.4 Remove the `buildCondition` exclusion from `.golangci.yml`; `golangci-lint run ./...` is clean
- [ ] 6.5 All built-in data loads and every existing round-trip test passes unchanged

## 7. Documentation and archive
- [ ] 7.1 Update the condition syntax and load errors in `openspec/project.md`, if described there
- [ ] 7.2 After merge, apply the spec deltas to `openspec/specs/` and move this change to `openspec/changes/archive/`
