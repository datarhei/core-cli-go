package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

var clusterEventsProcessCmd = &cobra.Command{
	Use:   "process [key=value] ... ; [key=value] ... ; ...",
	Short: "Process events",
	Long:  "Process events",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		filters := api.ProcessEventFilters{
			Filters: []api.ProcessEventFilter{},
		}

		for {
			var filter api.ProcessEventFilter
			var done bool

			args, filter, done = parseProcessFilter(args)
			if done {
				break
			}

			filters.Filters = append(filters.Filters, filter)
		}

		fmt.Printf("%+v\n", filters)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		events, err := client.ClusterProcessEvents(ctx, filters)
		if err != nil {
			return err
		}

		go func(ctx context.Context, events <-chan api.ProcessEvent) {
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
	clusterEventsCmd.AddCommand(clusterEventsProcessCmd)
}
