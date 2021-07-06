package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	sdk "github.com/pingcap/test-infra/sdk/core"
	"github.com/pingcap/test-infra/sdk/resource"
	_ "github.com/pingcap/test-infra/sdk/resource/impl/k8s"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/PingCAP-QE/metrics-checker/pkg/metrics"
)

func main() {
	_ = godotenv.Load()
	ctx, err := sdk.BuildContext()
	if err != nil {
		fmt.Println("cannot build TestContext", err)
		panic(1)
	}
	logger := ctx.Logger().WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))

	tc := ctx.Resource("tc").(resource.TiDBCluster)
	prometheus, err := tc.ServiceURL(resource.Prometheus)
	if err != nil {
		logger.Error("error when retrieving Prometheus ServiceURL of TC", zap.Error(err))
		panic(1)
	}
	checker, err := metrics.NewChecker(prometheus.String())
	if err != nil {
		logger.Error("error when constructing MetricsChecker", zap.Error(err))
		panic(1)
	}
	logger.Info("MetricsChecker created.")

	checker.AddRule(metrics.Rule{
		Name:              "alert_when_down_peer_more_than_30%",
		PromQL:            "sum(pd_regions_status{type=\"down-peer-region-count\"}) > bool 0.3 * sum(pd_cluster_status{type=\"region_count\"})",
		Interval:          5 * time.Second,
		EvaluatedCallback: EvaluatedWithLogger(logger),
		AlertCallback:     AlertWithLogger(logger),
	}).AddRule(metrics.Rule{
		Name:              "alert_when_up_peer_less_than_60%",
		PromQL:            "sum(pd_regions_status{type=\"up-peer-region-count\"}) < bool 0.6 * sum(pd_cluster_status{type=\"region_count\"})",
		Interval:          10 * time.Second,
		EvaluatedCallback: EvaluatedWithLogger(logger),
		AlertCallback:     AlertWithLogger(logger),
	})

	fmt.Println(len(checker.Rules))

	err = checker.RunBlocked()

	if err != nil {
		logger.Error("error when running MetricsChecker", zap.Error(err))
	}
}

func AlertWithLogger(logger *zap.Logger) metrics.AlertCallback {
	return func(rule metrics.Rule) {
		logger.Warn("Rule evaluated to True.", zap.String("rule_tag", rule.Name))
	}
}

func EvaluatedWithLogger(logger *zap.Logger) metrics.EvaluatedCallback {
	return func(rule metrics.Rule) {
		logger.Info("Rule evaluated.", zap.String("rule_tag", rule.Name))
	}
}
