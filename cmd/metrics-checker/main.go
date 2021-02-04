package main

import (
	"encoding/base64"
	"io/ioutil"
	"time"

	"github.com/pingcap/log"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/PingCAP-QE/metrics-checker/pkg/metrics"
)

// Variables which are passed in by flags.
var (
	prometheusAPIURL  string
	configFilePath    string
	configBase64      string
	grafanaAPIURL     string
	grafanaDataSource string
	config            Config
)

// Config represents information from config file.
type Config struct {
	startTime     time.Time
	StartAfter    time.Duration
	Interval      time.Duration
	Rules         []metrics.Rule
	MetricsToShow map[string]string
}

// AlertFunc do something when Rule failed.
func AlertFunction(rule metrics.Rule) {
	log.Fatal("Rule failed", zap.String("rule", rule.String()))
}

// NotifyFunc do something when Rule succeeded.
func NotifyFunction(rule metrics.Rule) {
	log.Info("Rule passed", zap.String("rule", rule.String()))
}

func (r *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp struct {
		StartAfter    string            `yaml:"start-after,omitempty"`
		Interval      string            `yaml:"interval,omitempty"`
		Rules         []metrics.Rule    `yaml:"rules"`
		MetricsToShow map[string]string `yaml:"metrics-to-show"`
	}
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	// TODO: Set some default value of config file here.
	// 		 Maybe not a good practice. Change it in the future.
	r.StartAfter = ParseDurationWithDefault(tmp.StartAfter, 0*time.Minute)
	r.Interval = ParseDurationWithDefault(tmp.Interval, 1*time.Minute)
	r.Rules = tmp.Rules
	r.MetricsToShow = tmp.MetricsToShow
	return nil
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

func LoadConfig(path string) Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	return LoadConfigFromBytes(file)
}

func LoadConfigFromBytes(b []byte) Config {
	config := Config{}
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return config
}

// TODO: There must be some better way to do config initialization. Ref: https://www.bookstack.cn/read/cobra/spilt.2.spilt.4.README.md

// InitConfig return a Config parsed from configBase64, if configBase64 is empty, load from configFilePath.
func InitConfig(configFilePath string, configBase64 string) Config {
	var c Config
	if configBase64 == "" {
		c = LoadConfig(configFilePath)
		log.Info("Load config from file", zap.String("file path", configFilePath))
	} else {
		configString, err := base64.StdEncoding.DecodeString(configBase64)
		if err != nil {
			log.Fatal(err.Error())
		}
		c = LoadConfigFromBytes(configString)
		log.Info("Load config from base64 string")
	}
	c.startTime = time.Now()

	if len(c.Rules) == 0 {
		log.Fatal("Number of rules == 0")
	}

	return c
}

func CreateGrafanaDashboard() {
	dashboardName := "Metrics Checker"

	reformedGrafanaURL, err := metrics.AddHTTPIfIP(grafanaAPIURL)
	if err != nil {
		log.Fatal("Grafana prometheusAPIURL invalid", zap.String("grafana", grafanaAPIURL))
	}
	grafanaAPIURL = reformedGrafanaURL

	if grafanaDataSource == "" {
		log.Fatal("Grafana datasource is not set.")
	}
	err = metrics.CreateMetricsDashboard(grafanaAPIURL, dashboardName, grafanaDataSource, config.MetricsToShow)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Created dashboard", zap.String("name", dashboardName), zap.String("grafana", grafanaAPIURL))
}

func main() {
	Execute()
}
