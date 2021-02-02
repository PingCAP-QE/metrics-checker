package main

import (
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Metrics checker",
	Short: "For checking prometheus metrics",
	Run: func(cmd *cobra.Command, args []string) {
		InitConfig(configFilePath, configBase64)

		address = AddHTTPIfIP(address)
		if grafanaAPIURL != "" {
			dashboardName := "Metrics Checker"
			grafanaAPIURL = AddHTTPIfIP(grafanaAPIURL)
			err := createMetricsDashboard(grafanaAPIURL, dashboardName, config.MetricsToShow)
			if err != nil {
				log.Fatalf("Create grafana metrics error: %s", err)
			}
			log.Printf("Created dashboard %s on %s", dashboardName, grafanaAPIURL)
		}

		log.Printf("Start checking metrics after %s", config.StartAfter)
		for time.Now().Before(config.startTime.Add(config.StartAfter)) {
			time.Sleep(time.Second)
		}
		log.Printf("Start checking metrics")
		log.Printf("Prometheus address: %s", address)

		// Create prometheus API
		client, err := api.NewClient(api.Config{
			Address: address,
		})
		if err != nil {
			log.Fatalf("Create prometheus api failed. address: %s, error: %s", address, err)
		}
		api := v1.NewAPI(client)

		for {
			// run rules
			for _, rule := range config.Rules {
				if !Check(api, rule.PromQL, time.Now()) {
					log.Fatalf("Rule %s failed.", rule)
				}
			}
			time.Sleep(config.Interval)
		}
	},
}

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
