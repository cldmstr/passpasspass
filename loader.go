package passpasspass

import (
	"fmt"
	"io/fs"

	"gopkg.in/yaml.v3"
)

// yamlConvention is a single (bid, conditions) entry in a YAML file.
type yamlConvention struct {
	Bid        string        `yaml:"bid"`
	Conditions conditionNode `yaml:"conditions"`
}

// yamlSubConvention is the top-level structure of a sub-convention file.
type yamlSubConvention struct {
	Name        string           `yaml:"name"`
	Conventions []yamlConvention `yaml:"conventions"`
}

// yamlCard is the top-level structure of a convention card file.
type yamlCard struct {
	Name        string           `yaml:"name"`
	Includes    []string         `yaml:"includes"`
	Conventions []yamlConvention `yaml:"conventions"`
}

// UnmarshalYAML handles conditions being either a mapping (single node)
// or a sequence (implicit all).
func (yc *yamlConvention) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("convention must be a mapping")
	}
	for i := 0; i+1 < len(value.Content); i += 2 {
		key := value.Content[i].Value
		val := value.Content[i+1]
		switch key {
		case "bid":
			yc.Bid = val.Value
		case "conditions":
			switch val.Kind {
			case yaml.SequenceNode:
				var children []conditionNode
				if err := val.Decode(&children); err != nil {
					return err
				}
				yc.Conditions = conditionNode{All: children}
			case yaml.MappingNode:
				if err := val.Decode(&yc.Conditions); err != nil {
					return err
				}
			default:
				return fmt.Errorf("conditions must be a mapping or sequence")
			}
		}
	}
	return nil
}

// conditionNode is a polymorphic YAML condition. Exactly one field is
// populated after unmarshaling.
type conditionNode struct {
	HCP        []int            `yaml:"hcp"`
	Balanced   *bool            `yaml:"balanced"`
	RuleOf20   *bool            `yaml:"rule_of_20"`
	SuitLength map[string][]int `yaml:"suit_length"`
	All        []conditionNode  `yaml:"all"`
	Any        []conditionNode  `yaml:"any"`

	// history-aware
	PartnerBid  *string `yaml:"partner_bid"`
	RHOBid      *string `yaml:"rho_bid"`
	LHOBid      *string `yaml:"lho_bid"`
	NoPriorBids *bool   `yaml:"no_prior_bids"`

	// partnership evaluation
	PartnershipHCP []int `yaml:"partnership_hcp"`
	MajorFit       []int `yaml:"major_fit"`
	MinorFit       []int `yaml:"minor_fit"`
	HeartFit       []int `yaml:"heart_fit"`
	SpadeFit       []int `yaml:"spade_fit"`

	// opponent HCP bounds
	LHOMinHCP *int `yaml:"lho_min_hcp"`
	RHOMinHCP *int `yaml:"rho_min_hcp"`
}

