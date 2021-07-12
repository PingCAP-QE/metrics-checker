package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/pingcap/log"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"go.uber.org/zap"
)

// Checker MetricsChecker checks Rules with given Interval.
type Checker struct {
	API   v1.API
	Rules []Rule
}

func NewChecker(promAddress string) (*Checker, error) {
	client, err := api.NewClient(api.Config{
		Address: promAddress,
	})
	if err != nil {
		return nil, err
	}
	m := &Checker{
		API:   v1.NewAPI(client),
		Rules: []Rule{},
	}
	return m, nil
}

func (m *Checker) AddRule(rule Rule) *Checker {
	m.Rules = append(m.Rules, rule)
	return m // return Checker itself enables the user to add rules as a chained method call, which will be more convenient
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

// RunBlocked starts processing of the MetricsChecker.
//   It is blocking.
//   If any rule returns an error, `RunBlocked` returns with this error,
//   stop all other rules immediately.
func (m *Checker) RunBlocked() error {
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	for _, rule := range m.Rules {
		go func(rule Rule) {
			defer wg.Done()
			wg.Add(1)
			err := m.RunRule(ctx, rule, rule.Interval)
			if err != nil {
				errChan <- err
				wg.Done()
			}
		}(rule)
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
			if rule.EvaluatedCallback != nil {
				rule.EvaluatedCallback(rule)
			}
			if ans == true {
				// Alert function is required.
				rule.AlertCallback(rule)
			}
			time.Sleep(interval)
		}
	}
}

// DefaultAlertCallback do something when Rule failed.
func DefaultAlertCallback(rule Rule) {
	log.Fatal("Rule failed", zap.String("rule", rule.String()))
}

// DefaultEvaluatedCallback do something when Rule succeeded.
func DefaultEvaluatedCallback(rule Rule) {
	log.Info("Rule evaluated", zap.String("rule", rule.String()))
}
