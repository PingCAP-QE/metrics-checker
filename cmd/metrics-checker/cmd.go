package main

import (
	"os"
	"time"

	"github.com/pingcap/log"
	"go.uber.org/zap"

	"github.com/spf13/cobra"

	"github.com/PingCAP-QE/metrics-checker/pkg/metrics"
)

var rootCmd = &cobra.Command{
	Use:   "Metrics checker",
	Short: "For checking prometheus metrics",
	Run: func(cmd *cobra.Command, args []string) {
		config := InitConfig(configFilePath, configBase64)

		reformedAddress, err := metrics.AddHTTPIfIP(prometheusAPIURL)
		if err != nil {
			log.Fatal("Prometheus prometheusAPIURL invalid", zap.String("prometheus", prometheusAPIURL))
		}
		prometheusAPIURL = reformedAddress

		if grafanaAPIURL != "" {
			CreateGrafanaDashboard()
		}

		log.Info("Waiting for checking metrics", zap.Duration("start after", config.StartAfter))
		for time.Now().Before(config.startTime.Add(config.StartAfter)) {
			time.Sleep(time.Second)
		}
		log.Info("Start checking metrics", zap.String("prometheus", prometheusAPIURL))

		for i := range config.Rules {
			config.Rules[i].NotifyFunc = NotifyFunction
			config.Rules[i].AlertFunc = AlertFunction
		}

		metricsChecker, err := metrics.NewMetricsChecker(prometheusAPIURL, config.Rules, config.Interval)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = metricsChecker.Run()
		if err != nil {
			log.Fatal("Metrics checker running error", zap.String("err", err.Error()))
		}
	},
}

// Execute ...
func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.PersistentFlags().StringVarP(&prometheusAPIURL, "address", "u", "http://127.0.0.1:9090", "Host and port of prometheus")
	rootCmd.PersistentFlags().StringVarP(&configFilePath, "config", "c", "./config.yaml", "Set config file path, overrided by --config-base64")
	rootCmd.PersistentFlags().StringVar(&configBase64, "config-base64", "", "Pass config file as base64 string, override --config")
	rootCmd.PersistentFlags().StringVar(&grafanaAPIURL, "grafana", "", "Pass config file as base64 string, override --config")
	rootCmd.PersistentFlags().StringVar(&grafanaDataSource, "grafana-datasource", "", "Datasource of grafana panels.")

	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println(err)
		os.Exit(1)
	}
}
