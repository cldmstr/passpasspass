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
	HCP        []int                `yaml:"hcp"`
	Balanced   *bool                `yaml:"balanced"`
	SuitLength map[string][]int     `yaml:"suit_length"`
	All        []conditionNode      `yaml:"all"`
	Any        []conditionNode      `yaml:"any"`
}

var suitIndex = map[string]int{
	"clubs": 0, "diamonds": 1, "hearts": 2, "spades": 3,
}

func buildCondition(n conditionNode) (Condition, error) {
	switch {
	case n.HCP != nil:
		return NewHCPIntervalCondition(n.HCP[0], n.HCP[1]), nil

	case n.Balanced != nil:
		return NewBalancedDistributionCondition(), nil

	case n.SuitLength != nil:
		r := DistributionRange{{0, 13}, {0, 13}, {0, 13}, {0, 13}}
		for suit, rng := range n.SuitLength {
			i, ok := suitIndex[suit]
			if !ok {
				return nil, fmt.Errorf("unknown suit %q in suit_length", suit)
			}
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
		conventions = append(conventions, NewConvention(bid, condition))
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
