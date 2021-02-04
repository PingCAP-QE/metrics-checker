package main

import (
	"os"
	"time"

	"github.com/pingcap/log"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/spf13/cobra"

	"github.com/PingCAP-QE/metrics-checker/pkg/metric"
)

var rootCmd = &cobra.Command{
	Use:   "Metrics checker",
	Short: "For checking prometheus metrics",
	Run: func(cmd *cobra.Command, args []string) {
		InitConfig(configFilePath, configBase64)

		reformedAddress, err := metric.AddHTTPIfIP(address)
		if err != nil {
			log.Fatal("Prometheus address invalid", zap.String("address", address))
		}
		address = reformedAddress

		if grafanaAPIURL != "" {
			dashboardName := "Metrics Checker"

			reformedGrafanaURL, err := metric.AddHTTPIfIP(grafanaAPIURL)
			if err != nil {
				log.Fatal("Grafana address invalid", zap.String("grafana", grafanaAPIURL))
			}
			grafanaAPIURL = reformedGrafanaURL

			err = metric.CreateMetricsDashboard(grafanaAPIURL, dashboardName, config.MetricsToShow)
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Info("Created dashboard", zap.String("name", dashboardName), zap.String("grafana url", grafanaAPIURL))
		}

		log.Info("Waiting for checking metrics", zap.Duration("start after", config.StartAfter))
		for time.Now().Before(config.startTime.Add(config.StartAfter)) {
			time.Sleep(time.Second)
		}
		log.Info("Start checking metrics", zap.String("prometheus address", address))

		// Create prometheus API
		client, err := api.NewClient(api.Config{
			Address: address,
		})
		if err != nil {
			log.Fatal(err.Error())
		}
		api := v1.NewAPI(client)

		for {
			// run rules
			for _, rule := range config.Rules {
				ans, err := metric.Check(api, rule.PromQL, time.Now())
				if err != nil {
					log.Warn(err.Error())
				}
				if ans == false {
					log.Fatal("Rule failed", zap.String("rule", rule.String()))
				}
				log.Info("Rule passed", zap.String("rule", rule.String()))
			}
			time.Sleep(config.Interval)
		}
	},
}

// Execute ...
func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.PersistentFlags().StringVarP(&address, "address", "u", "http://127.0.0.1:9090", "Host and port of prometheus")
	rootCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "./config.yaml", "Set config file path, overrided by --config-base64")
	rootCmd.PersistentFlags().StringVar(&configBase64, "config-base64", "", "Pass config file as base64 string, override --config")
	rootCmd.PersistentFlags().StringVar(&grafanaAPIURL, "grafana", "", "Pass config file as base64 string, override --config")

	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println(err)
		os.Exit(1)
	}
}
