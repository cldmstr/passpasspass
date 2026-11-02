# Tasks: update-partnership-evaluation-limits

## 1. Implementation
- [x] 1.1 `NewPartnershipEval` caps the combined HCP at `MaxDeckHCP` and each combined suit length at `HandSize`, minimum and maximum alike
- [x] 1.2 Tests: opposite an unconstrained partner the maximums are 40 HCP and 13 cards; minimums above the caps are capped too

## 2. Archive
- [x] 2.1 Apply the spec delta to `openspec/specs/` and move this change to `openspec/changes/archive/`