func buildCondition(n conditionNode) (Condition, error) {
	switch {
	case n.HCP != nil:
		return NewHCPIntervalCondition(n.HCP[0], n.HCP[1]), nil

	case n.Balanced != nil:
		return NewBalancedDistributionCondition(), nil

	case n.RuleOf20 != nil:
		return NewRuleOfTwentyCondition(), nil

	case n.SuitLength != nil:
		r := DistributionRange{{0, 13}, {0, 13}, {0, 13}, {0, 13}}
		for suit, rng := range n.SuitLength {
			s, err := NewSuitFromString(suit)
			if err != nil {
				return nil, fmt.Errorf("unknown suit %q in suit_length", suit)
			}
			i := int(s)
			if len(rng) != 2 {
				return nil, fmt.Errorf("suit_length %q must have exactly 2 values [min, max]", suit)
			}
			r[i] = SuitRange{rng[0], rng[1]}
		}
		return NewMinMaxDistributionCondition(r), nil

	case n.All != nil:
		conditions := make([]Condition, len(n.All))
		for i, child := range n.All {
			c, err := buildCondition(child)
			if err != nil {
				return nil, err
			}
			conditions[i] = c
		}
		return NewAllCondition(conditions...), nil

	case n.Any != nil:
		conditions := make([]Condition, len(n.Any))
		for i, child := range n.Any {
			c, err := buildCondition(child)
			if err != nil {
				return nil, err
			}
			conditions[i] = c
		}
		return NewMultiCondition(conditions...), nil

	case n.PartnerBid != nil:
		return NewPartnerBidCondition(*n.PartnerBid), nil

	case n.RHOBid != nil:
		return NewRHOBidCondition(*n.RHOBid), nil

	case n.LHOBid != nil:
		return NewLHOBidCondition(*n.LHOBid), nil

	case n.NoPriorBids != nil:
		return NewNoPriorBidsCondition(), nil

	case n.PartnershipHCP != nil:
		if len(n.PartnershipHCP) != 2 {
			return nil, fmt.Errorf("partnership_hcp must have exactly 2 values [min, max]")
		}
		return NewPartnershipHCPCondition(n.PartnershipHCP[0], n.PartnershipHCP[1]), nil

	case n.MajorFit != nil:
		if len(n.MajorFit) != 2 {
			return nil, fmt.Errorf("major_fit must have exactly 2 values [min, max]")
		}
		return NewMajorFitCondition(n.MajorFit[0], n.MajorFit[1]), nil

	case n.MinorFit != nil:
		if len(n.MinorFit) != 2 {
			return nil, fmt.Errorf("minor_fit must have exactly 2 values [min, max]")
		}
		return NewMinorFitCondition(n.MinorFit[0], n.MinorFit[1]), nil

	case n.HeartFit != nil:
		if len(n.HeartFit) != 2 {
			return nil, fmt.Errorf("heart_fit must have exactly 2 values [min, max]")
		}
		return NewHeartFitCondition(n.HeartFit[0], n.HeartFit[1]), nil

	case n.SpadeFit != nil:
		if len(n.SpadeFit) != 2 {
			return nil, fmt.Errorf("spade_fit must have exactly 2 values [min, max]")
		}
		return NewSpadeFitCondition(n.SpadeFit[0], n.SpadeFit[1]), nil

	case n.LHOMinHCP != nil:
		return NewLHOMinHCPCondition(*n.LHOMinHCP), nil

	case n.RHOMinHCP != nil:
		return NewRHOMinHCPCondition(*n.RHOMinHCP), nil

	default:
		return nil, fmt.Errorf("empty condition node")
	}
}

func buildConventions(entries []yamlConvention) ([]*Convention, error) {
	conventions := make([]*Convention, 0, len(entries))
	for _, entry := range entries {
		bid, err := NewBidFromString(entry.Bid)
		if err != nil {
			return nil, fmt.Errorf("invalid bid %q: %w", entry.Bid, err)
		}
		condition, err := buildCondition(entry.Conditions)
		if err != nil {
			return nil, fmt.Errorf("invalid conditions for bid %q: %w", entry.Bid, err)
		}
		c := NewConvention(bid, condition)
		c.node = entry.Conditions
		conventions = append(conventions, c)
	}
	return conventions, nil
}

// Loader reads YAML convention files and produces ConventionSheets.
type Loader struct {
	fsys fs.FS
}

// NewLoader creates a Loader that reads convention files from fsys.
// Sub-convention includes are resolved as "conventions/<name>.yaml"
// within the root of fsys.
func NewLoader(fsys fs.FS) *Loader {
	return &Loader{fsys: fsys}
}

// LoadCard reads a convention card file at cardPath (relative to the
// loader's FS root), resolves all includes, and returns a ready-to-use
// ConventionSheet.
func (l *Loader) LoadCard(cardPath string) (*ConventionSheet, error) {
	data, err := fs.ReadFile(l.fsys, cardPath)
	if err != nil {
		return nil, fmt.Errorf("reading card %q: %w", cardPath, err)
	}
	var card yamlCard
	if err := yaml.Unmarshal(data, &card); err != nil {
		return nil, fmt.Errorf("parsing card %q: %w", cardPath, err)
	}

	cs := NewConventionSheet()

	for _, name := range card.Includes {
		path := "conventions/" + name + ".yaml"
		data, err := fs.ReadFile(l.fsys, path)
		if err != nil {
			return nil, fmt.Errorf("loading sub-convention %q: %w", name, err)
		}
		var sub yamlSubConvention
		if err := yaml.Unmarshal(data, &sub); err != nil {
			return nil, fmt.Errorf("parsing sub-convention %q: %w", name, err)
		}
		conventions, err := buildConventions(sub.Conventions)
		if err != nil {
			return nil, fmt.Errorf("sub-convention %q: %w", name, err)
		}
		for _, c := range conventions {
			cs.Add(c)
		}
	}

	conventions, err := buildConventions(card.Conventions)
	if err != nil {
		return nil, fmt.Errorf("card conventions: %w", err)
	}
	for _, c := range conventions {
		cs.Add(c)
	}

	return cs, nil
}
