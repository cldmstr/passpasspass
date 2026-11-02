package passpasspass

import (
	"errors"
	"fmt"
	"maps"
	"slices"

	"gopkg.in/yaml.v3"
)

// conditionParser turns the value of one condition key into a Condition.
type conditionParser func(v *yaml.Node) (Condition, error)

// conditionParsers holds a parser for every condition key. It is filled in init because the
// parsers for `all` and `any` call parseCondition, which reads the table.
var conditionParsers map[string]conditionParser

// conditionKeys lists the keys of conditionParsers in order, for error messages.
var conditionKeys []string

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
	conditionKeys = slices.Sorted(maps.Keys(conditionParsers))
}

// parseCondition turns a YAML condition into a Condition. A list is an `all` of its items;
// a mapping is an `all` of its keys, each parsed by the parser registered for that key.
// A list with one item or a mapping with one key returns that condition directly.
// Errors carry the line of the innermost offending key or value and the keys leading to it.
func parseCondition(v *yaml.Node) (Condition, error) {
	switch {
	case v.Kind == yaml.SequenceNode:
		conditions, err := parseConditionList(v)
		if err != nil {
			return nil, err
		}
		return conjunction(conditions), nil

	case v.Kind == yaml.MappingNode && len(v.Content) > 0:
		return parseConditionKeys(v)

	case v.Kind == yaml.MappingNode:
		return nil, errorAt(v, errors.New("a condition mapping must not be empty"))

	default:
		return nil, errorAt(v, errors.New("a condition must be a mapping or a list"))
	}
}

// parseConditionKeys parses every key of a non-empty mapping, in file order.
func parseConditionKeys(v *yaml.Node) (Condition, error) {
	if err := checkKeys(v, "condition", conditionKeys); err != nil {
		return nil, err
	}
	conditions := make([]Condition, 0, len(v.Content)/2)
	for i := 0; i+1 < len(v.Content); i += 2 {
		key := v.Content[i]
		c, err := conditionParsers[key.Value](v.Content[i+1])
		if err != nil {
			return nil, keyError(key, err)
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
		return 0, 0, errors.New("must be [min, max] with two whole numbers")
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
		return nil, errors.New("must be a call such as 1NT or Pass")
	}
	return NewBidFromString(s)
}

// parseConditionList decodes a list of conditions, each parsed with parseCondition.
func parseConditionList(v *yaml.Node) ([]Condition, error) {
	if v.Kind != yaml.SequenceNode {
		return nil, errorAt(v, errors.New("must be a list of conditions"))
	}
	if len(v.Content) == 0 {
		return nil, errorAt(v, errors.New("a condition list must not be empty"))
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
	lo, hi, err := parseRange(v, MaxDeckHCP)
	if err != nil {
		return nil, err
	}
	return NewHCPIntervalCondition(lo, hi), nil
}

// parsePartnershipHCP parses `partnership_hcp: [min, max]`.
func parsePartnershipHCP(v *yaml.Node) (Condition, error) {
	lo, hi, err := parseRange(v, MaxDeckHCP)
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
		key := v.Content[i]
		s, err := NewSuitFromString(key.Value)
		if err != nil || s == NoTrump {
			return nil, keyError(key, errors.New("not a suit (clubs, diamonds, hearts or spades)"))
		}
		if seen[s] {
			return nil, keyError(key, errors.New("suit appears more than once"))
		}
		seen[s] = true
		lo, hi, err := parseRange(v.Content[i+1], HandSize)
		if err != nil {
			return nil, keyError(key, err)
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
	conditions, err := parseConditionList(v)
	if err != nil {
		return nil, err
	}
	return conjunction(conditions), nil
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
