CONFIG_FILE=./config.yaml
GO_TPC=go-tpc
CHECKER=./metrics-checker
PREPARE_CLEAN=0
GO_TPC_LOG=./go-tpc-log.txt

function clean_up {
  [ -n "$pid1" ] && kill $pid1
  [ -n "$pid2" ] && kill $pid2
}
function control_c {
  echo -en "\n## Caught SIGINT; Clean up and Exit \n"
  clean_up
  exit $?
}

trap control_c SIGINT
trap control_c SIGTERM

if [ -f "$CONFIG_FILE" ]; then
  mv $CONFIG_FILE $CONFIG_FILE.backup
  echo "backuped $CONFIG_FILE to $CONFIG_FILE.backup"
fi
echo "start-after: 0s" > "$CONFIG_FILE"
echo "interval: 10s" >> "$CONFIG_FILE"
echo "rules:    # 对 prometheus api 的 query" >> "$CONFIG_FILE"
echo "    - tag: tps" >> "$CONFIG_FILE"
echo "      promql: sum(rate(tidb_session_transaction_duration_seconds_count[1m])) > bool 2/3 * sum(rate(tidb_session_transaction_duration_seconds_count[5m]))" >> "$CONFIG_FILE"

echo ""
echo "Make sure tiup playground is up before we start. Run:"
echo "\$ tiup playground"
echo "Press enter to continue..."
read

if [ $PREPARE_CLEAN == 1 ]; then
  echo "Prepare go-tpc"
  $GO_TPC tpcc --warehouses 4 --parts 4 prepare
fi

$GO_TPC --threads 6 tpcc --warehouses 4 run > $GO_TPC_LOG &
pid1=$!

(sleep 5m; kill $pid1; pid1=""; $GO_TPC --threads 1 tpcc --warehouses 4 run >> $GO_TPC_LOG & echo "Switch thread from 6 to 1" ) &
pid2=$!

$CHECKER

if [ $PREPARE_CLEAN == 1 ]; then
  echo "Cleanup go-tpc"
  $GO_TPC tpcc --warehouses 4 cleanup
fi

clean_up
