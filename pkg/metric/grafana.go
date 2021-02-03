package metric

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/grafana-tools/sdk"
)

func CreateMetricsDashboard(apiURL, dashboardName string, metrics map[string]string) error {
	client := sdk.NewClient(apiURL, "admin:admin", http.DefaultClient)
	ctx := context.Background()
	boardFile, err := ioutil.ReadFile("templates/dashboard.json")
	if err != nil {
		return err
	}
	var board sdk.Board
	err = json.Unmarshal(boardFile, &board)
	if err != nil {
		return err
	}
	board.Title = dashboardName
	var id uint = 0
	for title, expr := range metrics {
		panel, err := panelFromFile("templates/panel.json")
		if err != nil {
			return err
		}
		panel.ID = id
		id++
		panel.CommonPanel.Title = title
		panel.GraphPanel.Targets[0].Expr = expr
		board.Panels = append(board.Panels, &panel)
	}
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
