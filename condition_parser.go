package passpasspass

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// conditionParser turns the value of one condition key into a Condition.
type conditionParser func(v *yaml.Node) (Condition, error)

// conditionParsers holds a parser for every condition key. It is filled in init because the
// parsers for `all` and `any` call parseCondition, which reads the table.
var conditionParsers map[string]conditionParser

func init() {
	conditionParsers = map[string]conditionParser{
		"hcp":             parseHCP,
		"suit_length":     parseSuitLength,
		"balanced":        parseBalanced,
		"rule_of_20":      parseRuleOf20,
		"no_prior_bids":   parseNoPriorBids,
		"all":             parseAll,
		"any":             parseAny,
		"partner_bid":     parsePartnerBid,
		"partnership_hcp": parsePartnershipHCP,
	}
}

// parseCondition turns a YAML condition into a Condition. A list is an `all` of its items;
// a mapping is an `all` of its keys, each parsed by the parser registered for that key.
// A list with one item or a mapping with one key returns that condition directly.
func parseCondition(v *yaml.Node) (Condition, error) {
	switch v.Kind {
	case yaml.SequenceNode:
		conditions, err := parseConditionList(v)
		if err != nil {
			return nil, err
		}
		return conjunction(conditions), nil

	case yaml.MappingNode:
		if len(v.Content) == 0 {
			return nil, errors.New("a condition mapping must not be empty")
		}
		return parseConditionKeys(v)

	default:
		return nil, errors.New("a condition must be a mapping or a list")
	}
}

// parseConditionKeys parses every key of a non-empty mapping, in file order.
func parseConditionKeys(v *yaml.Node) (Condition, error) {
	conditions := make([]Condition, 0, len(v.Content)/2)
	seen := make(map[string]bool, len(v.Content)/2)
	for i := 0; i+1 < len(v.Content); i += 2 {
		key, value := v.Content[i].Value, v.Content[i+1]
		if seen[key] {
			return nil, fmt.Errorf("%s: condition key appears more than once", key)
		}
		seen[key] = true
		parse, ok := conditionParsers[key]
		if !ok {
			known := slices.Sorted(maps.Keys(conditionParsers))
			return nil, fmt.Errorf("%s: unknown condition key (known keys: %s)", key, strings.Join(known, ", "))
		}
		c, err := parse(value)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}
		conditions = append(conditions, c)
	}
	return conjunction(conditions), nil
}

// conjunction returns the only condition of a single-element list, or an `all` of them.
func conjunction(conditions []Condition) Condition {
	if len(conditions) == 1 {
		return conditions[0]
	}
	return NewAllCondition(conditions...)
}

// parseRange decodes a range value `[min, max]` with 0 ≤ min ≤ max ≤ ceiling. Every range counts
// points or cards, so none goes below 0.
func parseRange(v *yaml.Node, ceiling int) (lo, hi int, err error) {
	var r []int
	if err := v.Decode(&r); err != nil {
		return 0, 0, err
	}
	if len(r) != 2 {
		return 0, 0, fmt.Errorf("must be [min, max], got %d values", len(r))
	}
	lo, hi = r[0], r[1]
	if lo > hi {
		return 0, 0, fmt.Errorf("min %d is greater than max %d", lo, hi)
	}
	if lo < 0 || hi > ceiling {
		return 0, 0, fmt.Errorf("[%d, %d] is outside 0 to %d", lo, hi, ceiling)
	}
	return lo, hi, nil
}

// parseFlag accepts only the YAML boolean `true`.
func parseFlag(v *yaml.Node) error {
	var b bool
	if v.Kind != yaml.ScalarNode || v.ShortTag() != "!!bool" || v.Decode(&b) != nil || !b {
		return errors.New("only true is accepted")
	}
	return nil
}

// parseBid decodes a call such as "1NT" or "Pass".
func parseBid(v *yaml.Node) (*Bid, error) {
	var s string
	if err := v.Decode(&s); err != nil {
		return nil, err
	}
	return NewBidFromString(s)
}

