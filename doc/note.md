## V0.0.1
设计（更简单的版本）

再进一步把 expression 的部分砍掉。所有的 query 只能返回一个布尔值。

- [ ] 要不要依赖 bench-toolset？
  - 算了吧，那我也太懒了？


```yaml
rule:    # 对 prometheus api 的 query
    tps:
        promql: sum(rate(tidb_session_transaction_duration_seconds_count[1m])) > bool sum(rate(tidb_session_transaction_duration_seconds_count[10m])) * 2 / 3
        start-after: 10m
```

## 问题

```bash
ofey@RX93 ~/Code/PingCAP-Internship/playground/metrics-fetching
$ ./metrics-fetching 
2021/01/25 10:44:48 checking query: sum(rate(tidb_session_transaction_duration_seconds_count[5m])) > bool sum(rate(tidb_session_transaction_duration_seconds_count[10m]))
panic: runtime error: index out of range [0] with length 0
```

在 prometheus 刚刚起来的时候，读不到数据会报这样的错。


## TODO
- [x] corner case 处理。
- [x] 造数据，验证是否能够工作。
- [x] 造数据的脚本
- [ ] 打包成 docker file。

## TPCC
1. 升级过程中，TPS 不能下降到 三分之一 一下，lat 不能升高 10% 以上

### 确认基本的正确性
造数据的目标：确认一下，TPS 降到标准以下的时候，能探测出来。

使用 tiup playground 起本地的测试集群。

使用 go-tpc 起一个 workload。

通过本地电脑观察得到，在 go-tpc 的 threading 从 2 变成 1 的时候，grafana 的 playground-overview/transaction-OPS 会下降大约一半，这就是我们需要监控的东西。

写一个 script 把这个自动测试的过程录下来。

### 通过命令行传递配置文件

