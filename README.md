# Metrics Checker
对测试的指标监控和校验工具

A quick and dirty way. I'll improve it.

详见：[大集群测试](https://docs.google.com/document/d/1EEFZVSifkDFwBzkzMKxhs3YmBJ_WHdkRXQDxrJfG_Pk/edit?ts=5ff6ee26#heading=h.crmsk8lqu128)

## Usage
配置文件为 config.yaml，示例如下：
```yaml
start-after: 30s  # optional, default = 0s
interval: 10s     # optional, default = 1m
rules:    # 对 prometheus api 的 query
    - tag: tps
      promql: sum(rate(tidb_session_transaction_duration_seconds_count[5m])) > bool sum(rate(tidb_session_transaction_duration_seconds_count[10m]))  # 现在需要一个返回 bool 值的 promQL 表达式。
metrics-to-show:  # 在 grafana 中生成图线
  tps_1m: sum(rate(tidb_session_transaction_duration_seconds_count[1m]))
  tps_10m: sum(rate(tidb_session_transaction_duration_seconds_count[10m]))
```

command usage:

```bash
$ make metrics-checker

$ ./metrics-checker
2021/01/25 15:08:47 Start checking metrics after 5s
2021/01/25 15:08:52 Start checking metrics
2021/01/25 15:08:52 checking query: sum(rate(tidb_session_transaction_duration_seconds_count[5m])) > bool sum(rate(tidb_session_transaction_duration_seconds_count[10m]))
2021/01/25 15:09:02 checking query: sum(rate(tidb_session_transaction_duration_seconds_count[5m])) > bool sum(rate(tidb_session_transaction_duration_seconds_count[10m]))
```

可以指定 prometheus 的地址
```bash
$ ./metrics-checker -u 127.0.0.1:9090
2021/01/25 15:21:07 Start checking metrics after 0s
2021/01/25 15:21:07 Start checking metrics
2021/01/25 15:21:07 Prometheus address: http://127.0.0.1:9090
2021/01/25 15:21:07 checking query: sum(rate(tidb_session_transaction_duration_seconds_count[5m])) > bool sum(rate(tidb_session_transaction_duration_seconds_count[10m]))
```

Config can also passed by base64 string
```bash
./metrics-checker --config-base64 c3RhcnQtYWZ0ZXI6IDEwMHMKaW50ZXJ2YWw6IDEwcwpydWxlczogICAgIyDlr7kgcHJvbWV0aGV1cyBhcGkg55qEIHF1ZXJ5CiAgICAtIHRhZzogdHBzCiAgICAgIHByb21xbDogc3VtKHJhdGUodGlkYl9zZXNzaW9uX3RyYW5zYWN0aW9uX2R1cmF0aW9uX3NlY29uZHNfY291bnRbMW1dKSkgPiBib29sIDIvMyAqIHN1bShyYXRlKHRpZGJfc2Vzc2lvbl90cmFuc2FjdGlvbl9kdXJhdGlvbl9zZWNvbmRzX2NvdW50WzVtXSkpCg==
2021/01/26 09:58:37 Load log from base64 string
2021/01/26 09:58:37 Start checking metrics after 1m40s
```

# docker image

Usage:
```bash
docker run -it --rm localhost/metrics-checker:latest metrics-checker --config-base64 c3RhcnQtYWZ0ZXI6IDBzCmludGVydmFsOiAxMHMKcnVsZXM6ICAgICMg5a+5IHByb21ldGhldXMgYXBpIOeahCBxdWVyeQogICAgLSB0YWc6IHRwcwogICAgICBwcm9tcWw6IHN1bShyYXRlKHRpZGJfc2Vzc2lvbl90cmFuc2FjdGlvbl9kdXJhdGlvbl9zZWNvbmRzX2NvdW50WzFtXSkpID4gYm9vbCAyLzMgKiBzdW0ocmF0ZSh0aWRiX3Nlc3Npb25fdHJhbnNhY3Rpb25fZHVyYXRpb25fc2Vjb25kc19jb3VudFs1bV0pKQo= --address 192.168.1.164:9090
```

## Test With Workload
[test-correctness-workload.sh](./tests/test-correctness-workload.sh) generates a mock workload and use our checker to observe it.

It generate custom [config.yaml](./config.yaml). Dependencies:
- [go-tpc](https://github.com/pingcap/go-tpc)
- [tiup](https://github.com/pingcap/tiup)

```bash
$ ./tests/test-correctness-workload.sh 
backuped ./config.yaml to ./config.yaml.backup

...

2021/01/25 17:33:13 Start checking metrics after 0s
2021/01/25 17:33:13 Start checking metrics
2021/01/25 17:33:13 Prometheus address: http://127.0.0.1:9090
2021/01/25 17:33:13 checking query: sum(rate(tidb_session_transaction_duration_seconds_count[1m])) > bool 2/3 * sum(rate(tidb_session_transaction_duration_seconds_count[5m]))
...
...
Switch thread from 6 to 1
2021/01/25 17:38:13 checking query: sum(rate(tidb_session_transaction_duration_seconds_count[1m])) > bool 2/3 * sum(rate(tidb_session_transaction_duration_seconds_count[5m]))
...
2021/01/25 17:38:43 Rule {tps sum(rate(tidb_session_transaction_duration_seconds_count[1m])) > bool 2/3 * sum(rate(tidb_session_transaction_duration_seconds_count[5m]))} failed.
```

## Generate grafana metric dashboard
Add metrics you want to show in config file:
```yaml
metrics-to-show:  # 在 grafana 中生成图线
  tps_1m: sum(rate(tidb_session_transaction_duration_seconds_count[1m]))
  tps_10m: sum(rate(tidb_session_transaction_duration_seconds_count[10m]))
```

Specify grafana api address with `--grafana` flag.
```bash
$ ./metrics-checker --grafana 127.0.0.1:3000
2021/01/28 19:09:03 Load config from file ./config.yaml
2021/01/28 19:09:03 Created dashboard Metrics Checker
```

Grafana Dashboard "Metrics Checker" will show the metrics in config file.
![Grafana Dashboard](doc/assets/grafana-metrics-checker.png)