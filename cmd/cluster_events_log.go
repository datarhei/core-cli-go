package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

var clusterEventsLogCmd = &cobra.Command{
	Use:   "log [component [key=value] [key=value] ...] ; [component [key=value] ...]",
	Short: "Retrieve events",
	Long:  "Retrieve events",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		filters := api.LogEventFilters{
			Filters: []api.LogEventFilter{},
		}

		for {
			var filter api.LogEventFilter
			var done bool

			args, filter, done = parseLogFilter(args)
			if done {
				break
			}

			filters.Filters = append(filters.Filters, filter)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		events, err := client.ClusterLogEvents(ctx, filters)
		if err != nil {
			return err
		}

		go func(ctx context.Context, events <-chan api.LogEvent) {
			for {
				select {
				case event, ok := <-events:
					if !ok {
						fmt.Printf("no more events, channel has been closed\n")
						cancel()
						return
					}
					writeJSON(os.Stdout, event, true)
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

		fmt.Printf("event stream finished\n")

		return nil
	},
}

func init() {
	clusterEventsCmd.AddCommand(clusterEventsLogCmd)
}
