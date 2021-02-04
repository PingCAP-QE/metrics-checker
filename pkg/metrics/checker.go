package metrics

import (
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// MetricsChecker checks Rules with given Interval.
type MetricsChecker struct {
	API      v1.API
	Interval time.Duration
	Rules    []Rule
}

// NewMetricsChecker creates a new instance of MetricsChecker.
func NewMetricsChecker(promAddress string, rules []Rule, interval time.Duration) (*MetricsChecker, error) {
	client, err := api.NewClient(api.Config{
		Address: promAddress,
	})
	if err != nil {
		return nil, err
	}
	m := &MetricsChecker{
		API:      v1.NewAPI(client),
		Rules:    rules,
		Interval: interval,
	}
	return m, nil
}

// TODO: Prometheus use channel and multiple goroutines to verify all rules.
// 		 We use single-threaded code here, should improve it.
// 		 Ref: https://github.com/prometheus/prometheus/blob/19c190b406c992278aaade63be92ecc7bb6a4921/rules/manager.go#L910

// CheckGiven checks whether a given promQL returns true or not.
func (m *MetricsChecker) CheckGiven(promQL string) (bool, error) {
	ans, err := Check(m.API, promQL, time.Now())
	if err != nil {
		return false, err
	}
	if ans == false {
		return false, nil
	}
	return true, nil
}

// Run starts processing of the MetricsChecker. It is blocking.
func (m *MetricsChecker) Run() error {
	for {
		for _, rule := range m.Rules {
			ans, err := m.CheckGiven(rule.PromQL)
			if err != nil {
				// TODO: I'm not sure whether Run() should return a error or not.
				return err
			}
			if ans == false {
				// Alert function is required.
				rule.AlertFunc(rule)
			}
			if rule.NotifyFunc != nil {
				rule.NotifyFunc(rule)
			}
		}
		time.Sleep(m.Interval)
	}
}
