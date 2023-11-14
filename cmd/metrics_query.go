package cmd

import (
	"fmt"
	"slices"

	"github.com/datarhei/core-client-go/v16/api"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var metricsQueryCmd = &cobra.Command{
	Use:   "query [name] ...",
	Short: "Query one or more metrics",
	Long:  "Query one or more metrics",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		query := api.MetricsQuery{
			Metrics: []api.MetricsQueryMetric{},
		}

		for _, name := range args {
			query.Metrics = append(query.Metrics, api.MetricsQueryMetric{
				Name:   name,
				Labels: map[string]string{},
			})
		}

		resp, err := client.Metrics(query)
		if err != nil {
			return fmt.Errorf("querying metrics failed: %w", err)
		}

		metrics := map[string][]api.MetricsResponseMetric{}

		for _, m := range resp.Metrics {
			metric, ok := metrics[m.Name]
			if !ok {
				metric = []api.MetricsResponseMetric{}
			}

			metric = append(metric, m)
			metrics[m.Name] = metric
		}

		metricNames := []string{}

		for name := range metrics {
			metricNames = append(metricNames, name)
		}

		slices.Sort(metricNames)

		for _, name := range metricNames {
			metric := metrics[name]

			t := table.NewWriter()

			labels := []string{}
			for label := range metric[0].Labels {
				labels = append(labels, label)
			}
			slices.Sort(labels)

			row := table.Row{"Metric"}
			for _, label := range labels {
				row = append(row, label)
			}
			row = append(row, "Value")
			t.AppendHeader(row)

			for _, m := range metric {
				row := table.Row{m.Name}
				for _, label := range labels {
					row = append(row, m.Labels[label])
				}
				row = append(row, m.Values[0].Value)
				t.AppendRow(row)
			}

			sort := []table.SortBy{}

			for i := 0; i < (len(labels) + 2); i++ {
				sort = append(sort, table.SortBy{Number: i + 1, Mode: table.Asc})
			}

			t.SortBy(sort)

			t.SetStyle(table.StyleLight)

			fmt.Println(t.Render())
		}

		return nil

	},
}

func init() {
	metricsCmd.AddCommand(metricsQueryCmd)
}
