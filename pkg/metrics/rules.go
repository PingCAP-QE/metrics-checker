package metrics

import (
	"fmt"
)

// Rule represents a bool PromQL expression, returning True means rule satisified.
//   Send alerts when a rule return true, which is similarly as prometheus
//   alerting rules.
type Rule struct {
	Tag        string `yaml:"tag"`
	PromQL     string `yaml:"promql"`
	NotifyFunc NotifyFunc
	AlertFunc  AlertFunc
}

// NotifyFunc do something when Rule failed.
type NotifyFunc func(rule Rule)

// AlertFunc do something when Rule satifified.
type AlertFunc func(rule Rule)

func (r *Rule) String() string {
	return fmt.Sprintf("Tag: %s; PromQL: %s", r.Tag, r.PromQL)
}
