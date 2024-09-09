package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

var clusterEventsCmd = &cobra.Command{
	Use:   "events [component [key=value] [key=value] ...] ; [component [key=value] ...]",
	Short: "Retrieve events",
	Long:  "Retrieve events",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		filters := api.EventFilters{
			Filters: []api.EventFilter{},
		}

		for {
			var filter api.EventFilter
			var done bool

			args, filter, done = parseFilter(args)
			if done {
				break
			}

			filters.Filters = append(filters.Filters, filter)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		events, err := client.ClusterEvents(ctx, filters)
		if err != nil {
			return err
		}

		go func(ctx context.Context, events <-chan api.Event) {
			for {
				select {
				case event, ok := <-events:
					if !ok {
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

		return nil
	},
}

func init() {
	clusterCmd.AddCommand(clusterEventsCmd)
}
