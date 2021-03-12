# Note of Development

## TODO
- [x] Provide some example config file:
  - [x] TPS and latency.
  - [x] Check deadlock of tidb.
    - Via heartbeat of tikvs.
  - [x] Check region health of tidb.
    - Via extra peer / miss peer.
- [ ] Check panic of tidb via logs.

From codes:
```
/metrics-checker/cmd/metricchecker/main.go
  41,5: 	// TODO: Set some default value of config file here.

/metrics-checker/pkg/metric/prometheus.go
  25,6: 		// TODO: When Prometheus doesn't up, the length of boolVector would be zero.
  28,7: 			// TODO: During startup of prometheus, query would return a zero-length
```

## Important external documentations
- [大集群测试](https://docs.google.com/document/d/1EEFZVSifkDFwBzkzMKxhs3YmBJ_WHdkRXQDxrJfG_Pk/edit?ts=5ff6ee26#heading=h.crmsk8lqu128): This internal documentation records the application of this checker.


