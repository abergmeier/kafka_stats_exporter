package types

import (
	"fmt"
	"regexp"
	"sort"
)

var (
	// See https://prometheus.io/docs/concepts/data_model/
	labelNameExp = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
)

// LabelNames are all possible names for Labels
type LabelNames struct {
	XXX []string
}

// AddLabelNames validates and adds all label names of the other collection
func (lns *LabelNames) AddLabelNames(other *LabelNames) error {
	return lns.AddStrings(other.XXX...)
}

// AddStrings validates and adds all the label names
func (lns *LabelNames) AddStrings(names ...string) error {
	for _, name := range names {
		ok := labelNameExp.MatchString(name)
		if !ok {
			return fmt.Errorf("not a valid Prometheus label name: %s", name)
		}
	}

	lns.XXX = append(lns.XXX, names...)
	return nil
}

func (lns *LabelNames) Sort() {
	sort.Strings(lns.XXX)
}

func (lns *LabelNames) Strings() []string {
	return lns.XXX
}
