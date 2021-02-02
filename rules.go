package main

type Rule struct {
	Tag    string `yaml:"tag"`
	PromQL string `yaml:"promql"`
}
