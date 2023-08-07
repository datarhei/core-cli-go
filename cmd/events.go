package cmd

import (
	"context"
	"os"
	"os/signal"
	"strings"

	"github.com/datarhei/core-client-go/v16/api"

	"github.com/spf13/cobra"
)

// eventsCmd represents the metrics command
var eventsCmd = &cobra.Command{
	Use:   "events",
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

		events, err := client.Events(ctx, filters)
		if err != nil {
			return err
		}

		go func(ctx context.Context, events <-chan api.Event) {
			for {
				select {
				case event := <-events:
					writeJSON(os.Stdout, event, true)
				case <-ctx.Done():
					return
				}
			}
		}(ctx, events)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		return nil
	},
}

func parseFilter(args []string) ([]string, api.EventFilter, bool) {
	if len(args) == 0 {
		return nil, api.EventFilter{}, true
	}

	filter := api.EventFilter{
		Component: args[0],
		Data:      map[string]string{},
	}

	for i, arg := range args[1:] {
		if arg == ";" {
			return args[i+1:], filter, false
		}

		key, value, found := strings.Cut(arg, "=")
		if !found {
			continue
		}

		if key == "level" {
			filter.Level = value
			continue
		}

		if key == "message" {
			filter.Message = value
			continue
		}

		filter.Data[key] = value
	}

	return nil, filter, false
}

func init() {
	rootCmd.AddCommand(eventsCmd)
}
