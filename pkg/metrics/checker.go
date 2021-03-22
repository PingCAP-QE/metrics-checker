package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// MetricsChecker checks Rules with given Interval.
type Checker struct {
	API      v1.API
	Interval time.Duration
	Rules    []Rule
}

// NewChecker creates a new instance of MetricsChecker.
func NewChecker(promAddress string, rules []Rule, interval time.Duration) (*Checker, error) {
	client, err := api.NewClient(api.Config{
		Address: promAddress,
	})
	if err != nil {
		return nil, err
	}
	m := &Checker{
		API:      v1.NewAPI(client),
		Rules:    rules,
		Interval: interval,
	}
	return m, nil
}

// CheckGiven checks whether a given promQL is satisfied.
func (m *Checker) CheckGiven(promQL string) (bool, error) {
	ans, err := Check(m.API, promQL, time.Now())
	if err != nil {
		return false, err
	}
	if ans == false {
		return false, nil
	}
	return true, nil
}

// Run starts processing of the MetricsChecker.
//   It is blocking.
//   If any rule returns an error, `Run` returns with this error,
//   stop all other rules immediately.
func (m *Checker) Run() error {
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	for _, rule := range m.Rules {
		go func(rule Rule, interval time.Duration) {
			defer wg.Done()
			wg.Add(1)
			err := m.RunRule(ctx, rule, interval)
			if err != nil {
				errChan <- err
			}
		}(rule, m.Interval)
	}
	err := <-errChan
	// If an error occurs, try to shut down all goroutines.
	cancel()
	wg.Wait()
	return err
}

// RunRule start to run checking on a certain rule.
//   Has nothing to do with the internal state of `Checker`
//   It will block the control flow.
func (m *Checker) RunRule(ctx context.Context, rule Rule, interval time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			ans, err := m.CheckGiven(rule.PromQL)
			if err != nil {
				return err
			}
			if ans == true {
				// Alert function is required.
				rule.AlertFunc(rule)
			}
			if rule.NotifyFunc != nil {
				rule.NotifyFunc(rule)
			}
			time.Sleep(interval)
		}
	}
}
