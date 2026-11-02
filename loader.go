package passpasspass

import (
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// yamlSubConvention is the top-level structure of a sub-convention file.
type yamlSubConvention struct {
	Name        string      `yaml:"name"`
	Conventions []yaml.Node `yaml:"conventions"`
}

// yamlCard is the top-level structure of a convention card file.
type yamlCard struct {
	Name        string      `yaml:"name"`
	Includes    []string    `yaml:"includes"`
	Conventions []yaml.Node `yaml:"conventions"`
}

// The keys accepted at each level of a file, in order, for checkKeys.
var (
	cardKeys          = []string{"conventions", "includes", "name"}
	subConventionKeys = []string{"conventions", "name"}
	conventionKeys    = []string{"bid", "conditions"}
)

// loadError is an error at a line of a convention file. It reads
// `file:line: bid 2S: any.suit_length.hearts: message`; the parts that are not known are left out.
type loadError struct {
	file string   // set by the loader once the error leaves the file
	line int      // the line of the innermost offending key or value
	bid  string   // the convention's bid, when the error is inside a convention
	path []string // the keys leading to the offending one, outermost first
	err  error
}

func (e *loadError) Error() string {
	var b strings.Builder
	if e.file != "" {
		b.WriteString(e.file + ":")
	}
	fmt.Fprintf(&b, "%d: ", e.line)
	if e.bid != "" {
		fmt.Fprintf(&b, "bid %s: ", e.bid)
	}
	if len(e.path) > 0 {
		b.WriteString(strings.Join(e.path, ".") + ": ")
	}
	b.WriteString(e.err.Error())
	return b.String()
}

func (e *loadError) Unwrap() error { return e.err }

// errorAt places err at the line of n.
func errorAt(n *yaml.Node, err error) error {
	return &loadError{line: n.Line, err: err}
}

// keyError puts key in front of err's key path. It places err at the key's line unless err
// already has a line from inside the key's value.
func keyError(key *yaml.Node, err error) error {
	var le *loadError
	if errors.As(err, &le) {
		le.path = append([]string{key.Value}, le.path...)
		return le
	}
	return &loadError{line: key.Line, path: []string{key.Value}, err: err}
}

// yamlSyntaxError matches the "yaml: line 3: message" form of yaml.v3's syntax errors.
var yamlSyntaxError = regexp.MustCompile(`^yaml: line (\d+): (.*)$`)

// syntaxError turns a yaml.v3 syntax error into a loadError at its line, so it reads like every
// other load error. Errors of any other form are returned unchanged.
func syntaxError(err error) error {
	m := yamlSyntaxError.FindStringSubmatch(err.Error())
	if m == nil {
		return err
	}
	line, _ := strconv.Atoi(m[1])
	return &loadError{line: line, err: errors.New(m[2])}
}

// inFile names the file at path in err.
func inFile(path string, err error) error {
	var le *loadError
	if errors.As(err, &le) {
		le.file = path
		return le
	}
	return fmt.Errorf("%s: %w", path, err)
}

// checkKeys rejects the first key of mapping n that is not in known or appears twice. kind names
// the mapping in the error, as in "unknown condition key".
func checkKeys(n *yaml.Node, kind string, known []string) error {
	seen := make(map[string]bool, len(n.Content)/2)
	for i := 0; i+1 < len(n.Content); i += 2 {
		key := n.Content[i]
		if !slices.Contains(known, key.Value) {
			return keyError(key, fmt.Errorf("unknown %s key (known keys: %s)", kind, strings.Join(known, ", ")))
		}
		if seen[key.Value] {
			return keyError(key, errors.New("key appears more than once"))
		}
		seen[key.Value] = true
	}
	return nil
}

// mappingValue returns the value of key in mapping n, or nil.
func mappingValue(n *yaml.Node, key string) *yaml.Node {
	if n.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i+1]
		}
	}
	return nil
}

// parseConvention turns a convention entry, a mapping with `bid` and `conditions`, into a
// Convention. Errors name the bid when the entry has one.
func parseConvention(n *yaml.Node) (*Convention, error) {
	c, err := parseConventionEntry(n)
	if err != nil {
		var le *loadError
		if bid := mappingValue(n, "bid"); bid != nil && errors.As(err, &le) {
			le.bid = bid.Value
		}
		return nil, err
	}
	return c, nil
}

func parseConventionEntry(n *yaml.Node) (*Convention, error) {
	if n.Kind != yaml.MappingNode {
		return nil, errorAt(n, errors.New("a convention must be a mapping with bid and conditions"))
	}
	if err := checkKeys(n, "convention", conventionKeys); err != nil {
		return nil, err
	}
	bidNode, conditionsNode := mappingValue(n, "bid"), mappingValue(n, "conditions")
	if bidNode == nil {
		return nil, errorAt(n, errors.New("convention has no bid"))
	}
	if conditionsNode == nil || conditionsNode.ShortTag() == "!!null" {
		return nil, errorAt(n, errors.New("convention has no conditions"))
	}
	bid, err := parseBid(bidNode)
	if err != nil {
		return nil, errorAt(bidNode, err)
	}
	condition, err := parseCondition(conditionsNode)
	if err != nil {
		return nil, err
	}
	if condition.maxCases(true) > maxConventionCases {
		return nil, errorAt(conditionsNode, fmt.Errorf(
			"conditions can show more than %d alternatives; split them into several conventions", maxConventionCases))
	}
	return NewConvention(bid, condition), nil
}

// addConventions parses every convention entry and adds it to cs, in file order.
func addConventions(cs *ConventionSheet, entries []yaml.Node) error {
	for i := range entries {
		c, err := parseConvention(&entries[i])
		if err != nil {
			return err
		}
		cs.Add(c)
	}
	return nil
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
// ConventionSheet. Errors start with the file and line, as in `cards/colin.yaml:4: bid 1NT: ...`.
func (l *Loader) LoadCard(cardPath string) (*ConventionSheet, error) {
	var card yamlCard
	if err := l.readFile(cardPath, "card", cardKeys, &card); err != nil {
		return nil, inFile(cardPath, err)
	}

	cs := NewConventionSheet()

	for _, name := range card.Includes {
		path := "conventions/" + name + ".yaml"
		var sub yamlSubConvention
		err := l.readFile(path, "sub-convention", subConventionKeys, &sub)
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("%s: include %s: %w", cardPath, name, inFile(path, err))
		}
		if err == nil {
			err = addConventions(cs, sub.Conventions)
		}
		if err != nil {
			return nil, inFile(path, err)
		}
	}

	if err := addConventions(cs, card.Conventions); err != nil {
		return nil, inFile(cardPath, err)
	}

	return cs, nil
}

// readFile reads the YAML file at path, rejects top-level keys not in known and decodes the
// file into v. kind names the file in the error, as in "unknown card key". Errors do not name
// the file; the caller does.
func (l *Loader) readFile(path, kind string, known []string, v any) error {
	data, err := fs.ReadFile(l.fsys, path)
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		return pathErr.Err
	}
	if err != nil {
		return err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return syntaxError(err)
	}
	if len(doc.Content) == 0 {
		return errors.New("file is empty")
	}
	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return errorAt(root, errors.New("must be a mapping"))
	}
	if err := checkKeys(root, kind, known); err != nil {
		return err
	}
	return root.Decode(v)
}
