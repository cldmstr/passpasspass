# Change: Cap the partnership evaluation at what the deck holds

## Why
The partnership evaluation adds partner's shown range to the player's own hand without limit.
Opposite a partner who has shown nothing, a player with 25 HCP and 7 hearts is told the
partnership holds up to 62 HCP and 20 hearts. The verbose output (`ppp -v`) prints these numbers
at every call, and a condition on the combined range would be judged against impossible values.

## What Changes
- The combined HCP range is capped at 40, the points in the deck.
- The combined length in each suit is capped at 13, the cards in a suit.
- Both the minimum and the maximum are capped.

## Impact
- Affected specs: `auction` (Partnership evaluation)
- Affected code: `auction.go` (`NewPartnershipEval`)
- `add-partner-aware-responses` modifies the same requirement and carries the caps.
