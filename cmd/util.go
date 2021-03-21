package cmd

import (
	"fmt"
	"time"

	"github.com/PingCAP-QE/metrics-checker/pkg/metrics"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

// AlertFunc do something when Rule failed.
func AlertFunction(rule metrics.Rule) {
	log.Fatal("Rule failed", zap.String("rule", rule.String()))
}

// NotifyFunc do something when Rule succeeded.
func NotifyFunction(rule metrics.Rule) {
	log.Info("Rule passed", zap.String("rule", rule.String()))
}

// ParseDurationWithDefault parse string `s` to duration, if s is an empty string, return `fallback`.
func ParseDurationWithDefault(s string, fallback time.Duration) time.Duration {
	if s == "" {
		return fallback
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatal(err.Error())
	}
	return d
}

// CreateGrafanaDashboard create dashboard on given Grafana.
func CreateGrafanaDashboard() {
	dashboardName := "Metrics Checker"

	reformedGrafanaURL, err := metrics.AddHTTPIfIP(Flag.grafanaAPIURL)
	if err != nil {
		log.Fatal("Grafana prometheusAPIURL invalid", zap.String("grafana", Flag.grafanaAPIURL))
	}
	Flag.grafanaAPIURL = reformedGrafanaURL

	if Flag.grafanaDataSource == "" {
		log.Fatal("Grafana datasource is not set.")
	}
	err = metrics.CreateMetricsDashboard(Flag.grafanaAPIURL, dashboardName, Flag.grafanaDataSource, Conf.MetricsToShow)
	fmt.Printf("MTS: %v\n", Conf.MetricsToShow)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Created dashboard", zap.String("name", dashboardName), zap.String("grafana", Flag.grafanaAPIURL))
}
