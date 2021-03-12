package metrics

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/grafana-tools/sdk"
)

const (
	dashboardJSON = `
{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "editable": true,
  "gnetId": null,
  "graphTooltip": 0,
  "links": [],
  "panels": [],
  "schemaVersion": 18,
  "style": "dark",
  "tags": [
    "templated"
  ],
  "templating": {
    "list": []
  },
  "time": {
    "from": "now-6h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": [
      "5s",
      "10s",
      "30s",
      "1m",
      "5m",
      "15m",
      "30m",
      "1h",
      "2h",
      "1d"
    ],
    "time_options": [
      "5m",
      "15m",
      "1h",
      "6h",
      "12h",
      "24h",
      "2d",
      "7d",
      "30d"
    ]
  },
  "timezone": "browser",
  "title": "Production Overview2",
  "version": 2
}
	`
	panelJSON = `
{
  "aliasColors": {},
  "bars": false,
  "dashLength": 10,
  "dashes": false,
  "datasource": "playground",
  "fill": 1,
  "gridPos": {
    "h": 9,
    "w": 12,
    "x": 0,
    "y": 0
  },
  "id": 2,
  "legend": {
    "avg": false,
    "current": false,
    "max": false,
    "min": false,
    "show": true,
    "total": false,
    "values": false
  },
  "lines": true,
  "linewidth": 1,
  "links": [],
  "nullPointMode": "null",
  "percentage": false,
  "pointradius": 2,
  "points": false,
  "renderer": "flot",
  "seriesOverrides": [],
  "spaceLength": 10,
  "stack": false,
  "steppedLine": false,
  "targets": [
    {
      "expr": "sum(rate(tidb_session_transaction_duration_seconds_count[10m]))",
      "format": "time_series",
      "intervalFactor": 1,
      "refId": "A"
    }
  ],
  "thresholds": [],
  "timeFrom": null,
  "timeRegions": [],
  "timeShift": null,
  "title": "Panel Title",
  "tooltip": {
    "shared": true,
    "sort": 0,
    "value_type": "individual"
  },
  "type": "graph",
  "xaxis": {
    "buckets": null,
    "mode": "time",
    "name": null,
    "show": true,
    "values": []
  },
  "yaxes": [
    {
      "format": "short",
      "label": null,
      "logBase": 1,
      "max": null,
      "min": null,
      "show": true
    },
    {
      "format": "short",
      "label": null,
      "logBase": 1,
      "max": null,
      "min": null,
      "show": true
    }
  ],
  "yaxis": {
    "align": false,
    "alignLevel": null
  }
}
	`
)

// CreateMetricsDashboard create a grafana dashboard on given api URL.
func CreateMetricsDashboard(apiURL, dashboardName string, datasource string, metrics map[string]string) error {
	client := sdk.NewClient(apiURL, "admin:admin", http.DefaultClient)
	ctx := context.Background()
	var board sdk.Board
	err := json.Unmarshal([]byte(dashboardJSON), &board)
	if err != nil {
		return err
	}
	board.Title = dashboardName
	var id uint = 0
	for title, expr := range metrics {
		var panel sdk.Panel
		err = json.Unmarshal([]byte(panelJSON), &panel)
		if err != nil {
			return err
		}
		panel.ID = id
		id++
		panel.CommonPanel.Title = title
		panel.GraphPanel.Targets[0].Expr = expr
		panel.Datasource = &datasource
		board.Panels = append(board.Panels, &panel)
	}

	log.Printf("Mark")
	_, err = client.SetDashboard(ctx, board, sdk.SetDashboardParams{
		FolderID:  0,
		Overwrite: true,
	})
	if err != nil {
		return err
	}
	return nil
}

func panelFromFile(filePath string) (sdk.Panel, error) {
	pFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return sdk.Panel{}, err
	}
	var p sdk.Panel
	err = json.Unmarshal(pFile, &p)
	if err != nil {
		return sdk.Panel{}, err
	}
	return p, nil
}
