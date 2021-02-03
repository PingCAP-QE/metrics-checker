package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/PingCAP-QE/metrics-checker/pkg/metric"
)

var (
	address        string
	configFilePath string
	configBase64   string
	grafanaAPIURL  string
	config         Config
)

type Config struct {
	startTime     time.Time
	StartAfter    time.Duration
	Interval      time.Duration
	Rules         []metric.Rule
	MetricsToShow map[string]string
}

func (r *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmp struct {
		StartAfter    string            `yaml:"start-after,omitempty"`
		Interval      string            `yaml:"interval,omitempty"`
		Rules         []metric.Rule     `yaml:"rules"`
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
		log.Fatalf("failed to parse '%s' to time.Duration: %v", s, err)
	}
	return d
}

func LoadConfig(path string) Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Read config file %s error: %v", path, err)
	}
	return LoadConfigFromBytes(file)
}

func LoadConfigFromBytes(b []byte) Config {
	config := Config{}
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("Load config error: %v", err)
	}
	return config
}

func InitConfig(configFilePath string, configBase64 string) {
	if configBase64 == "" {
		config = LoadConfig(configFilePath)
		log.Printf("Load config from file %s", configFilePath)
	} else {
		configString, err := base64.StdEncoding.DecodeString(configBase64)
		if err != nil {
			log.Fatalf("Base64 config decode error: %s", configBase64)
		}
		config = LoadConfigFromBytes(configString)
		log.Printf("Load config from base64 string")
	}
	config.startTime = time.Now()

	if len(config.Rules) == 0 {
		log.Fatalf("Number of rules == 0")
	}
}

func main() {
	Execute()
}
