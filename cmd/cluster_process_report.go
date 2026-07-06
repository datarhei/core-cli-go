package cmd

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"time"

	coreclient "github.com/datarhei/core-client-go/v16"
	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

var clusterProcessReportCmd = &cobra.Command{
	Use:   "report [-f|--follow]? [processid]",
	Short: "Show the report of the process with the given ID",
	Long:  "Show the report of the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		follow, _ := cmd.Flags().GetBool("follow")
		created, _ := cmd.Flags().GetString("created")
		exited, _ := cmd.Flags().GetString("exited")
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		if follow {
			filters := api.ProcessEventFilters{
				Filters: []api.ProcessEventFilter{
					{
						ProcessID: "^" + id.ID + "$",
						Domain:    "^" + id.Domain + "$",
						Type:      "line",
					},
				},
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			events, err := client.ClusterProcessEvents(ctx, filters)
			if err != nil {
				return err
			}

			start := time.Now()
			nevents := uint64(0)

			go func(ctx context.Context, events <-chan api.ProcessEvent) {
				for {
					select {
					case event, ok := <-events:
						if !ok {
							fmt.Printf("no more events, channel has been closed\n")
							cancel()
							return
						}
						fmt.Fprintf(os.Stdout, "%s\n", event.Line)
						nevents++
					case <-ctx.Done():
						return
					}
				}
			}(ctx, events)

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)

			select {
			case <-quit:
				cancel()
			case <-ctx.Done():
			}

			fmt.Printf("log stream finished (%.2f events/s)\n", float64(nevents)/time.Since(start).Seconds())
		} else {
			createdAt := int64(math.MinInt64)
			exitedAt := int64(math.MinInt64)

			if len(created) != 0 {
				t, err := time.Parse(time.DateTime, created)
				if err != nil {
					return fmt.Errorf("parsing created time: %w", err)
				}

				createdAt = t.Unix()
			}

			if len(exited) != 0 {
				t, err := time.Parse(time.DateTime, exited)
				if err != nil {
					return fmt.Errorf("parsing exited time: %w", err)
				}

				exitedAt = t.Unix()
			}

			report, err := client.ClusterProcessReport(id, createdAt, exitedAt)
			if err != nil {
				return err
			}

			if err := writeJSON(os.Stdout, report, true); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	clusterProcessCmd.AddCommand(clusterProcessReportCmd)

	clusterProcessReportCmd.Flags().BoolP("follow", "f", false, "Follow the current report")
	clusterProcessReportCmd.Flags().String("created", "", fmt.Sprintf("Show only reports that have been created at this time, format: %s", time.DateTime))
	clusterProcessReportCmd.Flags().String("exited", "", fmt.Sprintf("Show only report that have exited at this time, format: %s", time.DateTime))
}
