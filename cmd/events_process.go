package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

var eventsProcessCmd = &cobra.Command{
	Use:   "process [key=value] ... ; [key=value] ... ; ...",
	Short: "Process events",
	Long:  "Process events",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		render, _ := cmd.Flags().GetBool("render")

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

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		events, err := client.ProcessEvents(ctx, filters)
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
					if render {
						fmt.Println(renderProcessEvent(event))
					} else {
						writeJSON(os.Stdout, event, true)
					}
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

func parseProcessFilter(args []string) ([]string, api.ProcessEventFilter, bool) {
	if len(args) == 0 {
		return nil, api.ProcessEventFilter{}, true
	}

	filter := api.ProcessEventFilter{}

	for i, arg := range args {
		if arg == ";" {
			return args[i+1:], filter, false
		}

		key, value, found := strings.Cut(arg, "=")
		if !found {
			continue
		}

		if key == "pid" {
			filter.ProcessID = value
			continue
		}

		if key == "domain" {
			filter.Domain = value
			continue
		}

		if key == "type" {
			filter.Type = value
			continue
		}

		if key == "core_id" {
			filter.CoreID = value
			continue
		}
	}

	return nil, filter, false
}

func renderProcessEvent(event api.ProcessEvent) string {
	switch event.Type {
	case "line":
		return fmt.Sprintf("[%s@%s] %s", event.ProcessID, event.Domain, event.Line)
	case "progress":
		return fmt.Sprintf("[%s@%s] %v", event.ProcessID, event.Domain, event.Progress)
	}

	return ""
}

func init() {
	eventsCmd.AddCommand(eventsProcessCmd)
}
