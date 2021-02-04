package metric

import "fmt"

// Rule represents a bool PromQL expression, returning True means rule passed.
type Rule struct {
	Tag    string `yaml:"tag"`
	PromQL string `yaml:"promql"`
}

func (r *Rule) String() string {
	return fmt.Sprintf("Tag: %s; PromQL: %s", r.Tag, r.PromQL)
}
