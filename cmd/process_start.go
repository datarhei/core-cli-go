package cmd

import (
	"fmt"

	coreclient "github.com/datarhei/core-client-go/v16"

	"github.com/spf13/cobra"
)

// startProcessCmd represents the show command
var startProcessCmd = &cobra.Command{
	Use:   "start [processid]",
	Short: "Start the process with the given ID",
	Long:  "Start the process with the given ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid := args[0]

		client, err := connectSelectedCore()
		if err != nil {
			return err
		}

		id := coreclient.ParseProcessID(pid)

		if err := client.ProcessCommand(id, "start"); err != nil {
			return err
		}

		fmt.Printf("%s start\n", id)

		return nil
	},
}

func init() {
	processCmd.AddCommand(startProcessCmd)
}
