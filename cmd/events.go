package cmd

import (
	"github.com/spf13/cobra"
)

// eventsCmd represents the metrics command
var eventsCmd = &cobra.Command{
	Use:   "events ",
	Short: "Events related commands",
	Long:  "Events related commands",
}

func init() {
	rootCmd.AddCommand(eventsCmd)

	eventsCmd.PersistentFlags().Bool("render", false, "Display nicely rendered event, if available")
}
