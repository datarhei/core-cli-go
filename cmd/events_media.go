package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/datarhei/core-client-go/v16/api"
	"github.com/spf13/cobra"
)

var eventsMediaCmd = &cobra.Command{
	Use:   "media [fsname] [pattern]?",
	Short: "Media events",
	Long:  "Media events",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		storage := args[0]
		pattern := ""
		if len(args) == 2 {
			pattern = args[1]
		}

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		events, err := client.MediaEvents(ctx, storage, pattern)
		if err != nil {
			return err
		}

		go func(ctx context.Context, events <-chan api.MediaEvent) {
			for {
				select {
				case event, ok := <-events:
					if !ok {
						fmt.Printf("no more events, channel has been closed\n")
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
	eventsCmd.AddCommand(eventsMediaCmd)
}
