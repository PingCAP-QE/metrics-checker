package cmd

import (
	"encoding/base64"
	"io/ioutil"
	"time"

	"github.com/PingCAP-QE/metrics-checker/pkg/metrics"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Config is a global Config variable
var Conf Config

// Flag is a global flag variable
var Flag FlagConfig

// Config represents information from config file.
type Config struct {
	startTime     time.Time
	StartAfter    time.Duration
	Interval      time.Duration
	Rules         []metrics.Rule
	MetricsToShow map[string]string
}

// FlagConfig is a struct of variables which are passed in by flags.
type FlagConfig struct {
	prometheusAPIURL  string
	configFilePath    string
	configBase64      string
	grafanaAPIURL     string
	grafanaDataSource string
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

// UnmarshalYAML parse `Config` struct from yaml, set some default values.
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

	r.StartAfter = ParseDurationWithDefault(tmp.StartAfter, 0*time.Minute)
	r.Interval = ParseDurationWithDefault(tmp.Interval, 1*time.Minute)
	r.Rules = tmp.Rules
	r.MetricsToShow = tmp.MetricsToShow
	return nil
}

// InitConfig reads in config file from given/default place
func InitConfig() {
	if Flag.configBase64 == "" {
		Conf = LoadConfig(Flag.configFilePath)
		log.Info("Load config from file", zap.String("file path", Flag.configFilePath))
	} else {
		configString, err := base64.StdEncoding.DecodeString(Flag.configBase64)
		if err != nil {
			log.Fatal(err.Error())
		}
		Conf = LoadConfigFromBytes(configString)
		log.Info("Load config from base64 string")
	}
	Conf.startTime = time.Now()

	if len(Conf.Rules) == 0 {
		log.Fatal("Number of rules == 0")
	}
}
