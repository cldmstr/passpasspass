# Change: Print hands from spades down to clubs

## Why
Bridge hands are conventionally written top down from spades to clubs, the order of the suits'
rank. `ppp` prints clubs first, which makes hands harder to read at a glance.

## What Changes
- A hand prints one line per suit from spades down to clubs. Cards within a suit stay in
  ascending rank, and a void stays an empty line.
- A hand's shape, such as `[6,4,2,1]` next to each hand, lists the suits in the same order.
  Distributions stay indexed clubs first in code; only their printed form changes.

## Impact
- Affected specs: `cards-and-hands` (Suit order, Hand display)
- Affected code: `hand.go` (`Hand.String`, `Distribution.String`)
