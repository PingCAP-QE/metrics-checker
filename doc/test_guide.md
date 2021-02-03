# Test With Mock Workload
[test-correctness-workload.sh](../tests/test-correctness-workload.sh) generates a mock workload and use our checker to observe it.

It generate custom [config.yaml](../config.yaml). Dependencies:
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
