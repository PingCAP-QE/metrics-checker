package metrics

import (
	"fmt"
	"time"
)

// Rule represents a bool PromQL expression, returning True means rule satisified.
//   Send alerts when a rule return true, which is similarly as prometheus
//   alerting rules.
type Rule struct {
	Name              string `yaml:"tag"`
	PromQL            string `yaml:"promql"`
	Interval          time.Duration
	EvaluatedCallback EvaluatedCallback
	AlertCallback     AlertCallback
}

func BuildRuleWithDefaultCallback(tag string, interval time.Duration, promQL string) Rule {
	return Rule{
		Name:              tag,
		PromQL:            promQL,
		Interval:          interval,
		EvaluatedCallback: DefaultEvaluatedCallback,
		AlertCallback:     DefaultAlertCallback,
	}
}

// EvaluatedCallback will be called whenever the PromQL is evaluated.
type EvaluatedCallback func(rule Rule)

// AlertCallback will be called only the Rule evaluated to truthy.
type AlertCallback func(rule Rule)

func (r *Rule) String() string {
	return fmt.Sprintf("Name: %s; PromQL: %s", r.Name, r.PromQL)
}