// parseConditionList decodes a list of conditions, each parsed with parseCondition.
func parseConditionList(v *yaml.Node) ([]Condition, error) {
	if v.Kind != yaml.SequenceNode {
		return nil, errors.New("must be a list of conditions")
	}
	if len(v.Content) == 0 {
		return nil, errors.New("a condition list must not be empty")
	}
	conditions := make([]Condition, 0, len(v.Content))
	for _, item := range v.Content {
		c, err := parseCondition(item)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, c)
	}
	return conditions, nil
}

// parseHCP parses `hcp: [min, max]`.
func parseHCP(v *yaml.Node) (Condition, error) {
	lo, hi, err := parseRange(v, deckHCP)
	if err != nil {
		return nil, err
	}
	return NewHCPIntervalCondition(lo, hi), nil
}

// parsePartnershipHCP parses `partnership_hcp: [min, max]`.
func parsePartnershipHCP(v *yaml.Node) (Condition, error) {
	lo, hi, err := parseRange(v, deckHCP)
	if err != nil {
		return nil, err
	}
	return NewPartnershipHCPCondition(lo, hi), nil
}

// parseSuitLength parses `suit_length: {suit: [min, max], ...}`; unlisted suits are unconstrained.
func parseSuitLength(v *yaml.Node) (Condition, error) {
	if v.Kind != yaml.MappingNode || len(v.Content) == 0 {
		return nil, errors.New("must be a mapping of suits to [min, max]")
	}
	r := DistributionRange{{0, HandSize}, {0, HandSize}, {0, HandSize}, {0, HandSize}}
	seen := make(map[Suit]bool, len(r))
	for i := 0; i+1 < len(v.Content); i += 2 {
		name := v.Content[i].Value
		s, err := NewSuitFromString(name)
		if err != nil || s == NoTrump {
			return nil, fmt.Errorf("%q is not a suit (clubs, diamonds, hearts or spades)", name)
		}
		if seen[s] {
			return nil, fmt.Errorf("suit %q appears more than once", name)
		}
		seen[s] = true
		lo, hi, err := parseRange(v.Content[i+1], HandSize)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		r[s] = SuitRange{lo, hi}
	}
	return NewMinMaxDistributionCondition(r), nil
}

// parseBalanced parses `balanced: true`.
func parseBalanced(v *yaml.Node) (Condition, error) {
	if err := parseFlag(v); err != nil {
		return nil, err
	}
	return NewBalancedDistributionCondition(), nil
}

// parseRuleOf20 parses `rule_of_20: true`.
func parseRuleOf20(v *yaml.Node) (Condition, error) {
	if err := parseFlag(v); err != nil {
		return nil, err
	}
	return NewRuleOfTwentyCondition(), nil
}

// parseNoPriorBids parses `no_prior_bids: true`.
func parseNoPriorBids(v *yaml.Node) (Condition, error) {
	if err := parseFlag(v); err != nil {
		return nil, err
	}
	return NewNoPriorBidsCondition(), nil
}

// parseAll parses `all: [condition, ...]`, which passes when every condition passes.
func parseAll(v *yaml.Node) (Condition, error) {
	if v.Kind != yaml.SequenceNode {
		return nil, errors.New("must be a list of conditions")
	}
	return parseCondition(v)
}

// parseAny parses `any: [condition, ...]`, which passes when at least one condition passes.
func parseAny(v *yaml.Node) (Condition, error) {
	conditions, err := parseConditionList(v)
	if err != nil {
		return nil, err
	}
	return NewMultiCondition(conditions...), nil
}

// parsePartnerBid parses `partner_bid: <call>`.
func parsePartnerBid(v *yaml.Node) (Condition, error) {
	bid, err := parseBid(v)
	if err != nil {
		return nil, err
	}
	return partnerBidCondition{bid: bid}, nil
}
