package cmd

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
	Run:   root,
}

// root is the body of root command.
func root(cmd *cobra.Command, args []string) {
	reformedAddress, err := metrics.AddHTTPIfIP(Flag.prometheusAPIURL)
	if err != nil {
		log.Fatal("Prometheus prometheusAPIURL invalid", zap.String("prometheus", Flag.prometheusAPIURL))
	}
	Flag.prometheusAPIURL = reformedAddress

	if Flag.grafanaAPIURL != "" {
		CreateGrafanaDashboard()
	}

	log.Info("Waiting for checking metrics", zap.Duration("start after", Conf.StartAfter))
	for time.Now().Before(Conf.startTime.Add(Conf.StartAfter)) {
		time.Sleep(time.Second)
	}
	log.Info("Start checking metrics", zap.String("prometheus", Flag.prometheusAPIURL))

	metricsChecker, err := metrics.NewChecker(Flag.prometheusAPIURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	if metricsChecker == nil {
		return
	}

	for _, rule := range Conf.Rules {
		metricsChecker.AddRule(metrics.BuildRuleWithDefaultCallback(rule.Name, Conf.Interval, rule.PromQL))
	}

	err = metricsChecker.RunBlocked()
	if err != nil {
		log.Fatal("Metrics checker running error", zap.String("err", err.Error()))
	}
}

// Execute root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(InitConfig)
	rootCmd.SetOut(os.Stdout)

	rootCmd.PersistentFlags().StringVarP(&Flag.prometheusAPIURL, "address", "u", "http://127.0.0.1:9090", "Host and port of prometheus")
	rootCmd.PersistentFlags().StringVarP(&Flag.configFilePath, "config", "c", "./config.yaml", "Set config file path, overrided by --config-base64")
	rootCmd.PersistentFlags().StringVar(&Flag.configBase64, "config-base64", "", "Pass config file as base64 string, override --config")
	rootCmd.PersistentFlags().StringVar(&Flag.grafanaAPIURL, "grafana", "", "Pass config file as base64 string, override --config")
	rootCmd.PersistentFlags().StringVar(&Flag.grafanaDataSource, "grafana-datasource", "", "Datasource of grafana panels.")
}
